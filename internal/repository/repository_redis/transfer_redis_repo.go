package repositoryredis

import (
	"briefcash-transfer/internal/entity"
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	KeyFeeSettings string = "fee_settings"
	KeyBalance     string = "balance"
	FeePartner     string = "fee_partner"
	FeeService     string = "fee_service"
	FeeTax         string = "fee_tax"
	AdditionalFee  string = "additional_fee"
	TotalCharge    string = "total_charge"
)

type TransferRedisRepository interface {
	SetFee(ctx context.Context, settings []entity.FeeSettings) error
	SetBalance(ctx context.Context, balance []entity.MerchantBalance) error
	GetFeeByCodeAndChannel(ctx context.Context, merchantCode, channel string) (*entity.FeeSettings, error)
	GetBalanceByMerchantCode(ctx context.Context, merchantCode string) (float64, error)
	UpdateBalance(ctx context.Context, merchantCode, amount string) error
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

		data := map[string]string{
			FeePartner:    strconv.FormatFloat(float64(v.FeePartner), 'f', -1, 64),
			FeeService:    strconv.FormatFloat(float64(v.FeeService), 'f', -1, 64),
			FeeTax:        strconv.FormatFloat(float64(v.FeeTax), 'f', -1, 64),
			AdditionalFee: strconv.FormatFloat(float64(v.AdditionalFee), 'f', -1, 64),
			TotalCharge:   strconv.FormatFloat(float64(v.TotalCharge), 'f', -1, 64),
		}

		pipe.HSet(ctx, key, data)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to load fee settings to redis, with error: %w", err)
	}

	return nil
}

func (r *transferRedisRepository) SetBalance(ctx context.Context, balance []entity.MerchantBalance) error {
	if len(balance) == 0 {
		return fmt.Errorf("list merchant balance is empty")
	}

	pipe := r.client.TxPipeline()

	data := make(map[string]string, len(balance))
	for _, v := range balance {
		data[v.MerchantCode] = strconv.FormatFloat(v.Balance, 'f', -1, 64)
	}

	pipe.HSet(ctx, KeyBalance, data)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to load list merchant balance to redis, with error: %w", err)
	}

	return nil
}

func (r *transferRedisRepository) GetFeeByCodeAndChannel(ctx context.Context, merchantCode, channel string) (*entity.FeeSettings, error) {
	key := fmt.Sprintf("%s:%s:%s", KeyFeeSettings, merchantCode, channel)

	data, err := r.client.HGetAll(ctx, key).Result()

	if err != nil {
		return nil, fmt.Errorf("error on redis server while retrieving list fee merchant, with error %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("key not found or hash empty")
	}

	floatParser := func(data string) float64 {
		f, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return 0
		}
		return f
	}

	feeSetting := &entity.FeeSettings{}
	feeSetting.FeePartner = floatParser(data["fee_partner"])
	feeSetting.FeeService = floatParser(data["fee_service"])
	feeSetting.FeeTax = floatParser(data["fee_tax"])
	feeSetting.AdditionalFee = floatParser(data["additional_fee"])
	feeSetting.TotalCharge = floatParser(data["total_charge"])
	feeSetting.MerchantCode = merchantCode
	feeSetting.Channel = channel

	return feeSetting, nil
}

func (r *transferRedisRepository) GetBalanceByMerchantCode(ctx context.Context, merchantCode string) (float64, error) {

	amount, err := r.client.HGet(ctx, KeyBalance, merchantCode).Result()

	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("key or field not found in redis")
		}

		return 0, fmt.Errorf("error in redis server while retrieving merchant balance, with error: %w", err)
	}

	return strconv.ParseFloat(amount, 64)
}

func (r *transferRedisRepository) UpdateBalance(ctx context.Context, merchantCode, amount string) error {

	if err := r.client.HSet(ctx, KeyBalance, merchantCode, amount).Err(); err != nil {
		return fmt.Errorf("failed to update balance in redis, with error: %w", err)
	}

	return nil
}
