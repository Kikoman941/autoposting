package social_network_client

type SocialNetworkClient interface {
	GetAuthURL(string) (string, error)
	GetAccessToken(string, map[string][]string) (string, error)
	UploadImage()
	CreatePost(string, string, string) (string, error)
	DeletePost()
}
