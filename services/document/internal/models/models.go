// Package models — доменные модели Document Service.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentTemplate — шаблон документа.
type DocumentTemplate struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ID          string    `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type"`
	Description string    `gorm:"type:text" json:"description"`
	Content     string    `gorm:"type:text" json:"content"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}

// TableName возвращает имя таблицы шаблонов.
func (DocumentTemplate) TableName() string {
	return "document_templates"
}

// BeforeCreate заполняет ID перед созданием записи.
func (d *DocumentTemplate) BeforeCreate(_ *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

// GeneratedDocument — сгенерированный документ.
type GeneratedDocument struct {
	CreatedAt    time.Time `json:"created_at"`
	TemplateID   *string   `gorm:"type:uuid;index" json:"template_id,omitempty"`
	EntityID     *string   `gorm:"type:uuid;index" json:"entity_id,omitempty"`
	ID           string    `gorm:"type:uuid;primary_key" json:"id"`
	DocumentType string    `gorm:"type:varchar(50);not null" json:"document_type"`
	FileURL      string    `gorm:"type:varchar(500)" json:"file_url"`
	FileName     string    `gorm:"type:varchar(200)" json:"file_name"`
	Status       string    `gorm:"type:varchar(20);default:'generated'" json:"status"`
}

// TableName возвращает имя таблицы сгенерированных документов.
func (GeneratedDocument) TableName() string {
	return "generated_documents"
}

// BeforeCreate заполняет ID перед созданием записи.
func (g *GeneratedDocument) BeforeCreate(_ *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}
