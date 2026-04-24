package domains

import (
	"time"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID                string    `gorm:"primaryKey;type:uuid" json:"id"`
	Email             string    `gorm:"unique;not null" json:"email"`
	Role              UserRole  `gorm:"default:'user'" json:"role"`
	IsActive          bool      `gorm:"default:false" json:"is_active"`
	FullName          string    `json:"full_name"`
	PhoneNumber       string    `json:"phone_number"`
	AvatarURL         string    `json:"avatar_url"`
	IsVerified        bool      `gorm:"default:false" json:"is_verified"`
	Bio               string    `json:"bio"`
	TradingExperience string    `json:"trading_experience"` // 'Beginner', 'Intermediate', 'Pro'
	AuthProvider      string    `gorm:"column:auth_provider;default:'email'" json:"auth_provider"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
