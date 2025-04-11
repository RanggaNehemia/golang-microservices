package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/RanggaNehemia/golang-microservices/trade-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var SecretKey []byte
var expectedAud string

func init() {
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Panic("No .env file found", zap.Error(err))
	}
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		utils.Logger.Fatal("SECRET_KEY not set in environment")
	}
	SecretKey = []byte(secret)
	expectedAud = os.Getenv("WEB_CLIENT_ID")
}

func RequireUserToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Logger.Warn("Missing or bad auth header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or bad auth header"})
			return
		}
		tokenString := parts[1]

		// Local signature + exp check
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			utils.Logger.Warn("Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Validity Check
		resp, err := http.PostForm(
			os.Getenv("AUTH_URL")+"/oauth/introspect",
			url.Values{"token": {tokenString}},
		)
		if err != nil {
			utils.Logger.Warn("Introspection failed", zap.Error(err))
			c.AbortWithStatusJSON(500, gin.H{"error": "Introspection failed"})
			return
		}
		defer resp.Body.Close()

		var body struct {
			Active bool `json:"active"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			utils.Logger.Error("Bad introspection response", zap.Error(err))
			c.AbortWithStatusJSON(500, gin.H{"error": "Bad introspection response"})
			return
		}
		if !body.Active {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token revoked"})
			return
		}

		// Audience Check
		claims := token.Claims.(jwt.MapClaims)
		if aud, _ := claims["aud"].(string); aud != expectedAud {
			utils.Logger.Warn("Wrong audience", zap.String("Audience", aud))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Wrong audience"})
			return
		}
		c.Set("user_id", claims["sub"])
		c.Next()
	}
}
