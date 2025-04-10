package routes

import (
	"github.com/RanggaNehemia/golang-microservices/auth-service/controllers"
	"github.com/RanggaNehemia/golang-microservices/auth-service/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")

	// Public routes
	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)
	auth.POST("/client/token", controllers.ClientLogin)

	// Protected route
	auth.GET("/me", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		userID := c.MustGet("user_id")
		username := c.MustGet("username")

		c.JSON(200, gin.H{
			"user_id":  userID,
			"username": username,
		})
	})
}
