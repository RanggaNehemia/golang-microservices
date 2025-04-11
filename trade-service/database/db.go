package database

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/RanggaNehemia/golang-microservices/trade-service/models"
	"github.com/RanggaNehemia/golang-microservices/trade-service/utils"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Panic("No .env file found", zap.Error(err))
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Logger.Fatal("Failed to connect to database:", zap.Error(err))
	}

	db.AutoMigrate(&models.Trade{})
	DB = db
}
