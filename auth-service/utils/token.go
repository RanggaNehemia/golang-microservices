package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SecretKey []byte

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("SECRET_KEY not set in environment")
	}

	SecretKey = []byte(secret)
}

func GenerateJWT(userID uint, username string, subject string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"sub":      subject,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // expires in 24 hours
	})

	return token.SignedString(SecretKey)
}
