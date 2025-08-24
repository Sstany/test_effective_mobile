package main

import (
	"os"
	"subscrioption-service/internal/app/core"
	"subscrioption-service/internal/config"
	"subscrioption-service/migrations"
)

func main() {
	host := os.Getenv("HOST")
	connstr := os.Getenv("CONNECTION_STRING")

	logLevel := config.ErrorLevel

	debug := os.Getenv("DEBUG")
	if debug == "true" {
		logLevel = config.DebugLevel
	}
	cfg, err := config.New(host, connstr, logLevel, migrations.EmbedMigrations)
	if err != nil {
		panic(err)
	}

	core.Run(cfg)

}
