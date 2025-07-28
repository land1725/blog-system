package routes

import (
	"blog-system/controllers"
	"blog-system/middleware"
	"github.com/gin-gonic/gin"
)

// InitRoutes 初始化应用路由
// 参数: router - Gin引擎实例
func InitRoutes(router *gin.Engine) {
	// 认证相关路由
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// 文章相关路由
	posts := router.Group("/posts")
	{
		posts.GET("", controllers.GetPosts)                                           // 获取文章列表
		posts.GET("/:id", controllers.GetPost)                                        // 获取单篇文章
		posts.Use(middleware.AuthMiddleware())                                        // 以下路由需要认证
		posts.POST("", controllers.CreatePost)                                        // 创建文章
		posts.PUT("/:id", middleware.AuthorizePostOwner(), controllers.UpdatePost)    // 更新文章
		posts.DELETE("/:id", middleware.AuthorizePostOwner(), controllers.DeletePost) // 删除文章
	}

	// 评论相关路由
	comments := router.Group("/comments")
	{
		comments.GET("/post/:post_id", controllers.GetCommentsByPost) // 获取文章评论
		comments.Use(middleware.AuthMiddleware())                     // 以下路由需要认证
		comments.POST("", controllers.CreateComment)                  // 创建评论
	}
}
