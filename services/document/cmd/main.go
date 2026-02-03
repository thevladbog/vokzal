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
	"github.com/vokzal-tech/document-service/internal/config"
	"github.com/vokzal-tech/document-service/internal/handlers"
	"github.com/vokzal-tech/document-service/internal/models"
	"github.com/vokzal-tech/document-service/internal/pdf"
	"github.com/vokzal-tech/document-service/internal/repository"
	"github.com/vokzal-tech/document-service/internal/service"
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
	if cfg.Logger.Level == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	logger.Info("Starting Document Service", zap.String("version", "1.0.0"))

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

	if err := db.AutoMigrate(&models.DocumentTemplate{}, &models.GeneratedDocument{}); err != nil {
		logger.Warn("Auto-migration failed", zap.Error(err))
	}

	docRepo := repository.NewDocumentRepository(db)
	pdfGenerator := pdf.NewGenerator(logger)
	
	docService, err := service.NewDocumentService(docRepo, pdfGenerator, &cfg.MinIO, logger)
	if err != nil {
		logger.Fatal("Failed to create document service", zap.Error(err))
	}

	docHandler := handlers.NewDocumentHandler(docService, logger)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "document",
			"version": "1.0.0",
		})
	})

	v1 := router.Group("/v1")
	{
		doc := v1.Group("/document")
		{
			doc.POST("/ticket", docHandler.GenerateTicket)
			doc.POST("/pd2", docHandler.GeneratePD2)
			doc.GET("/:id", docHandler.GetDocument)
			doc.GET("/list", docHandler.ListDocuments)
		}
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Document service listening", zap.String("port", cfg.Server.Port))
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
