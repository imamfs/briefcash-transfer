package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type LedgerRepository interface {
	Save(ctx context.Context, statement *entity.AccountStatement) error
	Delete(ctx context.Context, transferId int64) error
	WithTransaction(trx *gorm.DB) LedgerRepository
}

type ledgerRepository struct {
	db *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) LedgerRepository {
	return &ledgerRepository{db}
}

func (a *ledgerRepository) Save(ctx context.Context, statement *entity.AccountStatement) error {
	if err := a.db.WithContext(ctx).Create(statement).Error; err != nil {
		return fmt.Errorf("failed to save account statement, with error: %w", err)
	}
	return nil
}

func (a *ledgerRepository) Delete(ctx context.Context, transferId int64) error {
	var accountStatement entity.AccountStatement
	if err := a.db.WithContext(ctx).Where("transaction_id = ?").Delete(&accountStatement).Error; err != nil {
		return fmt.Errorf("failed to delete record for transaction id %d, with error: %w", transferId, err)
	}
	return nil
}

func (a *ledgerRepository) WithTransaction(trx *gorm.DB) LedgerRepository {
	return &ledgerRepository{db: trx}
}
