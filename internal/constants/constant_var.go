package constants

type ResponseReturn struct {
	Code    string
	Message string
}

const (
	TransferSuccess        = "2004300"
	PendingTransfer        = "2024300"
	ErrBadRequest          = "4004300"
	ErrInsufficientFunds   = "4034314"
	ErrDataNotFound        = "4044301"
	ErrBalanceNotAvailable = "4044316"
	ErrInternalServerError = "5004301"
	ErrExternalServerError = "5004302"
	ErrTransferTimeout     = "5044300"
)

var ResponseMap = map[string]string{
	TransferSuccess:        "Successful",
	PendingTransfer:        "Transaction is being processed",
	ErrBadRequest:          "Invalid request",
	ErrInsufficientFunds:   "Insufficient funds",
	ErrDataNotFound:        "Data not found",
	ErrBalanceNotAvailable: "Merchant balance not found",
	ErrInternalServerError: "Internal server error",
	ErrExternalServerError: "External server error",
	ErrTransferTimeout:     "Timeout",
}

const (
	StatusCredit = "CR"
	StatusDebit  = "DB"
)

const (
	StatusFailedPublish = "FAILED_PUBLISH"
	StatusDone          = "DONE"
	StatusPending       = "PENDING"
	StatusRejected      = "REJECTED"
	StatusInProgress    = "PROGRESSING"
)
