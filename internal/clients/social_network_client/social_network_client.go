package social_network_client

type SocialNetworkClient interface {
	UploadImage()
	CreatePost(string, string, string) (string, error)
	DeletePost()
}
