// backend/main.go
package main

import (
	"online-book-management/backend/config"
	"online-book-management/backend/models"
	"online-book-management/backend/router"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 初始化 JWT 密钥
	config.InitJwtSecret()

	// 自动迁移
	config.DB.AutoMigrate(&models.User{}, &models.Book{})

	// 初始化路由
	r := router.SetupRouter()

	// 启动服务器
	if err := r.Run(":10086"); err != nil {
		panic("Failed to run server: " + err.Error())
	}
}
