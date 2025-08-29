package main

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"subscription-service/internal/app"
	"subscription-service/internal/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	address := os.Getenv("SERVER_PORT")
	dbConnectionString := os.Getenv("DATABASE_CONNECTION_STRING")
	maxOpenConnsRaw := os.Getenv("MAX_OPEN_CONNS")
	maxOpenConns, err := strconv.ParseInt(maxOpenConnsRaw, 10, 32)
	if err != nil {
		panic(err)
	}

	maxIdleTimeRaw := os.Getenv("MAX_IDLE_TIME")
	maxIdleTime, err := time.ParseDuration(maxIdleTimeRaw)
	if err != nil {
		panic(err)
	}

	maxLifetimeRaw := os.Getenv("MAX_LIFE_TIME")
	maxLifeTime, err := time.ParseDuration(maxLifetimeRaw)
	if err != nil {
		panic(err)
	}

	logLevel := config.InfoLevel

	debug := os.Getenv("DEBUG")
	if debug == "true" {
		logLevel = config.DebugLevel
	}

	cfg, err := config.New(
		address,
		dbConnectionString,
		logLevel,
		int32(maxOpenConns),
		maxIdleTime,
		maxLifeTime,
	)
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
