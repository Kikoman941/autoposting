package registry

import (
	"autoposting/internal/app/usecase"
	"log/slog"
)

type Container struct {
	Logger   *slog.Logger
	Usecases *Usecases
}

type Usecases struct {
	SocialNetwork *usecase.SocialNetworkUsecase
}
