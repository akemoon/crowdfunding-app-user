package http

import (
	"net/http"

	authHandler "github.com/akemoon/crowdfunding-app-user/modules/auth/api/handler"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/metrics"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/auth"
	userHandler "github.com/akemoon/crowdfunding-app-user/modules/user/api/handler"
	"github.com/akemoon/crowdfunding-app-user/modules/user/service/user"
	"github.com/akemoon/golib/httplib"
	"github.com/akemoon/golib/httplib/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	r *httplib.Router
}

func NewServer() *Server {
	return &Server{
		r: httplib.NewRouter().Use(
			middleware.BaseMetrics(),
		),
	}
}

func (s *Server) AddAuthHandlers(svc *auth.Service, m *metrics.AuthMetrics) {
	s.r.HandleFunc("POST /auth/signup", authHandler.SignUp(svc))
	s.r.HandleFunc("POST /auth/signin", authHandler.SignIn(svc, m))
	s.r.HandleFunc("POST /auth/signout", authHandler.SignOut(svc))
	s.r.HandleFunc("GET /auth/check", authHandler.CheckAccessToken(svc))
	s.r.HandleFunc("POST /auth/refresh", authHandler.Refresh(svc))
}

func (s *Server) AddUserHandlers(svc *user.Service) {
	s.r.HandleFunc("GET /users/{id}", userHandler.GetUserByID(svc))
	s.r.HandleFunc("GET /users/me", userHandler.GetMe(svc))
	s.r.HandleFunc("PATCH /users/me/profile", userHandler.UpdateProfile(svc))

	s.r.HandleFunc("POST /users/{id}/follow", userHandler.Follow(svc))
	s.r.HandleFunc("DELETE /users/{id}/follow", userHandler.Unfollow(svc))
}

func (s *Server) AddMetrics() {
	s.r.Handle("/metrics", promhttp.Handler())
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.r.Handler())
}
