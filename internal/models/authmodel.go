package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"time"
)

type SignInRequest struct {
	Email    string `json:"eMailAddress"`
	Password string `json:"password"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint
	Email     string
	Token     string
	IssuedAt  time.Time
	ExpiredAt time.Time
}

// TokenBlacklist represents the database model for revoked tokens
type TokenBlacklist struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"unique"`
	Reason    string
	ExpiresAt time.Time
	CreatedAt time.Time
}
