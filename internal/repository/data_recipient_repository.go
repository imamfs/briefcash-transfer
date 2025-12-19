package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type RecipientRepository interface {
	Save(ctx context.Context, recipient *entity.DataRecipient) error
	WithTransaction(trx *gorm.DB) RecipientRepository
}

type recipientRepository struct {
	db *gorm.DB
}

func NewRecipientRepository(db *gorm.DB) RecipientRepository {
	return &recipientRepository{db}
}

func (r *recipientRepository) Save(ctx context.Context, recipient *entity.DataRecipient) error {
	if err := r.db.WithContext(ctx).Create(recipient).Error; err != nil {
		return fmt.Errorf("failed to save data recipient, with error: %w", err)
	}

	return nil
}

func (r *recipientRepository) WithTransaction(trx *gorm.DB) RecipientRepository {
	return &recipientRepository{db: trx}
}
