package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string    `gorm:"not null" form:"title" json:"title" binding:"required"`     // 文章标题
	Content  string    `gorm:"not null" form:"content" json:"content" binding:"required"` // 文章内容
	UserID   uint      // 作者ID
	User     User      // 关联作者
	Comments []Comment // 文章关联的评论
}
