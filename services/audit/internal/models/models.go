// Package models — доменные модели Audit Service.
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONB — тип для хранения JSON в PostgreSQL.
type JSONB []byte

// Value реализует driver.Valuer для записи JSONB в БД.
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan реализует sql.Scanner для чтения JSONB из БД.
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

// AuditLog — модель лога аудита.
type AuditLog struct {
	ID         string    `gorm:"type:uuid;primary_key" json:"id"`
	EntityType string    `gorm:"type:varchar(50);not null;index" json:"entity_type"`
	EntityID   string    `gorm:"type:uuid;not null;index" json:"entity_id"`
	Action     string    `gorm:"type:varchar(50);not null" json:"action"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
	UserID     *string   `gorm:"type:uuid;index" json:"user_id,omitempty"`
	IPAddress  *string   `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  *string   `gorm:"type:varchar(500)" json:"user_agent,omitempty"`
	OldValue   JSONB     `gorm:"type:jsonb" json:"old_value,omitempty"`
	NewValue   JSONB     `gorm:"type:jsonb" json:"new_value,omitempty"`
}

// TableName возвращает имя таблицы для GORM.
func (AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate заполняет ID перед созданием записи.
func (a *AuditLog) BeforeCreate(_ *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// SetOldValue устанавливает old_value из интерфейса.
func (a *AuditLog) SetOldValue(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.OldValue = JSONB(data)
	return nil
}

// SetNewValue устанавливает new_value из интерфейса.
func (a *AuditLog) SetNewValue(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.NewValue = JSONB(data)
	return nil
}
