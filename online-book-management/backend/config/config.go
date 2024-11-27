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
	// dsn := "bookuser:password123@tcp(127.0.0.1:3306)/book_management?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:qaz123@tcp(127.0.0.1:3306)/book_management?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		fmt.Println("Connected to database successfully")
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
