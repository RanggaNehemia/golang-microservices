package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SecretKey []byte
var expectedAud string

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("SECRET_KEY not set in environment")
	}

	SecretKey = []byte(secret)

	expectedAud = os.Getenv("TRADE_SERVICE_CLIENT_ID")
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or bad auth header"})
			return
		}
		tokenString := parts[1]

		// Local signature + exp check
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Validity Check
		resp, err := http.PostForm(
			os.Getenv("AUTH_URL")+"/oauth/introspect",
			url.Values{"token": {tokenString}},
		)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Introspection failed"})
			return
		}
		defer resp.Body.Close()

		var body struct {
			Active bool `json:"active"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Wrong audience"})
			return
		}
		c.Next()
	}
}
