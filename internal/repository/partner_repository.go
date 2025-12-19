package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type PartnerRepository interface {
	FindAll(ctx context.Context) ([]entity.BankConfig, error)
}

type partnerRepository struct {
	db *gorm.DB
}

func NewPartnerRepository(db *gorm.DB) PartnerRepository {
	return &partnerRepository{db}
}

func (r *partnerRepository) FindAll(ctx context.Context) ([]entity.BankConfig, error) {
	var listConfig []entity.BankConfig

	err := r.db.WithContext(ctx).
		Select("partner.company_bank_code AS bank_code, domestic_bank.short_name AS bank_name, partner_url.kafka_topic, partner_url.kafka_topic_group").
		Joins("INNER JOIN partner_url ON partner.company_id = partner_url.company_id").
		Joins("INNER JOIN domestic_bank ON partner.company_id = domestic_bank.company_id").
		Scan(&listConfig).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch partner bank configuration: %w", err)
	}

	return listConfig, nil
}
