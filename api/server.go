package api

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/api/handler"
	"github.com/akemoon/crowdfunding-app-user/service/user"
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

func (s *Server) AddUserHandlers(svc *user.Service) {
	s.r.HandleFunc("POST /user", handler.CreateUser(svc))
	s.r.HandleFunc("GET /user/id/{id}", handler.GetUserByID(svc))
	s.r.HandleFunc("GET /user/{username}", handler.GetUserByUsername(svc))
	s.r.HandleFunc("GET /user/me", handler.GetMe(svc))
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
