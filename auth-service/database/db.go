package database

import (
	"os"

	"github.com/RanggaNehemia/golang-microservices/auth-service/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Panic("No .env file found")
	}

	dsn := os.Getenv("GORM_DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}

	DB = db
	utils.Logger.Info("Authentication Database Migrated")
}
