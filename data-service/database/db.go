package database

import (
	"os"

	"github.com/RanggaNehemia/golang-microservices/data-service/models"
	"github.com/RanggaNehemia/golang-microservices/data-service/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Panic("No .env file found", zap.Error(err))
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Logger.Error("Failed to connect to database", zap.Error(err))
	}

	db.AutoMigrate(&models.Price{})
	DB = db
	utils.Logger.Info("Data Database Migrated")
}
