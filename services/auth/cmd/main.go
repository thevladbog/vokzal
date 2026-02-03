// Package main — точка входа Auth Service (аутентификация и авторизация).
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vokzal-tech/auth-service/internal/config"
	"github.com/vokzal-tech/auth-service/internal/handlers"
	"github.com/vokzal-tech/auth-service/internal/middleware"
	"github.com/vokzal-tech/auth-service/internal/repository"
	"github.com/vokzal-tech/auth-service/internal/service"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const serverModeRelease = "release"

func main() {
	// Загрузить конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создать логгер
	var zapLogger *zap.Logger
	if cfg.Server.Mode == serverModeRelease {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		if syncErr := zapLogger.Sync(); syncErr != nil {
			log.Printf("logger sync: %v", syncErr)
		}
	}()

	zapLogger.Info("Вокзал.ТЕХ Auth Service starting...",
		zap.String("version", "1.0.0"),
		zap.String("mode", cfg.Server.Mode))

	// Подключиться к базе данных
	gormConfig := &gorm.Config{}
	if cfg.Server.Mode == serverModeRelease {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Error)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), gormConfig)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		zapLogger.Fatal("Failed to get database instance", zap.Error(err))
	}
	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			zapLogger.Error("failed to close database", zap.Error(closeErr))
		}
	}()

	// Настроить пул соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	zapLogger.Info("Database connected successfully")

	// Создать репозитории
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Создать сервисы
	authService := service.NewAuthService(userRepo, sessionRepo, cfg.JWT, zapLogger)

	// Создать handlers
	authHandler := handlers.NewAuthHandler(authService, zapLogger)

	// Создать middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, zapLogger)

	// Настроить Gin
	if cfg.Server.Mode == serverModeRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger(zapLogger))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "auth",
			"version": "1.0.0",
		})
	})

	// API routes
	v1 := router.Group("/v1")
	auth := v1.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/me", authMiddleware.RequireAuth(), authHandler.Me)

	// Users CRUD (admin only)
	users := v1.Group("/users")
	users.Use(authMiddleware.RequireAuth(), authMiddleware.RequireRole("admin"))
	users.GET("", authHandler.ListUsers)
	users.POST("", authHandler.CreateUser)
	users.GET("/:id", authHandler.GetUser)
	users.PUT("/:id", authHandler.UpdateUser)
	users.DELETE("/:id", authHandler.DeleteUser)

	// Создать сервер
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запустить сервер в горутине
	go func() {
		zapLogger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zapLogger.Info("Server stopped")
}
