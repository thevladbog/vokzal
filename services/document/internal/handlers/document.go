package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/document-service/internal/pdf"
	"github.com/vokzal-tech/document-service/internal/service"
	"go.uber.org/zap"
)

type DocumentHandler struct {
	service service.DocumentService
	logger  *zap.Logger
}

func NewDocumentHandler(service service.DocumentService, logger *zap.Logger) *DocumentHandler {
	return &DocumentHandler{
		service: service,
		logger:  logger,
	}
}

func (h *DocumentHandler) GenerateTicket(c *gin.Context) {
	var data pdf.TicketData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.service.GenerateTicket(c.Request.Context(), &data)
	if err != nil {
		h.logger.Error("Failed to generate ticket", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ticket"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": doc})
}

func (h *DocumentHandler) GeneratePD2(c *gin.Context) {
	var data pdf.PD2Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.service.GeneratePD2(c.Request.Context(), &data)
	if err != nil {
		h.logger.Error("Failed to generate PD-2", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PD-2"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": doc})
}

func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	doc, err := h.service.GetDocument(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": doc})
}

func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	docs, err := h.service.ListDocuments(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list documents", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": docs})
}
