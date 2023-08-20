package model

import "github.com/uptrace/bun"

type SocialNetworkAccount struct {
	bun.BaseModel `bun:"table:social_network_accounts"`
	ID            int               `bun:"id,pk,autoincrement"`
	SocialNetwork SocialNetworkName `bun:"social_network"`
	Credentials   string            `bun:"credentials"`
	AccessToken   *AccessToken      `bun:"access_token,nullzero"`
}

type AccessToken struct {
	Token     string
	ExpiresIn string
}
