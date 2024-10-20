package config

import (
	"errors"
	"log"
	"os"

	"github.com-Personal/go-fiber/internal/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	PORT        string
	HOST        string
}

// Load will load configuration from .env and Docker secrets.
func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using default values or environment variables.")
	}

	// databaseUrl, err := utils.getSecret("DATABASE_URL", "/run/secrets/DATABASE_URL")
	databaseUrl := utils.GetSecretOrEnv("DATABASE_URL")

	port := utils.GetSecretOrEnv("PORT")
	host := utils.GetSecretOrEnv("SERVER_HOST")

	if databaseUrl == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	return &Config{
		DatabaseURL: databaseUrl,
		PORT:        port,
		HOST:        host,
	}, nil
}

// getEnv fetches environment variables with a fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
