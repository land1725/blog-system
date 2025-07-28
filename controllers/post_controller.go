package controllers

import (
	"blog-system/database"
	"blog-system/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// CreatePost 创建新文章
// 参数: c - Gin上下文
func CreatePost(c *gin.Context) {
	var post models.Post
	// 绑定JSON数据
	var data map[string]interface{} // 定义一个 map 接收 JSON 数据

	// 绑定 JSON 数据到 map
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 标题类型断言和错误处理
	if title, ok := data["title"].(string); ok {
		post.Title = title
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title must be a string"})
		return
	}

	// 内容类型断言和错误处理
	if content, ok := data["content"].(string); ok {
		post.Content = content
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content must be a string"})
		return
	}

	// 从上下文中获取用户ID
	if userIdValue, exists := c.Get("userid"); exists {
		// 类型安全转换
		if userId, ok := userIdValue.(uint); ok {
			post.UserID = userId
		} else {
			// 记录详细的类型错误日志
			log.Printf("invalid userid type: expected uint, got %T", userIdValue)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user data"})
			return
		}
	}

	// 创建文章 - 添加错误处理
	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPosts 获取所有文章列表
// 参数: c - Gin上下文
func GetPosts(c *gin.Context) {
	// 查询所有文章
	var posts []models.Post

	// 添加数据库查询错误处理
	if err := database.DB.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost 获取单篇文章详情
// 参数: c - Gin上下文
func GetPost(c *gin.Context) {
	var post models.Post
	// 从URL参数获取文章ID
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 查询文章 - 添加错误处理
	if err := database.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost 更新文章
// 参数: c - Gin上下文
func UpdatePost(c *gin.Context) {
	// 获取文章ID
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var data map[string]interface{} // 定义一个 map 接收 JSON 数据

	// 绑定 JSON 数据到 map
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 准备更新数据
	updateData := make(map[string]interface{})

	// 标题类型断言和错误处理
	if title, ok := data["title"].(string); ok {
		updateData["title"] = title
	} else if data["title"] != nil { // 如果提供了title但不是字符串
		c.JSON(http.StatusBadRequest, gin.H{"error": "title must be a string"})
		return
	}

	// 内容类型断言和错误处理
	if content, ok := data["content"].(string); ok {
		updateData["content"] = content
	} else if data["content"] != nil { // 如果提供了content但不是字符串
		c.JSON(http.StatusBadRequest, gin.H{"error": "content must be a string"})
		return
	}

	// 如果没有提供任何有效更新字段
	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有提供有效的更新字段"})
		return
	}

	// 更新文章 - 添加错误处理
	result := database.DB.Model(&models.Post{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	// 检查是否成功更新了记录
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在或没有变化"})
		return
	}

	var updatedPost models.Post
	if err := database.DB.First(&updatedPost, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取更新后的文章失败"})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

// DeletePost 删除文章
// 参数: c - Gin上下文
func DeletePost(c *gin.Context) {
	// 获取文章ID
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 删除文章 - 添加错误处理
	result := database.DB.Where("id = ?", id).Delete(&models.Post{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	// 检查是否成功删除了记录
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
