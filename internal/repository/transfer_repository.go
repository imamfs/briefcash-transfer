package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TransferRepository interface {
	Save(ctx context.Context, trx *entity.Transaction) error
	FindByRefNo(ctx context.Context, partnerReferenceNo string) (*entity.TransferTemp, error)
	Update(ctx context.Context, id int64, status string) error
	WithTransaction(trx *gorm.DB) TransferRepository
}

type transferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) TransferRepository {
	return &transferRepository{db}
}

func (r *transferRepository) Save(ctx context.Context, transaction *entity.Transaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to save transaction with error %w", err)
	}
	return nil
}

func (r *transferRepository) FindByRefNo(ctx context.Context, partnerReferenceNo string) (*entity.TransferTemp, error) {
	var transferTemp entity.TransferTemp

	if err := r.db.WithContext(ctx).Select("id", "status").Where("partner_reference_no = ?", partnerReferenceNo).Scan(&transferTemp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query partner reference no %s: %w", partnerReferenceNo, err)
	}

	return &transferTemp, nil
}

func (r *transferRepository) Update(ctx context.Context, id int64, status string) error {
	if err := r.db.WithContext(ctx).Model(&entity.Transaction{}).
		Where("id = ? and status = 'PENDING'", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update transfer data in id %d:%w", id, err)
	}
	return nil
}

func (r *transferRepository) WithTransaction(trx *gorm.DB) TransferRepository {
	if trx == nil {
		return r
	}
	return &transferRepository{db: trx}
}
