package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/fiscal-service/internal/service"
	"go.uber.org/zap"
)

type FiscalHandler struct {
	service service.FiscalService
	logger  *zap.Logger
}

func NewFiscalHandler(service service.FiscalService, logger *zap.Logger) *FiscalHandler {
	return &FiscalHandler{
		service: service,
		logger:  logger,
	}
}

// Получить чек по ID
func (h *FiscalHandler) GetReceipt(c *gin.Context) {
	id := c.Param("id")
	receipt, err := h.service.GetReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": receipt})
}

// Получить чеки по билету
func (h *FiscalHandler) GetReceiptsByTicket(c *gin.Context) {
	ticketID := c.Query("ticket_id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is required"})
		return
	}

	receipts, err := h.service.GetReceiptsByTicket(c.Request.Context(), ticketID)
	if err != nil {
		h.logger.Error("Failed to get receipts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get receipts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": receipts})
}

// Создать Z-отчёт
func (h *FiscalHandler) CreateZReport(c *gin.Context) {
	var req struct {
		Date string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.service.CreateDailyZReport(c.Request.Context(), req.Date)
	if err != nil {
		h.logger.Error("Failed to create Z-report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": report})
}

// Получить Z-отчёт по дате
func (h *FiscalHandler) GetZReport(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required (YYYY-MM-DD)"})
		return
	}

	report, err := h.service.GetZReport(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Z-report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": report})
}

// Список Z-отчётов
func (h *FiscalHandler) ListZReports(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "30")
	limit, _ := strconv.Atoi(limitStr)

	reports, err := h.service.ListZReports(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list Z-reports", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list Z-reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reports})
}

// Статус ККТ
func (h *FiscalHandler) GetKKTStatus(c *gin.Context) {
	status, err := h.service.GetKKTStatus(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get KKT status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}
