package config

import "embed"

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
	Host              string
	Connectionstring  string
	Log               LogLevel
	EmbededMigrations embed.FS
}

func New(host string, connstr string, logLevel LogLevel, migrations embed.FS) (*Config, error) {
	return &Config{
		Host:              host,
		Connectionstring:  connstr,
		Log:               logLevel,
		EmbededMigrations: migrations,
	}, nil

}
