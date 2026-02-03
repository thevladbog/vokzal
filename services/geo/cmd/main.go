// Package main — точка входа Geo Service (Yandex Geocoder, расстояния).
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
	"github.com/redis/go-redis/v9"

	"github.com/vokzal-tech/geo-service/internal/config"
	"github.com/vokzal-tech/geo-service/internal/handlers"
	"github.com/vokzal-tech/geo-service/internal/service"
	"github.com/vokzal-tech/geo-service/internal/yandex"

	"go.uber.org/zap"
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

	logger.Info("Starting Geo Service", zap.String("version", "1.0.0"))

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

	// Создать Yandex Maps клиент
	yandexClient := yandex.NewClient(cfg.YandexMaps.APIKey, cfg.YandexMaps.BaseURL, logger)

	geoService := service.NewGeoService(yandexClient, redisClient, logger)
	geoHandler := handlers.NewGeoHandler(geoService, logger)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "geo",
			"version": "1.0.0",
		})
	})

	v1 := router.Group("/v1")
	geo := v1.Group("/geo")
	geo.GET("/geocode", geoHandler.Geocode)
	geo.GET("/reverse", geoHandler.ReverseGeocode)
	geo.GET("/distance", geoHandler.GetDistance)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Geo service listening", zap.String("port", cfg.Server.Port))
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
