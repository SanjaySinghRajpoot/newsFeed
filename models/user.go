package models

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// User Model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-"`
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

type Post struct {
}
