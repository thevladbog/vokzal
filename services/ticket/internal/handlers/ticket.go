// Package handlers содержит HTTP-обработчики API билетов.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/ticket-service/internal/service"
	"go.uber.org/zap"
)

// TicketHandler — обработчик HTTP-запросов для билетов и посадки.
type TicketHandler struct {
	service service.TicketService
	logger  *zap.Logger
}

// NewTicketHandler создаёт обработчик билетов.
func NewTicketHandler(service service.TicketService, logger *zap.Logger) *TicketHandler {
	return &TicketHandler{
		service: service,
		logger:  logger,
	}
}

// SellTicket продаёт билет.
func (h *TicketHandler) SellTicket(c *gin.Context) {
	var req service.SellTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.service.SellTicket(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to sell ticket", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": ticket})
}

// GetTicket возвращает билет по ID.
func (h *TicketHandler) GetTicket(c *gin.Context) {
	id := c.Param("id")
	ticket, err := h.service.GetTicket(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

// GetTicketByQR возвращает билет по QR-коду.
func (h *TicketHandler) GetTicketByQR(c *gin.Context) {
	qrCode := c.Query("qr_code")
	if qrCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "qr_code is required"})
		return
	}

	ticket, err := h.service.GetTicketByQR(c.Request.Context(), qrCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

// ListTicketsByTrip возвращает список билетов на рейс.
func (h *TicketHandler) ListTicketsByTrip(c *gin.Context) {
	tripID := c.Query("trip_id")
	if tripID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trip_id is required"})
		return
	}

	tickets, err := h.service.ListTicketsByTrip(c.Request.Context(), tripID)
	if err != nil {
		h.logger.Error("Failed to list tickets", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

// RefundTicket возвращает билет.
func (h *TicketHandler) RefundTicket(c *gin.Context) {
	ticketID := c.Param("id")

	// user_id из контекста (middleware) или заголовка X-User-ID (API Gateway после аутентификации)
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
	}
	if userID == "" {
		userID = "system"
	}

	result, err := h.service.RefundTicket(c.Request.Context(), ticketID, userID)
	if err != nil {
		h.logger.Error("Failed to refund ticket", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket refunded successfully",
		"data":    result,
	})
}

// StartBoarding начинает посадку.
func (h *TicketHandler) StartBoarding(c *gin.Context) {
	var req struct {
		TripID string `json:"trip_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// user_id из контекста (middleware) или заголовка X-User-ID (API Gateway после аутентификации)
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
	}
	if userID == "" {
		userID = "system"
	}

	if err := h.service.StartBoarding(c.Request.Context(), req.TripID, userID); err != nil {
		h.logger.Error("Failed to start boarding", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Boarding started successfully"})
}

// MarkBoarding отмечает посадку пассажира.
func (h *TicketHandler) MarkBoarding(c *gin.Context) {
	var req service.MarkBoardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.MarkBoarding(c.Request.Context(), &req); err != nil {
		h.logger.Error("Failed to mark boarding", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Boarding marked successfully"})
}

// GetBoardingStatus возвращает статус посадки по рейсу.
func (h *TicketHandler) GetBoardingStatus(c *gin.Context) {
	tripID := c.Query("trip_id")
	if tripID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trip_id is required"})
		return
	}

	status, err := h.service.GetBoardingStatus(c.Request.Context(), tripID)
	if err != nil {
		h.logger.Error("Failed to get boarding status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get boarding status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}
