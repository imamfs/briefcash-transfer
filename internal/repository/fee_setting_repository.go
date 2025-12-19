package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found in database")

type FeeSettingRepository interface {
	FindAll(ctx context.Context) ([]entity.FeeSettings, error)
	FindByCodeAndChannel(ctx context.Context, merchantCode, channel string) (entity.FeeSettings, error)
	WithTransaction(trx *gorm.DB) FeeSettingRepository
}

type feeSettingRepository struct {
	db *gorm.DB
}

func NewFeeSetting(db *gorm.DB) FeeSettingRepository {
	return &feeSettingRepository{db}
}

func (f *feeSettingRepository) FindAll(ctx context.Context) ([]entity.FeeSettings, error) {
	var feeSettings []entity.FeeSettings

	err := f.db.WithContext(ctx).Order("merchant_code ASC").Find(&feeSettings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get fee settings list %w", err)
	}

	return feeSettings, nil
}

func (f *feeSettingRepository) FindByCodeAndChannel(ctx context.Context, merchantCode, channel string) (entity.FeeSettings, error) {
	var feeSetting entity.FeeSettings
	err := f.db.WithContext(ctx).Where("merchant_code = ? AND channel = ?", merchantCode, channel).First(&feeSetting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.FeeSettings{}, nil
	}

	if err != nil {
		return entity.FeeSettings{}, fmt.Errorf("failed to get fee settings list, with error: %w", err)
	}
	return feeSetting, nil
}

func (f *feeSettingRepository) WithTransaction(trx *gorm.DB) FeeSettingRepository {
	return &feeSettingRepository{db: trx}
}
