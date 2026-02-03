// Package repository — слой доступа к данным Audit Service.
package repository

import (
	"context"
	"errors"

	"github.com/vokzal-tech/audit-service/internal/models"
	"gorm.io/gorm"
)

var (
	// ErrAuditLogNotFound возвращается, когда запись аудита не найдена.
	ErrAuditLogNotFound = errors.New("audit log not found")
)

// AuditRepository — интерфейс репозитория записей аудита.
type AuditRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	FindByID(ctx context.Context, id string) (*models.AuditLog, error)
	FindByEntity(ctx context.Context, entityType, entityID string) ([]*models.AuditLog, error)
	FindByUser(ctx context.Context, userID string, limit int) ([]*models.AuditLog, error)
	FindByDateRange(ctx context.Context, from, to string) ([]*models.AuditLog, error)
	List(ctx context.Context, limit int) ([]*models.AuditLog, error)
}

type auditRepository struct {
	db *gorm.DB
}

// NewAuditRepository создаёт новый AuditRepository.
func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(ctx context.Context, log *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditRepository) FindByID(ctx context.Context, id string) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := r.db.WithContext(ctx).First(&log, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuditLogNotFound
		}
		return nil, err
	}
	return &log, nil
}

func (r *auditRepository) FindByEntity(ctx context.Context, entityType, entityID string) ([]*models.AuditLog, error) {
	var logs []*models.AuditLog
	if err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *auditRepository) FindByUser(ctx context.Context, userID string, limit int) ([]*models.AuditLog, error) {
	var logs []*models.AuditLog
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *auditRepository) FindByDateRange(ctx context.Context, from, to string) ([]*models.AuditLog, error) {
	var logs []*models.AuditLog
	if err := r.db.WithContext(ctx).
		Where("created_at >= ? AND created_at <= ?", from, to).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *auditRepository) List(ctx context.Context, limit int) ([]*models.AuditLog, error) {
	var logs []*models.AuditLog
	query := r.db.WithContext(ctx).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
