// Package models содержит модели данных сервиса билетов.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ticket — модель билета.
type Ticket struct {
	ID            string     `gorm:"type:uuid;primary_key" json:"id"`
	TripID        string     `gorm:"type:uuid;not null;index" json:"trip_id"`
	SeatID        *string    `gorm:"type:uuid;index" json:"seat_id,omitempty"`
	PassengerName *string    `gorm:"type:varchar(100)" json:"passenger_name,omitempty"`
	PassengerDoc  *string    `gorm:"type:varchar(50)" json:"passenger_doc,omitempty"`
	Phone         *string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Email         *string    `gorm:"type:varchar(100)" json:"email,omitempty"`
	Price         float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Status        string     `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	PaymentMethod string     `gorm:"type:varchar(20)" json:"payment_method"`
	QRCode        string     `gorm:"type:varchar(255);unique" json:"qr_code"`
	BarCode       string     `gorm:"type:varchar(255);unique" json:"bar_code"`
	RefundedAt    *time.Time `json:"refunded_at,omitempty"`
	RefundAmount  *float64   `gorm:"type:decimal(10,2)" json:"refund_amount,omitempty"`
	RefundPenalty *float64   `gorm:"type:decimal(10,2)" json:"refund_penalty,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// BoardingEvent — модель события начала посадки.
type BoardingEvent struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	TripID    string    `gorm:"type:uuid;not null;unique;index" json:"trip_id"`
	StartedAt time.Time `gorm:"not null" json:"started_at"`
	StartedBy string    `gorm:"type:uuid;not null" json:"started_by"`
	CreatedAt time.Time `json:"created_at"`
}

// BoardingMark — модель отметки посадки.
type BoardingMark struct {
	ID         string    `gorm:"type:uuid;primary_key" json:"id"`
	TicketID   string    `gorm:"type:uuid;not null;index" json:"ticket_id"`
	MarkedAt   time.Time `gorm:"not null" json:"marked_at"`
	MarkedBy   string    `gorm:"type:uuid;not null" json:"marked_by"`
	ScanMethod string    `gorm:"type:varchar(20)" json:"scan_method"`
	CreatedAt  time.Time `json:"created_at"`
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
