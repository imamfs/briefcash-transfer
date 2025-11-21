package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type RecipientRepository interface {
	Save(ctx context.Context, sender *entity.DataRecipient) error
	WithTransaction(trx *gorm.DB) RecipientRepository
}

type recipientRepository struct {
	db *gorm.DB
}

func NewRecipientRepository(db *gorm.DB) RecipientRepository {
	return &recipientRepository{db}
}

func (r *recipientRepository) Save(ctx context.Context, sender *entity.DataRecipient) error {
	err := r.db.WithContext(ctx).Table("data_recipient").Create(sender).Error

	if err != nil {
		return fmt.Errorf("failed to save data recipient with error %w", err)
	}

	return nil
}

func (r *recipientRepository) WithTransaction(trx *gorm.DB) RecipientRepository {
	return &recipientRepository{db: trx}
}
