package repository

import (
	"briefcash-transfer/internal/entity"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BalanceRepository interface {
	Debit(ctx context.Context, merchantCode string, totalAmount float64) (float64, error)
	Credit(ctx context.Context, merchantCode string, amount float64) (float64, error)
	FindByCode(ctx context.Context, merchantCode string) (*entity.MerchantAccounts, error)
	WithTransaction(trx *gorm.DB) BalanceRepository
}

type balanceRepository struct {
	db *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepository{db}
}

func (a *balanceRepository) Debit(ctx context.Context, merchantCode string, totalAmount float64) (float64, error) {
	var balance float64
	result := a.db.WithContext(ctx).Model(&entity.MerchantAccounts{}).
		Where("merchant_code = ? and balance >= ?", merchantCode, totalAmount).
		Update("balance", gorm.Expr("balance - ?", totalAmount)).Scan(&balance)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return -1, fmt.Errorf("insufficient balance")
	}
	return balance, nil
}

func (a *balanceRepository) Credit(ctx context.Context, merchantCode string, amount float64) (float64, error) {
	var balance float64
	result := a.db.WithContext(ctx).Model(&entity.MerchantAccounts{}).
		Where("merchant_code = ?", merchantCode, amount).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "balance"}}}).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Scan(&balance)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return -1, fmt.Errorf("insufficient balance")
	}
	return balance, nil
}

func (a *balanceRepository) FindByCode(ctx context.Context, merchantCode string) (*entity.MerchantAccounts, error) {
	var accounts entity.MerchantAccounts
	if err := a.db.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE", Table: clause.Table{Name: clause.CurrentTable},
	}).Where("merchant_code = ?", merchantCode).First(&accounts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get merchant accounts, with error: %w", err)
	}
	return &accounts, nil
}

func (a *balanceRepository) WithTransaction(trx *gorm.DB) BalanceRepository {
	return &balanceRepository{db: trx}
}
