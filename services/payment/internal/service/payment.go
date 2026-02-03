package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/vokzal-tech/payment-service/internal/config"
	"github.com/vokzal-tech/payment-service/internal/models"
	"github.com/vokzal-tech/payment-service/internal/repository"
	"github.com/vokzal-tech/payment-service/internal/sbp"
	"github.com/vokzal-tech/payment-service/internal/tinkoff"
	"go.uber.org/zap"
)

type PaymentService interface {
	// Initialize payments
	InitTinkoffPayment(ctx context.Context, req *InitPaymentRequest) (*models.Payment, error)
	InitSBPPayment(ctx context.Context, req *InitPaymentRequest) (*models.Payment, error)
	InitCashPayment(ctx context.Context, req *InitPaymentRequest) (*models.Payment, error)

	// Check status
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
	GetPaymentByTicket(ctx context.Context, ticketID string) ([]*models.Payment, error)
	CheckPaymentStatus(ctx context.Context, id string) (*models.Payment, error)

	// Webhooks
	HandleTinkoffWebhook(ctx context.Context, data map[string]interface{}) error
	
	// List
	ListPayments(ctx context.Context, limit int) ([]*models.Payment, error)
}

type paymentService struct {
	repo           repository.PaymentRepository
	tinkoffClient  *tinkoff.TinkoffClient
	sbpClient      *sbp.SBPClient
	natsConn       *nats.Conn
	cfg            *config.Config
	logger         *zap.Logger
}

type InitPaymentRequest struct {
	TicketID    *string `json:"ticket_id"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

func NewPaymentService(
	repo repository.PaymentRepository,
	tinkoffClient *tinkoff.TinkoffClient,
	sbpClient *sbp.SBPClient,
	natsConn *nats.Conn,
	cfg *config.Config,
	logger *zap.Logger,
) PaymentService {
	return &paymentService{
		repo:          repo,
		tinkoffClient: tinkoffClient,
		sbpClient:     sbpClient,
		natsConn:      natsConn,
		cfg:           cfg,
		logger:        logger,
	}
}

// InitTinkoffPayment инициализирует платёж через Tinkoff
func (s *paymentService) InitTinkoffPayment(ctx context.Context, req *InitPaymentRequest) (*models.Payment, error) {
	// Создать запись в БД
	payment := &models.Payment{
		TicketID: req.TicketID,
		Amount:   req.Amount,
		Currency: "RUB",
		Method:   "card",
		Provider: "tinkoff",
		Status:   "pending",
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Инициализировать платёж в Tinkoff
	description := req.Description
	if description == "" {
		description = "Билет на автобус"
	}

	result, err := s.tinkoffClient.Init(payment.ID, req.Amount, description)
	if err != nil {
		payment.Status = "failed"
		errMsg := err.Error()
		payment.ErrorMsg = &errMsg
		s.repo.Update(ctx, payment)
		return nil, fmt.Errorf("failed to init tinkoff payment: %w", err)
	}

	// Обновить payment
	payment.ExternalID = &result.PaymentID
	payment.PaymentURL = &result.PaymentURL
	payment.Status = "processing"

	if err := s.repo.Update(ctx, payment); err != nil {
		s.logger.Error("Failed to update payment", zap.Error(err))
	}

	s.logger.Info("Tinkoff payment initialized",
		zap.String("payment_id", payment.ID),
		zap.String("external_id", result.PaymentID))

	return payment, nil
}

// InitSBPPayment инициализирует платёж через СБП
func (s *paymentService) InitSBPPayment(ctx context.Context, req *InitPaymentRequest) (*models.Payment, error) {
	// Создать запись в БД
	payment := &models.Payment{
		TicketID: req.TicketID,
		Amount:   req.Amount,
		Currency: "RUB",
		Method:   "sbp",
		Provider: "sbp",
		Status:   "pending",
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Сгенерировать QR код
	purpose := req.Description
	if purpose == "" {
		purpose = "Билет на автобус"
	}

	result, err := s.sbpClient.GenerateQR(req.Amount, purpose)
	if err != nil {
		payment.Status = "failed"
		errMsg := err.Error()
		payment.ErrorMsg = &errMsg
		s.repo.Update(ctx, payment)
		return nil, fmt.Errorf("failed to generate SBP QR: %w", err)
	}

	// Обновить payment
	payment.ExternalID = &result.PaymentID
	payment.QRCode = &result.QRString
	payment.Status = "processing"

	if err := s.repo.Update(ctx, payment); err != nil {
		s.logger.Error("Failed to update payment", zap.Error(err))
	}

	s.logger.Info("SBP payment initialized",
		zap.String("payment_id", payment.ID),
		zap.String("external_id", result.PaymentID))

	return payment, nil
}

// InitCashPayment создаёт запись о наличной оплате
func (s *paymentService) InitCashPayment(ctx context.Context, req *InitPaymentRequest) (*models.Payment, error) {
	payment := &models.Payment{
		TicketID: req.TicketID,
		Amount:   req.Amount,
		Currency: "RUB",
		Method:   "cash",
		Provider: "manual",
		Status:   "confirmed",
	}

	now := time.Now()
	payment.ConfirmedAt = &now

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Отправить событие подтверждения
	s.publishPaymentEvent("payment.confirmed", payment)

	s.logger.Info("Cash payment created", zap.String("payment_id", payment.ID))

	return payment, nil
}

func (s *paymentService) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *paymentService) GetPaymentByTicket(ctx context.Context, ticketID string) ([]*models.Payment, error) {
	return s.repo.FindByTicketID(ctx, ticketID)
}

// CheckPaymentStatus проверяет статус платежа у провайдера
func (s *paymentService) CheckPaymentStatus(ctx context.Context, id string) (*models.Payment, error) {
	payment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if payment.Status == "confirmed" || payment.Status == "failed" {
		return payment, nil
	}

	if payment.ExternalID == nil {
		return payment, nil
	}

	switch payment.Provider {
	case "tinkoff":
		result, err := s.tinkoffClient.GetState(*payment.ExternalID)
		if err != nil {
			s.logger.Error("Failed to check Tinkoff status", zap.Error(err))
			return payment, nil
		}

		if result.Status == "CONFIRMED" || result.Status == "AUTHORIZED" {
			payment.Status = "confirmed"
			now := time.Now()
			payment.ConfirmedAt = &now
			s.publishPaymentEvent("payment.confirmed", payment)
		} else if result.Status == "REJECTED" {
			payment.Status = "failed"
			errMsg := "Payment rejected"
			payment.ErrorMsg = &errMsg
		}

	case "sbp":
		result, err := s.sbpClient.GetStatus(*payment.ExternalID)
		if err != nil {
			s.logger.Error("Failed to check SBP status", zap.Error(err))
			return payment, nil
		}

		if result.Status == "paid" {
			payment.Status = "confirmed"
			payment.ConfirmedAt = result.PaidAt
			s.publishPaymentEvent("payment.confirmed", payment)
		} else if result.Status == "expired" || result.Status == "cancelled" {
			payment.Status = "failed"
			errMsg := fmt.Sprintf("Payment %s", result.Status)
			payment.ErrorMsg = &errMsg
		}
	}

	if err := s.repo.Update(ctx, payment); err != nil {
		s.logger.Error("Failed to update payment status", zap.Error(err))
	}

	return payment, nil
}

// HandleTinkoffWebhook обрабатывает webhook от Tinkoff
func (s *paymentService) HandleTinkoffWebhook(ctx context.Context, data map[string]interface{}) error {
	paymentID, ok := data["PaymentId"].(string)
	if !ok {
		return fmt.Errorf("invalid webhook data: no PaymentId")
	}

	status, ok := data["Status"].(string)
	if !ok {
		return fmt.Errorf("invalid webhook data: no Status")
	}

	payment, err := s.repo.FindByExternalID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	if status == "CONFIRMED" || status == "AUTHORIZED" {
		payment.Status = "confirmed"
		now := time.Now()
		payment.ConfirmedAt = &now
		s.publishPaymentEvent("payment.confirmed", payment)
	} else if status == "REJECTED" {
		payment.Status = "failed"
		errMsg := "Payment rejected"
		payment.ErrorMsg = &errMsg
	}

	if err := s.repo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	s.logger.Info("Webhook processed",
		zap.String("payment_id", payment.ID),
		zap.String("status", status))

	return nil
}

func (s *paymentService) ListPayments(ctx context.Context, limit int) ([]*models.Payment, error) {
	return s.repo.List(ctx, limit)
}

// NATS Events
func (s *paymentService) publishPaymentEvent(subject string, payment *models.Payment) {
	data, err := json.Marshal(payment)
	if err != nil {
		s.logger.Error("Failed to marshal payment event", zap.Error(err))
		return
	}

	if err := s.natsConn.Publish(subject, data); err != nil {
		s.logger.Error("Failed to publish payment event", zap.Error(err), zap.String("subject", subject))
	}
}
