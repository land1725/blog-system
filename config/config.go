package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver          string
	DBDSN             string
	JWTSecret         string
	ServerPort        string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

// LoadConfig 加载配置
func LoadConfig() Config {
	// 尝试加载 .env 文件
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("注意: 未找到 .env 文件，将使用系统环境变量")
	}

	// 获取数据库连接池配置，如果未设置则使用默认值
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetime := time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 30)) * time.Minute

	return Config{
		DBDriver:          getEnv("DB_DRIVER", "mysql"),
		DBDSN:             getEnv("DB_DSN", ""),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DBMaxOpenConns:    maxOpenConns,
		DBMaxIdleConns:    maxIdleConns,
		DBConnMaxLifetime: connMaxLifetime,
	}
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 获取整型环境变量
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("警告: 环境变量 %s 的值 '%s' 不是有效整数，使用默认值 %d", key, value, defaultValue)
		return defaultValue
	}
	return intValue
}
