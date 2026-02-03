// Package models содержит модели данных сервиса платежей.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment — модель платежа (карта, СБП, наличные).
type Payment struct {
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	RefundedAt   *time.Time `json:"refunded_at,omitempty"`
	ErrorMsg     *string    `gorm:"type:text" json:"error_msg,omitempty"`
	TicketID     *string    `gorm:"type:uuid;index" json:"ticket_id,omitempty"`
	RefundAmount *float64   `gorm:"type:decimal(10,2)" json:"refund_amount,omitempty"`
	ConfirmedAt  *time.Time `json:"confirmed_at,omitempty"`
	ExternalID   *string    `gorm:"type:varchar(100);index" json:"external_id,omitempty"`
	PaymentURL   *string    `gorm:"type:varchar(500)" json:"payment_url,omitempty"`
	QRCode       *string    `gorm:"type:text" json:"qr_code,omitempty"`
	Status       string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Currency     string     `gorm:"type:varchar(3);default:'RUB'" json:"currency"`
	ID           string     `gorm:"type:uuid;primary_key" json:"id"`
	Provider     string     `gorm:"type:varchar(20);not null" json:"provider"`
	Metadata     string     `gorm:"type:jsonb" json:"metadata,omitempty"`
	Method       string     `gorm:"type:varchar(20);not null" json:"method"`
	Amount       float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
}

// TableName возвращает имя таблицы для GORM.
func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate генерирует UUID для новой записи.
func (p *Payment) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
