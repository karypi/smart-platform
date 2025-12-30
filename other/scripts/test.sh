package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// type Alert struct {
// 	Status string `json: "status"`
// 	Labels struct {
// 		Severity  string `json:"severity"`
// 		AlertName string `json:"alertname"`
// 		Instance  string `json:"instance"`
// 	} `json:"labels"`
// 	Annotations struct {
// 		Summary     string `json:"summary"`
// 		Description string `json:"description"`
// 	} `json:"annotations"`
// 	StartsAt time.Time `json:"StartsAt"`
// 	EndsAt   time.Time `json:"endsAT"`
// }

// Alert 结构体
type Alert struct {
	Status      string      `json:"status"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
	StartsAt    time.Time   `json:"startsAt"`
	EndsAt      time.Time   `json:"endsAt"`
}

type Labels struct {
	Severity  string `json:"severity"`
	AlertName string `json:"alertname"`
	Instance  string `json:"instance"`
}

type Annotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type WebhookData struct {
	Alerts []Alert `json:"alerts"`
}

// Global var
var sendHistory []WebhookData
var mu sync.Mutex // 用于并发安全

// Generate Alert Message
func makeAlertData(data WebhookData) map[string]interface{} {
	sendData := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "This message is not in standard Prometheus format",
		},
	}

	if len(data.Alerts) == 0 {
		return sendData
	}

	for _, alert := range data.Alerts {
		severity := alert.Labels.Severity
		alertname := alert.Labels.AlertName
		instance := alert.Labels.Instance

		message := alert.Annotations.Summary
		if message == "" {
			message = alert.Annotations.Description
		}
		if message == "" {
			message = "null"
		}

		startAt := alert.StartsAt.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:05:05")
		endAt := alert.EndsAt.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:05:05")

		title := ""

		if alert.Status == "firing" {
			title = "【Problem】"

			sendData = map[string]interface{}{
				"msgtype": "markdown",
				"markdown": map[string]string{
					"content": fmt.Sprintf(
						"## %s\n\n"+
							">**Alert Severity**: %s\n\n"+
							">**Alert Type**: %s\n\n"+
							">**Alert Instance**: %s\n\n"+
							">**Alert Description**: %s\n\n"+
							">**Alert Status**: %s\n\n"+
							">**Alert Time**: %s\n\n",
						title, severity, alertname, instance, message, alert.Status, startAt,
					),
				},
			}
		} else if alert.Status == "resolved" {
			title = "【Resolve】"

			sendData = map[string]interface{}{
				"msgtype": "markdown",
				"markdown": map[string]string{
					"content": fmt.Sprintf(
						"## %s\n\n"+
							">**Alert Severity**: %s\n\n"
							">**Alert Type**: %s\n\n"+
							">**Alert Instance**: %s\n\n"+
							">**Alert Description**: %s\n\n"+
							">**Alert Status**: %s\n\n"+
							">**Alert StartTime**: %s\n\n"+
							">**Alert EndTime**: %s\n",
						title, severity, alertname, instance, message, alert.Status, startAt, endAt,
					),
				},
			}
		}
	}

	return sendData
}

func sendAlertWithToken(data WebhookData, token string) error {
	if token == "" {
		return fmt.Errorf("you must set ROBOT_TOKEN env")
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", token)
	body, _ := json.Marshal(makeAlertData(data))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Go 的 encoding/json 默认把 数字类型都解析成 float64。
	if result["errcode"].(float64) != 0 {
		return fmt.Errorf("notify wechat error:", result["errcode"])
	}

	// 保存发送任务, 确保多个用户同时发送告警时，历史记录不会被搞乱或崩溃
	mu.Lock()                               // 互斥锁，避免互斥
	defer mu.Unlock()                       // 当前函数执行完之后自动释放锁
	// 把新的 webhook 数据 data 添加到 sendHistory 里
	sendHistory = append(sendHistory, data) // sendHistory 是在内存里存储发送历史记录的 slice 切片, 在 go 里，slice 不是线程安全的   /

	return nil
}

func test() {
	r := gin.Default() // 类似 Python 的 Flask, 包含Logger（记录访问日志)、Recovery（避免程序崩溃）

	// Load Template
	r.LoadHTMLGlob("templates/*")

	// Home page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"history": sendHistory,
		})
	})

	// 表单提交发送告警
	r.POST("/send", func(c *gin.Context) {
		token := c.PostForm("token")

		alertname := c.PostForm("alertname")
		severity := c.PostForm("severity")
		instance := c.PostForm("instance")
		message := c.PostForm("message")

		// 用来接收 webhook 的JSON 数据
		data := WebhookData{
			Alerts: []Alert{
				{
					Status: "firing",
					Labels: Labels{
						Severity:  severity,
						AlertName: alertname,
						Instance:  instance,
					},
					Annotations: Annotations{
						Summary:     message,
						Description: message,
					},
					StartsAt: time.Now(),
					EndsAt:   time.Now().Add(5 * time.Minute),
				},
			},
		}

		if err := sendAlertWithToken(data, token); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Send failed: %v", err))
			return
		}

		// // c.BindJson(data) 把用户发送的 JSON 解析进 data 这个结构体里, 相当于 python 的 request.get_json(), 但 Go 是静态类型，必须把 JSON 填进结构体
		// if err := c.BindJSON(&data); err != nil {
		// 	c.String(400, "Bad request")
		// 	return
		// }
		// sendAlertWithToken(data)
		// c.String(200, "success")

		c.Redirect(http.StatusSeeOther, "/")
	})

	r.Run(":5000")
}
