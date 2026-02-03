// Package handlers содержит HTTP-обработчики API платежей.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vokzal-tech/payment-service/internal/service"
)

// PaymentHandler — обработчик HTTP-запросов для платежей.
type PaymentHandler struct {
	svc    service.PaymentService
	logger *zap.Logger
}

// NewPaymentHandler создаёт обработчик платежей.
func NewPaymentHandler(svc service.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		svc:    svc,
		logger: logger,
	}
}

// initPayment вызывает bind InitPaymentRequest, затем do(), при ошибке — 500, иначе 201.
func (h *PaymentHandler) initPayment(c *gin.Context, do func(*service.InitPaymentRequest) (interface{}, error), errMsg string) {
	var req service.InitPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payment, err := do(&req)
	if err != nil {
		h.logger.Error(errMsg, zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": payment})
}

// InitTinkoff инициализирует платёж через Tinkoff.
func (h *PaymentHandler) InitTinkoff(c *gin.Context) {
	h.initPayment(c, func(req *service.InitPaymentRequest) (interface{}, error) {
		return h.svc.InitTinkoffPayment(c.Request.Context(), req)
	}, "Failed to init Tinkoff payment")
}

// InitSBP инициализирует платёж через СБП.
func (h *PaymentHandler) InitSBP(c *gin.Context) {
	h.initPayment(c, func(req *service.InitPaymentRequest) (interface{}, error) {
		return h.svc.InitSBPPayment(c.Request.Context(), req)
	}, "Failed to init SBP payment")
}

// InitCash создаёт запись о наличной оплате.
func (h *PaymentHandler) InitCash(c *gin.Context) {
	h.initPayment(c, func(req *service.InitPaymentRequest) (interface{}, error) {
		return h.svc.InitCashPayment(c.Request.Context(), req)
	}, "Failed to create cash payment")
}

// GetPayment возвращает платёж по ID.
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.svc.GetPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

// CheckStatus проверяет статус платежа у провайдера.
func (h *PaymentHandler) CheckStatus(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.svc.CheckPaymentStatus(c.Request.Context(), id)
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

	payments, err := h.svc.GetPaymentByTicket(c.Request.Context(), ticketID)
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

	if err := h.svc.HandleTinkoffWebhook(c.Request.Context(), data); err != nil {
		h.logger.Error("Failed to handle Tinkoff webhook", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ListPayments возвращает список платежей.
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	payments, err := h.svc.ListPayments(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list payments", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}
