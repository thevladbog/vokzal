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
	"github.com/nats-io/nats.go"
	"github.com/vokzal-tech/fiscal-service/internal/atol"
	"github.com/vokzal-tech/fiscal-service/internal/config"
	"github.com/vokzal-tech/fiscal-service/internal/handlers"
	"github.com/vokzal-tech/fiscal-service/internal/models"
	"github.com/vokzal-tech/fiscal-service/internal/repository"
	"github.com/vokzal-tech/fiscal-service/internal/service"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загрузить конфигурацию
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Инициализировать логгер
	var logger *zap.Logger
	if cfg.Logger.Level == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	logger.Info("Starting Fiscal Service", zap.String("version", "1.0.0"))

	// Подключиться к БД
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Настроить connection pool
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate
	if err := db.AutoMigrate(&models.FiscalReceipt{}, &models.ZReport{}); err != nil {
		logger.Warn("Auto-migration failed", zap.Error(err))
	}

	// Подключиться к NATS
	natsConn, err := nats.Connect(cfg.NATS.URL,
		nats.UserInfo(cfg.NATS.User, cfg.NATS.Password),
		nats.Name("fiscal-service"))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsConn.Close()

	logger.Info("Connected to NATS", zap.String("url", cfg.NATS.URL))

	// Создать ATOL клиент
	atolClient := atol.NewATOLClient(cfg.LocalAgent.URL, logger)

	// Создать репозиторий
	fiscalRepo := repository.NewFiscalRepository(db)

	// Создать сервис
	fiscalService := service.NewFiscalService(fiscalRepo, atolClient, cfg, logger)

	// Подписаться на NATS события
	fiscalService.SubscribeToEvents(natsConn)

	// Создать handlers
	fiscalHandler := handlers.NewFiscalHandler(fiscalService, logger)

	// Настроить Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "fiscal",
			"version": "1.0.0",
		})
	})

	// API routes
	v1 := router.Group("/v1")
	{
		// Receipts
		receipts := v1.Group("/receipts")
		{
			receipts.GET("/:id", fiscalHandler.GetReceipt)
			receipts.GET("", fiscalHandler.GetReceiptsByTicket)
		}

		// Z-Reports
		reports := v1.Group("/z-reports")
		{
			reports.POST("", fiscalHandler.CreateZReport)
			reports.GET("", fiscalHandler.ListZReports)
			reports.GET("/date", fiscalHandler.GetZReport)
		}

		// KKT
		v1.GET("/kkt/status", fiscalHandler.GetKKTStatus)
	}

	// Создать HTTP сервер
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запустить сервер в горутине
	go func() {
		logger.Info("Fiscal service listening", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Запустить ежедневный cron для Z-отчётов (упрощённый вариант)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			date := time.Now().Format("2006-01-02")
			ctx := context.Background()
			if _, err := fiscalService.CreateDailyZReport(ctx, date); err != nil {
				logger.Error("Failed to create daily Z-report", zap.Error(err), zap.String("date", date))
			}
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
