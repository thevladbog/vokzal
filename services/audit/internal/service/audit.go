package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/vokzal-tech/audit-service/internal/models"
	"github.com/vokzal-tech/audit-service/internal/repository"
	"go.uber.org/zap"
)

type AuditService interface {
	CreateLog(ctx context.Context, req *CreateLogRequest) (*models.AuditLog, error)
	GetLog(ctx context.Context, id string) (*models.AuditLog, error)
	GetLogsByEntity(ctx context.Context, entityType, entityID string) ([]*models.AuditLog, error)
	GetLogsByUser(ctx context.Context, userID string, limit int) ([]*models.AuditLog, error)
	GetLogsByDateRange(ctx context.Context, from, to string) ([]*models.AuditLog, error)
	ListLogs(ctx context.Context, limit int) ([]*models.AuditLog, error)
	SubscribeToEvents(nc *nats.Conn)
}

type auditService struct {
	repo   repository.AuditRepository
	logger *zap.Logger
}

type CreateLogRequest struct {
	EntityType string      `json:"entity_type" binding:"required"`
	EntityID   string      `json:"entity_id" binding:"required"`
	Action     string      `json:"action" binding:"required"`
	UserID     *string     `json:"user_id"`
	OldValue   interface{} `json:"old_value"`
	NewValue   interface{} `json:"new_value"`
	IPAddress  *string     `json:"ip_address"`
	UserAgent  *string     `json:"user_agent"`
}

func NewAuditService(repo repository.AuditRepository, logger *zap.Logger) AuditService {
	return &auditService{
		repo:   repo,
		logger: logger,
	}
}

func (s *auditService) CreateLog(ctx context.Context, req *CreateLogRequest) (*models.AuditLog, error) {
	log := &models.AuditLog{
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		Action:     req.Action,
		UserID:     req.UserID,
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
	}

	if req.OldValue != nil {
		if err := log.SetOldValue(req.OldValue); err != nil {
			return nil, fmt.Errorf("failed to set old_value: %w", err)
		}
	}

	if req.NewValue != nil {
		if err := log.SetNewValue(req.NewValue); err != nil {
			return nil, fmt.Errorf("failed to set new_value: %w", err)
		}
	}

	if err := s.repo.Create(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create audit log: %w", err)
	}

	s.logger.Info("Audit log created",
		zap.String("entity_type", log.EntityType),
		zap.String("entity_id", log.EntityID),
		zap.String("action", log.Action))

	return log, nil
}

func (s *auditService) GetLog(ctx context.Context, id string) (*models.AuditLog, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *auditService) GetLogsByEntity(ctx context.Context, entityType, entityID string) ([]*models.AuditLog, error) {
	return s.repo.FindByEntity(ctx, entityType, entityID)
}

func (s *auditService) GetLogsByUser(ctx context.Context, userID string, limit int) ([]*models.AuditLog, error) {
	return s.repo.FindByUser(ctx, userID, limit)
}

func (s *auditService) GetLogsByDateRange(ctx context.Context, from, to string) ([]*models.AuditLog, error) {
	return s.repo.FindByDateRange(ctx, from, to)
}

func (s *auditService) ListLogs(ctx context.Context, limit int) ([]*models.AuditLog, error) {
	return s.repo.List(ctx, limit)
}

// SubscribeToEvents подписывается на NATS события для автоматического логирования
func (s *auditService) SubscribeToEvents(nc *nats.Conn) {
	nc.Subscribe("audit.log", func(msg *nats.Msg) {
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			s.logger.Error("Failed to unmarshal audit.log event", zap.Error(err))
			return
		}

		ctx := context.Background()
		req := &CreateLogRequest{
			EntityType: data["entity_type"].(string),
			EntityID:   data["entity_id"].(string),
			Action:     data["action"].(string),
			OldValue:   data["old_value"],
			NewValue:   data["new_value"],
		}

		if userID, ok := data["user_id"].(string); ok {
			req.UserID = &userID
		}

		if _, err := s.CreateLog(ctx, req); err != nil {
			s.logger.Error("Failed to create audit log from NATS", zap.Error(err))
		}
	})

	s.logger.Info("Subscribed to NATS events: audit.log")
}
