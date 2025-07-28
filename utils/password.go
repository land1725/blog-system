package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 密码加密
// 参数: password - 原始密码
// 返回值: 加密后的密码, 错误信息
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// CheckPassword 密码验证
// 参数: hashedPassword - 加密密码, password - 原始密码
// 返回值: 是否匹配
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
