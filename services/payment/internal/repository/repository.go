// Package repository содержит слой доступа к данным платежей.
package repository

import (
	"context"
	"errors"

	"github.com/vokzal-tech/payment-service/internal/models"
	"gorm.io/gorm"
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

func (r *paymentRepository) FindByID(ctx context.Context, id string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.WithContext(ctx).First(&payment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByExternalID(ctx context.Context, externalID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.WithContext(ctx).First(&payment, "external_id = ?", externalID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
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
