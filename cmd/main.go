package main

import (
	"autoposting/internal/app"
	"autoposting/internal/config"
	logging "autoposting/pkg"
)

func main() {
	conf := config.NewConfig()
	logger := logging.Init(conf.IsProd)

	logger.Info("Creating application")
	a := app.NewApp(conf, logger)
	logger.Info("Run application")
	a.Run()
}
