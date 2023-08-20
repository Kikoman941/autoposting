package main

import (
	"autoposting/internal/app"
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	baseContext, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP,
	)
	_ = cancel

	config, err := app.NewConfig()
	if err != nil {
		log.Fatal("Failed to init config app: ", err)
	}

	if err := app.Run(baseContext, config); err != nil {
		log.Fatal("failed to run app: ", err)
	}

	<-baseContext.Done()
}
