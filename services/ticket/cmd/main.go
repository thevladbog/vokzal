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
	"github.com/vokzal-tech/ticket-service/internal/config"
	"github.com/vokzal-tech/ticket-service/internal/handlers"
	"github.com/vokzal-tech/ticket-service/internal/models"
	"github.com/vokzal-tech/ticket-service/internal/repository"
	"github.com/vokzal-tech/ticket-service/internal/service"
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

	logger.Info("Starting Ticket Service", zap.String("version", "1.0.0"))

	// Подключиться к БД с оптимизированным connection pool
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

	// Auto-migrate (опционально для dev)
	if err := db.AutoMigrate(&models.Ticket{}, &models.BoardingEvent{}, &models.BoardingMark{}); err != nil {
		logger.Warn("Auto-migration failed", zap.Error(err))
	}

	// Подключиться к NATS
	natsConn, err := nats.Connect(cfg.NATS.URL,
		nats.UserInfo(cfg.NATS.User, cfg.NATS.Password),
		nats.Name("ticket-service"))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsConn.Close()

	logger.Info("Connected to NATS", zap.String("url", cfg.NATS.URL))

	// Создать репозитории
	ticketRepo := repository.NewTicketRepository(db)
	boardingRepo := repository.NewBoardingRepository(db)

	// Создать сервис
	ticketService := service.NewTicketService(ticketRepo, boardingRepo, natsConn, cfg, logger)

	// Создать handlers
	ticketHandler := handlers.NewTicketHandler(ticketService, logger)

	// Настроить Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ticket",
			"version": "1.0.0",
		})
	})

	// API routes
	v1 := router.Group("/v1")
	{
		// Tickets
		tickets := v1.Group("/tickets")
		{
			tickets.POST("/sell", ticketHandler.SellTicket)
			tickets.GET("", ticketHandler.ListTicketsByTrip)
			tickets.GET("/:id", ticketHandler.GetTicket)
			tickets.GET("/qr", ticketHandler.GetTicketByQR)
			tickets.POST("/:id/refund", ticketHandler.RefundTicket)
		}

		// Boarding
		boarding := v1.Group("/boarding")
		{
			boarding.POST("/start", ticketHandler.StartBoarding)
			boarding.POST("/mark", ticketHandler.MarkBoarding)
			boarding.GET("/status", ticketHandler.GetBoardingStatus)
		}
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
		logger.Info("Ticket service listening", zap.String("port", cfg.Server.Port))
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
