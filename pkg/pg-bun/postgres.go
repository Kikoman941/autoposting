package pg_bun

import (
	ewrap "autoposting/pkg/err-wrapper"
	"context"
	"database/sql"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DBConfig struct {
	Dsn                string
	ServiceName        string
	MaxOpenConnections int
	MaxIdleConnections int
	Hooks              []bun.QueryHook
}

func (c DBConfig) validate() error {
	err := validation.ValidateStruct(
		&c,
		validation.Field(&c.Dsn, validation.Required),
		validation.Field(&c.MaxOpenConnections, validation.Required, validation.Min(1)),
		validation.Field(&c.MaxIdleConnections, validation.Required, validation.Min(1)),
	)

	if err != nil {
		return ewrap.Errorf("failed to validate DB config: %w", err)
	}

	return nil
}

func NewDB(
	ctx context.Context,
	config DBConfig,
) (*bun.DB, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	sqlDB := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithDSN(config.Dsn),
			pgdriver.WithTLSConfig(nil),
			pgdriver.WithApplicationName(fmt.Sprintf("[%s]", config.ServiceName)),
		),
	)

	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)

	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)

	db := bun.NewDB(sqlDB, pgdialect.New())

	for _, hook := range config.Hooks {
		db.AddQueryHook(hook)
	}

	if err := db.Ping(); err != nil {
		return nil, ewrap.Errorf("failed to ping db: %w", err)
	}

	go func() {
		<-ctx.Done()
		if err := db.Close(); err != nil {
			fmt.Printf("Failed to close postgres connection: %s\n", err)
		}
		fmt.Println("\nPostgres connection successfully closed")
	}()

	return db, nil
}
