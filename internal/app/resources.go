package app

import (
	"autoposting/internal/adapters"
	"autoposting/internal/app/registry"
	"autoposting/internal/app/usecase"
	"autoposting/internal/domain/model"
	"autoposting/internal/domain/service"
	"autoposting/internal/infrastructure/social_network_client"
	"autoposting/internal/infrastructure/sqlc-pg/dao"
	ewrap "autoposting/pkg/err-wrapper"
	"autoposting/pkg/logger"
	"autoposting/pkg/pgx"
	"context"
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

	db, err := pgx.NewPool(&pgx.Config{
		DSN: config.PostgresDSN,
	}, logger)
	if err != nil {
		return nil, ewrap.Errorf("cannot get postgres connections pool: %w", err)
	}

	sqlcQueries := dao.New(db.Conn)

	socialNetworkClients := map[model.SocialNetworkName]social_network_client.SocialNetworkClient{
		"VK": social_network_client.NewVKClient(),
		"OK": social_network_client.NewOKClient(),
		"FB": social_network_client.NewFBClient(),
	}

	socialNetworkAccountService := service.NewService(
		logger,
		adapters.NewSocialNetworkAccountsRepository(sqlcQueries),
		adapters.NewSocialNetworkPagesRepository(sqlcQueries),
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
