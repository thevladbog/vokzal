// Package handlers содержит HTTP-обработчики API уведомлений.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/notify-service/internal/service"
	"go.uber.org/zap"
)

// NotifyHandler — обработчик HTTP-запросов для отправки и получения уведомлений.
type NotifyHandler struct {
	service service.NotifyService
	logger  *zap.Logger
}

// NewNotifyHandler создаёт обработчик уведомлений.
func NewNotifyHandler(service service.NotifyService, logger *zap.Logger) *NotifyHandler {
	return &NotifyHandler{
		service: service,
		logger:  logger,
	}
}

// SendSMS обрабатывает запрос на отправку SMS.
func (h *NotifyHandler) SendSMS(c *gin.Context) {
	var req struct {
		Phone   string `json:"phone" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.service.SendSMS(c.Request.Context(), req.Phone, req.Message)
	if err != nil {
		h.logger.Error("Failed to send SMS", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": notification})
}

// SendEmail обрабатывает запрос на отправку email.
func (h *NotifyHandler) SendEmail(c *gin.Context) {
	var req struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.service.SendEmail(c.Request.Context(), req.To, req.Subject, req.Body)
	if err != nil {
		h.logger.Error("Failed to send email", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": notification})
}

// SendTelegram обрабатывает запрос на отправку сообщения в Telegram.
func (h *NotifyHandler) SendTelegram(c *gin.Context) {
	var req struct {
		ChatID  int64  `json:"chat_id" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.service.SendTelegram(c.Request.Context(), req.ChatID, req.Message)
	if err != nil {
		h.logger.Error("Failed to send telegram", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": notification})
}

// SendTTS обрабатывает запрос на голосовое оповещение (TTS).
func (h *NotifyHandler) SendTTS(c *gin.Context) {
	var req struct {
		Text     string `json:"text" binding:"required"`
		Language string `json:"language"`
		Priority string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Language == "" {
		req.Language = "ru"
	}
	if req.Priority == "" {
		req.Priority = "normal"
	}

	notification, err := h.service.SendTTS(c.Request.Context(), req.Text, req.Language, req.Priority)
	if err != nil {
		h.logger.Error("Failed to send TTS", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": notification})
}

// GetNotification возвращает уведомление по ID.
func (h *NotifyHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")
	notification, err := h.service.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notification})
}

// ListNotifications возвращает список уведомлений с пагинацией.
func (h *NotifyHandler) ListNotifications(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	notifications, err := h.service.ListNotifications(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}
