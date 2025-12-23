package models

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 后续其他包用 models.DB 就能访问数据库连接, (类似 Django 的 form django.db import connection 的全局可访问对象)
var DB *gorm.DB

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

// Init DB
func InitDB() {
	username := "smartuser"
	password := "smartpass"
	host := "127.0.0.1"
	port := "3306"
	dbName := "smart_platform"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("MySQL DB Connect Failed!:", err)
	}

	// 自动迁移表结构, AutoMigrate 会自动根据模型结构在数据库中创建/修改表结构
	DB.AutoMigrate(&User{}, &Alert{})

	fmt.Println("MySQL DB Connect Successful!")
}
