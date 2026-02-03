package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/local-agent/internal/kkt"
	"github.com/vokzal-tech/local-agent/internal/printer"
	"github.com/vokzal-tech/local-agent/internal/scanner"
	"github.com/vokzal-tech/local-agent/internal/tts"
	"go.uber.org/zap"
)

type AgentHandler struct {
	kkt     *kkt.ATOLClient
	printer *printer.PrinterClient
	scanner *scanner.ScannerClient
	tts     *tts.TTSClient
	logger  *zap.Logger
}

func NewAgentHandler(
	kkt *kkt.ATOLClient,
	printer *printer.PrinterClient,
	scanner *scanner.ScannerClient,
	tts *tts.TTSClient,
	logger *zap.Logger,
) *AgentHandler {
	return &AgentHandler{
		kkt:     kkt,
		printer: printer,
		scanner: scanner,
		tts:     tts,
		logger:  logger,
	}
}

// Health check
func (h *AgentHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "local-agent",
	})
}

// ============= KKT Endpoints =============

// PrintReceipt печатает чек на ККТ
func (h *AgentHandler) PrintReceipt(c *gin.Context) {
	var req kkt.ReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid receipt request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.kkt.PrintReceipt(&req)
	if err != nil {
		h.logger.Error("Failed to print receipt", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateZReport создаёт Z-отчёт
func (h *AgentHandler) CreateZReport(c *gin.Context) {
	resp, err := h.kkt.CreateZReport()
	if err != nil {
		h.logger.Error("Failed to create Z-report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetKKTStatus возвращает статус ККТ
func (h *AgentHandler) GetKKTStatus(c *gin.Context) {
	status, err := h.kkt.GetStatus()
	if err != nil {
		h.logger.Error("Failed to get KKT status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// ============= Printer Endpoints =============

// PrintTicket печатает билет
func (h *AgentHandler) PrintTicket(c *gin.Context) {
	var req printer.TicketData
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid ticket data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.printer.PrintTicket(&req); err != nil {
		h.logger.Error("Failed to print ticket", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetPrinterStatus возвращает статус принтера
func (h *AgentHandler) GetPrinterStatus(c *gin.Context) {
	status, err := h.printer.GetStatus()
	if err != nil {
		h.logger.Error("Failed to get printer status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// ============= Scanner Endpoints =============

// ScanBarcode запускает чтение штрихкода
func (h *AgentHandler) ScanBarcode(c *gin.Context) {
	barcode, err := h.scanner.ReadBarcode()
	if err != nil {
		h.logger.Error("Failed to scan barcode", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"barcode": barcode})
}

// GetScannerStatus возвращает статус сканера
func (h *AgentHandler) GetScannerStatus(c *gin.Context) {
	status := h.scanner.GetStatus()
	c.JSON(http.StatusOK, status)
}

// ============= TTS Endpoints =============

// Announce добавляет голосовое оповещение
func (h *AgentHandler) Announce(c *gin.Context) {
	var req tts.Announcement
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid announcement request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Defaults
	if req.Language == "" {
		req.Language = "ru"
	}
	if req.Priority == "" {
		req.Priority = "normal"
	}

	if err := h.tts.Announce(req.Text, req.Language, req.Priority); err != nil {
		h.logger.Error("Failed to announce", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetTTSStatus возвращает статус TTS
func (h *AgentHandler) GetTTSStatus(c *gin.Context) {
	status := h.tts.GetStatus()
	c.JSON(http.StatusOK, status)
}

// ============= Common Status Endpoint =============

// GetAgentStatus возвращает статус всех компонентов агента
func (h *AgentHandler) GetAgentStatus(c *gin.Context) {
	kktStatus, _ := h.kkt.GetStatus()
	printerStatus, _ := h.printer.GetStatus()
	scannerStatus := h.scanner.GetStatus()
	ttsStatus := h.tts.GetStatus()

	c.JSON(http.StatusOK, gin.H{
		"agent":   "local-agent",
		"version": "1.0.0",
		"kkt":     kktStatus,
		"printer": printerStatus,
		"scanner": scannerStatus,
		"tts":     ttsStatus,
	})
}
