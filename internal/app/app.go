package app

import (
	"autoposting/internal/app/registry"
	"autoposting/internal/app/server"
	ewrap "autoposting/pkg/err-wrapper"
	"context"
	"log/slog"
)

type App struct {
	config    *Config
	container *registry.Container
	server    *server.Server
}

func Run(ctx context.Context, config *Config) error {
	app := App{
		config: config,
	}

	cnt, err := NewContainer(ctx, config)
	if err != nil {
		return ewrap.Errorf("failed to create container: %w", err)
	}
	app.container = cnt

	cnt.Logger.Debug("Init config", slog.Any("config", config))

	app.initAppServer()
	beforeShutdown := func() {}
	if err := app.server.Run(ctx, config.ServerAddr, config.IsProd, beforeShutdown); err != nil {
		return ewrap.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}

func (app *App) initAppServer() {
	app.server = server.NewServer(app.container.Logger)
	app.server.InitRoutes(app.container, app.config.IsProd)
}
