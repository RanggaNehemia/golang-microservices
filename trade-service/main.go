package main

import (
	"os"

	"github.com/RanggaNehemia/golang-microservices/trade-service/database"
	"github.com/RanggaNehemia/golang-microservices/trade-service/routes"
	"github.com/RanggaNehemia/golang-microservices/trade-service/tracing"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	shutdown := tracing.InitTracer()
	defer shutdown()

	_ = godotenv.Load()

	database.InitDB()

	router := gin.Default()

	router.Use(otelgin.Middleware("trade-service"))

	routes.RegisterTradeRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // default port if not specified
	}

	router.Run(":" + port)
}
