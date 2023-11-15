package model

type SocialNetworkAccount struct {
	ID            int
	SocialNetwork SocialNetworkName
	Credentials   string
	AccessToken   *AccessToken
}

type AccessToken struct {
	Token     string
	ExpiresIn string
}
