package model

import "github.com/uptrace/bun"

type SocialNetworkPage struct {
	bun.BaseModel `bun:"table:social_network_pages"`
	ID            int                    `bun:"id,pk,autoincrement"`
	AccountID     int                    `bun:"account_id"`
	Project       string                 `bun:"project"`
	PageID        string                 `bun:"page_id"`
	PageInfo      *SocialNetworkPageInfo `bun:"page_info"`
	AccessToken   *AccessToken           `bun:"access_token,nullzero"`
}

type SocialNetworkPageInfo struct {
	Name         string `json:"title"`
	Description  string `json:"description"`
	PreviewImage string `json:"previewImage"`
}
