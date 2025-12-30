package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"smart-platform/backend/models"
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

type AlertRequest struct {
	Token     string `json:"token"`
	AlertName string `json:"alertname"`
	Severity  string `json:"severity"`
	Instance  string `json:"instance"`
	Message   string `json:"message"`
	Status    string `json:"status"`
}

func buildWechatMessage(req AlertRequest) map[string]interface{} {
	title := "【Problem】"
	if req.Status == "resolved" {
		title = "【Resolve】"
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	content := fmt.Sprintf(
		"%s\n"+
			"Alert Severity: %s\n"+
			"Alert Type: %s\n"+
			"Alert Instance: %s\n"+
			"Alert Description: %s\n"+
			"Alert Status: %s\n"+
			"Alert Time: %s",
		title,
		req.Severity,
		req.AlertName,
		req.Instance,
		req.Message,
		req.Status,
		now,
	)

	return map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}
}

func SendAlert(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(401, gin.H{"error": "not logged in"})
		return
	}

	var req AlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Status == "" {
		req.Status = "firing"
	}

	// 将 sendData 序列化为 JSON 字节
	body, _ := json.Marshal(buildWechatMessage(req))

	_, err := http.Post(
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+req.Token,
		"application/json",
		// 将 body 包装成 io.Reader 数据流
		bytes.NewBuffer(body),
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 保存记录到数据库
	alertRecord := models.Alert{
		UserID:    userID.(uint),
		Token:     req.Token,
		AlertName: req.AlertName,
		Severity:  req.Severity,
		Instance:  req.Instance,
		Message:   req.Message,
		Status:    req.Status,
		CreatedAt: time.Now(),
	}
	// 通过 GORM 插入数据库, GORM 要求传入指针，直接修改对象里的字段，如果传入值类型，GORM修改的是副本，外部 alert.ID 不会更新
	models.DB.Create(&alertRecord)

	c.JSON(200, gin.H{
		"success": true,
	})
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
	c.JSON(200, gin.H{
		"success": true,
		"data":    alerts,
	})
}
