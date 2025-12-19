package service

import (
	"briefcash-transfer/internal/constants"
	"briefcash-transfer/internal/dto"
	"briefcash-transfer/internal/entity"
	"briefcash-transfer/internal/helper/kafkahelper"
	"briefcash-transfer/internal/helper/loghelper"
	"briefcash-transfer/internal/helper/timehelper"
	"briefcash-transfer/internal/manager"
	"briefcash-transfer/internal/protobuf"
	"briefcash-transfer/internal/repository"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type TransferService interface {
	TransferRequest(ctx context.Context, request dto.TransferRequest, merchantCode, externalId string) dto.TransferResponse
}

type transferService struct {
	recipientRepo  repository.RecipientRepository
	transferRepo   repository.TransferRepository
	feeSettingRepo repository.FeeSettingRepository
	ledgerRepo     repository.LedgerRepository
	merchantRepo   repository.BalanceRepository
	redisService   TransferRedisService
	partnerService BankPartner
	db             *gorm.DB
	kafkaProducer  *kafkahelper.KafkaProducer
}

func NewTransferService(recipientRepo repository.RecipientRepository, transferRepo repository.TransferRepository, feeSettingRepo repository.FeeSettingRepository,
	ledgerRepo repository.LedgerRepository, merchantRepo repository.BalanceRepository, redisService TransferRedisService,
	partnerService BankPartner, db *gorm.DB, kafkaProducer *kafkahelper.KafkaProducer) TransferService {
	return &transferService{recipientRepo, transferRepo, feeSettingRepo, ledgerRepo, merchantRepo, redisService, partnerService, db, kafkaProducer}
}

func (t *transferService) TransferRequest(ctx context.Context, request dto.TransferRequest, merchantCode, externalId string) dto.TransferResponse {
	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":   "transfer_service",
		"operation": "initiate_request",
		"bank_code": request.BeneficiaryBankCode,
		"trace_id":  externalId,
		"merchant":  merchantCode,
	})

	// get fee service charge from redis
	log.Info("Get fee setting configuration from redis")
	feeSetting, err := t.redisService.GetFeeSetting(ctx, merchantCode, request.AdditionalInfo.Channel, log)
	if err != nil {
		// fallback to databse if fee service charge not found in redis
		log.Warn("Fee not found in redis, fallback to DB")
		feeSettingDb, err := t.feeSettingRepo.FindByCodeAndChannel(ctx, merchantCode, request.AdditionalInfo.Channel)
		if err != nil {
			log.WithError(err).Error("Fee not found in redis and DB")
			return t.handleTransferResponse(constants.ErrDataNotFound, constants.ResponseMap[constants.ErrInternalServerError], "", request.PartnerReferenceNo, "0", &feeSettingDb)
		}

		// if fee service charge found in database, create goroutine to save data back into redis
		log.Infof("Fee setting found in database for merchant: %s and channel %s, processing cache in redis", feeSettingDb.MerchantCode, feeSettingDb.Channel)
		errorChannel := make(chan error, 1)
		var waitGroup sync.WaitGroup
		waitGroup.Go(func() {
			if err := t.redisService.SetFeeSetting(ctx, feeSettingDb, log); err != nil {
				errorChannel <- err
			}
		})

		waitGroup.Wait()
		close(errorChannel)

		if errUpdate := <-errorChannel; errUpdate != nil {
			log.WithError(<-errorChannel).Warn("Failed to cache fee setting in redis")
		}

		// assign fee service charge found in database to variabel feeSetting
		feeSetting = feeSettingDb
	}

	// Total amount transfer and service fee
	totalAmount := t.sumAmount(feeSetting, request.Amount.Value)

	// subtract balance in redis
	log.Info("Debit merchant balance in redis")
	balance, err := t.redisService.DebitBalance(ctx, merchantCode, totalAmount, log)

	if balance == 0 {
		return t.handleTransferResponse(constants.ErrBalanceNotAvailable, constants.ResponseMap[constants.ErrBalanceNotAvailable], "", request.PartnerReferenceNo, "0", &feeSetting)
	}

	if balance == -1 {
		return t.handleTransferResponse(constants.ErrInsufficientFunds, constants.ResponseMap[constants.ErrInsufficientFunds], "", request.PartnerReferenceNo, "0", &feeSetting)
	}

	if err != nil {
		return t.handleTransferResponse(constants.ErrInternalServerError, constants.ResponseMap[constants.ErrInternalServerError], "", request.PartnerReferenceNo, "0", &feeSetting)
	}

	// save transfer and account statement into database
	log.Info("Persist transfer, ledger, and updated balance to database")
	referenceNumber := t.generatedReferenceNumber(request)
	if err := t.PersistTransfer(ctx, request, feeSetting, merchantCode, referenceNumber, balance, totalAmount); err != nil {
		log.Warn("Persist failed, refund merchant balance in redis")
		if err := t.redisService.RefundBalance(ctx, merchantCode, totalAmount, log); err != nil {
			return t.handleTransferResponse(constants.ErrInternalServerError, constants.ResponseMap[constants.ErrInternalServerError], "", request.PartnerReferenceNo, "0", &feeSetting)
		}
		log.WithError(err).Error("Failed to persist transfer into database")
		return t.handleTransferResponse(constants.ErrInternalServerError, constants.ResponseMap[constants.ErrInternalServerError], "", request.PartnerReferenceNo, "0", &feeSetting)
	}

	// publish trigger transfer to kafka
	log.Info("Publish trigger transfer to kafka")
	if err := t.publishMessage(request, externalId, log); err != nil {
		log.Warn("Trigger transfer failed, refund merchant balance in redis")
		if err := t.redisService.RefundBalance(ctx, merchantCode, totalAmount, log); err != nil {
			return t.handleTransferResponse(constants.ErrInternalServerError, constants.ResponseMap[constants.ErrInternalServerError], "", request.PartnerReferenceNo, "0", &feeSetting)
		}

		log.Warn("Trigger transfer failed, update status transfer and refund balance in database")
		if err := t.handlePublishFailure(ctx, merchantCode, request, totalAmount); err != nil {
			return t.handleTransferResponse(constants.ErrInternalServerError, constants.ResponseMap[constants.ErrInternalServerError], "", request.PartnerReferenceNo, "0", &feeSetting)
		}

		log.WithError(err).Error("Failed publish trigger transfer to kafka")
		return t.handleTransferResponse(constants.ErrInternalServerError, constants.ResponseMap[constants.ErrInternalServerError], "", request.PartnerReferenceNo, "0", &feeSetting)
	}

	// return response to handler
	remainingBalance := strconv.FormatFloat(balance, 'f', 2, 64)
	return t.handleTransferResponse(constants.PendingTransfer, constants.ResponseMap[constants.PendingTransfer], referenceNumber, request.PartnerReferenceNo, remainingBalance, &feeSetting)
}

func (t *transferService) PersistTransfer(ctx context.Context, request dto.TransferRequest, feeCharge entity.FeeSettings, merchantCode, referenceNumber string, balance, totalAmount float64) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		transferTx := t.transferRepo.WithTransaction(tx)
		ledgerTx := t.ledgerRepo.WithTransaction(tx)
		recipientTx := t.recipientRepo.WithTransaction(tx)
		accountTx := t.merchantRepo.WithTransaction(tx)

		amountTransfer, err := parseAmount(request.Amount.Value)
		if err != nil {
			return err
		}

		pm := manager.NewTransferPersistenceManager(ledgerTx, transferTx, recipientTx, accountTx)

		// save recipient
		recipient, err := pm.CreateRecipient(ctx, request)
		if err != nil {
			return err
		}

		// save transfer
		transfer, err := pm.CreateTransfer(ctx, recipient, request, feeCharge, merchantCode, referenceNumber, amountTransfer)
		if err != nil {
			return err
		}

		// deduct balance
		newBalance, err := pm.DebitMerchant(ctx, merchantCode, totalAmount)
		if err != nil {
			return err
		}

		// save history transaction
		if err := pm.CreateTransferLedger(ctx, transfer.ID, request, newBalance, merchantCode, amountTransfer); err != nil {
			return err
		}

		// save history service fee
		if err := pm.CreateAdminFeeLedger(ctx, transfer.ID, request, newBalance, feeCharge, merchantCode); err != nil {
			return err
		}

		return nil
	})
}

func (t *transferService) handleTransferResponse(responseCode, responseMessage, referenceNumber, partnerReferenceNo string, balanceAfter string, fee *entity.FeeSettings) dto.TransferResponse {
	additionalInfo := map[string]string{}
	if fee != nil {
		additionalInfo["channel"] = fee.Channel
		additionalInfo["service_fee"] = strconv.FormatFloat(fee.TotalCharge, 'f', 2, 64)
		additionalInfo["balance_after"] = balanceAfter
	} else {
		additionalInfo = map[string]string{}
	}
	return dto.TransferResponse{
		ResponseCode:       responseCode,
		ResponseMessage:    responseMessage,
		ReferenceNumber:    referenceNumber,
		PartnerReferenceNo: partnerReferenceNo,
		TransactionDate:    timehelper.FormatTimeToISO7(time.Now()),
		AdditionalInfo:     additionalInfo,
	}
}

func (t *transferService) sumAmount(fee entity.FeeSettings, amount string) float64 {
	amountSend, _ := parseAmount(amount)
	totalCharge := fee.TotalCharge + amountSend
	return totalCharge
}

func (t *transferService) publishMessage(request dto.TransferRequest, externalId string, log *logrus.Entry) error {
	payload := &protobuf.TransferRequest{
		ExternalId:           externalId,
		PartnerRefNo:         request.PartnerReferenceNo,
		CustomerNumber:       request.CustomerNumber,
		AccountType:          request.AccountType,
		BeneficiaryAccountNo: request.BeneficiaryAccountNumber,
		BeneficiaryBankCode:  request.BeneficiaryBankCode,
		Amount:               request.Amount.Value,
		TransactionDate:      request.AdditionalInfo.TransactionDate,
		CustomerReference:    request.AdditionalInfo.CustomerReference,
		Channel:              request.AdditionalInfo.Channel,
		Remarks:              request.AdditionalInfo.Remarks,
		Email:                request.AdditionalInfo.Email,
		Address:              request.AdditionalInfo.Address,
		Citizenship:          request.AdditionalInfo.Citizenship,
		TransferPurpose:      request.AdditionalInfo.TransferPurpose,
		TransferActivity:     request.AdditionalInfo.TransferActivity,
		CustomerType:         request.AdditionalInfo.CustomerType,
	}

	protoBytes, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshal protobuf: %w", err)
	}

	topic := t.partnerService.GetBankConfig(request.BeneficiaryBankCode, log)
	key := request.PartnerReferenceNo
	return t.kafkaProducer.Publish(topic.KafkaTopic, key, protoBytes)
}

func (t *transferService) generatedReferenceNumber(request dto.TransferRequest) string {
	var reference string
	channel := request.AdditionalInfo.Channel
	destinationBank := request.BeneficiaryBankCode
	timeLayout := "20060102150405000"

	if len(channel) < 3 {
		reference = channel + "-" + destinationBank + "-" + time.Now().Format(timeLayout)
	} else {
		reference = channel[0:3] + "-" + destinationBank + "-" + time.Now().Format(timeLayout)
	}

	return reference
}

func (r *transferService) handlePublishFailure(ctx context.Context, merchantCode string, request dto.TransferRequest, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		ledgerTx := r.ledgerRepo.WithTransaction(tx)
		accountTx := r.merchantRepo.WithTransaction(tx)
		transferTx := r.transferRepo.WithTransaction(tx)

		pm := manager.NewTransferPersistenceManager(ledgerTx, transferTx, nil, accountTx)

		// get transfer Data
		transfer, err := pm.FindTransferByPartnerReference(ctx, request.PartnerReferenceNo)
		if err != nil {
			return err
		}

		// guard idempotency
		if transfer.Status == constants.StatusFailedPublish {
			return nil
		}

		if transfer.Status == constants.ResponseMap[constants.PendingTransfer] {
			return fmt.Errorf("invalid state transition from %s", transfer.Status)
		}

		// Update transfer status
		if err := pm.UpdateTransferStatus(ctx, transfer.ID, constants.StatusFailedPublish); err != nil {
			return err
		}

		// restore balance
		newBalance, err := pm.CreditMerchant(ctx, merchantCode, amount)
		if err != nil {
			return err
		}

		// refund in ledger with failed status
		if err := pm.CreateRefundLedger(ctx, transfer.ID, request, amount, newBalance, merchantCode); err != nil {
			return err
		}

		return nil
	})
}

func parseAmount(amount string) (float64, error) {
	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %w", err)
	}
	return value, err
}
