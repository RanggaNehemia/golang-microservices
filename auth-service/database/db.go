package database

import (
	"fmt"
	"log"
	"os"

	"github.com/RanggaNehemia/golang-microservices/auth-service/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Client{})
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	DB = db
	fmt.Println("Migrated")
}
