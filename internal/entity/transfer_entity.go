package entity

import "time"

type Transaction struct {
	ID                      int64     `gorm:"column:id;primaryKey"`
	Sender                  int64     `gorm:"column:data_sender_id"`
	Recipient               int64     `gorm:"column:data_recipient_id"`
	MerchantCode            string    `gorm:"column:merchant_code"`
	PartnerReferenceNo      string    `gorm:"column:partner_reference_no"`
	BankReferenceNo         string    `gorm:"column:bank_reference_no"`
	SystemReferenceNo       string    `gorm:"column:system_reference_no"`
	Amount                  float64   `gorm:"column:amount"`
	Currency                string    `gorm:"column:currency"`
	Remark                  string    `gorm:"column:remark"`
	TransactionType         string    `gorm:"column:transaction_type"`
	TransactionDate         time.Time `gorm:"column:transaction_date"`
	Status                  string    `gorm:"column:status"`
	IsReversal              bool      `gorm:"column:is_reversal"`
	CompanyCharge           float32   `gorm:"column:company_charge"`
	PartnerCharge           float32   `gorm:"column:partner_charge"`
	AdditionalPartnerCharge float32   `gorm:"column:additional_partner_charge"`
	TaxCharge               float32   `gorm:"column:tax_charge"`
	IsReconcile             bool      `gorm:"column:is_reconcile"`
	ReconcileDate           time.Time `gorm:"column:reconcile_date"`
}

type DataSender struct {
	ID                      int64     `gorm:"column:id;primaryKey"`
	IdType                  string    `gorm:"column:id_type"`
	IdNumber                string    `gorm:"column:id_number"`
	Name                    string    `gorm:"column:name"`
	BirthDate               time.Time `gorm:"column:birth_date"`
	BirthPlace              string    `gorm:"column:birth_place"`
	Address                 string    `gorm:"column:address"`
	Country                 string    `gorm:"column:country"`
	Email                   string    `gorm:"column:email"`
	Phone                   string    `gorm:"column:phone"`
	Profession              string    `gorm:"column:profession"`
	BeneficiaryRelationship string    `gorm:"column:beneficiary_relationship"`
	SourceOfFunds           string    `gorm:"column:source_of_funds"`
}

type DataRecipient struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	IdType        string    `gorm:"column:id_type"`
	IdNumber      string    `gorm:"column:id_number"`
	Name          string    `gorm:"column:name"`
	BirthDate     time.Time `gorm:"column:birth_date"`
	BirthPlace    string    `gorm:"column:birth_place"`
	Address       string    `gorm:"column:address"`
	Country       string    `gorm:"column:country"`
	Email         string    `gorm:"column:email"`
	Phone         string    `gorm:"column:phone"`
	Profession    string    `gorm:"column:profession"`
	AccountNumber string    `gorm:"column:account_number"`
	BankCode      string    `gorm:"column:bank_code"`
}

type FeeSettings struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	MerchantCode  string    `gorm:"column:merchant_code"`
	Channel       string    `gorm:"column:channel"`
	FeePartner    float32   `gorm:"column:fee_partner"`
	FeeService    float32   `gorm:"column:fee_service"`
	FeeTax        float32   `gorm:"column:fee_tax"`
	AdditionalFee float32   `gorm:"additional_fee"`
	TotalCharge   float32   `gorm:"total_charge"`
	CreatedAt     time.Time `gorm:"created_at"`
	LastUpdated   time.Time `gorm:"last_updated"`
}
