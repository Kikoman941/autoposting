package server

import (
	"autoposting/internal/app/registry"
	"autoposting/internal/app/server/handlers"
	"autoposting/internal/presentation/graphql/gen"
	"autoposting/internal/presentation/graphql/resolver"
	ewrap "autoposting/pkg/err-wrapper"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	logger *slog.Logger
	router *chi.Mux
}

func NewServer(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
		router: chi.NewRouter(),
	}
}

func (s *Server) InitRoutes(container *registry.Container, isProd bool) {
	graphqlHandler := handlers.NewGraphqlHandler(
		container.Logger,
		gen.NewExecutableSchema(
			gen.Config{
				Resolvers: resolver.NewResolver(
					container.Usecases,
				),
			},
		),
		isProd,
	)
	s.router.Handle("/graphql", graphqlHandler)
	s.router.Handle("/auth/get_token", handlers.GetAccessTokenHandler(container, container.Logger))
}

func (s *Server) Run(
	ctx context.Context,
	addr string,
	waitOfShutdown bool,
	shutdownInitiated func(),
) error {
	srv := http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	srv.RegisterOnShutdown(shutdownInitiated)

	go func() {
		<-ctx.Done()

		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Error("failed to shutdown server: %w", err)
		}

		const sleepShutdown = 5 * time.Second
		if waitOfShutdown {
			time.Sleep(sleepShutdown)
		}
	}()

	s.logger.Info("Run server", slog.String("addr", addr))

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return ewrap.Errorf("failed to Run: %w", err)
	}

	return nil
}
