// auth-service/main_test.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/RanggaNehemia/golang-microservices/auth-service/database"
	"github.com/RanggaNehemia/golang-microservices/auth-service/models"
	"github.com/RanggaNehemia/golang-microservices/auth-service/routes"
	"github.com/RanggaNehemia/golang-microservices/auth-service/utils"
	"github.com/gin-gonic/gin"
	oauth2Errors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauth2Models "github.com/go-oauth2/oauth2/v4/models"
	oauth2Server "github.com/go-oauth2/oauth2/v4/server"
	oauth2Store "github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	srv    *oauth2Server.Server
	router *gin.Engine
	ctx    = context.Background()
)

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env.test")

	// Init in‑memory GORM DB (only User table)
	database.InitTestDB()
	database.DB.AutoMigrate(&models.User{})
	defer database.CloseTestDB()

	// Build OAuth2 manager with in‑memory stores
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	// In‑memory token store
	memTokenStore, _ := oauth2Store.NewMemoryTokenStore()
	manager.MapTokenStorage(memTokenStore)

	// In‑memory client store + seed a client
	memClientStore := oauth2Store.NewClientStore()
	memClientStore.Set("webclient", &oauth2Models.Client{
		ID:     "webclient",
		Secret: "webclientsecret",
		Domain: "",
		UserID: "",
	})
	manager.MapClientStorage(memClientStore)

	// JWT generator
	manager.MapAccessGenerate(utils.NewCustomJWTAccessGenerate(
		[]byte(os.Getenv("SECRET_KEY")),
		jwt.SigningMethodHS512,
	))

	// OAuth2 server
	srv = oauth2Server.NewServer(oauth2Server.NewConfig(), manager)
	srv.SetClientInfoHandler(oauth2Server.ClientFormHandler)
	srv.SetPasswordAuthorizationHandler(passwordGrantHandler)
	srv.SetInternalErrorHandler(func(err error) *oauth2Errors.Response { return nil })
	srv.SetResponseErrorHandler(func(re *oauth2Errors.Response) {})

	// Gin router
	gin.SetMode(gin.TestMode)
	router = gin.New()
	router.Use(gin.Recovery())

	routes.RegisterAuthRoutes(router)

	oauth := router.Group("/oauth")
	{
		oauth.POST("/token", func(c *gin.Context) {
			srv.HandleTokenRequest(c.Writer, c.Request)
		})
		oauth.POST("/revoke", func(c *gin.Context) {
			token := c.PostForm("token")
			_ = memTokenStore.RemoveByAccess(ctx, token)
			_ = memTokenStore.RemoveByRefresh(ctx, token)
			c.Status(http.StatusOK)
		})
		oauth.POST("/introspect", func(c *gin.Context) {
			tkn := c.PostForm("token")
			ti, _ := memTokenStore.GetByAccess(ctx, tkn)
			active := ti != nil && ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).After(time.Now())
			resp := gin.H{"active": active}
			if active {
				resp["client_id"] = ti.GetClientID()
				resp["sub"] = ti.GetUserID()
				resp["scope"] = ti.GetScope()
				resp["iat"] = ti.GetAccessCreateAt().Unix()
				resp["exp"] = ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Unix()
			}
			c.JSON(http.StatusOK, resp)
		})
	}

	// Protected example
	router.GET("/auth/me", func(c *gin.Context) {
		ti, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"client_id":  ti.GetClientID(),
			"user_id":    ti.GetUserID(),
			"expires_in": int64(ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		})
	})

	os.Exit(m.Run())
}

// passwordGrantHandler checks username/password against GORM user table
func passwordGrantHandler(_ context.Context, clientID, username, password string) (string, error) {
	var u models.User
	if database.DB.First(&u, "username = ?", username).Error != nil {
		return "", oauth2Errors.ErrInvalidGrant
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", oauth2Errors.ErrInvalidGrant
	}
	return fmt.Sprint(u.ID), nil
}

// --- Now your tests follow exactly as before ---

func TestRegister_Success(t *testing.T) {
	// fresh user table
	database.DB.Exec("DELETE FROM users")

	payload := map[string]string{"username": "alice", "password": "pw123"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "User registered", resp["message"])
}

// ... and so on for Token, Me, Revoke, Introspect ...
