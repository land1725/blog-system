package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
)

// 自定义Claims结构体
type CustomClaims struct {
	jwt.RegisteredClaims
	Username string
	UserID   uint
}

// GenerateToken 生成JWT令牌
// 参数: userID - 用户ID, username - 用户名, secret - JWT密钥
// 返回值: 令牌字符串, 错误信息
func GenerateToken(userID uint, username, secret string) (string, error) {

	// 创建带有自定义Claims的token
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		Username:         username,
		UserID:           userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 用于签名的秘钥

	// 生成字符串形式的token
	strToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("出现错误: %v", err)
	}

	fmt.Println(strToken)
	return strToken, err
}

// ParseToken 解析验证JWT令牌
// 参数: tokenString - 令牌字符串, secret - JWT密钥
// 返回值: 解析后的声明, 错误信息
func ParseToken(tokenString, secret string) (*CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil // 使用传入的secret而不是全局变量
	})

	if err != nil {
		return nil, err
	}

	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
