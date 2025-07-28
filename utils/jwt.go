package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// CustomClaims 自定义声明结构（包含标准声明和用户信息）
type CustomClaims struct {
	jwt.RegisteredClaims
	Username string // 用户名（展示用）
	UserID   uint   // 用户ID（核心身份标识）
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username, secret string) (string, error) {
	claims := CustomClaims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// 关键：设置24小时有效期（安全性与用户体验平衡点）
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// 关键：使用HS256算法创建令牌（防篡改签名）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 关键：使用密钥签名（密钥泄露=系统安全崩溃）
	return token.SignedString([]byte(secret))
}

// ParseToken 解析验证JWT令牌
func ParseToken(tokenString, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	// 关键修复：先处理错误情况
	if err != nil {
		return nil, fmt.Errorf("令牌解析失败: %w", err)
	}

	// 确保token非空后再访问其属性
	if token == nil {
		return nil, fmt.Errorf("无效的令牌")
	}

	// 安全检查：类型断言+有效性验证
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的令牌声明")
}
