package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type TransferRepository interface {
	Save(ctx context.Context, trx *entity.Transaction) error
	WithTransaction(trx *gorm.DB) TransferRepository
}

type transferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) TransferRepository {
	return &transferRepository{db}
}

func (r *transferRepository) Save(ctx context.Context, trx *entity.Transaction) error {
	err := r.db.WithContext(ctx).Table("transaction").Create(trx).Error
	if err != nil {
		return fmt.Errorf("failed to save transaction with error %w", err)
	}
	return nil
}

func (r *transferRepository) WithTransaction(trx *gorm.DB) TransferRepository {
	return &transferRepository{db: trx}
}
