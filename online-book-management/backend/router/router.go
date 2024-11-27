// backend/router/router.go
package router

import (
	"online-book-management/backend/controllers"
	"online-book-management/backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 设置CORS配置，允许当前域
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:10086"} // 调整为前端实际运行地址
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true // 允许发送 Cookie
	r.Use(cors.New(corsConfig))

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 公共路由
	r.GET("/", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)

	// 受保护的路由
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/books", controllers.GetBooks)
		authorized.GET("/books/:id", controllers.GetBook)
		authorized.POST("/books", controllers.CreateBook)
		authorized.PUT("/books/:id", controllers.UpdateBook)
		authorized.DELETE("/books/:id", controllers.DeleteBook)

		// 前端页面路由
		authorized.GET("/books_page", controllers.ShowBooksPage)
		authorized.GET("/add_book", controllers.ShowAddBookPage)
	}

	return r
}
