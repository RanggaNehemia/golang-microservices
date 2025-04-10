package main

import (
	"os"

	"github.com/RanggaNehemia/golang-microservices/trade-service/database"
	"github.com/RanggaNehemia/golang-microservices/trade-service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	router := gin.Default()

	routes.RegisterTradeRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // default port if not specified
	}

	router.Run(":" + port)
}
