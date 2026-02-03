// Package models — доменные модели Fiscal Service.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FiscalReceipt — модель фискального чека.
type FiscalReceipt struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ErrorMsg   *string   `gorm:"type:text" json:"error_msg,omitempty"`
	ID         string    `gorm:"type:uuid;primary_key" json:"id"`
	TicketID   string    `gorm:"type:uuid;not null;index" json:"ticket_id"`
	Type       string    `gorm:"type:varchar(20);not null" json:"type"`
	OFDURL     string    `gorm:"type:varchar(500)" json:"ofd_url"`
	KKTSerial  string    `gorm:"type:varchar(50)" json:"kkt_serial"`
	FiscalSign string    `gorm:"type:varchar(100)" json:"fiscal_sign"`
	Status     string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Amount     float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
}

// ZReport — модель Z-отчёта.
type ZReport struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ID           string    `gorm:"type:uuid;primary_key" json:"id"`
	Date         string    `gorm:"type:date;not null;index" json:"date"`
	KKTSerial    string    `gorm:"type:varchar(50);not null" json:"kkt_serial"`
	Status       string    `gorm:"type:varchar(20);not null" json:"status"`
	FiscalSign   string    `gorm:"type:varchar(100)" json:"fiscal_sign"`
	ShiftNumber  int       `gorm:"not null" json:"shift_number"`
	TotalSales   float64   `gorm:"type:decimal(10,2)" json:"total_sales"`
	TotalRefunds float64   `gorm:"type:decimal(10,2)" json:"total_refunds"`
	SalesCount   int       `json:"sales_count"`
	RefundsCount int       `json:"refunds_count"`
}

// TableName возвращает имя таблицы чеков.
func (FiscalReceipt) TableName() string {
	return "fiscal_receipts"
}

// TableName возвращает имя таблицы Z-отчётов.
func (ZReport) TableName() string {
	return "z_reports"
}

// BeforeCreate заполняет ID перед созданием записи.
func (f *FiscalReceipt) BeforeCreate(_ *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate заполняет ID перед созданием записи.
func (z *ZReport) BeforeCreate(_ *gorm.DB) error {
	if z.ID == "" {
		z.ID = uuid.New().String()
	}
	return nil
}
