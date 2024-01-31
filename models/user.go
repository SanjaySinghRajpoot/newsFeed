package models

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-"`
	Posts    []Post `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
}

type Post struct {
	gorm.Model
	ID      uint   `json:"id" gorm:"primary_key"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `gorm:"foreignKey:UserID"`
}

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type UserJWT struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
