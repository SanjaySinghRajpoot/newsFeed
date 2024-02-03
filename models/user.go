package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-"`
	Posts    []Post
}

// type Post struct {
// 	gorm.Model
// 	ID      uint   `json:"id" gorm:"primary_key"`
// 	Content string `json:"content"`
// 	UserID  uint   `gorm:"foreignKey:UserID" json:"user_id"`
// 	User    User   `gorm:"foreignKey:UserID"`
// 	Comment []Comment
// }

type Post struct {
	gorm.Model
	Content  string `json:"content"`
	UserID   uint   `gorm:"foreignKey:UserID" json:"userID"`
	User     User   `gorm:"foreignKey:UserID"`
	Comments []Comment
}

type Comment struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey"`
	PostID uint   `gorm:"foreignKey:PostID" json:"post_id" binding:"required,gt=0"`
	UserID uint   `gorm:"foreignKey:UserID"`
	Body   string `gorm:"type:text" json:"body"`
	User   User
}

type Follower struct {
	gorm.Model
	FollowerUserID  uint `json:"follower_user_id"`
	FollowingUserID uint `json:"following_user_id"`
}

type NotificationPayload struct {
	// Type -> email, sms, inapp
	UserID      []uint `json:"user_id"`
	Type        string `json:"type"  validate:"required,oneof=sms email inapp"`
	Description string `json:"description"`
}
