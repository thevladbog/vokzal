// Package main — точка входа Audit Service (логирование операций по 152-ФЗ).
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

	"github.com/vokzal-tech/audit-service/internal/config"
	"github.com/vokzal-tech/audit-service/internal/handlers"
	"github.com/vokzal-tech/audit-service/internal/models"
	"github.com/vokzal-tech/audit-service/internal/repository"
	"github.com/vokzal-tech/audit-service/internal/service"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	var logger *zap.Logger
	var errLog error
	if cfg.Logger.Level == "production" {
		logger, errLog = zap.NewProduction()
	} else {
		logger, errLog = zap.NewDevelopment()
	}
	if errLog != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", errLog))
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			fmt.Fprintf(os.Stderr, "logger sync: %v\n", syncErr)
		}
	}()

	logger.Info("Starting Audit Service", zap.String("version", "1.0.0"))

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if migErr := db.AutoMigrate(&models.AuditLog{}); migErr != nil {
		logger.Warn("Auto-migration failed", zap.Error(migErr))
	}

	natsConn, err := nats.Connect(cfg.NATS.URL,
		nats.UserInfo(cfg.NATS.User, cfg.NATS.Password),
		nats.Name("audit-service"))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsConn.Close()

	logger.Info("Connected to NATS", zap.String("url", cfg.NATS.URL))

	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo, logger)

	// Подписаться на события
	auditService.SubscribeToEvents(natsConn)

	auditHandler := handlers.NewAuditHandler(auditService, logger)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "audit",
			"version": "1.0.0",
		})
	})

	v1 := router.Group("/v1")
	audit := v1.Group("/audit")
	audit.POST("/log", auditHandler.CreateLog)
	audit.GET("/:id", auditHandler.GetLog)
	audit.GET("/entity", auditHandler.GetLogsByEntity)
	audit.GET("/user", auditHandler.GetLogsByUser)
	audit.GET("/date-range", auditHandler.GetLogsByDateRange)
	audit.GET("/list", auditHandler.ListLogs)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Audit service listening", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

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
