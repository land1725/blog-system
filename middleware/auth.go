package middleware

import (
	"blog-system/config"
	"blog-system/database"
	"blog-system/models"
	"blog-system/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			c.JSON(http.StatusUnauthorized, gin.H{ // 改为401状态码
				"code": 2003,
				"msg":  "Authorization header is required",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusBadRequest, gin.H{ // 改为400状态码
				"code": 2004,
				"msg":  "Invalid Authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{ // 改为400状态码
				"code": 2004,
				"msg":  "Token string is empty",
			})
			c.Abort()
			return
		}

		cfg := config.LoadConfig()
		if cfg.JWTSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{ // 添加500状态码
				"code": 5000,
				"msg":  "Server configuration error",
			})
			c.Abort()
			return
		}

		// 解析JWT令牌
		mc, err := utils.ParseToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{ // 改为401状态码
				"code": 2005,
				"msg":  "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 验证用户ID有效性
		if mc.UserID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{ // 改为401状态码
				"code": 2006,
				"msg":  "Invalid user in token",
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
		id, err := strconv.ParseUint(param, 10, 32) // 更安全的转换方式
		if err != nil || id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			c.Abort()
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
		if !ok || userId == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// 4. 查询数据库验证文章所有者
		var post models.Post
		if err := database.DB.First(&post, id).Error; err != nil {
			// 更精确的错误处理
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
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
