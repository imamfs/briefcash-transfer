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
	TransactionDate   string `json:"transactionDate"`
	CustomerReference string `json:"customerReference"`
	Channel           string `json:"channel"` // online, bifast, sknbi, rtgs, va, wallet
	Remarks           string `json:"remarks"`
	Email             string `json:"email"`
	Address           string `json:"address"`
	Citizenship       string `json:"citizenship"` // wna, wni
	TransferPurpose   string `json:"transferPurpose"`
	TransferActivity  string `json:"transferActivity"` // only mandatory for non indonesian citizen
	CustomerType      string `json:"customerType"`     // 01 - individu, 02 - corporate, 03 - others
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

type BRITransferExternalRequest struct {
	PartnerReferenceNo     string                                   `json:"partnerReferenceNo"`
	Amount                 TransferAmountData                       `json:"amount"`
	BeneficiaryAccountName string                                   `json:"beneficiaryAccountName"`
	BeneficiaryAccountNo   string                                   `json:"beneficiaryAccountNo"`
	BeneficiaryAddress     string                                   `json:"beneficiaryAddress"`
	BeneficiaryBankCode    string                                   `json:"beneficiaryBankCode"`
	BeneficiaryBankName    string                                   `json:"beneficiaryBankName"`
	BeneficiaryEmail       string                                   `json:"beneficiaryEmail"`
	SourceAccountNo        string                                   `json:"sourceAccountNo"`
	TransactionDate        string                                   `json:"transactionDate"`
	AdditionalInfo         BRITransferExternalAdditionalInfoRequest `json:"additionalInfo"`
}

type BRITransferExternalAdditionalInfoRequest struct {
	ServiceCode          string `json:"serviceCode"`
	ReferenceNo          string `json:"referenceNo"`
	ExternalId           string `json:"externalId"`
	SenderIdentityNumber string `json:"senderIdentityNumber"`
	SenderType           string `json:"senderType"`           // 01 - individual, 02 - corporate, 03 - government, 04 - remittance, 99 - others
	SenderResidentStatus string `json:"senderResidentStatus"` // 01 - resident, 02 - non resident
}

type BRITransferExternalResponse struct {
	ResponseCode         string                                    `json:"responseCode"`
	ResponseMessage      string                                    `json:"responseMessage"`
	ReferenceNo          string                                    `json:"referenceNo"`
	PartnerReferenceNo   string                                    `json:"partnerReferenceNo"`
	Amount               TransferAmountData                        `json:"amount"`
	BeneficiaryAccountNo string                                    `json:"beneficiaryAccountNo"`
	BeneficiaryBankCode  string                                    `json:"beneficiaryBankCode"`
	SourceAccountNo      string                                    `json:"sourceAccountNo"`
	AdditionalInfo       BRITransferExternalAdditionalInfoResponse `json:"additionalInfo"`
}

type BRITransferExternalAdditionalInfoResponse struct {
	OriginalReferenceNo string `json:"originalReferenceNo"`
	JournalSequence     string `json:"journalSequence"`
	ExternalId          string `json:"externalId"`
}

type BRITransferVARequest struct {
	PartnerServiceId   string             `json:"partnerServiceId"`
	CustomerNo         string             `json:"customerNo"`
	VirtualAccountNo   string             `json:"virtualAccountNo"`
	VirtualAccountName string             `json:"virtualAccountName"`
	SourceAccountNo    string             `json:"sourceAccountNo"`
	PartnerReferenceNo string             `json:"partnerReferenceNo"`
	PaidAmount         TransferAmountData `json:"paidAmount"`
	TrxDateTime        string             `json:"trxDateTime"`
}

type BRITransferVAResponse struct {
	ResponseCode       string            `json:"responseCode"`
	ResponseMessage    string            `json:"responseMessage"`
	VirtualAccountData BRITransferVAData `json:"virtualAccountData"`
}

type BRITransferVAData struct {
	PartnerServiceId   string             `json:"partnerServiceId"`
	CustomerNo         string             `json:"customerNo"`
	VirtualAccountNo   string             `json:"virtualAccountNo"`
	VirtualAccountName string             `json:"virtualAccountName"`
	PartnerReferenceNo string             `json:"partnerReferenceNo"`
	PaymentRequestId   string             `json:"paymentRequestId"`
	PaidAmount         TransferAmountData `json:"paidAmount"`
	TrxDateTime        string             `json:"trxDateTime"`
}

type CIMBTransferInternalRequest struct {
	PartnerReferenceNo   string                          `json:"partnerReferenceNo"`
	Amount               TransferAmountData              `json:"amount"`
	BeneficiaryAccountNo string                          `json:"beneficiaryAccountNo"`
	Remark               string                          `json:"remark"`
	SourceAccountNo      string                          `json:"sourceAccountNo"`
	TransactionDate      string                          `json:"transactionDate"`
	AdditionalData       CIMBTransferInternalInfoRequest `json:"additionalInfo"`
}

type CIMBTransferInternalInfoRequest struct {
	BeneficiaryAccountName string `json:"beneficiaryAccountName"`
}

type CIMBTransferInternaResponse struct {
	ResponseCode         string             `json:"responseCode"`
	ResponseMessage      string             `json:"responseMessage"`
	PartnerReferenceNo   string             `json:"partnerReferenceNo"`
	Amount               TransferAmountData `json:"amount"`
	BeneficiaryAccountNo string             `json:"beneficiaryAccountNo"`
	Currency             string             `json:"currency"`
	SourceAccountNo      string             `json:"sourceAccountNo"`
}

type CIMBTransferInternalInfoResponse struct {
	Remark                 string `json:"remark"`
	BeneficiaryAccountName string `json:"beneficiaryAccountName"`
}

type CIMBTransferExternalRequest struct {
	PartnerReferenceNo     string                          `json:"partnerReferenceNo"`
	Amount                 TransferAmountData              `json:"amount"`
	BeneficiaryAccountName string                          `json:"beneficiaryAccountName"`
	BeneficiaryAccountNo   string                          `json:"beneficiaryAccountNo"`
	BeneficiaryBankCode    string                          `json:"beneificaryBankCode"`
	SourceAccountNo        string                          `json:"sourceAccountNo"`
	TransactionDate        string                          `json:"transactionDate"`
	AdditionalInfo         CIMBTransferExternalInfoRequest `json:"additionalInfo"`
}

type CIMBTransferExternalInfoRequest struct {
	Remark         string `json:"remark"`
	TrxType        string `json:"trxType"`
	ProxyValue     string `json:"proxyValue"`
	ProxyType      string `json:"proxyType"`
	TrxPurposeCode string `json:"trxPurposeCode"`
}

type CIMBTransferExternalResponse struct {
	ResponseCode           string                           `json:"responseCode"`
	ResponseMessage        string                           `json:"responseMessage"`
	ReferenceNo            string                           `json:"referenceNo"`
	PartnerReferenceNo     string                           `json:"partnerReferenceNo"`
	Amount                 TransferAmountData               `json:"amount"`
	BeneficiaryAccountNo   string                           `json:"beneficiaryAccountNo"`
	BeneficiaryAccountName string                           `json:"beneficiaryAccountName"`
	BeneficiaryBankCode    string                           `json:"beneficiaryBankCode"`
	SourceAccountNo        string                           `json:"sourceAccountNo"`
	TransactionDate        string                           `json:"transactionDate"`
	AdditionalInfo         CIMBTransferExternalInfoResponse `json:"additionalInfo"`
}

type CIMBTransferExternalInfoResponse struct {
	Remark         string `json:"remark"`
	TrxType        string `json:"trxType"`
	ProxyValue     string `json:"proxyValue"`
	ProxyType      string `json:"proxyType"`
	TrxPurposeCode string `json:"trxPurposeCode"`
}

type CIMBTransferVARequest struct {
	PartnerServiceId   string             `json:"partnerServiceId"`
	CustomerNo         string             `json:"customerNo"`
	VirtualAccountNo   string             `json:"virtualAccountNo"`
	VirtualAccountName string             `json:"virtualAccountName"`
	PartnerReferenceNo string             `json:"partnerReferenceNo"`
	PaidAmount         TransferAmountData `json:"paidAmount"`
	TotalAmount        TransferAmountData `json:"totalAmount"`
	TrxDateTime        string             `json:"trxDateTime"`
	AdditionalInfo     map[string]string  `json:"additionalInfo"`
}

type CIMBTransferVAResponse struct {
	ResponseCode       string     `json:"responseCode"`
	ResponseMessage    string     `json:"responseMessage"`
	VirtualAccountData CIMBVAData `json:"virtualAccountData"`
}

type CIMBVAData struct {
	VirtualAccountName string             `json:"virtualAccountName"`
	PartnerServiceId   string             `json:"partnerServiceId"`
	CustomerNo         string             `json:"customerNo"`
	InquiryRequestId   string             `json:"InquiryRequestId"`
	PaymentRequestId   string             `json:"paymentRequestId"`
	PaidAmount         TransferAmountData `json:"paidAmount"`
	TotalAmount        TransferAmountData `json:"totalAmount"`
	TrxDateTime        string             `json:"trxDateTime"`
	PartnerReferenceNo string             `json:"partnerReferenceNo"`
	VirtualAccountNo   string             `json:"virtualAccountNo"`
	ReferenceNo        string             `json:"referenceNo"`
	PaymentFlagReason  CIMBVAFlagReason   `json:"paymentFlagReason"`
}

type CIMBVAFlagReason struct {
	English   string `json:"english"`
	Indonesia string `json:"indonesia"`
}

type PermataMessageHeaderRequest struct {
	RequestTimestamp    string `json:"RequestTimestamp"`
	CustomerReferenceId string `json:"CustRefID"`
}

type PermataInternalMessageBodyRequest struct {
	FromAccount            string `json:"FromAccount"`
	ToAccount              string `json:"ToAccount"`
	Amount                 int64  `json:"amount"`
	CurrencyCode           string `json:"CurrencyCode"`
	ChargeTo               string `json:"ChargeTo"`
	TrxDesc                string `json:"TrxDesc"`
	TrxDesc2               string `json:"TrxDesc2"`
	BeneficiaryEmail       string `json:"BenefEmail"`
	BeneficiaryAccountName string `json:"BenefAccName"`
	BeneficiaryPhoneNo     string `json:"BenefPhoneNo"`
	FromAccountName        string `json:"FromAcctName"`
	TkiFlag                string `json:"TkiFlag"`
}

type PermataTransferInternalRequest struct {
	MessageHeader PermataMessageHeaderRequest       `json:"MsgRqHdr"`
	MessageBody   PermataInternalMessageBodyRequest `json:"XferInfo"`
}

type PermataMessageHeaderResponse struct {
	ResponseTimestamp   string `json:"ResponseTimestamp"`
	CustomerReferenceId string `json:"CustRefID"`
	StatusCode          string `json:"StatusCode"`
	StatusDesc          string `json:"StatusDesc"`
}

type PermataExternalMessageBodyRequest struct {
	FromAccount            string `json:"FromAccount"`
	ToAccount              string `json:"ToAccount"`
	ToBankId               string `json:"ToBankId"`
	ToBankName             string `json:"ToBankName"`
	Amount                 int64  `json:"amount"`
	ChargeTo               string `json:"ChargeTo"`
	TrxDesc                string `json:"TrxDesc"`
	TrxDesc2               string `json:"TrxDesc2"`
	BeneficiaryEmail       string `json:"BenefEmail"`
	BeneficiaryAccountName string `json:"BenefAccName"`
	BeneficiaryPhoneNo     string `json:"BenefPhoneNo"`
	FromAccountName        string `json:"FromAcctName"`
	DatiII                 string `json:"DatiII"`
	TkiFlag                string `json:"TkiFlag"`
}

type PermataTransferExternalRequest struct {
	MessageHeader PermataMessageHeaderRequest       `json:"MsgRqHdr"`
	MessageBody   PermataExternalMessageBodyRequest `json:"XferInfo"`
}

type PermataSKNMessageBodyRequest struct {
	FromAccount               string `json:"FromAccount"`
	ToAccount                 string `json:"ToAccount"`
	ToBankId                  string `json:"ToBankId"`
	ToBankName                string `json:"ToBankName"`
	Amount                    int64  `json:"amount"`
	CurrencyCode              string `json:"CurrencyCode"`
	ChargeTo                  string `json:"ChargeTo"`
	TrxDesc                   string `json:"TrxDesc"`
	TrxDesc2                  string `json:"TrxDesc2"`
	ResidentStatus            string `json:"ResidentStatus"`
	BeneficiaryType           string `json:"BenefType"`
	BeneficiaryEmail          string `json:"BenefEmail"`
	BeneficiaryAccountName    string `json:"BenefAccName"`
	BeneficiaryPhoneNo        string `json:"BenefPhoneNo"`
	BeneficiaryBankAddress    string `json:"BenefBankAddress"`
	BeneficiaryBankBranchName string `json:"BenefBankBranchName"`
	BeneficiaryBankCity       string `json:"BenefBankCity"`
	FromAccountName           string `json:"FromAcctName"`
	FromCurrencyCode          string `json:"FromCurrencyCode"`
	Filler1                   string `json:"Filler1"`
	Filler2                   string `json:"Filler2"`
	Filler3                   string `json:"Filler3"`
	BeneficiaryAddress1       string `json:"BenefAddress1"`
	BeneficiaryAddress2       string `json:"BenefAddress2"`
	BeneficiaryAddress3       string `json:"BenefAddress3"`
	DatiII                    string `json:"DatiII"`
	TkiFlag                   string `json:"TkiFlag"`
}

type PermataTransferSKNRequest struct {
	MessageHeader PermataMessageHeaderRequest  `json:"MsgRqHdr"`
	MessageBody   PermataSKNMessageBodyRequest `json:"XferInfo"`
}

type PermataRTGSMessageBodyRequest struct {
	FromAccount               string `json:"FromAccount"`
	ToAccount                 string `json:"ToAccount"`
	ToBankId                  string `json:"ToBankId"`
	ToBankName                string `json:"ToBankName"`
	Amount                    int64  `json:"amount"`
	CurrencyCode              string `json:"CurrencyCode"`
	ChargeTo                  string `json:"ChargeTo"`
	TrxDesc                   string `json:"TrxDesc"`
	TrxDesc2                  string `json:"TrxDesc2"`
	CitizenStatus             string `json:"CitizenStatus"`
	ResidentStatus            string `json:"ResidentStatus"`
	BeneficiaryEmail          string `json:"BenefEmail"`
	BeneficiaryAccountName    string `json:"BenefAccName"`
	BeneficiaryPhoneNo        string `json:"BenefPhoneNo"`
	BeneficiaryBankAddress    string `json:"BenefBankAddress"`
	BeneficiaryBankBranchName string `json:"BenefBankBranchName"`
	BeneficiaryBankCity       string `json:"BenefBankCity"`
	FromAccountName           string `json:"FromAcctName"`
	FromCurrencyCode          string `json:"FromCurrencyCode"`
	Filler1                   string `json:"Filler1"`
	Filler2                   string `json:"Filler2"`
	Filler3                   string `json:"Filler3"`
	BeneficiaryAddress1       string `json:"BenefAddress1"`
	BeneficiaryAddress2       string `json:"BenefAddress2"`
	BeneficiaryAddress3       string `json:"BenefAddress3"`
	DatiII                    string `json:"DatiII"`
	TkiFlag                   string `json:"TkiFlag"`
}

type PermataTransferRTGSRequest struct {
	MessageHeader PermataMessageHeaderRequest   `json:"MsgRqHdr"`
	MessageBody   PermataRTGSMessageBodyRequest `json:"XferInfo"`
}

type PermataTransferResponse struct {
	MessageHeader    PermataMessageHeaderResponse `json:"MsgRsHdr"`
	TransactionRefNo string                       `json:"TrxReffNo"`
}

type PermataVAMessageBodyRequest struct {
	BillType             string `json:"BillType"` // MOBILEVOUCHER, VIRTUALACCOUNT, CREDIT CARD
	InstitutionCode      string `json:"InstCode"`
	BillNumber           string `json:"BillNumber"`
	TransactionAmount    string `json:"TrxAmount"`
	Currency             string `json:"Currency"`
	UserId               string `json:"UserId"`
	DebitAccountNumber   string `json:"DebAccNumber"`
	DebitAccountName     string `json:"DebAccName"`
	DebitAccountCurrency string `json:"DebAccCur"`
}

type PermataTransferVARequest struct {
	MessageHeader   PermataMessageHeaderRequest `json:"MsgRqHdr"`
	BillPaymentInfo PermataVAMessageBodyRequest `json:"BillPaymentInfo"`
}

type PermataTransferVAResponse struct {
	MessageHeader   PermataMessageHeaderResponse `json:"MsgRsHdr"`
	BillReferenceNo string                       `json:"BillRefNo"`
	InstitutionCode string                       `json:"InstCode"`
}
