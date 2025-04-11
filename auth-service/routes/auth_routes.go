package routes

import (
	"github.com/RanggaNehemia/golang-microservices/auth-service/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")

	// Public routes
	auth.POST("/register", controllers.Register)
}
