package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/akemoon/crowdfunding-app-user/modules/user/repo/user/postgres"
	"github.com/akemoon/crowdfunding-app-user/modules/user/service/user"
	platformHTTP "github.com/akemoon/crowdfunding-app-user/platform/http"
	pgLib "github.com/akemoon/golib/pglib"
)

type AppConfig struct {
	PostgresDSN string
	HTTPPort    string
}

type App struct {
	config AppConfig

	db *sql.DB

	userSvc *user.Service

	server platformHTTP.Server
}

func NewApp(config AppConfig) *App {
	return &App{config: config}
}

func (a *App) Init() error {
	var err error

	a.db, err = pgLib.Connect(context.Background(), a.config.PostgresDSN)
	if err != nil {
		return fmt.Errorf("postgres connection: %w", err)
	}

	repo := postgres.NewUserRepo(a.db)
	a.userSvc = user.NewService(repo)

	a.server = *platformHTTP.NewServer()
	a.server.AddUserHandlers(a.userSvc)

	return nil
}

func (a *App) Run() error {
	addr := ":" + a.config.HTTPPort
	log.Printf("http server listening on %s", addr)

	err := a.server.ListenAndServe(addr)
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
