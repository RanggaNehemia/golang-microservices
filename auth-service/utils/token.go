package utils

import (
	"os"

	"github.com/joho/godotenv"
)

var SecretKey []byte

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		Logger.Panic("No .env file found")
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		Logger.Error("SECRET_KEY not set in environment")
	}

	SecretKey = []byte(secret)
}
