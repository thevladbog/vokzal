// Package service содержит бизнес-логику продажи, возврата и посадки по билетам.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/vokzal-tech/ticket-service/internal/config"
	"github.com/vokzal-tech/ticket-service/internal/models"
	"github.com/vokzal-tech/ticket-service/internal/repository"
)

// TicketService — интерфейс сервиса билетов (продажа, возврат, посадка).
type TicketService interface {
	// Продажа
	SellTicket(ctx context.Context, req *SellTicketRequest) (*models.Ticket, error)
	GetTicket(ctx context.Context, id string) (*models.Ticket, error)
	GetTicketByQR(ctx context.Context, qrCode string) (*models.Ticket, error)
	ListTicketsByTrip(ctx context.Context, tripID string) ([]*models.Ticket, error)

	// Возврат
	RefundTicket(ctx context.Context, ticketID string, userID string) (*RefundResult, error)

	// Посадка
	StartBoarding(ctx context.Context, tripID string, userID string) error
	MarkBoarding(ctx context.Context, req *MarkBoardingRequest) error
	GetBoardingStatus(ctx context.Context, tripID string) (*BoardingStatus, error)
}

type ticketService struct {
	ticketRepo   repository.TicketRepository
	boardingRepo repository.BoardingRepository
	natsConn     *nats.Conn
	cfg          *config.Config
	logger       *zap.Logger
}

// SellTicketRequest — запрос на продажу билета.
type SellTicketRequest struct {
	SeatID        *string `json:"seat_id"`
	PassengerName *string `json:"passenger_name"`
	PassengerDoc  *string `json:"passenger_doc"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	TripID        string  `json:"trip_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Price         float64 `json:"price" binding:"required,gt=0"`
}

// RefundResult — результат возврата билета.
type RefundResult struct {
	OriginalAmount float64 `json:"original_amount"`
	Penalty        float64 `json:"penalty"`
	RefundAmount   float64 `json:"refund_amount"`
}

// MarkBoardingRequest — запрос на отметку посадки.
type MarkBoardingRequest struct {
	TicketID   string `json:"ticket_id" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
	ScanMethod string `json:"scan_method"`
}

// BoardingStatus — статус посадки по рейсу.
type BoardingStatus struct {
	StartedAt      *time.Time `json:"started_at,omitempty"`
	TripID         string     `json:"trip_id"`
	TotalTickets   int        `json:"total_tickets"`
	BoardedCount   int        `json:"boarded_count"`
	BoardingActive bool       `json:"boarding_active"`
}

// NewTicketService создаёт сервис билетов.
func NewTicketService(
	ticketRepo repository.TicketRepository,
	boardingRepo repository.BoardingRepository,
	natsConn *nats.Conn,
	cfg *config.Config,
	logger *zap.Logger,
) TicketService {
	return &ticketService{
		ticketRepo:   ticketRepo,
		boardingRepo: boardingRepo,
		natsConn:     natsConn,
		cfg:          cfg,
		logger:       logger,
	}
}

// SellTicket продаёт билет.
func (s *ticketService) SellTicket(ctx context.Context, req *SellTicketRequest) (*models.Ticket, error) {
	// Проверить доступность места
	if req.SeatID != nil {
		available, err := s.ticketRepo.CheckSeatAvailability(ctx, req.TripID, *req.SeatID)
		if err != nil {
			return nil, fmt.Errorf("failed to check seat availability: %w", err)
		}
		if !available {
			return nil, repository.ErrSeatAlreadyTaken
		}
	}

	// Создать билет
	ticket := &models.Ticket{
		TripID:        req.TripID,
		SeatID:        req.SeatID,
		PassengerName: req.PassengerName,
		PassengerDoc:  req.PassengerDoc,
		Phone:         req.Phone,
		Email:         req.Email,
		Price:         req.Price,
		Status:        "active",
		PaymentMethod: req.PaymentMethod,
	}

	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Отправить событие в NATS для фискализации
	s.publishTicketEvent("ticket.sold", ticket)

	s.logger.Info("Ticket sold",
		zap.String("ticket_id", ticket.ID),
		zap.String("trip_id", ticket.TripID),
		zap.Float64("price", ticket.Price))

	return ticket, nil
}

func (s *ticketService) GetTicket(ctx context.Context, id string) (*models.Ticket, error) {
	return s.ticketRepo.FindByID(ctx, id)
}

func (s *ticketService) GetTicketByQR(ctx context.Context, qrCode string) (*models.Ticket, error) {
	return s.ticketRepo.FindByQRCode(ctx, qrCode)
}

func (s *ticketService) ListTicketsByTrip(ctx context.Context, tripID string) ([]*models.Ticket, error) {
	return s.ticketRepo.FindByTripID(ctx, tripID)
}

// RefundTicket возвращает билет.
func (s *ticketService) RefundTicket(ctx context.Context, ticketID, userID string) (*RefundResult, error) {
	// Получить билет
	ticket, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	// Проверить статус
	if ticket.Status != "active" {
		return nil, fmt.Errorf("ticket is not active, current status: %s", ticket.Status)
	}

	// Проверить, не началась ли посадка
	boardingEvent, err := s.boardingRepo.FindEventByTripID(ctx, ticket.TripID)
	if err != nil {
		return nil, fmt.Errorf("failed to check boarding status: %w", err)
	}
	if boardingEvent != nil {
		return nil, repository.ErrBoardingAlreadyStarted
	}

	departureTime, err := s.ticketRepo.GetTripDepartureTime(ctx, ticket.TripID)
	if err != nil {
		s.logger.Warn("GetTripDepartureTime failed, using nil for penalty", zap.Error(err), zap.String("trip_id", ticket.TripID))
	}
	penalty := s.calculateRefundPenalty(ticket.Price, departureTime)

	// Обновить билет
	now := time.Now()
	refundAmount := ticket.Price - penalty
	ticket.Status = "returned"
	ticket.RefundedAt = &now
	ticket.RefundAmount = &refundAmount
	ticket.RefundPenalty = &penalty

	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	// Отправить событие для фискализации возврата
	s.publishTicketEvent("ticket.returned", ticket)

	// Логировать в audit (через NATS)
	s.publishAuditEvent(ctx, "ticket", ticketID, "refund", userID, ticket.Price, refundAmount)

	s.logger.Info("Ticket refunded",
		zap.String("ticket_id", ticketID),
		zap.Float64("original", ticket.Price),
		zap.Float64("penalty", penalty),
		zap.Float64("refund", refundAmount))

	return &RefundResult{
		OriginalAmount: ticket.Price,
		Penalty:        penalty,
		RefundAmount:   refundAmount,
	}, nil
}

// Расчёт штрафа за возврат (время отправления рейса из БД trips+schedules).
func (s *ticketService) calculateRefundPenalty(price float64, departureTime *time.Time) float64 {
	hoursUntilDeparture := 25.0
	if departureTime != nil && !departureTime.IsZero() {
		hoursUntilDeparture = time.Until(*departureTime).Hours()
		if hoursUntilDeparture < 0 {
			hoursUntilDeparture = 0
		}
	}

	var penaltyRate float64
	switch {
	case hoursUntilDeparture > 24:
		penaltyRate = s.cfg.Business.RefundPenalty.Over24Hours
	case hoursUntilDeparture >= 12:
		penaltyRate = s.cfg.Business.RefundPenalty.Between12_24
	default:
		penaltyRate = s.cfg.Business.RefundPenalty.Under12Hours
	}
	return price * penaltyRate
}

// StartBoarding начинает посадку (блокировка возвратов).
func (s *ticketService) StartBoarding(ctx context.Context, tripID, userID string) error {
	// Проверить, не началась ли уже посадка
	existing, err := s.boardingRepo.FindEventByTripID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to check boarding status: %w", err)
	}
	if existing != nil {
		return repository.ErrBoardingAlreadyStarted
	}

	// Создать событие начала посадки
	event := &models.BoardingEvent{
		TripID:    tripID,
		StartedAt: time.Now(),
		StartedBy: userID,
	}

	if err := s.boardingRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to create boarding event: %w", err)
	}

	// Отправить событие в NATS
	s.publishBoardingEvent("boarding.started", map[string]interface{}{
		"trip_id":    tripID,
		"started_at": event.StartedAt,
		"started_by": userID,
	})

	s.logger.Info("Boarding started", zap.String("trip_id", tripID), zap.String("user_id", userID))

	return nil
}

// MarkBoarding отмечает посадку пассажира.
func (s *ticketService) MarkBoarding(ctx context.Context, req *MarkBoardingRequest) error {
	// Получить билет
	ticket, err := s.ticketRepo.FindByID(ctx, req.TicketID)
	if err != nil {
		return err
	}

	// Проверить статус билета
	if ticket.Status != "active" {
		return fmt.Errorf("ticket is not active, current status: %s", ticket.Status)
	}

	// Проверить, началась ли посадка
	boardingEvent, err := s.boardingRepo.FindEventByTripID(ctx, ticket.TripID)
	if err != nil {
		return fmt.Errorf("failed to check boarding: %w", err)
	}
	if boardingEvent == nil {
		return repository.ErrBoardingNotStarted
	}

	// Проверить, не отмечен ли уже
	marked, err := s.boardingRepo.CheckIfMarked(ctx, req.TicketID)
	if err != nil {
		return fmt.Errorf("failed to check if marked: %w", err)
	}
	if marked {
		return fmt.Errorf("ticket already marked for boarding")
	}

	// Создать отметку посадки
	mark := &models.BoardingMark{
		TicketID:   req.TicketID,
		MarkedAt:   time.Now(),
		MarkedBy:   req.UserID,
		ScanMethod: req.ScanMethod,
	}

	if err := s.boardingRepo.CreateMark(ctx, mark); err != nil {
		return fmt.Errorf("failed to create boarding mark: %w", err)
	}

	s.logger.Info("Boarding marked",
		zap.String("ticket_id", req.TicketID),
		zap.String("user_id", req.UserID))

	return nil
}

// GetBoardingStatus возвращает статус посадки.
func (s *ticketService) GetBoardingStatus(ctx context.Context, tripID string) (*BoardingStatus, error) {
	// Проверить событие посадки
	event, err := s.boardingRepo.FindEventByTripID(ctx, tripID)
	if err != nil {
		return nil, err
	}

	status := &BoardingStatus{
		TripID:         tripID,
		BoardingActive: event != nil,
	}

	if event != nil {
		status.StartedAt = &event.StartedAt
	}

	// Получить статистику
	tickets, err := s.ticketRepo.FindByTripID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	status.TotalTickets = len(tickets)

	if event != nil {
		marks, err := s.boardingRepo.FindMarksByTripID(ctx, tripID)
		if err != nil {
			return nil, err
		}
		status.BoardedCount = len(marks)
	}

	return status, nil
}

// publishTicketEvent публикует событие по билету в NATS.
func (s *ticketService) publishTicketEvent(subject string, ticket *models.Ticket) {
	data, err := json.Marshal(ticket)
	if err != nil {
		s.logger.Error("Failed to marshal ticket event", zap.Error(err))
		return
	}

	if err := s.natsConn.Publish(subject, data); err != nil {
		s.logger.Error("Failed to publish ticket event", zap.Error(err), zap.String("subject", subject))
	}
}

func (s *ticketService) publishBoardingEvent(subject string, data map[string]interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("Failed to marshal boarding event", zap.Error(err))
		return
	}

	if err := s.natsConn.Publish(subject, payload); err != nil {
		s.logger.Error("Failed to publish boarding event", zap.Error(err), zap.String("subject", subject))
	}
}

func (s *ticketService) publishAuditEvent(_ context.Context, entityType, entityID, action, userID string, oldValue, newValue interface{}) {
	event := map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
		"action":      action,
		"user_id":     userID,
		"old_value":   oldValue,
		"new_value":   newValue,
		"timestamp":   time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		s.logger.Error("Failed to marshal audit event", zap.Error(err))
		return
	}

	if err := s.natsConn.Publish("audit.log", data); err != nil {
		s.logger.Error("Failed to publish audit event", zap.Error(err))
	}
}
