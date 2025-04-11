package controllers

import (
	"time"

	"github.com/RanggaNehemia/golang-microservices/data-service/database"
	"github.com/RanggaNehemia/golang-microservices/data-service/models"
	"github.com/RanggaNehemia/golang-microservices/data-service/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetLatestPrice(c *gin.Context) {
	var latestPrice models.Price
	result := database.DB.Order("created_at DESC").First(&latestPrice)
	if result.Error != nil {
		utils.Logger.Error("Could not get price", zap.Error(result.Error))
		c.JSON(500, gin.H{"error": "Could not get price"})
		return
	}

	c.JSON(200, latestPrice)
}

func GetLowestPrice(c *gin.Context) {
	var lowestPrice models.Price
	timeLimit := time.Now().Add(-24 * time.Hour)
	result := database.DB.Where("created_at > ?", timeLimit).Order("value ASC").First(&lowestPrice)
	if result.Error != nil {
		utils.Logger.Error("Could not get price", zap.Error(result.Error))
		c.JSON(500, gin.H{"error": "Could not get price"})
		return
	}

	c.JSON(200, lowestPrice)
}
