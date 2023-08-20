package repository

import (
	"autoposting/internal/domain/model"
	"autoposting/internal/infrastructure/postgres"
	"context"
)

type SocialNetworkAccountsRepository interface {
	CreateAccount(context.Context, *model.SocialNetworkAccount) error
	FindAccounts(context.Context, postgres.FindSocialNetworkAccountQuery) ([]model.SocialNetworkAccount, error)
	UpdateAccount(context.Context, *model.SocialNetworkAccount) (*model.SocialNetworkAccount, error)
	FindBySocialNetwork(context.Context, model.SocialNetworkName) (*model.SocialNetworkAccount, error)
}
