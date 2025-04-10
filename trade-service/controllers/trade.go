package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/RanggaNehemia/golang-microservices/trade-service/database"
	"github.com/RanggaNehemia/golang-microservices/trade-service/models"

	"github.com/gin-gonic/gin"
)

type TradeInput struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func fetchLowestPrice() (float64, error) {
	machineToken := os.Getenv("TRADE_SERVICE_TOKEN")

	req, err := http.NewRequest("GET", "http://localhost:8081/data/lowest", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+machineToken)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("error on fetching lowest price: %d", resp.StatusCode)
	}

	var body struct {
		Value float64 `json:"Value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return body.Value, nil
}

func PlaceTrade(c *gin.Context) {
	var input TradeInput
	if err := c.ShouldBindJSON(&input); err != nil || input.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trade price"})
		return
	}

	lowestPrice, err := fetchLowestPrice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if input.Price < lowestPrice/2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Price must be more than or equals to %.2f", lowestPrice/2)})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User token required"})
		return
	}
	userID := userIDVal.(uint)

	trade := models.Trade{
		UserID:   userID,
		Price:    input.Price,
		Quantity: input.Quantity,
	}

	if err := database.DB.Create(&trade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save trade"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trade placed", "trade": trade})
}
