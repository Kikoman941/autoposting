package repository

import (
	"autoposting/internal/domain/model"
	"context"
)

type SocialNetworkPagesRepository interface {
	CreatePage(context.Context, *model.SocialNetworkPage) error
	FindPage()
}
