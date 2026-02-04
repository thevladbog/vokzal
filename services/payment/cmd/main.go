// Package main — точка входа Payment Service.
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
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/vokzal-tech/payment-service/internal/config"
	"github.com/vokzal-tech/payment-service/internal/handlers"
	"github.com/vokzal-tech/payment-service/internal/models"
	"github.com/vokzal-tech/payment-service/internal/repository"
	"github.com/vokzal-tech/payment-service/internal/sbp"
	"github.com/vokzal-tech/payment-service/internal/service"
	"github.com/vokzal-tech/payment-service/internal/tinkoff"
)

func initLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.Logger.Level == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	logger, errLog := initLogger(cfg)
	if errLog != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", errLog))
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			fmt.Fprintf(os.Stderr, "logger sync: %v\n", syncErr)
		}
	}()

	logger.Info("Starting Payment Service", zap.String("version", "1.0.0"))

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

	if migErr := db.AutoMigrate(&models.Payment{}); migErr != nil {
		logger.Warn("Auto-migration failed", zap.Error(migErr))
	}

	// Подключиться к NATS
	natsConn, err := nats.Connect(cfg.NATS.URL,
		nats.UserInfo(cfg.NATS.User, cfg.NATS.Password),
		nats.Name("payment-service"))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsConn.Close()

	logger.Info("Connected to NATS", zap.String("url", cfg.NATS.URL))

	// Создать клиенты для провайдеров
	tinkoffClient := tinkoff.NewTinkoffClient(
		cfg.Tinkoff.TerminalKey,
		cfg.Tinkoff.Password,
		cfg.Tinkoff.APIUrl,
		logger,
	)

	sbpClient := sbp.NewSBPClient(
		cfg.SBP.MerchantID,
		cfg.SBP.APIUrl,
		cfg.SBP.APIKey,
		logger,
	)

	// Создать репозиторий
	paymentRepo := repository.NewPaymentRepository(db)

	// Создать сервис
	paymentService := service.NewPaymentService(
		paymentRepo,
		tinkoffClient,
		sbpClient,
		natsConn,
		cfg,
		logger,
	)

	// Создать handlers
	paymentHandler := handlers.NewPaymentHandler(paymentService, logger)

	// Настроить Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "payment",
			"version": "1.0.0",
		})
	})

	// Traefik strips /v1/payment prefix, service receives /, /:id, /tinkoff/init, etc.
	router.POST("/tinkoff/init", paymentHandler.InitTinkoff)
	router.POST("/sbp/init", paymentHandler.InitSBP)
	router.POST("/cash/init", paymentHandler.InitCash)
	router.GET("/list", paymentHandler.ListPayments)
	router.GET("", paymentHandler.GetPaymentsByTicket)
	router.GET("/:id", paymentHandler.GetPayment)
	router.GET("/:id/status", paymentHandler.CheckStatus)
	router.POST("/webhooks/tinkoff", paymentHandler.TinkoffWebhook)

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
		logger.Info("Payment service listening", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
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
