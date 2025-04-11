package utils

import (
	"log"
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
		log.Println("No .env file found, relying on system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("SECRET_KEY is required")
	}

	pgxURL := os.Getenv("PGX_DATABASE_URL")
	if pgxURL == "" {
		log.Fatal("PGX_DATABASE_URL is required")
	}

	// You can also make these configurable via env
	return &Config{
		Port:            port,
		SecretKey:       secret,
		PGXDatabaseURL:  pgxURL,
		TokenTTL:        time.Hour,      // 1h
		RefreshTokenTTL: 24 * time.Hour, // 24h
	}
}
