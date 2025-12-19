package main

import (
	"briefcash-transfer/config"
	"briefcash-transfer/internal/controller"
	"briefcash-transfer/internal/helper/dbhelper"
	"briefcash-transfer/internal/helper/kafkahelper"
	"briefcash-transfer/internal/helper/loghelper"
	"briefcash-transfer/internal/helper/redishelper"
	"briefcash-transfer/internal/repository"
	repositoryredis "briefcash-transfer/internal/repository/repository-redis"
	"briefcash-transfer/internal/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	loghelper.InitLogger("./resource/app.log", logrus.InfoLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
		sign := <-signalChannel
		loghelper.Logger.Infof("Received shutdown signal %s", sign.String())
		cancel()
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		loghelper.Logger.Fatal("Failed to load credential configuration")
	}

	dbConfig := dbhelper.DBConfig{
		Hostname: cfg.DBHost, Port: cfg.DBPort, DBname: cfg.DBName,
		Username: cfg.DBUsername, Password: cfg.DBPassword, SslMode: "disable",
	}

	dbCon, err := dbhelper.NewDBConfig(dbConfig)
	if err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to established connection to databases")
	}
	defer dbCon.Close()

	redisClient, err := redishelper.NewRedisHelper(cfg)
	if err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to established connection to redis server")
	}
	defer redisClient.Close()
	redsync := redishelper.NewRedsync(redisClient.Client)

	kafkaAddres := fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)
	kafkaService, err := kafkahelper.NewKafkaProducer([]string{kafkaAddres})
	if err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to establish kafka server")
	}

	redisRepo := repositoryredis.NewRedisRepository(redisClient.Client)
	balanceRepo := repository.NewBalanceRepository(dbCon.DB)
	recipientRepo := repository.NewRecipientRepository(dbCon.DB)
	feeSettingRepo := repository.NewFeeSetting(dbCon.DB)
	ledgerRepo := repository.NewLedgerRepository(dbCon.DB)
	merchantRepo := repository.NewMerchantBalanceRepository(dbCon.DB)
	partnerRepo := repository.NewPartnerRepository(dbCon.DB)
	transferRepo := repository.NewTransferRepository(dbCon.DB)

	partnerService := service.NewPartnerService(partnerRepo)
	if err := partnerService.LoadAllBankPartner(ctx); err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to load all bank partner configuration to memory")
	}

	redisService := service.NewRedisService(feeSettingRepo, merchantRepo, redisRepo, redsync)

	if err := redisService.LoadFeeSetting(ctx); err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to load fee setting to redis")
	}

	if err := redisService.LoadBalance(ctx); err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to load merchant balance to redis")
	}

	transferService := service.NewTransferService(recipientRepo, transferRepo, feeSettingRepo, ledgerRepo, balanceRepo, redisService, partnerService, dbCon.DB, kafkaService)

	transferController := controller.NewTransferController(transferService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(RequestLoggerMiddleware())

	api := router.Group("/api/v1")
	api.POST("/transfer", transferController.Transfer)

	server := &http.Server{
		Addr:    cfg.AppPort,
		Handler: router,
	}

	go func() {
		loghelper.Logger.Info("Transfer service is running...")
		if err := router.Run(cfg.AppPort); err != nil {
			loghelper.Logger.WithError(err).Fatal("Failed to start Transfer Service")
		}
	}()

	<-ctx.Done()

	loghelper.Logger.Info("Shutting down apps properly...")

	shutDownctx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutDownctx); err != nil {
		loghelper.Logger.WithError(err).Error("Forced shutdown due to timeout")
	} else {
		loghelper.Logger.Info("Transfer service shutdown completed")
	}

}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		duration := time.Since(start)
		status := c.Writer.Status()

		loghelper.Logger.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.FullPath(),
			"status":   status,
			"duration": duration.String(),
			"clientIp": c.ClientIP(),
		}).Info("Handled request")
	}
}
