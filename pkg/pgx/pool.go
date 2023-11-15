package pgx

import (
	ewrap "autoposting/pkg/err-wrapper"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Pool struct {
	Conn *pgxpool.Pool
}

type Config struct {
	DSN            string
	LatencyObserve latencyObserveFn
	WithDebug      bool
}

func NewPool(config *Config, logger *slog.Logger) (*Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, ewrap.Errorf("failed to parse database dsn: %w", err)
	}

	pgxConfig.ConnConfig.Tracer = NewInterceptor(logger, config.LatencyObserve, config.WithDebug)

	// QueryExecModeSimpleProtocol need for PG connections through PGBouncer to avoid PreparedStatements cache issues
	pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, ewrap.Errorf("failed to create database connection: %w", err)
	}

	return &Pool{conn}, nil
}
