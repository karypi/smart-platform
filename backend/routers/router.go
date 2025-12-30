package routers

import (
	"smart-platform/backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/logout", controllers.Logout)
		api.POST("/alert/send", middleware.Auth(), controllers.SendAlert)
		api.GET("/alert/history", middleware.Auth(), controllers.AlertHistory)
	}
}
