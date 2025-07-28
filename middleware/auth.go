package middleware

import (
	"blog-system/config"
	"blog-system/database"
	"blog-system/models"
	"blog-system/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// AuthMiddleware JWT认证中间件
// 返回值: Gin处理函数
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		cfg := config.LoadConfig()
		jwt := cfg.JWTSecret
		mc, err := utils.ParseToken(parts[1], jwt)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Set("userid", mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

// AuthorizePostOwner 验证文章所有者中间件
// 返回值: Gin处理函数
func AuthorizePostOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从URL参数获取文章ID
		param := c.Param("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			c.Abort() // 终止后续处理
			return
		}

		// 2. 从上下文中获取用户ID
		userIdValue, exists := c.Get("userid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// 3. 类型断言确保userid是uint类型
		userId, ok := userIdValue.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			c.Abort()
			return
		}

		// 4. 查询数据库验证文章所有者
		var post models.Post
		if err := database.DB.First(&post, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			c.Abort()
			return
		}

		// 5. 验证当前用户是否是文章作者
		if post.UserID != userId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this post"})
			c.Abort()
			return
		}

		// 继续执行下一个处理程序
		c.Next()
	}
}
