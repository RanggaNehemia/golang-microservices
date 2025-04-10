package routes

import (
	"github.com/RanggaNehemia/golang-microservices/trade-service/controllers"
	"github.com/RanggaNehemia/golang-microservices/trade-service/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterTradeRoutes(router *gin.Engine) {
	trade := router.Group("/trade")
	trade.Use(middleware.RequireUserToken())

	trade.POST("/place", controllers.PlaceTrade)
}
