package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification модель уведомления
type Notification struct {
	ID        string     `gorm:"type:uuid;primary_key" json:"id"`
	Type      string     `gorm:"type:varchar(20);not null;index" json:"type"` // sms, email, telegram, tts
	Recipient string     `gorm:"type:varchar(100);not null" json:"recipient"`
	Message   string     `gorm:"type:text;not null" json:"message"`
	Subject   *string    `gorm:"type:varchar(200)" json:"subject,omitempty"`
	Status    string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, sent, failed
	SentAt    *time.Time `json:"sent_at,omitempty"`
	ErrorMsg  *string    `gorm:"type:text" json:"error_msg,omitempty"`
	Metadata  string     `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
