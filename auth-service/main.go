package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RanggaNehemia/golang-microservices/auth-service/database"
	"github.com/RanggaNehemia/golang-microservices/auth-service/models"
	"github.com/RanggaNehemia/golang-microservices/auth-service/routes"
	"github.com/RanggaNehemia/golang-microservices/auth-service/utils"
	"github.com/gin-gonic/gin"
	oauth2Errors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauth2Server "github.com/go-oauth2/oauth2/v4/server"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := utils.Load()

	database.ConnectDatabase()

	ctx := context.Background()

	pgxConn, err := pgx.Connect(ctx, cfg.PGXDatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pgxConn.Close(ctx)

	adapter := pgx4adapter.NewConn(pgxConn)

	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	// Token
	tokenStore, err := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	if err != nil {
		log.Fatalf("Failed to create token store: %v\n", err)
	}
	defer tokenStore.Close()
	manager.MapTokenStorage(tokenStore)

	//Client
	clientStore, err := pg.NewClientStore(adapter)
	if err != nil {
		log.Fatalf("Failed to create client store: %v\n", err)
	}
	manager.MapClientStorage(clientStore)

	// JWT token generator
	manager.MapAccessGenerate(utils.NewCustomJWTAccessGenerate(
		[]byte(cfg.SecretKey),
		jwt.SigningMethodHS512,
	))

	utils.SeedOAuthClients(ctx, pgxConn)

	// Create the OAuth2 server with default configuration.
	srv := oauth2Server.NewServer(oauth2Server.NewConfig(), manager)

	srv.SetClientInfoHandler(func(r *http.Request) (id, secret string, err error) {
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm error: %v", err)
		}
		id = r.Form.Get("client_id")
		secret = r.Form.Get("client_secret")
		log.Printf("ClientInfoHandler sees client_id=%q, client_secret=%q", id, secret)
		return id, secret, nil
	})

	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		ctx = context.WithValue(ctx, utils.ClientIDKey, clientID)

		var user models.User
		if result := database.DB.First(&user, "username = ?", username); result.Error != nil {
			return "", oauth2Errors.ErrInvalidGrant
		}
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			return "", oauth2Errors.ErrInvalidGrant
		}
		userID = fmt.Sprint(user.ID)
		return userID, nil
	})
	srv.SetInternalErrorHandler(func(err error) *oauth2Errors.Response {
		log.Printf("OAuth2 Internal Error: %v", err)
		return nil
	})
	srv.SetResponseErrorHandler(func(re *oauth2Errors.Response) {
		log.Printf("OAuth2 Response Error: %v", re.Error)
	})

	// HTTP
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// — Health check —
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	routes.RegisterAuthRoutes(r)

	oauth := r.Group("/oauth")
	{
		oauth.POST("/token", func(c *gin.Context) {
			srv.HandleTokenRequest(c.Writer, c.Request)
		})
		oauth.GET("/authorize", func(c *gin.Context) {
			srv.HandleAuthorizeRequest(c.Writer, c.Request)
		})

		// — Revocation endpoint (RFC 7009) —
		oauth.POST("/revoke", func(c *gin.Context) {
			token := c.PostForm("token")
			hint := c.PostForm("token_type_hint")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
				return
			}
			// Try to revoke by the hinted type, or fall back
			var err error
			switch hint {
			case "access_token":
				err = tokenStore.RemoveByAccess(c.Request.Context(), token)
			case "refresh_token":
				err = tokenStore.RemoveByRefresh(c.Request.Context(), token)
			default:
				err = tokenStore.RemoveByAccess(c.Request.Context(), token)
				if err != nil {
					err = tokenStore.RemoveByRefresh(c.Request.Context(), token)
				}
			}
			// Always return 200 even if token wasn't found per spec
			c.Status(http.StatusOK)
		})

		// — Introspection endpoint (RFC 7662) —
		oauth.POST("/introspect", func(c *gin.Context) {
			token := c.PostForm("token")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
				return
			}
			ti, err := tokenStore.GetByAccess(c.Request.Context(), token)
			if err != nil || ti == nil {
				// inactive token
				c.JSON(http.StatusOK, gin.H{"active": false})
				return
			}
			// Check expiry
			active := ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).After(time.Now())
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

	// — Protected example —
	r.GET("/auth/me", func(c *gin.Context) {
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

	r.Run(":" + cfg.Port)
}
