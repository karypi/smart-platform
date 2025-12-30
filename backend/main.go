package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// func main() {
// 	// Init MySql DB
// 	models.InitDB()

// 	r := gin.Default()
// 	// Session store, 32字节密钥用于加密签名, "mysession" 是 cookie 名称, super-secret-key 是服务器密钥，用于 cookie 签名
// 	store := cookie.NewStore([]byte("super-secret-key"))
// 	store.Options(sessions.Options{
// 		Path:     "/",  // 决定 Cookie 在哪些 URL 生效
// 		MaxAge:   3600, // Cookie 保留时长
// 		HttpOnly: true,
// 		Secure:   false, // ❗ HTTP 一定要 false
// 	})
// 	r.Use(sessions.Sessions("mysession", store))

// 	// User Login/Register
// 	r.POST("/api/login", controllers.Login)
// 	r.POST("/api/register", controllers.Register)
// 	r.POST("/api/logout", controllers.Logout)

// 	// 受保护路由
// 	// r.Group("/") 创建一个Route Group,并增加中间件 AuthMiddleware(), 这个中间件用于检查用户是否登录
// 	auth := r.Group("/")
// 	// 任何访问这些路由都会先通过 AuthMiddleware 进行检查 (类似 Flask 的装饰器或 Django 的 @login_required )
// 	auth.Use(controllers.AuthMiddleware())
// 	{
// 		auth.POST("/api/alert/send", controllers.SendAlert)
// 		auth.GET("/api/alert/history", controllers.AlertHistory)
// 	}

// 	r.Run(":5000")
// }

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("smart-session", store))

	routers.RegisterRoutes(r)
	r.Run(":5000")
}
