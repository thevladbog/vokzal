// Package models содержит модели данных сервиса уведомлений.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification — модель уведомления (SMS, email, Telegram, TTS).
type Notification struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Subject   *string    `gorm:"type:varchar(200)" json:"subject,omitempty"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	ErrorMsg  *string    `gorm:"type:text" json:"error_msg,omitempty"`
	ID        string     `gorm:"type:uuid;primary_key" json:"id"`
	Type      string     `gorm:"type:varchar(20);not null;index" json:"type"`
	Recipient string     `gorm:"type:varchar(100);not null" json:"recipient"`
	Message   string     `gorm:"type:text;not null" json:"message"`
	Status    string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Metadata  string     `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// TableName возвращает имя таблицы для GORM.
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate генерирует UUID для новой записи.
func (n *Notification) BeforeCreate(_ *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
