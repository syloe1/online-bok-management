// backend/config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JwtSecret []byte

func InitDB() {
	// 数据库连接字符串，请根据实际情况修改
	dsn := "root:qaz123@tcp(127.0.0.1:3306)/book_management?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	} else {
		fmt.Println("成功连接到数据库")
	}
}

func InitJwtSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 生产环境中应通过环境变量设置，避免硬编码
		secret = "your_secure_secret_key"
	}
	JwtSecret = []byte(secret)
}
