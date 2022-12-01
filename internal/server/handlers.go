package server

import (
	"autoposting/internal/social_account"
	logging "autoposting/pkg"
	"net/http"
)

func GetAuthUrlHandler(socialAccountService *social_account.SocialAccountService, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		network := r.URL.Query().Get("network")
		authURL, err := socialAccountService.GetAuthURL(network)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}
		_, _ = w.Write([]byte(authURL))
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
