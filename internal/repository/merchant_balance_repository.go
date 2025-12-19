package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type MerchantBalanceRepository interface {
	FindAll(ctx context.Context) ([]entity.MerchantBalance, error)
	WithTransaction(trx *gorm.DB) MerchantBalanceRepository
}

type merchantBalanceRepository struct {
	db *gorm.DB
}

func NewMerchantBalanceRepository(db *gorm.DB) MerchantBalanceRepository {
	return &merchantBalanceRepository{db}
}

func (b *merchantBalanceRepository) FindAll(ctx context.Context) ([]entity.MerchantBalance, error) {
	var merchant []entity.MerchantBalance
	err := b.db.WithContext(ctx).Select("merchant_code", "balance").
		Order("merchant_code ASC").Find(&merchant).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve list of merchant balance: %w", err)
	}

	return merchant, nil
}

func (b *merchantBalanceRepository) WithTransaction(trx *gorm.DB) MerchantBalanceRepository {
	return &merchantBalanceRepository{db: trx}
}
