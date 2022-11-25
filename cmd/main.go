package main

import (
	"amplifr/internal/app"
	"amplifr/internal/config"
	logging "amplifr/pkg"
)

func main() {
	conf := config.NewConfig()
	logger := logging.Init(conf.IsProd)

	logger.Info("Creating application")
	a := app.NewApp(conf, logger)
	logger.Info("Run application")
	a.Run()
}
