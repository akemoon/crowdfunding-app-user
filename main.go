// @title           User Service API
// @version         1.0
// @description     CRUD API for user management
// @host            localhost:10000
// @BasePath        /

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	envPostgresDSN = "POSTGRES_DSN"
	envHTTPPort    = "HTTP_PORT"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp(cfg)

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

func loadConfig() (AppConfig, error) {
	dsn, err := getRequiredEnv(envPostgresDSN)
	if err != nil {
		return AppConfig{}, err
	}

	port := strings.TrimSpace(os.Getenv(envHTTPPort))
	if port == "" {
		port = "10000"
	}

	return AppConfig{
		PostgresDSN: dsn,
		HTTPPort:    port,
	}, nil
}
