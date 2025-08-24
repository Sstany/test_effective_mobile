package core

import (
	"subscrioption-service/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

}
