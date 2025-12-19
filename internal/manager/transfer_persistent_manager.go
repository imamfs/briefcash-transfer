package manager

import (
	"briefcash-transfer/internal/constants"
	"briefcash-transfer/internal/dto"
	"briefcash-transfer/internal/entity"
	"briefcash-transfer/internal/repository"
	"context"
	"fmt"
	"time"
)

type AccountStatementManager interface {
	CreateRecipient(ctx context.Context, request dto.TransferRequest) (*entity.DataRecipient, error)
	CreateTransfer(ctx context.Context, recipient *entity.DataRecipient, request dto.TransferRequest, adminFee entity.FeeSettings, partnerId, referenceNumber string, amountTransfer float64) (*entity.Transaction, error)
	FindTransferByPartnerReference(ctx context.Context, partnerReferenceNo string) (*entity.TransferTemp, error)
	DebitMerchant(ctx context.Context, merchantCode string, totalAmount float64) (float64, error)
	CreditMerchant(ctx context.Context, merchantCode string, amount float64) (float64, error)
	CreateTransferLedger(ctx context.Context, transferId int64, request dto.TransferRequest, balance float64, merchantCode string, amountTransfer float64) error
	CreateAdminFeeLedger(ctx context.Context, transferId int64, request dto.TransferRequest, balance float64, feeSetting entity.FeeSettings, merchantCode string) error
	CreateRefundLedger(ctx context.Context, transferId int64, request dto.TransferRequest, refundAmount, balance float64, merchantCode string) error
	UpdateTransferStatus(ctx context.Context, transferId int64, status string) error
}

type transferPersistenceService struct {
	ledgerRepo    repository.LedgerRepository
	transferRepo  repository.TransferRepository
	recipientRepo repository.RecipientRepository
	merchantRepo  repository.BalanceRepository
}

func NewTransferPersistenceManager(ledgerRepo repository.LedgerRepository, transferRepo repository.TransferRepository,
	recipientRepo repository.RecipientRepository, merchantRepo repository.BalanceRepository) AccountStatementManager {
	return &transferPersistenceService{ledgerRepo, transferRepo, recipientRepo, merchantRepo}

}

// data recipient
func (tp *transferPersistenceService) CreateRecipient(ctx context.Context, request dto.TransferRequest) (*entity.DataRecipient, error) {
	recipient := &entity.DataRecipient{
		AccountNumber: request.BeneficiaryAccountNumber,
		BankCode:      request.BeneficiaryBankCode,
	}

	if err := tp.recipientRepo.Save(ctx, recipient); err != nil {
		return nil, err
	}

	return recipient, nil
}

// data transfer
func (tp *transferPersistenceService) CreateTransfer(ctx context.Context, recipient *entity.DataRecipient, request dto.TransferRequest, adminFee entity.FeeSettings, partnerId, referenceNumber string, amountTransfer float64) (*entity.Transaction, error) {
	transfer := &entity.Transaction{
		MerchantCode:            partnerId,
		PartnerReferenceNo:      request.PartnerReferenceNo,
		BankReferenceNo:         nil,
		SystemReferenceNo:       &referenceNumber,
		Amount:                  amountTransfer,
		Currency:                "IDR",
		Remark:                  request.AdditionalInfo.Remarks,
		TransactionType:         request.AdditionalInfo.Channel,
		TransactionDate:         time.Now(),
		Status:                  constants.ResponseMap[constants.PendingTransfer],
		IsReversal:              false,
		IsReconcile:             false,
		ReconcileDate:           nil,
		CompanyCharge:           float32(adminFee.FeeService),
		PartnerCharge:           float32(adminFee.FeePartner),
		AdditionalPartnerCharge: float32(adminFee.AdditionalFee),
		TaxCharge:               float32(adminFee.FeeTax),
		Recipient:               recipient.ID,
	}

	if err := tp.transferRepo.Save(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (tp *transferPersistenceService) FindTransferByPartnerReference(ctx context.Context, partnerReferenceNo string) (*entity.TransferTemp, error) {
	transfer, err := tp.transferRepo.FindByRefNo(ctx, partnerReferenceNo)
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (tp *transferPersistenceService) UpdateTransferStatus(ctx context.Context, id int64, status string) error {
	return tp.transferRepo.Update(ctx, id, status)
}

// data merchant account
func (tp *transferPersistenceService) DebitMerchant(ctx context.Context, partnerId string, totalAmount float64) (float64, error) {
	newBalance, err := tp.merchantRepo.Debit(ctx, partnerId, totalAmount)
	if err != nil {
		return 0, err
	}
	return newBalance, err
}

func (tp *transferPersistenceService) CreditMerchant(ctx context.Context, merchantCode string, amount float64) (float64, error) {
	newBalance, err := tp.merchantRepo.Credit(ctx, merchantCode, amount)
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}

// Data account statement
func (tp *transferPersistenceService) CreateTransferLedger(ctx context.Context, transferId int64, request dto.TransferRequest, balance float64, merchantCode string, amountTransfer float64) error {
	statement := tp.buildStatement(transferId, request, balance, -amountTransfer, request.AdditionalInfo.Remarks, constants.StatusDebit, merchantCode)
	return tp.ledgerRepo.Save(ctx, statement)
}

func (tp *transferPersistenceService) CreateAdminFeeLedger(ctx context.Context, transferId int64, request dto.TransferRequest, balance float64, feeSetting entity.FeeSettings, merchantCode string) error {
	statement := tp.buildStatement(transferId, request, balance, -feeSetting.TotalCharge, "Service Fee", constants.StatusDebit, merchantCode)
	return tp.ledgerRepo.Save(ctx, statement)
}

func (tp *transferPersistenceService) CreateRefundLedger(ctx context.Context, transferId int64, request dto.TransferRequest, refundAmount, balance float64, merchantCode string) error {
	if refundAmount < 0 {
		refundAmount = -refundAmount
	}

	statement := tp.buildStatement(transferId, request, balance, refundAmount, fmt.Sprintf("Refund fee: system error for ref: %s", request.PartnerReferenceNo), constants.StatusCredit, merchantCode)
	return tp.ledgerRepo.Save(ctx, statement)
}

func (tp *transferPersistenceService) buildStatement(transferId int64, requestDto dto.TransferRequest, balance, amount float64, description, status, merchantCode string) *entity.AccountStatement {
	return &entity.AccountStatement{
		TransactionId:       transferId,
		TransctionReference: requestDto.PartnerReferenceNo,
		MerchantCode:        merchantCode,
		Status:              status,
		Channel:             requestDto.AdditionalInfo.Channel,
		Description:         description,
		Amount:              amount,
		BalanceAfter:        balance,
		CreatedAt:           time.Now(),
	}
}
