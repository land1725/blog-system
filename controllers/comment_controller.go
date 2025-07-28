package controllers

import (
	"blog-system/database"
	"blog-system/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// CreateComment 创建评论
// 参数: c - Gin上下文
func CreateComment(c *gin.Context) {
	var comment models.Comment
	// 绑定JSON数据
	var data map[string]interface{} // 定义一个 map 接收 JSON 数据

	// 绑定 JSON 数据到 map
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理post_id - 添加更严格的类型检查
	if postID, ok := data["post_id"]; ok {
		switch v := postID.(type) {
		case float64:
			comment.PostID = uint(v)
		case int:
			comment.PostID = uint(v)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "post_id must be a number"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}

	// 处理content - 添加类型检查和错误处理
	if content, ok := data["content"].(string); ok {
		comment.Content = content
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content must be a string"})
		return
	}

	// 从上下文中获取用户ID
	if userIdValue, exists := c.Get("userid"); exists {
		// 类型安全转换
		if userId, ok := userIdValue.(uint); ok {
			comment.UserID = userId
		} else {
			// 记录详细的类型错误日志
			log.Printf("invalid userid type: expected uint, got %T", userIdValue)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user data"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// 创建评论 - 添加错误处理
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetCommentsByPost 获取文章的所有评论
// 参数: c - Gin上下文
func GetCommentsByPost(c *gin.Context) {
	// 从URL参数获取文章ID
	param := c.Param("post_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var post models.Post
	// 查询该文章的所有评论 - 添加错误处理
	if err := database.DB.Preload("Comments").Where("id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, post.Comments)
}
