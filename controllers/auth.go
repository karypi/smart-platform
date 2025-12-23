package controllers

import (
	"net/http"
	"smart-platform/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(c *gin.Context) {
	// nil 表示不传递任何数据给模板, 不需要从后端获取动态数据. 常用于纯静态页面: 登录页, 关于页, 帮助页等
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	// 获取前端提交的Form 表单取字段用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 使用 GORM 在 users 表中查找用户名匹配的第一条记录, First(&user) 会把结果 populate 到 &user 结构体里
	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "User doesn't exist.")
		return
	}

	// 对比是否和数据库中存的哈希密码匹配 (Bcrypt 加密验证), 数据库里存的是加密后的密码，而不是明文，因此需要使用 bcrypt 来验证
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.String(http.StatusBadRequest, "Password is wrong.")
		return
	}

	// 登录成功，设置 cookie (模拟session), 有效期 3600秒. "/" 表示整个网站都能访问这个 cookie, 这个 cookie 用来判断用户是否已登录,如果没 cookie → 说明没登录 → 跳到登录页
	// c.SetCookie("user_id", strconv.Itoa(int(user.ID)), 3600, "/", "", false, true)

	// 登录时, 用 session 保存用户登录信息, 服务器会把 user_id 写入加密的 cookie, 用户即使查看 cookie，也看不到明文 ID，无法篡改。
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	// 登录成功，跳到主页
	c.Redirect(http.StatusSeeOther, "/")
}

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 对生成的密码加密, GenerateFromPassword 自动加 salt，非常安全
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 构造 User Struct
	user := models.User{Username: username, Password: string(hashed)}

	// 写入数据库. 类似于 INSERT INTO users (username, password) VALUES (?, ?)  Error 用于检查是否出错(重复用户名等)
	if err := models.DB.Create(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "Register failed: %v", err)
		return
	}

	// 注册成功, 跳转到登录页
	c.Redirect(http.StatusSeeOther, "/login")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}

// AuthMiddleware 用于保护那些需要登录才能访问的页面
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 cookie 中读取 user_id ,只有登录成功后，才会拥有这个 cookie, 未登录用户不会有
		// userID, err := c.Cookie("user_id")

		// 从 session 中读取
		session := sessions.Default(c)
		userID := session.Get("user_id")

		// 没有 cookie,说明没登录，自动跳转到登录页面
		if userID == nil {
			c.Redirect(http.StatusSeeOther, "/login") //Http 状态码，值是 303, 请求成功，但请客户端用 GET 请求去获取目标 URL, 常用场景：表单提交后重定向到另一个页面
			// 阻止后续的 handler 执行
			c.Abort()
			return
		}

		// 如果登录成功 -> 放行
		c.Next()
	}
}
