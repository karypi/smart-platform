package models

import "time"

// GORM 的模型,对应数据库表 alerts, 表名默认是结构体名的复数形式
type Alert struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint // 外键, 指向 users.id, 用于记录哪位用户发出的告警
	Token     string
	AlertName string
	Severity  string
	Instance  string
	Message   string
	Status    string
	CreatedAt time.Time
}
