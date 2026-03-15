package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	tokenRepo "github.com/akemoon/crowdfunding-app-user/modules/auth/repo/token/redis"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/metrics"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/auth"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/token"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/tool/hasher/bcrypt"
	"github.com/akemoon/crowdfunding-app-user/modules/user/repo/user/postgres"
	"github.com/akemoon/crowdfunding-app-user/modules/user/service/user"
	"github.com/akemoon/crowdfunding-app-user/platform/http"
	platformRedis "github.com/akemoon/crowdfunding-app-user/platform/redis"
	userPublisher "github.com/akemoon/crowdfunding-app-user/publisher/user"
	pgLib "github.com/akemoon/golib/postgres"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

const defaultHTTPAddr = ":80"

type AppConfig struct {
	PostgresURL           string
	PostgresMigrationsDir string
	RedisURL              string
	JWTSecret             string
	AvatarsBaseURL        string
	KafkaBrokers          []string
	UserTopic             string
}

type App struct {
	ctx context.Context

	config AppConfig

	db *sql.DB

	redisClient *redis.Client

	publisher *userPublisher.Publisher

	authSvc *auth.Service

	userSvc *user.Service

	server http.Server
}

func NewApp(ctx context.Context, config AppConfig) *App {
	return &App{
		ctx:    ctx,
		config: config,
	}
}

func (a *App) InitDB() error {
	var err error

	a.db, err = pgLib.Connect(a.ctx, a.config.PostgresURL)
	if err != nil {
		return fmt.Errorf("postgres connection: %w", err)
	}

	err = pgLib.Migrate(a.ctx, a.db, a.config.PostgresMigrationsDir)
	if err != nil {
		return fmt.Errorf("postgres migration: %w", err)
	}

	a.redisClient, err = platformRedis.NewClient(a.ctx, a.config.RedisURL)
	if err != nil {
		return fmt.Errorf("redis connection: %w", err)
	}

	return nil
}

func (a *App) InitServices() error {
	publisher, err := userPublisher.NewPublisher(a.config.KafkaBrokers, a.config.UserTopic)
	if err != nil {
		return fmt.Errorf("init kafka publisher: %w", err)
	}
	a.publisher = publisher

	repo := postgres.NewUserRepo(a.db)

	tokenRepo := tokenRepo.NewRefreshTokenRepo(a.redisClient)
	tokenSvc := token.NewService(tokenRepo, a.config.JWTSecret)

	hasher := bcrypt.NewHasher(10)

	a.authSvc = auth.NewService(repo, hasher, tokenSvc, publisher)

	a.userSvc = user.NewService(repo, a.config.AvatarsBaseURL)

	return nil
}

func (a *App) InitServer() {
	authMetrics := metrics.NewAuthMetrics(prometheus.DefaultRegisterer)
	a.server = *http.NewServer()
	a.server.AddAuthHandlers(a.authSvc, authMetrics)
	a.server.AddUserHandlers(a.userSvc)
	a.server.AddMetrics()
}


func (a *App) Init() error {
	err := a.InitDB()
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}

	err = a.InitServices()
	if err != nil {
		return fmt.Errorf("init services: %w", err)
	}

	a.InitServer()

	return nil
}

func (a *App) Run() error {
	log.Printf("http server listening on %s", defaultHTTPAddr)

	err := a.server.ListenAndServe(defaultHTTPAddr)
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
