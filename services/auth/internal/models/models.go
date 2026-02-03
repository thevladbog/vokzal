package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User модель пользователя
type User struct {
	ID           string    `gorm:"type:uuid;primary_key" json:"id"`
	Username     string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string    `gorm:"type:varchar(100);not null" json:"full_name"`
	Role         string    `gorm:"type:varchar(20);not null" json:"role"`
	StationID    *string   `gorm:"type:uuid" json:"station_id,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Session модель сессии
type Session struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string    `gorm:"type:varchar(255);unique;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (Session) TableName() string {
	return "sessions"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
