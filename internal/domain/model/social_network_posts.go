package model

import "time"

type SocialNetworkPost struct {
	ID          int
	Page        int
	PostData    *PostData
	PublishedAt time.Time
}

type PostData struct {
	Text  string
	Image string
	Url   string
}
