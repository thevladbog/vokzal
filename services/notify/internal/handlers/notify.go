// Package handlers содержит HTTP-обработчики API уведомлений.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vokzal-tech/notify-service/internal/service"
)

// NotifyHandler — обработчик HTTP-запросов для отправки и получения уведомлений.
type NotifyHandler struct {
	svc    service.NotifyService
	logger *zap.Logger
}

// NewNotifyHandler создаёт обработчик уведомлений.
func NewNotifyHandler(svc service.NotifyService, logger *zap.Logger) *NotifyHandler {
	return &NotifyHandler{
		svc:    svc,
		logger: logger,
	}
}

// sendCreated вызывает send(), при ошибке логирует и отвечает 500, иначе 201 с data.
func (h *NotifyHandler) sendCreated(c *gin.Context, send func() (interface{}, error), errMsg string) {
	notification, err := send()
	if err != nil {
		h.logger.Error(errMsg, zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": notification})
}

// SendSMS обрабатывает запрос на отправку SMS.
//
//nolint:dupl // дублирует шаблон SendTelegram (bind + sendCreated), разный тип запроса.
func (h *NotifyHandler) SendSMS(c *gin.Context) {
	var req struct {
		Phone   string `json:"phone" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.sendCreated(c, func() (interface{}, error) {
		return h.svc.SendSMS(c.Request.Context(), req.Phone, req.Message)
	}, "Failed to send SMS")
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

	notification, err := h.svc.SendEmail(c.Request.Context(), req.To, req.Subject, req.Body)
	if err != nil {
		h.logger.Error("Failed to send email", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": notification})
}

// SendTelegram обрабатывает запрос на отправку сообщения в Telegram.
//
//nolint:dupl // дублирует шаблон SendSMS (bind + sendCreated), разный тип запроса.
func (h *NotifyHandler) SendTelegram(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
		ChatID  int64  `json:"chat_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.sendCreated(c, func() (interface{}, error) {
		return h.svc.SendTelegram(c.Request.Context(), req.ChatID, req.Message)
	}, "Failed to send telegram")
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

	notification, err := h.svc.SendTTS(c.Request.Context(), req.Text, req.Language, req.Priority)
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
	notification, err := h.svc.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notification})
}

// ListNotifications возвращает список уведомлений с пагинацией.
func (h *NotifyHandler) ListNotifications(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	notifications, err := h.svc.ListNotifications(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}
