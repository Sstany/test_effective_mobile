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
	Address          string
	ConnectionString string
	Log              LogLevel
}

type DatabaseConfig struct {
	Connectionstring string
	MaxOpenConns     int32
	MaxLifetime      time.Duration
	MaxIdleTime      time.Duration
}

func New(
	address string,
	connStr string,
	logLevel LogLevel,
) (*Config, error) {
	return &Config{
		Address:          address,
		ConnectionString: connStr,
		Log:              logLevel,
	}, nil
}
