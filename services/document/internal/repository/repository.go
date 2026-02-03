// Package repository — слой доступа к данным Document Service.
package repository

import (
	"context"
	"errors"

	"github.com/vokzal-tech/document-service/internal/models"
	"gorm.io/gorm"
)

var (
	// ErrTemplateNotFound возвращается, когда шаблон не найден.
	ErrTemplateNotFound = errors.New("template not found")
	// ErrDocumentNotFound возвращается, когда документ не найден.
	ErrDocumentNotFound = errors.New("document not found")
)

// DocumentRepository — интерфейс репозитория документов.
type DocumentRepository interface {
	CreateTemplate(ctx context.Context, template *models.DocumentTemplate) error
	FindTemplateByID(ctx context.Context, id string) (*models.DocumentTemplate, error)
	FindTemplateByName(ctx context.Context, name string) (*models.DocumentTemplate, error)
	ListTemplates(ctx context.Context) ([]*models.DocumentTemplate, error)
	UpdateTemplate(ctx context.Context, template *models.DocumentTemplate) error
	
	CreateDocument(ctx context.Context, doc *models.GeneratedDocument) error
	FindDocumentByID(ctx context.Context, id string) (*models.GeneratedDocument, error)
	FindDocumentsByEntity(ctx context.Context, entityID string) ([]*models.GeneratedDocument, error)
	ListDocuments(ctx context.Context, limit int) ([]*models.GeneratedDocument, error)
}

type documentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository создаёт новый DocumentRepository.
func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) CreateTemplate(ctx context.Context, template *models.DocumentTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *documentRepository) FindTemplateByID(ctx context.Context, id string) (*models.DocumentTemplate, error) {
	var template models.DocumentTemplate
	if err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	return &template, nil
}

func (r *documentRepository) FindTemplateByName(ctx context.Context, name string) (*models.DocumentTemplate, error) {
	var template models.DocumentTemplate
	if err := r.db.WithContext(ctx).First(&template, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	return &template, nil
}

func (r *documentRepository) ListTemplates(ctx context.Context) ([]*models.DocumentTemplate, error) {
	var templates []*models.DocumentTemplate
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *documentRepository) UpdateTemplate(ctx context.Context, template *models.DocumentTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *documentRepository) CreateDocument(ctx context.Context, doc *models.GeneratedDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *documentRepository) FindDocumentByID(ctx context.Context, id string) (*models.GeneratedDocument, error) {
	var doc models.GeneratedDocument
	if err := r.db.WithContext(ctx).First(&doc, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) FindDocumentsByEntity(ctx context.Context, entityID string) ([]*models.GeneratedDocument, error) {
	var docs []*models.GeneratedDocument
	if err := r.db.WithContext(ctx).Where("entity_id = ?", entityID).Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *documentRepository) ListDocuments(ctx context.Context, limit int) ([]*models.GeneratedDocument, error) {
	var docs []*models.GeneratedDocument
	query := r.db.WithContext(ctx).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}
