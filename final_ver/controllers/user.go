// backend/controllers/user.go
package controllers

import (
	"net/http"
	"time"

	"online-book-management/backend/config"
	"online-book-management/backend/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ShowLoginPage 显示登录页面
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Register 注册新用户
func Register(c *gin.Context) {
	var input struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "用户名和密码为必填项"})
		return
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "用户名已存在"})
		return
	}

	// 注册成功后，提示用户登录
	c.HTML(http.StatusOK, "register.html", gin.H{"success": "注册成功，请登录"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var input struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "用户名和密码为必填项"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "无效的用户名或密码"})
		return
	}

	// 比较密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "无效的用户名或密码"})
		return
	}

	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   user.Username,
	})

	tokenString, err := token.SignedString(config.JwtSecret) // 使用公共的JwtSecret
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "无法生成令牌"})
		return
	}

	// 设置令牌到Cookie
	c.SetCookie("Authorization", "Bearer "+tokenString, 3600*72, "/", "", false, true) // 修正域名为空字符串

	// 登录成功后，重定向到图书列表页面
	c.Redirect(http.StatusSeeOther, "/books_page")
}
