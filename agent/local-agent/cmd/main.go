package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vokzal-tech/local-agent/internal/config"
	"github.com/vokzal-tech/local-agent/internal/handlers"
	"github.com/vokzal-tech/local-agent/internal/kkt"
	"github.com/vokzal-tech/local-agent/internal/printer"
	"github.com/vokzal-tech/local-agent/internal/scanner"
	"github.com/vokzal-tech/local-agent/internal/tts"
)

func main() {
	// Загрузить конфигурацию
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Инициализировать logger
	logger, err := initLogger(cfg.Logger.Level)
	if err != nil {
		panic(fmt.Sprintf("Failed to init logger: %v", err))
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "logger sync: %v\n", err)
		}
	}()

	logger.Info("Starting Vokzal.TECH Local Agent",
		zap.String("version", "1.0.0"),
		zap.String("port", cfg.Server.Port),
	)

	// Инициализировать клиенты
	kktClient := kkt.NewATOLClient(
		cfg.KKT.ATOLDriver,
		cfg.KKT.INN,
		cfg.KKT.OFDUrl,
		cfg.KKT.Enabled,
		logger,
	)

	printerClient := printer.NewPrinterClient(
		cfg.Printer.DevicePath,
		cfg.Printer.Type,
		cfg.Printer.Enabled,
		logger,
	)

	scannerClient := scanner.NewScannerClient(
		cfg.Scanner.DevicePath,
		cfg.Scanner.Mode,
		cfg.Scanner.Enabled,
		logger,
	)

	ttsClient := tts.NewTTSClient(
		cfg.TTS.Engine,
		cfg.TTS.Voice,
		cfg.TTS.Volume,
		cfg.TTS.Enabled,
		logger,
	)

	// Инициализировать handlers
	handler := handlers.NewAgentHandler(kktClient, printerClient, scannerClient, ttsClient, logger)

	// Настроить Gin
	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/status", handler.GetAgentStatus)

	kktGroup := router.Group("/kkt")
	kktGroup.POST("/receipt", handler.PrintReceipt)
	kktGroup.POST("/z-report", handler.CreateZReport)
	kktGroup.GET("/status", handler.GetKKTStatus)
	printerGroup := router.Group("/printer")
	printerGroup.POST("/ticket", handler.PrintTicket)
	printerGroup.GET("/status", handler.GetPrinterStatus)
	scannerGroup := router.Group("/scanner")
	scannerGroup.POST("/scan", handler.ScanBarcode)
	scannerGroup.GET("/status", handler.GetScannerStatus)
	ttsGroup := router.Group("/tts")
	ttsGroup.POST("/announce", handler.Announce)
	ttsGroup.GET("/status", handler.GetTTSStatus)

	// Запустить сервер
	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("Local Agent started", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Local Agent...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Local Agent stopped")
}

func initLogger(level string) (*zap.Logger, error) {
	var zapLevel zap.AtomicLevel
	switch level {
	case "debug":
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	zapCfg := zap.Config{
		Level:            zapLevel,
		Encoding:         "json",
		OutputPaths:      []string{"stdout", "/var/log/vokzal/agent.log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}

	return zapCfg.Build()
}
