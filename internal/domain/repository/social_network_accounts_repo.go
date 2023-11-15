package repository

import (
	"autoposting/internal/adapters"
	"autoposting/internal/domain/model"
	"context"
)

type SocialNetworkAccountsRepository interface {
	CreateAccount(context.Context, *model.SocialNetworkAccount) error
	FindAccounts(context.Context, adapters.FindSocialNetworkAccountQuery) ([]model.SocialNetworkAccount, error)
	UpdateAccount(context.Context, *model.SocialNetworkAccount) (*model.SocialNetworkAccount, error)
	FindBySocialNetwork(context.Context, model.SocialNetworkName) (*model.SocialNetworkAccount, error)
}
