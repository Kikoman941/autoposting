package handlers

import (
	"autoposting/internal/app/registry"
	"context"
	"log/slog"
	"net/http"
)

func GetAccessTokenHandler(container *registry.Container, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		if err := container.Usecases.SocialNetwork.GetAccessToken(context.TODO(), queryParams); err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte("done"))
	}
}
