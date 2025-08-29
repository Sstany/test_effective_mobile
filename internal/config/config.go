package config

import "time"

type LogLevel int8

const (
	DebugLevel LogLevel = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

type Config struct {
	Address  string
	Log      LogLevel
	DBConfig DatabaseConfig
}

type DatabaseConfig struct {
	ConnectionString string
	MaxOpenConns     int32
	MaxLifetime      time.Duration
	MaxIdleTime      time.Duration
}

func New(
	address string,
	connStr string,
	logLevel LogLevel,
	maxOpenConns int32,
	maxLifeTime time.Duration,
	maxIdleTime time.Duration,
) (*Config, error) {
	return &Config{
		Address: address,
		Log:     logLevel,
		DBConfig: DatabaseConfig{
			ConnectionString: connStr,
			MaxOpenConns:     maxOpenConns,
			MaxLifetime:      maxLifeTime,
			MaxIdleTime:      maxIdleTime,
		},
	}, nil
}
