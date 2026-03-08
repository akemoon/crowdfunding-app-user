package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	envPostgresURL           = "POSTGRES_URL"
	envPostgresMigrationsDir = "POSTGRES_MIGRATIONS_DIR"
	envRedisURL              = "REDIS_URL"
	envJWTSecret             = "JWT_SECRET"
	envAvatarsBaseURL        = "AVATARS_BASE_URL"
)

// @title User service API
// @version 1.0
// @description User service API for a crowdfunding app.
func main() {
	mainCtx := context.Background()

	cfg, err := loadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp(mainCtx, cfg)

	err = app.Init()
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func getRequiredEnv(key string) (string, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return "", fmt.Errorf("env %s is empty", key)
	}
	return val, nil
}

func loadConfigFromEnv() (AppConfig, error) {
	postgresURL, err := getRequiredEnv(envPostgresURL)
	if err != nil {
		return AppConfig{}, err
	}

	migrationsDir, err := getRequiredEnv(envPostgresMigrationsDir)
	if err != nil {
		return AppConfig{}, err
	}

	redisURL, err := getRequiredEnv(envRedisURL)
	if err != nil {
		return AppConfig{}, err
	}

	jwtSecret, err := getRequiredEnv(envJWTSecret)
	if err != nil {
		return AppConfig{}, err
	}

	avatarsBaseURL, err := getRequiredEnv(envAvatarsBaseURL)
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		PostgresURL:           postgresURL,
		PostgresMigrationsDir: migrationsDir,
		RedisURL:              redisURL,
		JWTSecret:             jwtSecret,
		AvatarsBaseURL:        avatarsBaseURL,
	}, nil
}
