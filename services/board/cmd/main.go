// Package main — точка входа Board Service (табло отправлений, WebSocket).
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/vokzal-tech/board-service/internal/cache"
	"github.com/vokzal-tech/board-service/internal/config"
	"github.com/vokzal-tech/board-service/internal/handlers"
	"github.com/vokzal-tech/board-service/internal/websocket"
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
	defer func() { _ = logger.Sync() }()

	logger.Info("Starting Board Service", zap.String("version", "1.0.0"))

	// Подключиться к БД
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

	// Подключиться к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Info("Connected to Redis", zap.String("addr", cfg.Redis.Address()))

	redisCache := cache.NewRedisCache(redisClient)

	// Подключиться к NATS
	natsConn, err := nats.Connect(cfg.NATS.URL,
		nats.UserInfo(cfg.NATS.User, cfg.NATS.Password),
		nats.Name("board-service"))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsConn.Close()

	logger.Info("Connected to NATS", zap.String("url", cfg.NATS.URL))

	// Создать WebSocket hub
	hub := websocket.NewHub(logger)
	go hub.Run()

	// Подписаться на NATS события
	_, err = natsConn.Subscribe("trip.created", func(msg *nats.Msg) {
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			logger.Error("Failed to unmarshal trip.created", zap.Error(err))
			return
		}

		// Инвалидировать кэш
		if date, ok := data["date"].(string); ok {
			_ = redisCache.InvalidateTrips(ctx, date)
		}

		// Отправить обновление через WebSocket
		hub.Broadcast(&websocket.Message{
			Type:   "trip_created",
			TripID: data["id"].(string),
			Data:   data,
		})
	})
	if err != nil {
		logger.Error("Failed to subscribe to trip.created", zap.Error(err))
	}

	_, err = natsConn.Subscribe("trip.status_changed", func(msg *nats.Msg) {
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			logger.Error("Failed to unmarshal trip.status_changed", zap.Error(err))
			return
		}

		// Инвалидировать кэш
		if date, ok := data["date"].(string); ok {
			_ = redisCache.InvalidateTrips(ctx, date)
		}

		// Отправить обновление через WebSocket
		message := &websocket.Message{
			Type:   "trip_update",
			TripID: data["id"].(string),
			Status: data["status"].(string),
		}

		if delay, ok := data["delay_minutes"].(float64); ok {
			message.DelayMinutes = int(delay)
		}

		hub.Broadcast(message)
	})
	if err != nil {
		logger.Error("Failed to subscribe to trip.status_changed", zap.Error(err))
	}

	logger.Info("Subscribed to NATS events: trip.created, trip.status_changed")

	// Разрешить origins для WebSocket (из конфига, через запятую)
	var allowedOrigins []string
	for _, o := range strings.Split(cfg.WebSocket.AllowedOrigins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			allowedOrigins = append(allowedOrigins, trimmed)
		}
	}

	// Создать handlers
	boardHandler := handlers.NewBoardHandler(db, hub, logger, allowedOrigins, cfg.WebSocket.AllowAllOriginsInDev)

	// Настроить Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "board",
			"version": "1.0.0",
		})
	})

	// API routes
	v1 := router.Group("/v1")
	{
		board := v1.Group("/board")
		{
			board.GET("/ws", boardHandler.HandleWebSocket)
			board.GET("/public", boardHandler.GetPublicBoard)
			board.GET("/platform/:platform", boardHandler.GetPlatformBoard)
			board.GET("/stats", boardHandler.GetWebSocketStats)
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
		logger.Info("Board service listening", zap.String("port", cfg.Server.Port))
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
