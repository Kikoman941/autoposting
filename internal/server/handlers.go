package server

import (
	"amplifr/internal/social_account"
	logging "amplifr/pkg"
	"net/http"
)

func NewSocialAccountHandler(socialAccountService *social_account.SocialAccountService, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		network := r.URL.Query().Get("network")
		credentials := r.URL.Query().Get("credentials")
		if err := socialAccountService.CreateAccount(network, credentials); err != nil {
			logger.Errorln(err)
			_, _ = w.Write([]byte("Cannot create new account"))
			return
		}
		_, _ = w.Write([]byte("Done"))
		return
	}
}

func CreatePostHandler(socialAccountService *social_account.SocialAccountService, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		network := r.URL.Query().Get("network")
		project := r.URL.Query().Get("project")
		post := r.URL.Query().Get("post")
		if err := socialAccountService.CreatePost(network, project, post); err != nil {
			logger.Errorln(err)
			_, _ = w.Write([]byte("Cannot create new post"))
			return
		}
		_, _ = w.Write([]byte("Done"))
		return
	}
}
