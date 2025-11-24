package repositoryredis

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var KeyFeeSettings string = "feeSettings"
var KeyBalance string = "balance"

type TransferRedisRepository interface {
	SetFee(ctx context.Context, settings []entity.FeeSettings) error
	SetBalance(ctx context.Context, balance []entity.MerchantBalance) error
	GetFeeByCodeAndChannel(ctx context.Context, merchantCode, chanel string) (float32, error)
	GetBalanceByCode(ctx context.Context, merchantCode string) (float64, error)
}

type transferRedisRepository struct {
	client *redis.Client
}

func NewTransferRedisRepository(client *redis.Client) TransferRedisRepository {
	return &transferRedisRepository{client}
}

func (r *transferRedisRepository) SetFee(ctx context.Context, settings []entity.FeeSettings) error {
	if len(settings) == 0 {
		return fmt.Errorf("list fee settings is empty")
	}

	pipe := r.client.TxPipeline()
	for _, v := range settings {
		key := fmt.Sprintf("%s:%s:%s", KeyFeeSettings, v.MerchantCode, v.Channel)

		data := map[string]float32{
			"fee_partner":    v.FeePartner,
			"fee_service":    v.FeeService,
			"fee_tax":        v.FeeTax,
			"additional_fee": v.AdditionalFee,
			"total_charge":   v.TotalCharge,
		}

		pipe.HSet(ctx, key, data)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to load fee settings to redis, with error: %w", err)
	}

	return nil
}

func (r *transferRedisRepository) SetBalance(ctx context.Context, balance []entity.MerchantBalance) error {
	return nil
}

func (r *transferRedisRepository) GetFeeByCodeAndChannel(ctx context.Context, merchantCode, chanel string) (float32, error) {
	return 0, nil
}

func (r *transferRedisRepository) GetBalanceByCode(ctx context.Context, merchantCode string) (float64, error) {
	return 0, nil
}
