package pg_bun

import (
	"context"
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
	"log/slog"
	"time"
)

type LoggerHook struct {
	logger *slog.Logger
}

func NewLoggerHook(
	logger *slog.Logger,
) *LoggerHook {
	return &LoggerHook{
		logger: logger,
	}
}

func (h *LoggerHook) BeforeQuery(
	ctx context.Context,
	_ *bun.QueryEvent,
) context.Context {
	return ctx
}

func (h *LoggerHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	err := event.Err
	logFields := slog.Group(
		"postgres",
		slog.Any("err", err),
		slog.String("query", event.Query),
		slog.String("duration", time.Since(event.StartTime).String()),
		slog.String("operation", event.Operation()),
	)

	if event.Err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.logger.Error("SQL error ", logFields)
	} else {
		h.logger.Debug("SQL ", logFields)
	}
}
