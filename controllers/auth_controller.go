package controllers

import (
	"blog-system/config"
	"blog-system/database"
	"blog-system/models"
	"blog-system/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Register 用户注册
// 参数: c - Gin上下文
func Register(c *gin.Context) {
	var user models.User
	// 绑定JSON数据到user结构体
	err := c.ShouldBind(&user)
	if err != nil {
		log.Fatalf("绑定用户错误: %v", err)
		return
	}
	// 密码加密
	password, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	} else {
		user.Password = password
	}
	// 创建用户
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败: " + err.Error()})
		return
	}
	// 返回响应
	c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})
}

// Login 用户登录
// 参数: c - Gin上下文
// 返回值: 包含JWT的响应
func Login(c *gin.Context) {
	// 使用 map 接收 JSON 数据
	var data map[string]interface{}

	// 绑定 JSON 数据到 map
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 从 map 中提取用户名和密码
	username, usernameOk := data["username"].(string)
	password, passwordOk := data["password"].(string)

	// 验证字段存在性和类型
	if !usernameOk || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}
	if !passwordOk || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}

	var user models.User
	// 数据库查询错误处理
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证用户凭据
	if !utils.CheckPassword(user.Password, password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	cfg := config.LoadConfig()
	jwtSecret := cfg.JWTSecret

	// 生成JWT错误处理
	token, err := utils.GenerateToken(user.ID, user.Username, jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
