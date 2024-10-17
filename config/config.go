package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	PORT        string
	HOST        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using default values or environment variables.")
	}

	databaseUrl := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres")
	port := getEnv("PORT", "8080")
	host := getEnv("SERVER_HOST", "localhost")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	return &Config{
		DatabaseURL: databaseUrl,
		PORT:        port,
		HOST:        host,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
