package config

import (
	"briefcash-transfer/internal/helper/loghelper"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	DBName     string
	AppPort    string
	RedisHost  string
	RedisPort  string
	KafkaHost  string
	KafkaPort  string
	KafkaTopic string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		loghelper.Logger.WithError(err).Error("Error while reading .env file config, now using system environment variables")
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_URL"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		RedisHost:  os.Getenv("REDIS_ADDRESS"),
		RedisPort:  os.Getenv("REDIS_PORT"),
		KafkaHost:  os.Getenv("KAFKA_HOST"),
		KafkaPort:  os.Getenv("KAFKA_PORT"),
		KafkaTopic: os.Getenv("KAFKA_TOPIC"),
		AppPort: func() string {
			if value := os.Getenv("APP_PORT"); value != "" {
				return value
			}
			return ":8080"
		}(),
	}

	if cfg.DBHost == "" {
		loghelper.Logger.Error("DB_HOST not found in environment")
		return nil, fmt.Errorf("DB_HOST not found in environment")
	}

	if cfg.RedisHost == "" {
		loghelper.Logger.Error("REDIS_HOST not found in environment")
		return nil, fmt.Errorf("REDIS_HOST not found in environment")
	}

	if cfg.KafkaHost == "" {
		loghelper.Logger.Error("KAFKA_HOST not found in environment")
		return nil, fmt.Errorf("KAFKA_HOST not found in environment")
	}

	return cfg, nil
}
