package main

import (
	"smart-platform/routers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("smart-session", store))

	routers.RegisterRoutes(r)
	r.Run(":5000")
}
