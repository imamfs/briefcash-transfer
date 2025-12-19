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

type RedisRepository interface {
	SetListFee(ctx context.Context, settings []entity.FeeSettings) error
	SetFee(ctx context.Context, feeSetting entity.FeeSettings) error
	SetBalance(ctx context.Context, balance []entity.MerchantBalance) error
	SetPendingStatus(ctx context.Context, externalId string) error
	FindByCodeAndChannel(ctx context.Context, merchantCode, channel string) (entity.FeeSettings, error)
	FindByMerchantCode(ctx context.Context, merchantCode string) (float64, error)
	UpdateBalance(ctx context.Context, merchantCode, amount string) error
	RefundBalance(ctx context.Context, merchantCode string, amount float64) error
	DeletePendingStatus(ctx context.Context, externalId string) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client}
}

func (r *redisRepository) SetListFee(ctx context.Context, listFee []entity.FeeSettings) error {
	if len(listFee) == 0 {
		return fmt.Errorf("list fee setting is empty")
	}

	pipe := r.client.TxPipeline()
	for _, v := range listFee {
		key := fmt.Sprintf("%s:%s:%s", KeyFeeSettings, v.MerchantCode, v.Channel)

		data := map[string]string{
			FeePartner:    strconv.FormatFloat(float64(v.FeePartner), 'f', 2, 64),
			FeeService:    strconv.FormatFloat(float64(v.FeeService), 'f', 2, 64),
			FeeTax:        strconv.FormatFloat(float64(v.FeeTax), 'f', 2, 64),
			AdditionalFee: strconv.FormatFloat(float64(v.AdditionalFee), 'f', 2, 64),
			TotalCharge:   strconv.FormatFloat(float64(v.TotalCharge), 'f', 2, 64),
		}

		pipe.HSet(ctx, key, data)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to cache list fee setting to redis, with error: %w", err)
	}

	return nil
}

func (r *redisRepository) SetFee(ctx context.Context, settings entity.FeeSettings) error {
	key := fmt.Sprintf("%s:%s:%s", KeyFeeSettings, settings.MerchantCode, settings.Channel)

	data := map[string]string{
		FeePartner:    strconv.FormatFloat(float64(settings.FeePartner), 'f', 2, 64),
		FeeService:    strconv.FormatFloat(float64(settings.FeeService), 'f', 2, 64),
		FeeTax:        strconv.FormatFloat(float64(settings.FeeTax), 'f', 2, 64),
		AdditionalFee: strconv.FormatFloat(float64(settings.AdditionalFee), 'f', 2, 64),
		TotalCharge:   strconv.FormatFloat(float64(settings.TotalCharge), 'f', 2, 64),
	}

	if err := r.client.HSet(ctx, key, data).Err(); err != nil {
		return fmt.Errorf("failed to cache fee setting to redis, with error: %w", err)
	}

	return nil
}

func (r *redisRepository) SetBalance(ctx context.Context, balances []entity.MerchantBalance) error {
	if len(balances) == 0 {
		return fmt.Errorf("list balance is empty")
	}

	pipe := r.client.TxPipeline()

	data := make(map[string]string, len(balances))
	for _, v := range balances {
		data[v.MerchantCode] = strconv.FormatFloat(v.Balance, 'f', 2, 64)
	}

	pipe.HSet(ctx, KeyBalance, data)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to cache list balance to redis, with error: %w", err)
	}

	return nil
}

func (r *redisRepository) SetPendingStatus(ctx context.Context, externalId string) error {
	key := "pending_transaction"

	if err := r.client.SAdd(ctx, key, externalId).Err(); err != nil {
		return fmt.Errorf("failed to add pending transaction in redis, with error %w", err)
	}

	return nil
}

func (r *redisRepository) FindByCodeAndChannel(ctx context.Context, merchantCode, channel string) (entity.FeeSettings, error) {
	key := fmt.Sprintf("%s:%s:%s", KeyFeeSettings, merchantCode, channel)

	data, err := r.client.HGetAll(ctx, key).Result()

	if err != nil {
		return entity.FeeSettings{}, fmt.Errorf("error on redis server while retrieving list fee merchant, with error %w", err)
	}

	if len(data) == 0 {
		return entity.FeeSettings{}, fmt.Errorf("key not found or hash empty")
	}

	floatParser := func(data string) float64 {
		f, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return 0
		}
		return f
	}

	feeSetting := entity.FeeSettings{}
	feeSetting.FeePartner = floatParser(data["fee_partner"])
	feeSetting.FeeService = floatParser(data["fee_service"])
	feeSetting.FeeTax = floatParser(data["fee_tax"])
	feeSetting.AdditionalFee = floatParser(data["additional_fee"])
	feeSetting.TotalCharge = floatParser(data["total_charge"])
	feeSetting.MerchantCode = merchantCode
	feeSetting.Channel = channel

	return feeSetting, nil
}

func (r *redisRepository) FindByMerchantCode(ctx context.Context, merchantCode string) (float64, error) {

	amount, err := r.client.HGet(ctx, KeyBalance, merchantCode).Result()

	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("key or field not found in redis")
		}

		return 0, fmt.Errorf("error in redis server while retrieving merchant balance, with error: %w", err)
	}

	return strconv.ParseFloat(amount, 64)
}

func (r *redisRepository) UpdateBalance(ctx context.Context, merchantCode, amount string) error {
	if err := r.client.HSet(ctx, KeyBalance, merchantCode, amount).Err(); err != nil {
		return fmt.Errorf("failed to update balance in redis, with error: %w", err)
	}
	return nil
}

func (r *redisRepository) RefundBalance(ctx context.Context, merchantCode string, amount float64) error {
	decrement := r.client.HIncrByFloat(ctx, KeyBalance, merchantCode, amount)
	if decrement.Err() != nil {
		return fmt.Errorf("failed to update balance in redis, with error: %w", decrement.Err())
	}
	return nil
}

func (r *redisRepository) DeletePendingStatus(ctx context.Context, externalId string) error {
	key := "pending_transaction"

	if err := r.client.SRem(ctx, key, externalId).Err(); err != nil {
		return fmt.Errorf("failed to remove pending transaction in redis, with error %w", err)
	}

	return nil
}
