package repository

import (
	"context"
	"errors"

	"github.com/vokzal-tech/fiscal-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrReceiptNotFound = errors.New("receipt not found")
	ErrZReportNotFound = errors.New("z-report not found")
)

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

func NewFiscalRepository(db *gorm.DB) FiscalRepository {
	return &fiscalRepository{db: db}
}

// Receipts
func (r *fiscalRepository) CreateReceipt(ctx context.Context, receipt *models.FiscalReceipt) error {
	return r.db.WithContext(ctx).Create(receipt).Error
}

func (r *fiscalRepository) FindReceiptByID(ctx context.Context, id string) (*models.FiscalReceipt, error) {
	var receipt models.FiscalReceipt
	if err := r.db.WithContext(ctx).First(&receipt, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReceiptNotFound
		}
		return nil, err
	}
	return &receipt, nil
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

// Z-Reports
func (r *fiscalRepository) CreateZReport(ctx context.Context, report *models.ZReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *fiscalRepository) FindZReportByDate(ctx context.Context, date string) (*models.ZReport, error) {
	var report models.ZReport
	if err := r.db.WithContext(ctx).First(&report, "date = ?", date).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZReportNotFound
		}
		return nil, err
	}
	return &report, nil
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
