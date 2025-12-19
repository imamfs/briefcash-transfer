package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type SenderRepository interface {
	Save(ctx context.Context, sender *entity.DataSender) error
	WithTransaction(trx *gorm.DB) SenderRepository
}

type senderRepository struct {
	db *gorm.DB
}

func NewSenderRepository(db *gorm.DB) SenderRepository {
	return &senderRepository{db}
}

func (s *senderRepository) Save(ctx context.Context, sender *entity.DataSender) error {
	if err := s.db.WithContext(ctx).Table("data_sender").Create(sender).Error; err != nil {
		return fmt.Errorf("failed to save data sender, with error: %w", err)
	}
	return nil
}

func (s *senderRepository) WithTransaction(trx *gorm.DB) SenderRepository {
	return &senderRepository{db: trx}
}
