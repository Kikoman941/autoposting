package server

import (
	"autoposting/internal/social_account"
	logging "autoposting/pkg"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	logger *logging.Logger
	router *chi.Mux
}

func NewServer(logger *logging.Logger) *Server {
	return &Server{
		logger: logger,
		router: chi.NewRouter(),
	}
}

func (s *Server) ListenAndServe() {
	s.logger.Fatal(
		http.ListenAndServe(":8080", s.router),
	)
}

func (s *Server) InitRoutes(socialAccountService *social_account.SocialAccountService) {
	s.router.Get("/auth/get-auth-url", NewSocialAccountHandler(socialAccountService, s.logger))
	s.router.Get("/create-post", CreatePostHandler(socialAccountService, s.logger))
}
