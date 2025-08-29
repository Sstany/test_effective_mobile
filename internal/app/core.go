package app

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"subscription-service/internal/adapter/db"
	"subscription-service/internal/adapter/repo"
	"subscription-service/internal/app/usecase"
	"subscription-service/internal/config"
	handler "subscription-service/internal/controller/http"
)

const (
	MaxOpenConns = 25
	MaxIdleTime  = 10 * time.Minute
	MaxLifetime  = 5 * time.Minute
)

func Run(cfg *config.Config) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zapcore.Level(cfg.Log))
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	logger, err := loggerConfig.Build(zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCaller())
	if err != nil {
		panic(err)
	}

	logger.Info("start service")

	defer func() {
		if lErr := logger.Sync(); lErr != nil {
			panic(lErr)
		}
	}()

	pool, err := db.NewPostgresPool(context.Background(), cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	subRepo, err := repo.NewSubscription(pool, logger.Named("subscription-repo"))
	if err != nil {
		panic(err)
	}

	transactionController := repo.NewTransactionSQL(pool, logger.Named("transaction-ctrl"))

	subUsecase, err := usecase.NewSubscription(subRepo, transactionController, logger.Named("subscription-usecase"))
	if err != nil {
		panic(err)
	}

	handler.NewServer(cfg.Address, subUsecase, pool, logger.Named("http")).Start()
}
