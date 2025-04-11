package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/RanggaNehemia/golang-microservices/trade-service/database"
	"github.com/RanggaNehemia/golang-microservices/trade-service/models"
	"github.com/RanggaNehemia/golang-microservices/trade-service/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type TradeInput struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func fetchLowestPrice() (float64, error) {
	token, err := utils.GetMachineToken()
	if err != nil {
		return 0, err
	}

	req, _ := http.NewRequest("GET", os.Getenv("DATA_SERVICE_URL")+"/data/lowest", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		utils.Logger.Error("Error on fetching lowest price")
		return 0, fmt.Errorf("error on fetching lowest price: %d", resp.StatusCode)
	}

	var body struct {
		Value float64 `json:"Value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		utils.Logger.Error("Failed to parse JSON", zap.Error(err))
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
		utils.Logger.Error("Error on fetching lowest price", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if input.Price < lowestPrice/2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Price must be more than or equals to %.2f", lowestPrice/2)})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.Logger.Warn("Missing user token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User token required"})
		return
	}

	// Convert string to uint
	userIDStr, ok := userIDVal.(string)
	if !ok {
		utils.Logger.Warn("Invalid user ID format")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse user ID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	userID := uint(userIDUint64)

	// Now use userID (as uint) safely
	trade := models.Trade{
		UserID:   userID,
		Price:    input.Price,
		Quantity: input.Quantity,
	}

	if err := database.DB.Create(&trade).Error; err != nil {
		utils.Logger.Error("Failed to save trade", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save trade"})
		return
	}

	utils.Logger.Info("Trade placed", zap.Uint("trade", trade.ID))
	c.JSON(http.StatusOK, gin.H{"message": "Trade placed", "trade": trade})
}
