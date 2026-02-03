// Package service — бизнес-логика Fiscal Service.
package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/vokzal-tech/fiscal-service/internal/atol"
	"github.com/vokzal-tech/fiscal-service/internal/config"
	"github.com/vokzal-tech/fiscal-service/internal/models"
	"github.com/vokzal-tech/fiscal-service/internal/repository"
	"go.uber.org/zap"
)

const receiptStatusFailed = "failed"

// FiscalService — интерфейс сервиса фискализации.
type FiscalService interface {
	// Receipts
	ProcessTicketSold(ctx context.Context, ticketData map[string]interface{}) error
	ProcessTicketRefund(ctx context.Context, ticketData map[string]interface{}) error
	GetReceipt(ctx context.Context, id string) (*models.FiscalReceipt, error)
	GetReceiptsByTicket(ctx context.Context, ticketID string) ([]*models.FiscalReceipt, error)

	// Z-Reports
	CreateDailyZReport(ctx context.Context, date string) (*models.ZReport, error)
	GetZReport(ctx context.Context, date string) (*models.ZReport, error)
	ListZReports(ctx context.Context, limit int) ([]*models.ZReport, error)

	// KKT Status
	GetKKTStatus(ctx context.Context) (map[string]interface{}, error)

	// NATS Events
	SubscribeToEvents(nc *nats.Conn)
}

type fiscalService struct {
	repo       repository.FiscalRepository
	atolClient *atol.ATOLClient
	cfg        *config.Config
	logger     *zap.Logger
}

// NewFiscalService создаёт новый FiscalService.
func NewFiscalService(
	repo repository.FiscalRepository,
	atolClient *atol.ATOLClient,
	cfg *config.Config,
	logger *zap.Logger,
) FiscalService {
	return &fiscalService{
		repo:       repo,
		atolClient: atolClient,
		cfg:        cfg,
		logger:     logger,
	}
}

// ProcessTicketSold обрабатывает продажу билета (фискализация чека).
//
//nolint:dupl // структура совпадает с ProcessTicketRefund (продажа/возврат)
func (s *fiscalService) ProcessTicketSold(ctx context.Context, ticketData map[string]interface{}) error {
	ticketID, _ := ticketData["id"].(string)
	price, _ := ticketData["price"].(float64)

	// Создать запись чека
	receipt := &models.FiscalReceipt{
		TicketID: ticketID,
		Type:     "sale",
		Amount:   price,
		Status:   "pending",
	}

	if err := s.repo.CreateReceipt(ctx, receipt); err != nil {
		return fmt.Errorf("failed to create receipt: %w", err)
	}

	// Отправить на ККТ
	req := &atol.ReceiptRequest{
		Operation: "sell",
		Items: []atol.ReceiptItem{
			{
				Name:     "Билет на автобус",
				Quantity: 1,
				Price:    price,
				VAT:      "none",
			},
		},
		Payment: atol.Payment{
			Type:   "card",
			Amount: price,
		},
		Company: atol.Company{
			INN:       s.cfg.ATOL.CompanyINN,
			Name:      s.cfg.ATOL.CompanyName,
			TaxSystem: s.cfg.ATOL.TaxSystem,
		},
	}

	result, err := s.atolClient.PrintReceipt(req)
	if err != nil {
		receipt.Status = receiptStatusFailed
		errMsg := err.Error()
		receipt.ErrorMsg = &errMsg
		_ = s.repo.UpdateReceipt(ctx, receipt)
		return fmt.Errorf("failed to print receipt: %w", err)
	}

	if result.Success {
		receipt.Status = "confirmed"
		receipt.OFDURL = result.OFDURL
		receipt.FiscalSign = result.FiscalSign
		receipt.KKTSerial = result.KKTSerial
	} else {
		receipt.Status = receiptStatusFailed
		receipt.ErrorMsg = &result.ErrorMsg
	}

	if updateErr := s.repo.UpdateReceipt(ctx, receipt); updateErr != nil {
		s.logger.Error("Failed to update receipt", zap.Error(updateErr))
	}

	s.logger.Info("Receipt processed",
		zap.String("receipt_id", receipt.ID),
		zap.String("ticket_id", ticketID),
		zap.String("status", receipt.Status))

	return nil
}

// ProcessTicketRefund обрабатывает возврат билета (фискализация чека возврата).
//
//nolint:dupl // структура совпадает с ProcessTicketSold (продажа/возврат), рефакторинг усложнит чтение
func (s *fiscalService) ProcessTicketRefund(ctx context.Context, ticketData map[string]interface{}) error {
	ticketID, _ := ticketData["id"].(string)
	refundAmount, _ := ticketData["refund_amount"].(float64)

	// Создать запись чека возврата
	receipt := &models.FiscalReceipt{
		TicketID: ticketID,
		Type:     "refund",
		Amount:   refundAmount,
		Status:   "pending",
	}

	if err := s.repo.CreateReceipt(ctx, receipt); err != nil {
		return fmt.Errorf("failed to create refund receipt: %w", err)
	}

	// Отправить на ККТ
	req := &atol.ReceiptRequest{
		Operation: "refund",
		Items: []atol.ReceiptItem{
			{
				Name:     "Возврат билета на автобус",
				Quantity: 1,
				Price:    refundAmount,
				VAT:      "none",
			},
		},
		Payment: atol.Payment{
			Type:   "card",
			Amount: refundAmount,
		},
		Company: atol.Company{
			INN:       s.cfg.ATOL.CompanyINN,
			Name:      s.cfg.ATOL.CompanyName,
			TaxSystem: s.cfg.ATOL.TaxSystem,
		},
	}

	result, err := s.atolClient.PrintReceipt(req)
	if err != nil {
		receipt.Status = receiptStatusFailed
		errMsg := err.Error()
		receipt.ErrorMsg = &errMsg
		_ = s.repo.UpdateReceipt(ctx, receipt)
		return fmt.Errorf("failed to print refund receipt: %w", err)
	}

	if result.Success {
		receipt.Status = "confirmed"
		receipt.OFDURL = result.OFDURL
		receipt.FiscalSign = result.FiscalSign
		receipt.KKTSerial = result.KKTSerial
	} else {
		receipt.Status = receiptStatusFailed
		receipt.ErrorMsg = &result.ErrorMsg
	}

	if updateErr := s.repo.UpdateReceipt(ctx, receipt); updateErr != nil {
		s.logger.Error("Failed to update refund receipt", zap.Error(updateErr))
	}

	s.logger.Info("Refund receipt processed",
		zap.String("receipt_id", receipt.ID),
		zap.String("ticket_id", ticketID),
		zap.String("status", receipt.Status))

	return nil
}

func (s *fiscalService) GetReceipt(ctx context.Context, id string) (*models.FiscalReceipt, error) {
	return s.repo.FindReceiptByID(ctx, id)
}

func (s *fiscalService) GetReceiptsByTicket(ctx context.Context, ticketID string) ([]*models.FiscalReceipt, error) {
	return s.repo.FindReceiptByTicketID(ctx, ticketID)
}

// CreateDailyZReport создаёт дневной Z-отчёт.
func (s *fiscalService) CreateDailyZReport(ctx context.Context, date string) (*models.ZReport, error) {
	// Проверить, не существует ли уже отчёт
	existing, err := s.repo.FindZReportByDate(ctx, date)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Создать Z-отчёт на ККТ
	result, err := s.atolClient.CreateZReport()
	if err != nil {
		return nil, fmt.Errorf("failed to create Z-report on KKT: %w", err)
	}

	// Серийный номер ККТ из ответа агента или значение по умолчанию
	kktSerial := result.KKTSerial
	if kktSerial == "" {
		kktSerial = "KKT001"
	}

	// Сохранить в БД
	report := &models.ZReport{
		Date:         date,
		KKTSerial:    kktSerial,
		ShiftNumber:  result.ShiftNumber,
		TotalSales:   result.TotalSales,
		TotalRefunds: result.TotalRefunds,
		SalesCount:   result.SalesCount,
		RefundsCount: result.RefundsCount,
		FiscalSign:   result.FiscalSign,
		Status:       "completed",
	}

	if err := s.repo.CreateZReport(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to save Z-report: %w", err)
	}

	s.logger.Info("Z-report created",
		zap.String("date", date),
		zap.Int("shift", result.ShiftNumber),
		zap.Float64("sales", result.TotalSales))

	return report, nil
}

func (s *fiscalService) GetZReport(ctx context.Context, date string) (*models.ZReport, error) {
	return s.repo.FindZReportByDate(ctx, date)
}

func (s *fiscalService) ListZReports(ctx context.Context, limit int) ([]*models.ZReport, error) {
	return s.repo.FindAllZReports(ctx, limit)
}

func (s *fiscalService) GetKKTStatus(_ context.Context) (map[string]interface{}, error) {
	return s.atolClient.GetKKTStatus()
}

// SubscribeToEvents подписывается на NATS-события для фискализации.
func (s *fiscalService) SubscribeToEvents(nc *nats.Conn) {
	_, _ = nc.Subscribe("ticket.sold", func(msg *nats.Msg) {
		var ticketData map[string]interface{}
		if err := json.Unmarshal(msg.Data, &ticketData); err != nil {
			s.logger.Error("Failed to unmarshal ticket.sold event", zap.Error(err))
			return
		}

		ctx := context.Background()
		if err := s.ProcessTicketSold(ctx, ticketData); err != nil {
			s.logger.Error("Failed to process ticket.sold", zap.Error(err))
		}
	})

	_, _ = nc.Subscribe("ticket.returned", func(msg *nats.Msg) {
		var ticketData map[string]interface{}
		if err := json.Unmarshal(msg.Data, &ticketData); err != nil {
			s.logger.Error("Failed to unmarshal ticket.returned event", zap.Error(err))
			return
		}

		ctx := context.Background()
		if err := s.ProcessTicketRefund(ctx, ticketData); err != nil {
			s.logger.Error("Failed to process ticket.returned", zap.Error(err))
		}
	})

	s.logger.Info("Subscribed to NATS events: ticket.sold, ticket.returned")
}
