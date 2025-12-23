package main

import (
	"smart-platform/controllers"
	"smart-platform/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Init MySql DB
	models.InitDB()

	r := gin.Default()
	// Session store, 32字节密钥用于加密签名, "mysession" 是 cookie 名称, super-secret-key 是服务器密钥，用于 cookie 签名
	store := cookie.NewStore([]byte("super-secret-key"))
	store.Options(sessions.Options{
		Path:     "/",  // 决定 Cookie 在哪些 URL 生效
		MaxAge:   3600, // Cookie 保留时长
		HttpOnly: true,
		Secure:   false, // ❗ HTTP 一定要 false
	})
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("templates/*")
	// 将本地 ./static 目录映射到 HTTP 的 /static, 前端可以通过 /static 访问静态文件(JS/CSS 图片)
	r.Static("/static", "./static")

	// User Login/Register
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.Login)
	r.GET("/register", controllers.RegisterPage)
	r.POST("/register", controllers.Register)
	r.POST("/logout", controllers.Logout)

	// 受保护路由
	// r.Group("/") 创建一个Route Group,并增加中间件 AuthMiddleware(), 这个中间件用于检查用户是否登录
	auth := r.Group("/")

	// 任何访问这些路由都会先通过 AuthMiddleware 进行检查 (类似 Flask 的装饰器或 Django 的 @login_required )
	auth.Use(controllers.AuthMiddleware())
	{
		auth.GET("/", controllers.DashboardPage)
		auth.POST("/alert/send", controllers.SendAlert)
		auth.GET("/alert/history", controllers.AlertHistory)
	}

	r.Run(":5000")
}
