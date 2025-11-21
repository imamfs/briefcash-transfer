package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found in database")

type FeeRepository interface {
	FindAll(ctx context.Context) ([]entity.FeeSettings, error)
	WithTransaction(trx *gorm.DB) FeeRepository
}

type feeRepository struct {
	db *gorm.DB
}

func NewFeeSettings(db *gorm.DB) FeeRepository {
	return &feeRepository{db}
}

func (f *feeRepository) FindAll(ctx context.Context) ([]entity.FeeSettings, error) {
	var feeSettings []entity.FeeSettings

	err := f.db.WithContext(ctx).Table("fee_settings").Find(&feeSettings).Order("merchant_code ASC").Error

	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get fee settings list, with error %w", err)
	}

	return feeSettings, nil
}

func (f *feeRepository) WithTransaction(trx *gorm.DB) FeeRepository {
	return &feeRepository{db: trx}
}
