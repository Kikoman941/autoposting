package handlers

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"log/slog"
)

func NewGraphqlHandler(
	logger *slog.Logger,
	schema graphql.ExecutableSchema,
	isProd bool,
) *handler.Server {
	srv := handler.NewDefaultServer(schema)
	return srv
}
