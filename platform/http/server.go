package http

import (
	"net/http"

	authHandler "github.com/akemoon/crowdfunding-app-user/modules/auth/api/handler"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/auth"
	userHandler "github.com/akemoon/crowdfunding-app-user/modules/user/api/handler"
	"github.com/akemoon/crowdfunding-app-user/modules/user/service/user"
	"github.com/akemoon/golib/myhttp"
	"github.com/akemoon/golib/myhttp/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	r *myhttp.Router
}

func NewServer() *Server {
	return &Server{
		r: myhttp.NewRouter().Use(
			middleware.BaseMetrics(),
		),
	}
}

func (s *Server) AddAuthHandlers(svc *auth.Service) {
	s.r.HandleFunc("POST /signup", authHandler.SignUp(svc))
	s.r.HandleFunc("POST /signin", authHandler.SignIn(svc))
	s.r.HandleFunc("POST /signout", authHandler.SignOut(svc))
	s.r.HandleFunc("GET /check", authHandler.CheckAccessToken(svc))
}

func (s *Server) AddUserHandlers(svc *user.Service) {
	s.r.HandleFunc("GET /users/{id}", userHandler.GetUserByID(svc))
	s.r.HandleFunc("GET /users/me", userHandler.GetMe(svc))
	s.r.HandleFunc("PATCH /users/me/profile", userHandler.UpdateProfile(svc))

	s.r.HandleFunc("POST /users/{id}/follow", userHandler.Follow(svc))
	s.r.HandleFunc("DELETE /users/{id}/follow", userHandler.Unfollow(svc))
}

func (s *Server) AddSwaggerUI() {
	s.r.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}

func (s *Server) AddMetrics() {
	s.r.Handle("/metrics", promhttp.Handler())
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.r.Handler())
}
