// Package main — точка входа Board Service (табло отправлений, WebSocket).
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

func initLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.Logger.Level == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

func subscribeNATS(natsConn *nats.Conn, redisCache *cache.RedisCache, hub *websocket.Hub, logger *zap.Logger) {
	ctx := context.Background()
	_, err := natsConn.Subscribe("trip.created", func(msg *nats.Msg) {
		var data map[string]interface{}
		if unmarshalErr := json.Unmarshal(msg.Data, &data); unmarshalErr != nil {
			logger.Error("Failed to unmarshal trip.created", zap.Error(unmarshalErr))
			return
		}
		if date, ok := data["date"].(string); ok {
			if invErr := redisCache.InvalidateTrips(ctx, date); invErr != nil {
				logger.Warn("failed to invalidate trips cache", zap.Error(invErr), zap.String("date", date))
			}
		}
		var tripID string
		if id, ok := data["id"].(string); ok {
			tripID = id
		}
		hub.Broadcast(&websocket.Message{
			Type:   "trip_created",
			TripID: tripID,
			Data:   data,
		})
	})
	if err != nil {
		logger.Error("Failed to subscribe to trip.created", zap.Error(err))
	}
	_, err = natsConn.Subscribe("trip.status_changed", func(msg *nats.Msg) {
		var data map[string]interface{}
		if unmarshalErr := json.Unmarshal(msg.Data, &data); unmarshalErr != nil {
			logger.Error("Failed to unmarshal trip.status_changed", zap.Error(unmarshalErr))
			return
		}

		// Инвалидировать кэш
		if date, ok := data["date"].(string); ok {
			if invErr := redisCache.InvalidateTrips(ctx, date); invErr != nil {
				logger.Warn("failed to invalidate trips cache", zap.Error(invErr), zap.String("date", date))
			}
		}

		// Отправить обновление через WebSocket
		var tripID, status string
		if id, ok := data["id"].(string); ok {
			tripID = id
		}
		if s, ok := data["status"].(string); ok {
			status = s
		}
		message := &websocket.Message{
			Type:   "trip_update",
			TripID: tripID,
			Status: status,
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

	logger.Info("Starting Board Service", zap.String("version", "1.0.0"))

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	ctx := context.Background()
	if pingErr := redisClient.Ping(ctx).Err(); pingErr != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(pingErr))
	}
	logger.Info("Connected to Redis", zap.String("addr", cfg.Redis.Address()))
	redisCache := cache.NewRedisCache(redisClient)

	natsConn, err := nats.Connect(cfg.NATS.URL,
		nats.UserInfo(cfg.NATS.User, cfg.NATS.Password),
		nats.Name("board-service"))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsConn.Close()
	logger.Info("Connected to NATS", zap.String("url", cfg.NATS.URL))

	hub := websocket.NewHub(logger)
	go hub.Run()
	subscribeNATS(natsConn, redisCache, hub, logger)

	var allowedOrigins []string
	for _, o := range strings.Split(cfg.WebSocket.AllowedOrigins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			allowedOrigins = append(allowedOrigins, trimmed)
		}
	}
	boardHandler := handlers.NewBoardHandler(db, hub, logger, allowedOrigins, cfg.WebSocket.AllowAllOriginsInDev)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "board",
			"version": "1.0.0",
		})
	})
	// Traefik strips /v1/board prefix, service receives /ws, /public, /platform/:platform, /stats
	router.GET("/ws", boardHandler.HandleWebSocket)
	router.GET("/public", boardHandler.GetPublicBoard)
	router.GET("/platform/:platform", boardHandler.GetPlatformBoard)
	router.GET("/stats", boardHandler.GetWebSocketStats)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		logger.Info("Board service listening", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server exited")
}
