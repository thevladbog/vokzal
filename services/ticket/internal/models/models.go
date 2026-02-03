// Package models содержит модели данных сервиса билетов.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ticket — модель билета.
type Ticket struct {
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	RefundedAt    *time.Time `json:"refunded_at,omitempty"`
	RefundAmount  *float64   `gorm:"type:decimal(10,2)" json:"refund_amount,omitempty"`
	PassengerDoc  *string    `gorm:"type:varchar(50)" json:"passenger_doc,omitempty"`
	Phone         *string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Email         *string    `gorm:"type:varchar(100)" json:"email,omitempty"`
	SeatID        *string    `gorm:"type:uuid;index" json:"seat_id,omitempty"`
	RefundPenalty *float64   `gorm:"type:decimal(10,2)" json:"refund_penalty,omitempty"`
	PassengerName *string    `gorm:"type:varchar(100)" json:"passenger_name,omitempty"`
	PaymentMethod string     `gorm:"type:varchar(20)" json:"payment_method"`
	BarCode       string     `gorm:"type:varchar(255);unique" json:"bar_code"`
	ID            string     `gorm:"type:uuid;primary_key" json:"id"`
	QRCode        string     `gorm:"type:varchar(255);unique" json:"qr_code"`
	Status        string     `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	TripID        string     `gorm:"type:uuid;not null;index" json:"trip_id"`
	Price         float64    `gorm:"type:decimal(10,2);not null" json:"price"`
}

// BoardingEvent — модель события начала посадки.
type BoardingEvent struct {
	StartedAt time.Time `gorm:"not null" json:"started_at"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	TripID    string    `gorm:"type:uuid;not null;unique;index" json:"trip_id"`
	StartedBy string    `gorm:"type:uuid;not null" json:"started_by"`
}

// BoardingMark — модель отметки посадки.
type BoardingMark struct {
	MarkedAt   time.Time `gorm:"not null" json:"marked_at"`
	CreatedAt  time.Time `json:"created_at"`
	ID         string    `gorm:"type:uuid;primary_key" json:"id"`
	TicketID   string    `gorm:"type:uuid;not null;index" json:"ticket_id"`
	MarkedBy   string    `gorm:"type:uuid;not null" json:"marked_by"`
	ScanMethod string    `gorm:"type:varchar(20)" json:"scan_method"`
}

// TableName возвращает имя таблицы для GORM (Ticket).
func (Ticket) TableName() string {
	return "tickets"
}

// TableName возвращает имя таблицы для GORM (BoardingEvent).
func (BoardingEvent) TableName() string {
	return "boarding_events"
}

// TableName возвращает имя таблицы для GORM (BoardingMark).
func (BoardingMark) TableName() string {
	return "boarding_marks"
}

// BeforeCreate генерирует UUID и коды для новой записи (Ticket).
func (t *Ticket) BeforeCreate(_ *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	if t.QRCode == "" {
		t.QRCode = "TK" + uuid.New().String()[:8]
	}
	if t.BarCode == "" {
		t.BarCode = "BC" + uuid.New().String()[:12]
	}
	return nil
}

// BeforeCreate генерирует UUID для новой записи (BoardingEvent).
func (b *BoardingEvent) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate генерирует UUID для новой записи (BoardingMark).
func (b *BoardingMark) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
