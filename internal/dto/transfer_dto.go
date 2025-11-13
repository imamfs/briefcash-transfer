package dto

type TransferRequest struct {
	PartnerReferenceNo       string                 `json:"partnerReferenceNo"`
	CustomerNumber           string                 `json:"customerNumber"` // phone number
	AccountType              string                 `json:"accountType"`
	BeneficiaryAccountNumber string                 `json:"beneficiaryAccountNumber"`
	BeneficiaryBankCode      string                 `json:"beneficiaryBankCode"`
	Amount                   TransferAmountData     `json:"amount"`
	AdditionalInfo           TransferAdditionalInfo `json:"additionalInfo"`
}

type TransferAmountData struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type TransferAdditionalInfo struct {
	TransactionDate  string `json:"transactionDate"`
	Channel          string `json:"channel"` // online, bifast, sknbi, rtgs, va, wallet
	Remarks          string `json:"remarks"`
	Email            string `json:"email"`
	Address          string `json:"address"`
	Citizenship      string `json:"citizenship"` // wna, wni
	TransferPurpose  string `json:"transferPurpose"`
	TransferActivity string `json:"transferActivity"` // only mandatory for non indonesian citizen
	CustomerType     string `json:"customerType"`     // 01 - individu, 02 - corporate, 03 - others
}

type TransferResponse struct {
	ResponseCode       string            `json:"responseCode"`
	ResponseMessage    string            `json:"responseMessage"`
	ReferenceNumber    string            `json:"referenceNumber"`
	PartnerReferenceNo string            `json:"partnerReferenceNo"`
	TransactionDate    string            `json:"transactionDate"`
	AdditionalInfo     map[string]string `json:"additionalInfo"`
}

type BCATransferExternalRequest struct {
	PartnerReferenceNo     string                  `json:"partnerReferenceNo"`
	Amount                 TransferAmountData      `json:"amount"`
	BeneficiaryAccountName string                  `json:"beneficiaryAccountName"`
	BeneficiaryAccountNo   string                  `json:"beneficiaryAccountNo"`
	BeneficiaryBankCode    string                  `json:"beneficiaryBankCode"`
	BeneficiaryEmail       string                  `json:"beneficiaryEmail"`
	SourceAccountNo        string                  `json:"sourceAccountNo"`
	TransactionDate        string                  `json:"transactionDate"`
	AdditionalInfo         BCATransferExternalInfo `json:"additionalInfo"`
}

type BCATransferExternalInfo struct {
	TransferType string `json:"transferType"`
	PurposeCode  string `json:"purposeCode"`
}

type BCATransferExternalResponse struct {
	ResponseCode         string             `json:"responseCode"`
	ResponseMessage      string             `json:"responseMessage"`
	PartnerReferenceNo   string             `json:"partnerReferenceNo"`
	ReferenceNo          string             `json:"referenceNo"`
	Amount               TransferAmountData `json:"amount"`
	BeneficiaryAccountNo string             `json:"beneficiaryAccountNo"`
	BeneficiaryBankCode  string             `json:"beneficiaryBankCode"`
	SourceAccountNo      string             `json:"sourceAccountNo"`
	AdditionalInfo       map[string]string  `json:"additionalInfo"`
}

type BCATransferInternalRequest struct {
	PartnerReferenceNo   string                  `json:"partnerReferenceNo"`
	BeneficiaryEmail     string                  `json:"beneficiaryEmail"`
	Amount               TransferAmountData      `json:"amount"`
	BeneficiaryAccountNo string                  `json:"beneficiaryAccountNo"`
	Remark               string                  `json:"remark"`
	SourceAccountNo      string                  `json:"sourceAccountNo"`
	TransactionDate      string                  `json:"transactionDate"`
	AdditionalInfo       BCATransferInternalInfo `json:"additionalInfo"`
}

type BCATransferInternalInfo struct {
	EconomicActivity   string `json:"economicActivity"`
	TransactionPurpose string `json:"transactionPurpose"`
}

type BCATransferInternalResponse struct {
	ResponseCode         string                  `json:"responseCode"`
	ResponseMessage      string                  `json:"responseMessage"`
	PartnerReferenceNo   string                  `json:"partnerReferenceNo"`
	ReferenceNo          string                  `json:"referenceNo"`
	Amount               TransferAmountData      `json:"amount"`
	BeneficiaryAccountNo string                  `json:"beneficiaryAccountNo"`
	SourceAccountNo      string                  `json:"sourceAccountNo"`
	TransactionDate      string                  `json:"transactionDate"`
	AdditionalInfo       BCATransferInternalInfo `json:"additionalInfo"`
}

type BCATransferToVARequest struct {
	VirtualAccountNo    string `json:"virtualAccountNo"`
	VirtualAccountEmail string `json:"virtualAccountEmail"`
	SourceAccountNo     string `json:"sourceAccountNo"`
	PartnerReferenceNo  string `json:"partnerReferenceNo"`
	PaidAmount          string `json:"paidAmount"`
	TrxDateTime         string `json:"trxDateTime"`
}

type BCATransferToVAResponse struct {
	ResponseCode       string              `json:"responseCode"`
	ResponseMessage    string              `json:"responseMessage"`
	VirtualAccountData BCATransferToVAData `json:"virtualAccountData"`
}

type BCATransferToVAData struct {
	VirtualAccountNo    string             `json:"virtualAccountNo"`
	VirtualAccountName  string             `json:"virtualAccountName"`
	VirtualAccountEmail string             `json:"virtualAccountEmail"`
	SourceAccountNo     string             `json:"sourceAccountNo"`
	PartnerReferenceNo  string             `json:"partnerReferenceNo"`
	ReferenceNo         string             `json:"referenceNo"`
	PaidAmount          TransferAmountData `json:"paidAmount"`
	TotalAmount         TransferAmountData `json:"TotalAmount"`
	TrxDateTime         string             `json:"trxDateTime"`
	BillDetails         []BCAVABillDetails `json:"billDetails"`
	FreeTexts           []BCAVADescription `json:"freeTexts"`
	FeeAmount           TransferAmountData `json:"feeAmount"`
	ProductName         string             `json:"productName"`
}

type BCAVADescription struct {
	English   string `json:"english"`
	Indonesia string `json:"indonesia"`
}

type BCAVABillDetails struct {
	BillDescription BCAVADescription     `json:"billDescription"`
	BillAmount      []TransferAmountData `json:"billAmount"`
}

type BRITransferInternalRequest struct {
	PartnerReferenceNo   string             `json:"partnerReferenceNo"`
	Amount               TransferAmountData `json:"amount"`
	BeneficiaryAccountNo string             `json:"beneficiaryAccountNo"`
	FeeType              string             `json:"feeType"`
	Remark               string             `json:"remark"`
	SourceAccountNo      string             `json:"sourceAccountNo"`
	TransactionDate      string             `json:"transactionDate"`
	AdditionalInfo       map[string]string  `json:"additionalInfo"`
}

type BRITransferInternalResponse struct {
	ResponseCode         string             `json:"responseCode"`
	ResponseMessage      string             `json:"responseMessage"`
	PartnerReferenceNo   string             `json:"partnerReferenceNo"`
	ReferenceNo          string             `json:"referenceNo"`
	Amount               TransferAmountData `json:"amount"`
	BeneficiaryAccountNo string             `json:"beneficiaryAccountNo"`
	SourceAccountNo      string             `json:"sourceAccountNo"`
	TransactionDate      string             `json:"transactionDate"`
}
