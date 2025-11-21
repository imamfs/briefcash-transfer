package redishelper

import (
	"briefcash-transfer/config"
	"briefcash-transfer/internal/helper/loghelper"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHelper struct {
	Client *redis.Client
}

func NewRedisHelper(cfg config.Config) (*RedisHelper, error) {
	if cfg.RedisHost == "" || cfg.RedisPort == "" {
		loghelper.Logger.Info("Redis host or port is not configured in environment")
		return nil, fmt.Errorf("redis host or port is empty")
	}
	address := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:         address,
		DB:           0,
		PoolSize:     50,
		MinIdleConns: 10,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		loghelper.Logger.WithError(err).Error("Failed to connect to redis")
		return nil, fmt.Errorf("failed to ping redis server")
	}

	loghelper.Logger.Info("Redis server connected successfully")

	return &RedisHelper{Client: client}, nil
}

func (r *RedisHelper) Close() error {
	return r.Client.Close()
}
