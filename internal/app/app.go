package app

import (
	"amplifr/internal/clients/db"
	"amplifr/internal/clients/social_network_client"
	"amplifr/internal/config"
	"amplifr/internal/server"
	"amplifr/internal/social_account"
	socialAccountStorage "amplifr/internal/social_account/storage"
	logging "amplifr/pkg"
	"context"
)

type App struct {
	config *config.Config
	logger *logging.Logger
}

func NewApp(config *config.Config, logger *logging.Logger) *App {
	return &App{
		config: config,
		logger: logger,
	}
}

func (a *App) Run() {
	postgresqlClient, err := db.NewClient(context.TODO(), a.config.PostgresqlDSN)
	if err != nil {
		a.logger.Fatalf("Cannot get postgresql client:\n%s", err)
	}
	socialNetworkAccounts := map[string]social_network_client.SocialNetworkClient{
		"VK": social_network_client.NewVKClient(),
		"OK": social_network_client.NewOKClient(),
		"FB": social_network_client.NewFBClient(),
	}
	socialAccountRepository := socialAccountStorage.NewRepository(postgresqlClient)
	socialAccountService := social_account.NewService(
		socialAccountRepository,
		socialNetworkAccounts,
	)
	s := server.NewServer(a.logger)
	s.InitRoutes(socialAccountService)
	s.ListenAndServe()
}
