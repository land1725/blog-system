package main

import (
	"blog-system/config"
	"blog-system/database"
	"blog-system/middleware"
	"blog-system/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 自动迁移数据库表结构
	if err := database.AutoMigrate(database.DB); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建Gin引擎实例
	router := gin.Default()

	// 设置中间件
	router.Use(middleware.Logger())

	// 初始化路由
	routes.InitRoutes(router)

	// 启动服务器
	cfg := config.LoadConfig()
	addr := ":" + cfg.ServerPort
	log.Printf("服务器启动，监听地址: http://localhost%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
