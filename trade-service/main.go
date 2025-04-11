package main

import (
	"os"

	"github.com/RanggaNehemia/golang-microservices/trade-service/database"
	"github.com/RanggaNehemia/golang-microservices/trade-service/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	database.InitDB()

	router := gin.Default()

	routes.RegisterTradeRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // default port if not specified
	}

	router.Run(":" + port)
}
