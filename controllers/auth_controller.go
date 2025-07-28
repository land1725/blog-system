package controllers

import (
	"blog-system/config"
	"blog-system/database"
	"blog-system/models"
	"blog-system/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 用户注册
// 参数: c - Gin上下文
func Register(c *gin.Context) {
	var user models.User
	// 绑定JSON数据到user结构体
	err := c.ShouldBind(&user)
	if err != nil {
		return
	}
	// 密码加密
	password, err := utils.HashPassword(user.Password)
	if err != nil {
		return
	} else {
		user.Password = password
	}
	// 创建用户
	database.DB.Create(&user)
	// 返回响应
	c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})
}

// Login 用户登录
// 参数: c - Gin上下文
// 返回值: 包含JWT的响应
func Login(c *gin.Context) {
	var credentials struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}
	// 绑定JSON数据
	err := c.ShouldBind(&credentials)
	if err != nil {
		return
	}

	var user models.User
	err = database.DB.Where("username = ? ", credentials.Username).First(&user).Error
	if err != nil {
		return
	}
	// 验证用户凭据
	res := utils.CheckPassword(user.Password, credentials.Password)
	if res == false {
		return
	}
	cfg := config.LoadConfig()
	jwt := cfg.JWTSecret
	// 生成JWT
	token, err := utils.GenerateToken(user.ID, user.Username, jwt)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
