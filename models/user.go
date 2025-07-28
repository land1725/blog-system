package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string    `gorm:"unique;not null" form:"username" json:"username" binding:"required"` // 添加 form 和 json 标签
	Password string    `gorm:"not null" form:"password" json:"password" binding:"required"`        // 添加 form 和 json 标签
	Email    string    `gorm:"unique;not null" form:"email" json:"email" binding:"required,email"` // 添加 form 和 json 标签
	Posts    []Post    `form:"-" json:"posts,omitempty"`                                           // 禁用 form 绑定
	Comments []Comment `form:"-" json:"comments,omitempty"`                                        // 禁用 form 绑定
}
