package entity

type BankConfig struct {
	BankCode        string `gorm:"column:bank_code"`
	BankName        string `gorm:"column:bank_name"`
	KafkaTopic      string `gorm:"column:kafka_topic"`
	KafkaTopicGroup string `gorm:"column:kafka_topic_group"`
}
