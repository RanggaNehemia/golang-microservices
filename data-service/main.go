package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/RanggaNehemia/golang-microservices/data-service/controllers"
	"github.com/RanggaNehemia/golang-microservices/data-service/database"
	"github.com/RanggaNehemia/golang-microservices/data-service/middleware"
	"github.com/RanggaNehemia/golang-microservices/data-service/models"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()

	router := gin.Default()

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

	router.Run(":5081")
}
