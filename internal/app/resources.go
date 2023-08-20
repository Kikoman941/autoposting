package app

import (
	"autoposting/internal/app/registry"
	"autoposting/internal/app/usecase"
	"autoposting/internal/domain/model"
	"autoposting/internal/domain/service"
	"autoposting/internal/infrastructure/postgres"
	"autoposting/internal/infrastructure/social_network_client"
	ewrap "autoposting/pkg/err-wrapper"
	"autoposting/pkg/logger"
	pg_bun "autoposting/pkg/pg-bun"
	"context"
	"github.com/uptrace/bun"
)

func NewContainer(ctx context.Context, config *Config) (*registry.Container, error) {
	logger, err := logger.NewLogger(
		&logger.Options{
			IsProd:           config.IsProd,
			LogLevel:         config.LogLevel,
			EnableStacktrace: false,
		},
	)

	if err != nil {
		return nil, ewrap.Errorf("failed to init logger: %w", err)
	}

	postgresClient, err := pg_bun.NewDB(
		ctx,
		pg_bun.DBConfig{
			Dsn:                config.PostgresDSN,
			MaxOpenConnections: 100,
			MaxIdleConnections: 100,
			Hooks: []bun.QueryHook{
				pg_bun.NewLoggerHook(logger),
			},
		},
	)
	if err != nil {
		return nil, ewrap.Errorf("cannot get postgres client: %w", err)
	}

	socialNetworkClients := map[model.SocialNetworkName]social_network_client.SocialNetworkClient{
		"VK": social_network_client.NewVKClient(),
		"OK": social_network_client.NewOKClient(),
		"FB": social_network_client.NewFBClient(),
	}

	socialNetworkAccountService := service.NewService(
		logger,
		postgres.NewSocialNetworkAccountsRepository(postgresClient),
		postgres.NewSocialNetworkPagesRepository(postgresClient),
		socialNetworkClients,
	)

	container := registry.Container{
		Logger: logger,
		Usecases: &registry.Usecases{
			SocialNetwork: usecase.NewSocialNetworkUsecase(socialNetworkAccountService),
		},
	}

	return &container, nil
}
