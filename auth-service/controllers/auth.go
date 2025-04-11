package controllers

import (
	"net/http"

	"github.com/RanggaNehemia/golang-microservices/auth-service/database"
	"github.com/RanggaNehemia/golang-microservices/auth-service/models"
	"github.com/RanggaNehemia/golang-microservices/auth-service/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Logger.Error("Error on registering user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		utils.Logger.Error("Password hashing failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	user := models.User{Username: input.Username, Password: string(hashedPassword)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		utils.Logger.Error("Error on registering user", zap.Error(result.Error))
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	utils.Logger.Info("User registered", zap.String("username", user.Username))
	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}
