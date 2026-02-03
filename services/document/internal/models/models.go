package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentTemplate шаблон документа
type DocumentTemplate struct {
	ID          string    `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type"` // pd2, ticket, invoice, custom
	Description string    `gorm:"type:text" json:"description"`
	Content     string    `gorm:"type:text" json:"content"` // HTML или простой шаблон
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (DocumentTemplate) TableName() string {
	return "document_templates"
}

func (d *DocumentTemplate) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

// GeneratedDocument сгенерированный документ
type GeneratedDocument struct {
	ID           string    `gorm:"type:uuid;primary_key" json:"id"`
	TemplateID   *string   `gorm:"type:uuid;index" json:"template_id,omitempty"`
	DocumentType string    `gorm:"type:varchar(50);not null" json:"document_type"`
	EntityID     *string   `gorm:"type:uuid;index" json:"entity_id,omitempty"` // ID билета, рейса и т.д.
	FileURL      string    `gorm:"type:varchar(500)" json:"file_url"`
	FileName     string    `gorm:"type:varchar(200)" json:"file_name"`
	Status       string    `gorm:"type:varchar(20);default:'generated'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func (GeneratedDocument) TableName() string {
	return "generated_documents"
}

func (g *GeneratedDocument) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}
