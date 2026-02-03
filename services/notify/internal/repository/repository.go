// Package repository содержит слой доступа к данным уведомлений.
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/vokzal-tech/notify-service/internal/models"
)

// ErrNotificationNotFound возвращается, когда уведомление не найдено.
var ErrNotificationNotFound = errors.New("notification not found")

// NotificationRepository — интерфейс репозитория уведомлений.
type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	FindByID(ctx context.Context, id string) (*models.Notification, error)
	FindByType(ctx context.Context, notifType string, limit int) ([]*models.Notification, error)
	Update(ctx context.Context, notification *models.Notification) error
	List(ctx context.Context, limit int) ([]*models.Notification, error)
}

type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository создаёт репозиторий уведомлений.
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) FindByID(ctx context.Context, id string) (*models.Notification, error) {
	var notification models.Notification
	if err := r.db.WithContext(ctx).First(&notification, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) FindByType(ctx context.Context, notifType string, limit int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := r.db.WithContext(ctx).Where("type = ?", notifType).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

func (r *notificationRepository) List(ctx context.Context, limit int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := r.db.WithContext(ctx).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
