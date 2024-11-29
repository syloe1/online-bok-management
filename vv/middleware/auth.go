// backend/middleware/auth.go
package middleware

import (
	"log"
	"net/http"
	"strings"

	"online-book-management/backend/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Cookie 中获取 Authorization
		authCookie, err := c.Cookie("Authorization")
		if err != nil || authCookie == "" {
			log.Println("Authorization cookie missing")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		log.Printf("Authorization cookie: %s", authCookie)

		parts := strings.SplitN(authCookie, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Authorization cookie format invalid")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法是否是预期的
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return config.JwtSecret, nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !token.Valid {
			log.Println("Invalid token")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		log.Println("Token is valid")
		c.Next()
	}
}
