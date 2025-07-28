package controllers

import (
	"blog-system/database"
	"blog-system/models"
	"fmt"
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
	post.Title = data["title"].(string)
	post.Content = data["content"].(string)
	//err := c.ShouldBindJSON(&post)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// 从上下文中获取用户ID
	// 设置文章作者
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
	// 创建文章
	err := database.DB.Create(&post).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, post)
}

// GetPosts 获取所有文章列表
// 参数: c - Gin上下文
func GetPosts(c *gin.Context) {
	// 查询所有文章（带分页）
	var posts []models.Post
	database.DB.Find(&posts)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		fmt.Println(id)
	}
	database.DB.Where("id = ?", id).Find(&post)
	// 查询文章
	c.JSON(http.StatusOK, post)
}

// UpdatePost 更新文章
// 参数: c - Gin上下文
func UpdatePost(c *gin.Context) {
	// 获取文章ID
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		fmt.Println(id)
	}
	//var updatedPost models.Post
	var data map[string]interface{} // 定义一个 map 接收 JSON 数据

	// 绑定 JSON 数据到 map
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 更新文章
	database.DB.Model(&models.Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":   data["title"].(string),
		"content": data["content"].(string),
	})
	var updatedPost models.Post
	database.DB.First(&updatedPost, id)
	c.JSON(http.StatusOK, updatedPost)
}

// DeletePost 删除文章
// 参数: c - Gin上下文
func DeletePost(c *gin.Context) {
	// 获取文章ID
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		fmt.Println(id)
	}
	// 删除文章
	database.DB.Where("id = ?", id).Delete(&models.Post{})
	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
