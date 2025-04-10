package main

import (
	"github.com/RanggaNehemia/golang-microservices/trade-service/database"
	"github.com/RanggaNehemia/golang-microservices/trade-service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	router := gin.Default()

	routes.RegisterTradeRoutes(router)

	router.Run(":8082")
}
