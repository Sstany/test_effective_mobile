package main

import (
	"fmt"
	"os"

	"subscrioption-service/internal/app"
	"subscrioption-service/internal/config"
)

func main() {
	address := os.Getenv("SERVER_PORT")
	dbHost := os.Getenv("DATABASE_HOST")
	dbName := os.Getenv("DATABASE_NAME")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbUser := os.Getenv("DATABASE_USER")
	dbPort := os.Getenv("DATABASE_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	logLevel := config.ErrorLevel

	debug := os.Getenv("DEBUG")
	if debug == "true" {
		logLevel = config.DebugLevel
	}

	cfg, err := config.New(address, connStr, logLevel)
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
