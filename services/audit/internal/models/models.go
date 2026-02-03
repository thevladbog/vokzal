package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JSONB []byte

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONB value")
	}
	*j = s
	return nil
}

// AuditLog модель лога аудита
type AuditLog struct {
	ID         string    `gorm:"type:uuid;primary_key" json:"id"`
	EntityType string    `gorm:"type:varchar(50);not null;index" json:"entity_type"`
	EntityID   string    `gorm:"type:uuid;not null;index" json:"entity_id"`
	Action     string    `gorm:"type:varchar(50);not null" json:"action"`
	UserID     *string   `gorm:"type:uuid;index" json:"user_id,omitempty"`
	OldValue   JSONB     `gorm:"type:jsonb" json:"old_value,omitempty"`
	NewValue   JSONB     `gorm:"type:jsonb" json:"new_value,omitempty"`
	IPAddress  *string   `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  *string   `gorm:"type:varchar(500)" json:"user_agent,omitempty"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// SetOldValue устанавливает old_value из интерфейса
func (a *AuditLog) SetOldValue(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.OldValue = JSONB(data)
	return nil
}

// SetNewValue устанавливает new_value из интерфейса
func (a *AuditLog) SetNewValue(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.NewValue = JSONB(data)
	return nil
}
