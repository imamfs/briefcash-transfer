package dbhelper

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBHelper struct {
	DB *gorm.DB
}

type DBConfig struct {
	Hostname string
	Port     string
	DBname   string
	Username string
	Password string
	SslMode  string
}

func NewDBConfig(cfg DBConfig) (*DBHelper, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Hostname, cfg.Port, cfg.Username, cfg.Password, cfg.DBname, cfg.SslMode,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to establish database connection %w", err)
	}

	sqlDb, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("failed to get generic database %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := sqlDb.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database %w", err)
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return &DBHelper{DB: db}, err
}

func (conn *DBHelper) Close() error {
	sqlDb, err := conn.DB.DB()

	if err == nil {
		return sqlDb.Close()
	}

	return fmt.Errorf("failed to close database connection %w", err)
}
