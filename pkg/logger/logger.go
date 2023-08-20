package logger

import (
	ewrap "autoposting/pkg/err-wrapper"
	"log/slog"
	"os"
)

type Options struct {
	IsProd           bool
	LogLevel         string
	EnableStacktrace bool
}

func NewLogger(options *Options) (*slog.Logger, error) {
	var handler interface{}
	var level slog.Level
	err := level.UnmarshalText([]byte(options.LogLevel))
	if err != nil {
		return nil, ewrap.Errorf("failed to build logger: %w", err)
	}

	if options.IsProd {
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: options.EnableStacktrace,
				Level:     level,
			},
		)
	} else {
		handler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: options.EnableStacktrace,
				Level:     level,
			},
		)
	}

	return slog.New(handler.(slog.Handler)), nil
}
