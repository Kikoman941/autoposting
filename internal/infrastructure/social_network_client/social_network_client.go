package social_network_client

type SocialNetworkClient interface {
	GetAuthURL(string) (string, error)
	GetAccessToken(string, map[string][]string) (string, error)
	GetAccountPages(string, string) ([]SocialNetworkPage, error)
	UploadImage()
	CreatePost(string, string, string) (string, error)
	DeletePost()
}

type SocialNetworkPage struct {
	ID          string
	Name        string
	Description string
	Image       string
}
