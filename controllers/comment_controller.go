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

// CreateComment 创建评论
// 参数: c - Gin上下文
func CreateComment(c *gin.Context) {
	var comment models.Comment
	// 绑定JSON数据
	// 绑定JSON数据
	var data map[string]interface{} // 定义一个 map 接收 JSON 数据

	// 绑定 JSON 数据到 map
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if postID, ok := data["post_id"].(float64); ok {
		comment.PostID = uint(postID) // 将 float64 转换为 uint
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id must be a number"})
		return
	}
	comment.Content = data["content"].(string)
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
	}
	// 创建评论
	database.DB.Create(&comment)
	c.JSON(http.StatusCreated, comment)
}

// GetCommentsByPost 获取文章的所有评论
// 参数: c - Gin上下文
func GetCommentsByPost(c *gin.Context) {
	var post models.Post
	// 从URL参数获取文章ID
	param := c.Param("post_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		fmt.Println(id)
	}
	// 查询该文章的所有评论
	database.DB.Preload("Comments").Where("id =?", id).Find(&post)

	c.JSON(http.StatusOK, post.Comments)
}
