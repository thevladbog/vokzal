// Package main — точка входа Schedule Service.
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

	"github.com/vokzal-tech/schedule-service/internal/config"
	"github.com/vokzal-tech/schedule-service/internal/handlers"
	"github.com/vokzal-tech/schedule-service/internal/models"
	"github.com/vokzal-tech/schedule-service/internal/repository"
	"github.com/vokzal-tech/schedule-service/internal/service"
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

	logger.Info("Starting Schedule Service", zap.String("version", "1.0.0"))

	// Подключиться к БД
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if migErr := db.AutoMigrate(&models.Station{}, &models.Route{}, &models.Schedule{}, &models.Trip{}); migErr != nil {
		logger.Warn("Auto-migration failed", zap.Error(migErr))
	}

	// Подключиться к NATS
	natsConn, err := nats.Connect(cfg.NATS.URL,
		nats.UserInfo(cfg.NATS.User, cfg.NATS.Password),
		nats.Name("schedule-service"))
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsConn.Close()

	logger.Info("Connected to NATS", zap.String("url", cfg.NATS.URL))

	// Создать репозитории
	stationRepo := repository.NewStationRepository(db)
	routeRepo := repository.NewRouteRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	tripRepo := repository.NewTripRepository(db)

	// Создать сервис
	scheduleService := service.NewScheduleService(stationRepo, routeRepo, scheduleRepo, tripRepo, natsConn, logger)

	// Создать handlers
	scheduleHandler := handlers.NewScheduleHandler(scheduleService, logger)

	// Настроить Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "schedule",
			"version": "1.0.0",
		})
	})

	v1 := router.Group("/v1")
	stations := v1.Group("/stations")
	stations.POST("", scheduleHandler.CreateStation)
	stations.GET("", scheduleHandler.ListStations)
	stations.GET("/:id", scheduleHandler.GetStation)
	stations.PATCH("/:id", scheduleHandler.UpdateStation)
	stations.DELETE("/:id", scheduleHandler.DeleteStation)
	routes := v1.Group("/routes")
	routes.POST("", scheduleHandler.CreateRoute)
	routes.GET("", scheduleHandler.ListRoutes)
	routes.GET("/:id", scheduleHandler.GetRoute)
	routes.PATCH("/:id", scheduleHandler.UpdateRoute)
	routes.DELETE("/:id", scheduleHandler.DeleteRoute)
	schedules := v1.Group("/schedules")
	schedules.POST("", scheduleHandler.CreateSchedule)
	schedules.GET("", scheduleHandler.ListSchedulesByRoute)
	schedules.GET("/:id", scheduleHandler.GetSchedule)
	schedules.PATCH("/:id", scheduleHandler.UpdateSchedule)
	schedules.DELETE("/:id", scheduleHandler.DeleteSchedule)
	trips := v1.Group("/trips")
	trips.POST("", scheduleHandler.CreateTrip)
	trips.GET("", scheduleHandler.ListTripsByDate)
	trips.GET("/:id", scheduleHandler.GetTrip)
	trips.PATCH("/:id/status", scheduleHandler.UpdateTripStatus)
	trips.POST("/generate", scheduleHandler.GenerateTrips)

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
		logger.Info("Schedule service listening", zap.String("port", cfg.Server.Port))
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
