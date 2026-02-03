// Package repository содержит слой доступа к данным платежей.
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/vokzal-tech/payment-service/internal/models"
)

// ErrPaymentNotFound возвращается, когда платёж не найден.
var ErrPaymentNotFound = errors.New("payment not found")

// PaymentRepository — интерфейс репозитория платежей.
type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	FindByID(ctx context.Context, id string) (*models.Payment, error)
	FindByExternalID(ctx context.Context, externalID string) (*models.Payment, error)
	FindByTicketID(ctx context.Context, ticketID string) ([]*models.Payment, error)
	Update(ctx context.Context, payment *models.Payment) error
	List(ctx context.Context, limit int) ([]*models.Payment, error)
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository создаёт репозиторий платежей.
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func findFirstBy[T any](db *gorm.DB, ctx context.Context, query string, arg any, notFoundErr error) (*T, error) {
	var t T
	if err := db.WithContext(ctx).First(&t, query, arg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFoundErr
		}
		return nil, err
	}
	return &t, nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id string) (*models.Payment, error) {
	return findFirstBy[models.Payment](r.db, ctx, "id = ?", id, ErrPaymentNotFound)
}

func (r *paymentRepository) FindByExternalID(ctx context.Context, externalID string) (*models.Payment, error) {
	return findFirstBy[models.Payment](r.db, ctx, "external_id = ?", externalID, ErrPaymentNotFound)
}

func (r *paymentRepository) FindByTicketID(ctx context.Context, ticketID string) ([]*models.Payment, error) {
	var payments []*models.Payment
	if err := r.db.WithContext(ctx).Where("ticket_id = ?", ticketID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) List(ctx context.Context, limit int) ([]*models.Payment, error) {
	var payments []*models.Payment
	query := r.db.WithContext(ctx).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
