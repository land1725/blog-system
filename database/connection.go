package database

import (
	"blog-system/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.LoadConfig()

	var dialector gorm.Dialector

	switch cfg.DBDriver {
	case "mysql":
		if cfg.DBDSN == "" {
			return fmt.Errorf("MySQL DSN 未配置")
		}
		dialector = mysql.Open(cfg.DBDSN)
	case "sqlite":
		dsn := cfg.DBDSN
		if dsn == "" {
			dsn = "blog.db"
		}
		dialector = sqlite.Open(dsn)
	default:
		return fmt.Errorf("不支持的数据库驱动: %s", cfg.DBDriver)
	}

	// 创建数据库连接
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接池失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)       // 最大打开连接数
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)       // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime) // 连接最大生存时间

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	log.Printf("成功连接到 %s 数据库", cfg.DBDriver)
	DB = db
	return nil
}
