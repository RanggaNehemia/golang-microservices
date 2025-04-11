package database

import (
	"log"
	"os"

	"github.com/RanggaNehemia/golang-microservices/auth-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func InitTestDB() {
	dsn := os.Getenv("GORM_TEST_DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test Postgres: %v", err)
	}

	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	db.AutoMigrate(&models.User{})

	TestDB = db
	DB = db
}

func CloseTestDB() {
	sqlDB, _ := TestDB.DB()
	sqlDB.Close()
}
