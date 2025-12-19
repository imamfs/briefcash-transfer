package service

import (
	"briefcash-transfer/internal/entity"
	"briefcash-transfer/internal/helper/loghelper"
	"briefcash-transfer/internal/repository"
	repositoryredis "briefcash-transfer/internal/repository/repository-redis"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type TransferRedisService interface {
	LoadFeeSetting(ctx context.Context) error
	LoadBalance(ctx context.Context) error
	GetFeeSetting(ctx context.Context, merchantCode, channel string, log *logrus.Entry) (entity.FeeSettings, error)
	SetFeeSetting(ctx context.Context, feeSetting entity.FeeSettings, log *logrus.Entry) error
	DebitBalance(ctx context.Context, merchantCode string, amount float64, log *logrus.Entry) (float64, error)
	RefundBalance(ctx context.Context, merchantCode string, amount float64, log *logrus.Entry) error
}

type transferRedisService struct {
	feeRepository      repository.FeeSettingRepository
	merchantRepository repository.MerchantBalanceRepository
	redisRepository    repositoryredis.RedisRepository
	locker             *redsync.Redsync
}

func NewRedisService(feeRepository repository.FeeSettingRepository, merchantBalanceRepository repository.MerchantBalanceRepository, transferRedisRepository repositoryredis.RedisRepository, locker *redsync.Redsync) TransferRedisService {
	return &transferRedisService{feeRepository, merchantBalanceRepository, transferRedisRepository, locker}
}

func (r *transferRedisService) LoadFeeSetting(ctx context.Context) error {
	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":   "redis_service",
		"operation": "cache_fee_setting",
	})

	// find all service fee in database
	log.Info("Fetching all fee setting from database")
	feeList, err := r.feeRepository.FindAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to fetch fee setting from database")
		return err
	}

	// load all service fee to redis
	log.Info("Caching fee setting to redis")
	if err := r.redisRepository.SetListFee(ctx, feeList); err != nil {
		log.WithError(err).Error("Failed to cache fee setting to redis")
		return err
	}

	log.Info("Fee settings succesfully cached")
	return nil
}

func (r *transferRedisService) LoadBalance(ctx context.Context) error {
	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":   "redis_service",
		"operation": "cache_balance",
	})

	// Query all active merchants balance
	log.Info("Fetching merchant balance from database")
	balanceList, err := r.merchantRepository.FindAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to fetch merchant balance from database")
		return err
	}

	// Load all active merchants balance to redis
	log.Info("Caching merchant balance to redis")
	if err := r.redisRepository.SetBalance(ctx, balanceList); err != nil {
		log.WithError(err).Error("Failed to cache merchant balance to redis")
		return err
	}

	log.Info("Merchant balance successfully cached")
	return nil
}

func (r *transferRedisService) GetFeeSetting(ctx context.Context, merchantCode, channel string, log *logrus.Entry) (entity.FeeSettings, error) {
	// Get service fee based on merchant code and payment channel
	log.Infof("Fetch fee setting for merchant: %s and channel: %s", merchantCode, channel)
	fee, err := r.redisRepository.FindByCodeAndChannel(ctx, merchantCode, channel)
	if err != nil {
		log.WithError(err).Error("Failed to fetch fee setting from redis")
		return entity.FeeSettings{}, err
	}
	return fee, nil
}

func (r *transferRedisService) SetFeeSetting(ctx context.Context, feeSetting entity.FeeSettings, log *logrus.Entry) error {
	// Create new service fee to redis
	log.Infof("Caching fee setting for merchant: %s and channel: %s", feeSetting.MerchantCode, feeSetting.Channel)
	if err := r.redisRepository.SetFee(ctx, feeSetting); err != nil {
		log.WithError(err).Error("Failed to cache fee setting in redis")
		return err
	}
	return nil
}

func (r *transferRedisService) DebitBalance(ctx context.Context, merchantCode string, amount float64, log *logrus.Entry) (float64, error) {
	// Set redis lock
	keyLock := "lock:balance:" + merchantCode
	mutex := r.locker.NewMutex(
		keyLock,
		redsync.WithExpiry(300*time.Millisecond),
		redsync.WithTries(2),
	)

	if err := mutex.Lock(); err != nil {
		return 0, fmt.Errorf("failed to start lock balance in redis %v", err)
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	// get current merchant balance from redis
	log.Info("Fetching merchant balance from redis")
	currentBalance, err := r.redisRepository.FindByMerchantCode(ctx, merchantCode)
	if err == redis.Nil {
		log.WithError(err).Errorf("Merchant %s not found in redis", merchantCode)
		return 0, err
	}

	if err != nil {
		log.WithError(err).Error("Failed to fetch merchant balance from redis")
		return 0, err
	}

	// check sufficiency merchant balance
	log.Info("Checking sufficiency merchant balance")
	if currentBalance < amount {
		log.Warnf("Insufficient balance: merchant: %s, balance: %.2f, needed %.2f", merchantCode, currentBalance, amount)
		return -1, fmt.Errorf("insufficient merchant balance, current balance is %f", currentBalance)
	}

	// calculate balance
	newBalance := currentBalance - amount

	if newBalance < 0 {
		log.Warnf("Insufficient balance: merchant: %s, balance: %.2f", merchantCode, currentBalance)
		return -1, fmt.Errorf("insufficient merchant balance, current balance is %f", currentBalance)
	}

	// update merchant CalculateBalanceRedis
	log.Infof("Updating new balance to redis for merchant: %s", merchantCode)
	if err := r.redisRepository.UpdateBalance(ctx, merchantCode, strconv.FormatFloat(newBalance, 'f', 2, 64)); err != nil {
		log.WithError(err).Errorf("Failed to update new balance to redis for merchant: %s", merchantCode)
		return 0, err
	}

	return newBalance, nil
}

func (r *transferRedisService) RefundBalance(ctx context.Context, merchantCode string, amount float64, log *logrus.Entry) error {
	keyLock := "lock:balance:" + merchantCode
	mutex := r.locker.NewMutex(
		keyLock,
		redsync.WithExpiry(300*time.Millisecond),
		redsync.WithTries(2),
	)

	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to start lock balance in redis %v", err)
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	log.Infof("Refund merchant %s balance in redis, with %f returned to balance", merchantCode, amount)
	if err := r.redisRepository.RefundBalance(ctx, merchantCode, amount); err != nil {
		log.WithError(err).Error("Failed to refund merchant balance in redis")
		return err
	}
	log.Info("Refund balance successfully executed")
	return nil
}
