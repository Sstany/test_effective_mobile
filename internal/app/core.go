package app

import (
	"context"

	"subscrioption-service/internal/adapter/db"
	"subscrioption-service/internal/adapter/repo"
	"subscrioption-service/internal/app/usecase"
	"subscrioption-service/internal/config"
	handler "subscrioption-service/internal/controller/http"

	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	MaxOpenConns = 25
	MaxIdleTime  = 10
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

	logger.Info("start")

	defer func() {
		if lErr := logger.Sync(); lErr != nil {
			panic(err)
		}
	}()

	dbCfg := config.DatabaseConfig{
		Connectionstring: cfg.ConnectionString,
		MaxOpenConns:     MaxOpenConns,
		MaxIdleTime:      MaxIdleTime,
		MaxLifetime:      MaxLifetime,
	}

	pool, err := db.NewPostgresPool(context.Background(), dbCfg)
	if err != nil {
		panic(err)
	}

	subRepo, err := repo.NewSubscription(pool, logger.Named("subscription-repo"))
	if err != nil {
		panic(err)
	}

	subUsecase, err := usecase.NewSubscription(subRepo, pool, logger.Named("subscription-usecase"))
	if err != nil {
		panic(err)
	}

	handler.NewServer(cfg.Address, subUsecase, pool, logger.Named("http")).Start()
}
