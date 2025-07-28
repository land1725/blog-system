package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"` // 评论内容
	UserID  uint   // 评论者ID
	User    User   // 关联评论者
	PostID  uint   // 关联文章ID
	Post    Post   // 关联文章
}
