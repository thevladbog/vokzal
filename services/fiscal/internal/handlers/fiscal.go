// Package handlers — HTTP-обработчики Fiscal Service.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vokzal-tech/fiscal-service/internal/service"
)

// FiscalHandler обрабатывает HTTP-запросы к API фискализации.
type FiscalHandler struct {
	svc    service.FiscalService
	logger *zap.Logger
}

// NewFiscalHandler создаёт новый FiscalHandler.
func NewFiscalHandler(svc service.FiscalService, logger *zap.Logger) *FiscalHandler {
	return &FiscalHandler{
		svc:    svc,
		logger: logger,
	}
}

// GetReceipt возвращает чек по ID.
func (h *FiscalHandler) GetReceipt(c *gin.Context) {
	id := c.Param("id")
	receipt, err := h.svc.GetReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": receipt})
}

// GetReceiptsByTicket возвращает чеки по билету.
func (h *FiscalHandler) GetReceiptsByTicket(c *gin.Context) {
	ticketID := c.Query("ticket_id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is required"})
		return
	}

	receipts, err := h.svc.GetReceiptsByTicket(c.Request.Context(), ticketID)
	if err != nil {
		h.logger.Error("Failed to get receipts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get receipts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": receipts})
}

// CreateZReport создаёт Z-отчёт за дату.
func (h *FiscalHandler) CreateZReport(c *gin.Context) {
	var req struct {
		Date string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.svc.CreateDailyZReport(c.Request.Context(), req.Date)
	if err != nil {
		h.logger.Error("Failed to create Z-report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": report})
}

// GetZReport возвращает Z-отчёт по дате.
func (h *FiscalHandler) GetZReport(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required (YYYY-MM-DD)"})
		return
	}

	report, err := h.svc.GetZReport(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Z-report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": report})
}

// ListZReports возвращает список Z-отчётов.
func (h *FiscalHandler) ListZReports(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "30")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 30
	}

	reports, err := h.svc.ListZReports(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list Z-reports", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list Z-reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reports})
}

// GetKKTStatus возвращает статус ККТ.
func (h *FiscalHandler) GetKKTStatus(c *gin.Context) {
	status, err := h.svc.GetKKTStatus(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get KKT status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}
