// Package handlers — HTTP-обработчики Document Service.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vokzal-tech/document-service/internal/pdf"
	"github.com/vokzal-tech/document-service/internal/service"
)

// DocumentHandler обрабатывает HTTP-запросы к API документов.
type DocumentHandler struct {
	svc    service.DocumentService
	logger *zap.Logger
}

// NewDocumentHandler создаёт новый DocumentHandler.
func NewDocumentHandler(svc service.DocumentService, logger *zap.Logger) *DocumentHandler {
	return &DocumentHandler{
		svc:    svc,
		logger: logger,
	}
}

// generateDoc — общий шаг: bind JSON, вызов generate, ответ 201 или ошибка.
func (h *DocumentHandler) generateDoc(c *gin.Context, bindTarget interface{}, generate func() (interface{}, error), errMsg string) {
	if err := c.ShouldBindJSON(bindTarget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doc, err := generate()
	if err != nil {
		h.logger.Error(errMsg, zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": doc})
}

// GenerateTicket генерирует PDF-билет.
func (h *DocumentHandler) GenerateTicket(c *gin.Context) {
	var data pdf.TicketData
	h.generateDoc(c, &data, func() (interface{}, error) {
		return h.svc.GenerateTicket(c.Request.Context(), &data)
	}, "Failed to generate ticket")
}

// GeneratePD2 генерирует проездной документ ПД-2.
func (h *DocumentHandler) GeneratePD2(c *gin.Context) {
	var data pdf.PD2Data
	h.generateDoc(c, &data, func() (interface{}, error) {
		return h.svc.GeneratePD2(c.Request.Context(), &data)
	}, "Failed to generate PD-2")
}

// GetDocument возвращает документ по ID.
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	doc, err := h.svc.GetDocument(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": doc})
}

// ListDocuments возвращает список документов.
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	docs, err := h.svc.ListDocuments(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list documents", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": docs})
}
