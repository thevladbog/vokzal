// Package handlers — HTTP-обработчики Audit Service.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/audit-service/internal/service"
	"go.uber.org/zap"
)

// AuditHandler обрабатывает HTTP-запросы к API аудита.
type AuditHandler struct {
	service service.AuditService
	logger  *zap.Logger
}

// NewAuditHandler создаёт новый AuditHandler.
func NewAuditHandler(service service.AuditService, logger *zap.Logger) *AuditHandler {
	return &AuditHandler{
		service: service,
		logger:  logger,
	}
}

// CreateLog создаёт запись в журнале аудита.
func (h *AuditHandler) CreateLog(c *gin.Context) {
	var req service.CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получить IP и User-Agent
	ip := c.ClientIP()
	req.IPAddress = &ip
	ua := c.GetHeader("User-Agent")
	if ua != "" {
		req.UserAgent = &ua
	}

	log, err := h.service.CreateLog(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create audit log", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create log"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": log})
}

// GetLog возвращает запись аудита по ID.
func (h *AuditHandler) GetLog(c *gin.Context) {
	id := c.Param("id")
	log, err := h.service.GetLog(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": log})
}

// GetLogsByEntity возвращает записи аудита по типу и ID сущности.
func (h *AuditHandler) GetLogsByEntity(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")

	if entityType == "" || entityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_type and entity_id are required"})
		return
	}

	logs, err := h.service.GetLogsByEntity(c.Request.Context(), entityType, entityID)
	if err != nil {
		h.logger.Error("Failed to get logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// GetLogsByUser возвращает записи аудита по ID пользователя.
func (h *AuditHandler) GetLogsByUser(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.service.GetLogsByUser(c.Request.Context(), userID, limit)
	if err != nil {
		h.logger.Error("Failed to get user logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// GetLogsByDateRange возвращает записи аудита за период дат.
func (h *AuditHandler) GetLogsByDateRange(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to dates are required (YYYY-MM-DD)"})
		return
	}

	logs, err := h.service.GetLogsByDateRange(c.Request.Context(), from, to)
	if err != nil {
		h.logger.Error("Failed to get logs by date range", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// ListLogs возвращает список записей аудита с пагинацией.
func (h *AuditHandler) ListLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.service.ListLogs(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}
