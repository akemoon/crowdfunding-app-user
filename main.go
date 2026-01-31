package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/akemoon/crowdfunding-app-user/api"
	_ "github.com/akemoon/crowdfunding-app-user/docs"
	pgRepo "github.com/akemoon/crowdfunding-app-user/repo/user/postgres"
	"github.com/akemoon/crowdfunding-app-user/service/user"
	"github.com/akemoon/golib/postgres"
)

const (
	migrationsDir = "/app/migrations/postgres"

	envPostgresDSN = "POSTGRES_DSN"
)

// @title User Service API
// @version 1.0
// @description User Service API for a Crowdfunding App.
func main() {
	mainCtx := context.Background()

	db, err := initDB(mainCtx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	repo := pgRepo.NewUserRepo(db)
	svc := user.NewService(repo)

	srv := api.NewServer()
	srv.AddUserHandlers(svc)
	srv.AddSwaggerUI()
	srv.AddMetrics()

	log.Printf("listen on port 80")

	srv.ListenAndServe(":80")
}

func initDB(ctx context.Context) (*sql.DB, error) {
	dsn := os.Getenv(envPostgresDSN)
	if dsn == "" {
		log.Fatalf("env var %s is empty", envPostgresDSN)
	}

	db, err := postgres.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	err = postgres.Migrate(ctx, db, migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("db migration error: %w", err)
	}

	return db, nil
}
