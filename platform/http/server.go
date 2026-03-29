package http

import (
	"net/http"

	_ "github.com/akemoon/crowdfunding-app-user/docs"
	userHandler "github.com/akemoon/crowdfunding-app-user/modules/user/api/handler"
	"github.com/akemoon/crowdfunding-app-user/modules/user/service/user"
	"github.com/akemoon/golib/httplib"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	r *httplib.Router
}

func NewServer() *Server {
	return &Server{
		r: httplib.NewRouter(),
	}
}

func (s *Server) AddUserHandlers(svc *user.Service) {
	s.r.HandleFunc("POST /users", userHandler.CreateUser(svc))
	s.r.HandleFunc("GET /users", userHandler.ListUsers(svc))
	s.r.HandleFunc("GET /users/{id}", userHandler.GetUserByID(svc))
	s.r.HandleFunc("PUT /users/{id}", userHandler.UpdateProfile(svc))
	s.r.HandleFunc("POST /users/{id}/follow", userHandler.Follow(svc))
	s.r.HandleFunc("GET /users/{id}/follows", userHandler.GetFollows(svc))
	s.r.HandleFunc("DELETE /users/{id}/follow", userHandler.Unfollow(svc))

	s.r.Handle("/swagger/", httpSwagger.WrapHandler)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.r.Handler())
}
