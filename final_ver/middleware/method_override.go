package middleware

import "github.com/gin-gonic/gin"

// MethodOverride middleware allows method overriding via a form parameter _method.
func MethodOverride() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果请求方法是 POST 且表单中有 "_method" 参数，则重写请求方法
		if method := c.DefaultPostForm("_method", ""); method != "" {
			c.Request.Method = method
		}
		c.Next()
	}
}
