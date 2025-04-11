package utils

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	SecretKey       string
	PGXDatabaseURL  string
	TokenTTL        time.Duration
	RefreshTokenTTL time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		Logger.Panic("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		Logger.Error("SECRET_KEY is required")
	}

	pgxURL := os.Getenv("PGX_DATABASE_URL")
	if pgxURL == "" {
		Logger.Error("PGX_DATABASE_URL is required")
	}

	return &Config{
		Port:            port,
		SecretKey:       secret,
		PGXDatabaseURL:  pgxURL,
		TokenTTL:        time.Hour,      // 1h
		RefreshTokenTTL: 24 * time.Hour, // 24h
	}
}
