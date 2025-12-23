package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"smart-platform/models"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 把来自表单的字段组织成企业微信机器人期望的 JSON 格式（markdown 类型）
func makeAlertDataFromForm(title, severity, alertname, instance, message, status, startAt string) map[string]interface{} {

	// map[string]interface{} （可被 json.Marshal 直接序列化)
	return map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": fmt.Sprintf(
				"## %s\n\n"+
					">**Alert Severity**: %s\n\n"+
					">**Alert Type**: %s\n\n"+
					">**Alert Instance**: %s\n\n"+
					">**Alert Description**: %s\n\n"+
					">**Alert Status**: %s\n\n"+
					">**Alert Time**: %s\n",
				title, severity, alertname, instance, message, status, startAt,
			),
		},
	}
}

func SendAlert(c *gin.Context) {
	// 读取 user_id
	userIDstr, _ := c.Cookie("user_id")
	// 把 cookie(字符串)转为整数 userID, 用 _ 忽略错误
	userID, _ := strconv.Atoi(userIDstr)

	token := c.PostForm("token")
	alertname := c.PostForm("alertname")
	severity := c.PostForm("severity")
	instance := c.PostForm("instance")
	message := c.PostForm("message")
	status := c.DefaultPostForm("status", "firing")

	title := "【Problem】"
	if status == "resolved" {
		title = "【Resolve】"
	}

	// In(time.FixedZone("CST", 8*3600)) 转到 +08:00 时区（中国时区）
	startAt := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")
	sendData := makeAlertDataFromForm(title, severity, alertname, instance, message, status, startAt)

	// 将 sendData 序列化为 JSON 字节
	body, _ := json.Marshal(sendData)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+token,
		"application/json",

		// 将 body 包装成 io.Reader 数据流
		bytes.NewBuffer(body),
	)

	if err != nil {
		c.String(http.StatusInternalServerError, "Send failed: %v", err)
		return
	}

	// 在函数返回时关闭响应体，防止资源泄露
	defer resp.Body.Close()

	// 保存记录到数据库
	alertRecord := models.Alert{
		UserID:    uint(userID),
		Token:     token,
		AlertName: alertname,
		Severity:  severity,
		Instance:  instance,
		Message:   message,
		Status:    status,
		CreatedAt: time.Now(),
	}

	// 通过 GORM 插入数据库, GORM 要求传入指针，直接修改对象里的字段，如果传入值类型，GORM修改的是副本，外部 alert.ID 不会更新
	models.DB.Create(&alertRecord)

	// 请求结束后重定向到 Dashboard 主页, 常用于表单提交后重定向可以避免刷新重复提交
	c.Redirect(http.StatusSeeOther, "/")
}

// 获取历史记录
func AlertHistory(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not logged in",
		})
		return
	}

	// 声明一个 Alert 结构体切片,用来保存数据库查到的记录
	var alerts []models.Alert

	// 查询数据库
	models.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&alerts)
	c.JSON(200, alerts)
}

// Home Page
func DashboardPage(c *gin.Context) {
	// 获取 cookie 中的 user_id, 渲染模板, 并把 userID 注入模板上下文 (在模板里可以通过 {{ .userID }} 访问)
	// userIDStr, _ := c.Cookie("user_id")
	session := sessions.Default(c)
	userID := session.Get("user_id")

	// gin.H 是 mapp[string]interface{}的别名,方便传数据给模板
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"userID": userID,
	})
}
