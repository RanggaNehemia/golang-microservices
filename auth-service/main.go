package main

import (
	"log"
	"os"

	"github.com/RanggaNehemia/golang-microservices/auth-service/database"
	"github.com/RanggaNehemia/golang-microservices/auth-service/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	database.ConnectDatabase()

	routes.RegisterAuthRoutes(r)

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}

	r.Run(":" + port)
}
