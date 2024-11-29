// backend/middleware/user_context.go
package middleware

import (
	"strings"

	"online-book-management/backend/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("User", nil) // 默认无用户
		authCookie, err := c.Cookie("Authorization")
		if err != nil || authCookie == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authCookie, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return config.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username, ok := claims["sub"].(string)
			if ok {
				// 设置用户信息到上下文
				c.Set("User", username)
			}
		}

		c.Next()
	}
}
