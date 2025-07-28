package database

import (
	"blog-system/models"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
// 参数: db - 数据库实例
func AutoMigrate(db *gorm.DB) error {
	// 执行自动迁移
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)

	return err
}
