// Package repository — слой доступа к данным Fiscal Service.
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/vokzal-tech/fiscal-service/internal/models"
)

var (
	// ErrReceiptNotFound возвращается, когда чек не найден.
	ErrReceiptNotFound = errors.New("receipt not found")
	// ErrZReportNotFound возвращается, когда Z-отчёт не найден.
	ErrZReportNotFound = errors.New("z-report not found")
)

// FiscalRepository — интерфейс репозитория фискальных данных.
type FiscalRepository interface {
	// Receipts
	CreateReceipt(ctx context.Context, receipt *models.FiscalReceipt) error
	FindReceiptByID(ctx context.Context, id string) (*models.FiscalReceipt, error)
	FindReceiptByTicketID(ctx context.Context, ticketID string) ([]*models.FiscalReceipt, error)
	UpdateReceipt(ctx context.Context, receipt *models.FiscalReceipt) error

	// Z-Reports
	CreateZReport(ctx context.Context, report *models.ZReport) error
	FindZReportByDate(ctx context.Context, date string) (*models.ZReport, error)
	FindAllZReports(ctx context.Context, limit int) ([]*models.ZReport, error)
	UpdateZReport(ctx context.Context, report *models.ZReport) error
}

type fiscalRepository struct {
	db *gorm.DB
}

// NewFiscalRepository создаёт новый FiscalRepository.
func NewFiscalRepository(db *gorm.DB) FiscalRepository {
	return &fiscalRepository{db: db}
}

// CreateReceipt создаёт запись чека.
func (r *fiscalRepository) CreateReceipt(ctx context.Context, receipt *models.FiscalReceipt) error {
	return r.db.WithContext(ctx).Create(receipt).Error
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

func (r *fiscalRepository) FindReceiptByID(ctx context.Context, id string) (*models.FiscalReceipt, error) {
	return findFirstBy[models.FiscalReceipt](r.db, ctx, "id = ?", id, ErrReceiptNotFound)
}

func (r *fiscalRepository) FindReceiptByTicketID(ctx context.Context, ticketID string) ([]*models.FiscalReceipt, error) {
	var receipts []*models.FiscalReceipt
	if err := r.db.WithContext(ctx).Where("ticket_id = ?", ticketID).Find(&receipts).Error; err != nil {
		return nil, err
	}
	return receipts, nil
}

func (r *fiscalRepository) UpdateReceipt(ctx context.Context, receipt *models.FiscalReceipt) error {
	return r.db.WithContext(ctx).Save(receipt).Error
}

// CreateZReport создаёт запись Z-отчёта.
func (r *fiscalRepository) CreateZReport(ctx context.Context, report *models.ZReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *fiscalRepository) FindZReportByDate(ctx context.Context, date string) (*models.ZReport, error) {
	return findFirstBy[models.ZReport](r.db, ctx, "date = ?", date, ErrZReportNotFound)
}

func (r *fiscalRepository) FindAllZReports(ctx context.Context, limit int) ([]*models.ZReport, error) {
	var reports []*models.ZReport
	query := r.db.WithContext(ctx).Order("date DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *fiscalRepository) UpdateZReport(ctx context.Context, report *models.ZReport) error {
	return r.db.WithContext(ctx).Save(report).Error
}
