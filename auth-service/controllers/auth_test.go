package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RanggaNehemia/golang-microservices/auth-service/database"
	"github.com/RanggaNehemia/golang-microservices/auth-service/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {
	err := godotenv.Load("../.env.test")
	if err != nil {
		println("No .env.test file found, continuing...")
	}
	database.InitTestDB()
	defer database.CloseTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/register", Register)

	// Prepare payload
	payload := models.User{Username: "foo", Password: "bar123"}
	body, _ := json.Marshal(payload)

	// Perform request
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "User registered", resp["message"])
}
