package model

type SocialNetworkPage struct {
	ID          int
	AccountID   int
	Project     string
	PageID      string
	PageInfo    *SocialNetworkPageInfo
	AccessToken *AccessToken
}

type SocialNetworkPageInfo struct {
	Name         string `json:"title"`
	Description  string `json:"description"`
	PreviewImage string `json:"previewImage"`
}
