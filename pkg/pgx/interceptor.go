package pgx

import (
	"context"
	"git.sholding.ru/gm/go-packages/pkg/logger"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type latencyObserveFn func(duration time.Duration, status int, label string)

type sqlQueryPayload struct {
	timeStart time.Time
	sql       string
	args      []any
	label     string
}

type (
	ctxKey int
)

const (
	ctxKeySQLPayload ctxKey = iota
	ctxKeyLabel
)

type Interceptor struct {
	logger         logger.ILogger
	latencyObserve latencyObserveFn
	withDebug      bool
}

func NewInterceptor(
	logger logger.ILogger,
	latencyObserveFn latencyObserveFn,
	withDebug bool,
) *Interceptor {
	return &Interceptor{
		logger,
		latencyObserveFn,
		withDebug,
	}
}

func (i *Interceptor) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	return context.WithValue(ctx, ctxKeySQLPayload, &sqlQueryPayload{ // nolint:staticcheck,revive
		timeStart: time.Now(),
		sql:       data.SQL,
		args:      data.Args,
		label:     getLabel(ctx),
	})
}

func (i *Interceptor) TraceQueryEnd(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	queryPayload, ok := ctx.Value(ctxKeySQLPayload).(*sqlQueryPayload)
	if !ok {
		i.logger.ErrorContext(ctx, "failed to trac sql: sqlQueryPayload is null", slog.Any("data", data))

		return
	}

	duration := time.Since(queryPayload.timeStart)

	if data.Err != nil {
		i.logger.ErrorContext(
			ctx,
			"failed to do sql",
			slog.String("query", strings.ReplaceAll(queryPayload.sql, "\n", " ")),
			slog.Any("args", queryPayload.args),
			slog.Duration("duration", duration),
			slog.Any("err", data.Err),
		)

		i.latencyObserve(duration, http.StatusInternalServerError, queryPayload.label)

		return
	}

	i.latencyObserve(duration, http.StatusOK, queryPayload.label)

	if !i.withDebug {
		return
	}

	i.logger.DebugContext(
		ctx,
		"SQL",
		slog.String("query", strings.ReplaceAll(queryPayload.sql, "\n", " ")),
		slog.Any("args", queryPayload.args),
		slog.Duration("duration", duration),
	)
}

func AddLabel(ctx context.Context, label string) context.Context {
	return context.WithValue(ctx, ctxKeyLabel, label)
}

func getLabel(ctx context.Context) string {
	if v := ctx.Value(ctxKeyLabel); v != nil {
		if l, ok := v.(string); ok {
			return l
		}
	}

	return "NONE"
}
