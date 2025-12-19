package entity

import "time"

type AccountStatement struct {
	ID                  int64     `gorm:"column:id;primaryKey"`
	TransactionId       int64     `gorm:"column:transaction_id"`
	TransctionReference string    `gorm:"column:transaction_reference"`
	MerchantCode        string    `gorm:"column:merchant_code"`
	Status              string    `gorm:"column:type"`
	Channel             string    `gorm:"column:channel"`
	Description         string    `gorm:"column:description"`
	Amount              float64   `gorm:"column:amount"`
	BalanceAfter        float64   `gorm:"column:balance_after"`
	CreatedAt           time.Time `gorm:"column:created_at"`
}

type Transaction struct {
	ID                      int64      `gorm:"column:id;primaryKey"`
	Sender                  int64      `gorm:"column:data_sender_id"`
	Recipient               int64      `gorm:"column:data_recipient_id"`
	MerchantCode            string     `gorm:"column:merchant_code"`
	PartnerReferenceNo      string     `gorm:"column:partner_reference_no"`
	BankReferenceNo         *string    `gorm:"column:bank_reference_no"`
	SystemReferenceNo       *string    `gorm:"column:system_reference_no"`
	Amount                  float64    `gorm:"column:amount"`
	Currency                string     `gorm:"column:currency"`
	Remark                  string     `gorm:"column:remark"`
	TransactionType         string     `gorm:"column:transaction_type"`
	TransactionDate         time.Time  `gorm:"column:transaction_date"`
	Status                  string     `gorm:"column:status"`
	IsReversal              bool       `gorm:"column:is_reversal"`
	CompanyCharge           float32    `gorm:"column:company_charge"`
	PartnerCharge           float32    `gorm:"column:partner_charge"`
	AdditionalPartnerCharge float32    `gorm:"column:additional_partner_charge"`
	TaxCharge               float32    `gorm:"column:tax_charge"`
	IsReconcile             bool       `gorm:"column:is_reconcile"`
	ReconcileDate           *time.Time `gorm:"column:reconcile_date"`
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
	FeePartner    float64   `gorm:"column:fee_partner"`
	FeeService    float64   `gorm:"column:fee_service"`
	FeeTax        float64   `gorm:"column:fee_tax"`
	AdditionalFee float64   `gorm:"additional_fee"`
	TotalCharge   float64   `gorm:"total_charge"`
	CreatedAt     time.Time `gorm:"created_at"`
	LastUpdated   time.Time `gorm:"last_updated"`
}

type MerchantBalance struct {
	MerchantCode string  `gorm:"column:merchant_code"`
	Balance      float64 `gorm:"column:balance"`
}

type MerchantAccounts struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	MerchantCode  string    `gorm:"column:merchant_code"`
	AccountNumber string    `gorm:"column:account_number"`
	Balance       float64   `gorm:"column:balance"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	LastUpdated   time.Time `gorm:"column:last_updated"`
}

type TransferTemp struct {
	ID     int64  `gorm:"column:id;primaryKey"`
	Status string `gorm:"column:status"`
}
