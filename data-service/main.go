package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/RanggaNehemia/golang-microservices/data-service/controllers"
	"github.com/RanggaNehemia/golang-microservices/data-service/database"
	"github.com/RanggaNehemia/golang-microservices/data-service/middleware"
	"github.com/RanggaNehemia/golang-microservices/data-service/models"
	"github.com/RanggaNehemia/golang-microservices/data-service/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	shutdown := tracing.InitTracer()
	defer shutdown()

	database.ConnectDatabase()

	router := gin.Default()
	router.Use(otelgin.Middleware("data-service"))

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				price := models.Price{
					Value: float64(rand.Intn(10000)) + rand.Float64(), // Random 0–10000.x
				}
				database.DB.Create(&price)
				fmt.Println("Generated price:", price.Value)
			}
		}
	}()

	protected := router.Group("/data")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.GET("/latest", controllers.GetLatestPrice)
	protected.GET("/lowest", controllers.GetLowestPrice)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // default port if not specified
	}

	router.Run(":" + port)
}
