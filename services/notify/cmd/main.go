// Package main — точка входа Notify Service (SMS, Email, Telegram, TTS).
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

	"github.com/vokzal-tech/notify-service/internal/config"
	"github.com/vokzal-tech/notify-service/internal/email"
	"github.com/vokzal-tech/notify-service/internal/handlers"
	"github.com/vokzal-tech/notify-service/internal/models"
	"github.com/vokzal-tech/notify-service/internal/repository"
	"github.com/vokzal-tech/notify-service/internal/service"
	"github.com/vokzal-tech/notify-service/internal/sms"
	"github.com/vokzal-tech/notify-service/internal/telegram"
	"github.com/vokzal-tech/notify-service/internal/tts"
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

	logger.Info("Starting Notify Service", zap.String("version", "1.0.0"))

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

	if migErr := db.AutoMigrate(&models.Notification{}); migErr != nil {
		logger.Warn("Auto-migration failed", zap.Error(migErr))
	}

	natsConn, err := nats.Connect(cfg.NATS.URL,
		nats.UserInfo(cfg.NATS.User, cfg.NATS.Password),
		nats.Name("notify-service"))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsConn.Close()

	smsClient := sms.NewSMSRuClient(cfg.SMS.APIID, cfg.SMS.URL, logger)
	emailClient := email.NewEmailClient(cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.Username, cfg.Email.Password, cfg.Email.From, logger)

	var telegramClient *telegram.TelegramClient
	if cfg.Telegram.BotToken != "" {
		telegramClient, err = telegram.NewTelegramClient(cfg.Telegram.BotToken, logger)
		if err != nil {
			logger.Warn("Failed to create Telegram client", zap.Error(err))
		}
	}

	ttsClient := tts.NewTTSClient(cfg.LocalAgent.URL, logger)

	notifyRepo := repository.NewNotificationRepository(db)
	notifyService := service.NewNotifyService(notifyRepo, smsClient, emailClient, telegramClient, ttsClient, logger)
	notifyHandler := handlers.NewNotifyHandler(notifyService, logger)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "notify",
			"version": "1.0.0",
		})
	})

	v1 := router.Group("/v1")
	notify := v1.Group("/notify")
	notify.POST("/sms", notifyHandler.SendSMS)
	notify.POST("/email", notifyHandler.SendEmail)
	notify.POST("/telegram", notifyHandler.SendTelegram)
	notify.POST("/tts", notifyHandler.SendTTS)
	notify.GET("/:id", notifyHandler.GetNotification)
	notify.GET("/list", notifyHandler.ListNotifications)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Notify service listening", zap.String("port", cfg.Server.Port))
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
