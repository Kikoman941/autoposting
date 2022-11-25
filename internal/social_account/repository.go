package social_account

import "context"

type SocialAccountRepository interface {
	CreateAccount(context.Context, *SocialAccount) error
	GetOneAccount(context.Context, int) (SocialAccount, error)
	FindAccountByNetwork(context.Context, string) (SocialAccount, error)
	GetGroup(context.Context, int, string) (Group, error)
}
