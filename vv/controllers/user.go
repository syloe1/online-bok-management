
package controllers

import (
	"log"
	"net/http"
	"time"

	"online-book-management/backend/config"
	"online-book-management/backend/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ShowAuthPage 显示统一认证页面（登录和注册）
func ShowAuthPage(c *gin.Context) {
	user, exists := c.Get("User")
	if !exists || user == nil {
		c.HTML(http.StatusOK, "auth.html", nil)
		return
	}
	// 如果用户已登录，重定向到图书列表页面
	c.Redirect(http.StatusSeeOther, "/books_page")
}

// Register 注册新用户
func Register(c *gin.Context) {
	var input struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		log.Printf("注册失败: 绑定表单数据错误: %v", err)
		c.HTML(http.StatusBadRequest, "auth.html", gin.H{"register_error": "用户名和密码为必填项"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		log.Printf("注册失败: 用户名 %s 已存在", input.Username)
		c.HTML(http.StatusBadRequest, "auth.html", gin.H{"register_error": "用户名已存在"})
		return
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("注册失败: 密码哈希错误: %v", err)
		c.HTML(http.StatusInternalServerError, "auth.html", gin.H{"register_error": "密码加密失败"})
		return
	}

	// 创建用户
	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		log.Printf("注册失败: 数据库插入错误: %v", result.Error)
		c.HTML(http.StatusInternalServerError, "auth.html", gin.H{"register_error": "无法创建用户"})
		return
	}

	log.Printf("用户注册成功: %s", user.Username)
	// 注册成功后，提示用户登录
	c.HTML(http.StatusOK, "auth.html", gin.H{"register_success": "注册成功，请登录"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var input struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		log.Printf("登录失败: 绑定表单数据错误: %v", err)
		c.HTML(http.StatusBadRequest, "auth.html", gin.H{"login_error": "用户名和密码为必填项"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		log.Printf("登录失败: 用户名 %s 不存在", input.Username)
		c.HTML(http.StatusUnauthorized, "auth.html", gin.H{"login_error": "无效的用户名或密码"})
		return
	}

	// 比较密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("登录失败: 密码不匹配 for 用户 %s", input.Username)
		c.HTML(http.StatusUnauthorized, "auth.html", gin.H{"login_error": "无效的用户名或密码"})
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
		log.Printf("登录失败: 生成JWT错误: %v", err)
		c.HTML(http.StatusInternalServerError, "auth.html", gin.H{"login_error": "无法生成令牌"})
		return
	}

	// 设置令牌到Cookie
	c.SetCookie("Authorization", "Bearer "+tokenString, 3600*72, "/", "", false, true) // 修正域名为空字符串

	// 登录成功后，重定向到图书列表页面
	c.Redirect(http.StatusSeeOther, "/books_page")
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// 清除认证 Cookie
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	// 重定向到认证页面
	c.Redirect(http.StatusSeeOther, "/")
}