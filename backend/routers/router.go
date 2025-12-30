package routers

import (
	"smart-platform/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/logout", controllers.Logout)
		api.POST("/alert/send", controllers.SendAlert)
		api.GET("/alert/history", controllers.AlertHistory)
	}
}
