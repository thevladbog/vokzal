// Package handlers содержит HTTP-обработчики API платежей.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/payment-service/internal/service"
	"go.uber.org/zap"
)

// PaymentHandler — обработчик HTTP-запросов для платежей.
type PaymentHandler struct {
	service service.PaymentService
	logger  *zap.Logger
}

// NewPaymentHandler создаёт обработчик платежей.
func NewPaymentHandler(service service.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		logger:  logger,
	}
}

// InitTinkoff инициализирует платёж через Tinkoff.
func (h *PaymentHandler) InitTinkoff(c *gin.Context) {
	var req service.InitPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.InitTinkoffPayment(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to init Tinkoff payment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": payment})
}

// InitSBP инициализирует платёж через СБП.
func (h *PaymentHandler) InitSBP(c *gin.Context) {
	var req service.InitPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.InitSBPPayment(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to init SBP payment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": payment})
}

// InitCash создаёт запись о наличной оплате.
func (h *PaymentHandler) InitCash(c *gin.Context) {
	var req service.InitPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.InitCashPayment(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create cash payment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": payment})
}

// GetPayment возвращает платёж по ID.
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.service.GetPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

// CheckStatus проверяет статус платежа у провайдера.
func (h *PaymentHandler) CheckStatus(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.service.CheckPaymentStatus(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to check payment status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

// GetPaymentsByTicket возвращает платежи по билету.
func (h *PaymentHandler) GetPaymentsByTicket(c *gin.Context) {
	ticketID := c.Query("ticket_id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is required"})
		return
	}

	payments, err := h.service.GetPaymentByTicket(c.Request.Context(), ticketID)
	if err != nil {
		h.logger.Error("Failed to get payments", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}

// TinkoffWebhook обрабатывает webhook от Tinkoff.
func (h *PaymentHandler) TinkoffWebhook(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.HandleTinkoffWebhook(c.Request.Context(), data); err != nil {
		h.logger.Error("Failed to handle Tinkoff webhook", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ListPayments возвращает список платежей.
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	payments, err := h.service.ListPayments(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list payments", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}
