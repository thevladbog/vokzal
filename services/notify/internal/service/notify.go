// Package service содержит бизнес-логику отправки и хранения уведомлений.
package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/vokzal-tech/notify-service/internal/email"
	"github.com/vokzal-tech/notify-service/internal/models"
	"github.com/vokzal-tech/notify-service/internal/repository"
	"github.com/vokzal-tech/notify-service/internal/sms"
	"github.com/vokzal-tech/notify-service/internal/telegram"
	"github.com/vokzal-tech/notify-service/internal/tts"
	"go.uber.org/zap"
)

const (
	statusFailed = "failed"
	statusSent   = "sent"
)

// NotifyService — интерфейс сервиса уведомлений (SMS, email, Telegram, TTS).
type NotifyService interface {
	SendSMS(ctx context.Context, phone, message string) (*models.Notification, error)
	SendEmail(ctx context.Context, to, subject, body string) (*models.Notification, error)
	SendTelegram(ctx context.Context, chatID int64, message string) (*models.Notification, error)
	SendTTS(ctx context.Context, text, language, priority string) (*models.Notification, error)
	GetNotification(ctx context.Context, id string) (*models.Notification, error)
	ListNotifications(ctx context.Context, limit int) ([]*models.Notification, error)
}

type notifyService struct {
	repo           repository.NotificationRepository
	smsClient      *sms.SMSRuClient
	emailClient    *email.EmailClient
	telegramClient *telegram.TelegramClient
	ttsClient      *tts.TTSClient
	logger         *zap.Logger
}

// NewNotifyService создаёт сервис уведомлений.
func NewNotifyService(
	repo repository.NotificationRepository,
	smsClient *sms.SMSRuClient,
	emailClient *email.EmailClient,
	telegramClient *telegram.TelegramClient,
	ttsClient *tts.TTSClient,
	logger *zap.Logger,
) NotifyService {
	return &notifyService{
		repo:           repo,
		smsClient:      smsClient,
		emailClient:    emailClient,
		telegramClient: telegramClient,
		ttsClient:      ttsClient,
		logger:         logger,
	}
}

func (s *notifyService) SendSMS(ctx context.Context, phone, message string) (*models.Notification, error) {
	notification := &models.Notification{
		Type:      "sms",
		Recipient: phone,
		Message:   message,
		Status:    "pending",
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Отправить SMS
	if err := s.smsClient.Send(phone, message); err != nil {
		notification.Status = statusFailed
		errMsg := err.Error()
		notification.ErrorMsg = &errMsg
		_ = s.repo.Update(ctx, notification)
		return nil, err
	}

	now := time.Now()
	notification.Status = statusSent
	notification.SentAt = &now

	if err := s.repo.Update(ctx, notification); err != nil {
		s.logger.Error("Failed to update notification", zap.Error(err))
	}

	return notification, nil
}

func (s *notifyService) SendEmail(ctx context.Context, to, subject, body string) (*models.Notification, error) {
	notification := &models.Notification{
		Type:      "email",
		Recipient: to,
		Subject:   &subject,
		Message:   body,
		Status:    "pending",
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Отправить Email
	if err := s.emailClient.Send(to, subject, body); err != nil {
		notification.Status = statusFailed
		errMsg := err.Error()
		notification.ErrorMsg = &errMsg
		_ = s.repo.Update(ctx, notification)
		return nil, err
	}

	now := time.Now()
	notification.Status = statusSent
	notification.SentAt = &now

	if err := s.repo.Update(ctx, notification); err != nil {
		s.logger.Error("Failed to update notification", zap.Error(err))
	}

	return notification, nil
}

func (s *notifyService) SendTelegram(ctx context.Context, chatID int64, message string) (*models.Notification, error) {
	notification := &models.Notification{
		Type:      "telegram",
		Recipient: strconv.FormatInt(chatID, 10),
		Message:   message,
		Status:    "pending",
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Отправить Telegram
	if err := s.telegramClient.Send(chatID, message); err != nil {
		notification.Status = statusFailed
		errMsg := err.Error()
		notification.ErrorMsg = &errMsg
		_ = s.repo.Update(ctx, notification)
		return nil, err
	}

	now := time.Now()
	notification.Status = statusSent
	notification.SentAt = &now

	if err := s.repo.Update(ctx, notification); err != nil {
		s.logger.Error("Failed to update notification", zap.Error(err))
	}

	return notification, nil
}

func (s *notifyService) SendTTS(ctx context.Context, text, language, priority string) (*models.Notification, error) {
	notification := &models.Notification{
		Type:      "tts",
		Recipient: "local-speaker",
		Message:   text,
		Status:    "pending",
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Отправить TTS
	if err := s.ttsClient.Announce(text, language, priority); err != nil {
		notification.Status = statusFailed
		errMsg := err.Error()
		notification.ErrorMsg = &errMsg
		_ = s.repo.Update(ctx, notification)
		return nil, err
	}

	now := time.Now()
	notification.Status = statusSent
	notification.SentAt = &now

	if err := s.repo.Update(ctx, notification); err != nil {
		s.logger.Error("Failed to update notification", zap.Error(err))
	}

	return notification, nil
}

func (s *notifyService) GetNotification(ctx context.Context, id string) (*models.Notification, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *notifyService) ListNotifications(ctx context.Context, limit int) ([]*models.Notification, error) {
	return s.repo.List(ctx, limit)
}
