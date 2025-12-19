package controller

import (
	"briefcash-transfer/internal/constants"
	"briefcash-transfer/internal/dto"
	"briefcash-transfer/internal/helper/loghelper"
	"briefcash-transfer/internal/helper/timehelper"
	"briefcash-transfer/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type transferController struct {
	svc service.TransferService
}

func NewTransferController(svc service.TransferService) *transferController {
	return &transferController{svc}
}

func (t *transferController) Transfer(ctx *gin.Context) {
	start := time.Now()

	var request dto.TransferRequest
	externalId := ctx.GetHeader("X-EXTERNAL-ID")
	merchantCode := ctx.GetHeader("X-PARTNER-ID")

	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":     "transfer_controller",
		"trace_id":    externalId,
		"merchant_id": merchantCode,
	})

	defer func() {
		log.WithField("processing_time", time.Since(start).Milliseconds())
	}()

	log.Info("Parsing payload request")
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.TransferResponse{
			ResponseCode:       constants.ErrBadRequest,
			ResponseMessage:    constants.ResponseMap[constants.ErrBadRequest],
			ReferenceNumber:    "",
			PartnerReferenceNo: "",
			TransactionDate:    timehelper.FormatTimeToISO7(time.Now()),
		})
		return
	}

	response := t.svc.TransferRequest(ctx, request, merchantCode, externalId)

	httpStatus := map[string]int{
		constants.ErrDataNotFound:        http.StatusNotFound,
		constants.ErrInsufficientFunds:   http.StatusForbidden,
		constants.ErrInternalServerError: http.StatusInternalServerError,
		constants.PendingTransfer:        http.StatusAccepted,
	}

	log.Info("Populate response")
	ctx.JSON(httpStatus[response.ResponseCode], response)
}
