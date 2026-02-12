// Package config загружает конфигурацию Ticket Service.
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Insecure JWT secret placeholders; app must not start in release mode when secret is one of these.
var insecureJWTSecrets = []string{
	"vokzal_jwt_secret_change_in_production",
	"vokzal-tech-jwt-secret-change-me-in-production",
	"vokzal-tech-jwt-secret-change-me",
}

// Config — корневая конфигурация сервиса.
//
//nolint:govet // fieldalignment: keep mapstructure/readability order
type Config struct {
	NATS     NATSConfig     `mapstructure:"nats"`
	Server   ServerConfig   `mapstructure:"server"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
	Business BusinessConfig `mapstructure:"business"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// JWTConfig — настройки JWT для проверки токенов (тот же секрет, что в Auth Service).
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// ServerConfig — настройки HTTP-сервера.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig — настройки БД.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Port     int    `mapstructure:"port"`
}

// NATSConfig — настройки подключения к NATS.
type NATSConfig struct {
	URL      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// LoggerConfig — настройки логгера.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// BusinessConfig — бизнес-настройки (штрафы за возврат и т.п.).
type BusinessConfig struct {
	RefundPenalty RefundPenaltyConfig `mapstructure:"refund_penalty"`
}

// RefundPenaltyConfig — коэффициенты штрафа за возврат по времени до отправления.
type RefundPenaltyConfig struct {
	Over24Hours  float64 `mapstructure:"over_24_hours"`
	Between12_24 float64 `mapstructure:"between_12_24"`
	Under12Hours float64 `mapstructure:"under_12_hours"`
}

// Load загружает конфигурацию из файла и переменных окружения.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vokzal/ticket")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VOKZAL_TICKET")

	viper.SetDefault("server.port", "8083")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "admin")
	viper.SetDefault("database.password", "vokzal_secret_2026")
	viper.SetDefault("database.dbname", "vokzal")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("nats.url", "nats://localhost:4222")
	viper.SetDefault("nats.user", "vokzal")
	viper.SetDefault("nats.password", "nats_secret_2026")
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("business.refund_penalty.over_24_hours", 0.10)
	viper.SetDefault("business.refund_penalty.between_12_24", 0.20)
	viper.SetDefault("business.refund_penalty.under_12_hours", 0.30)
	viper.SetDefault("jwt.secret", "vokzal_jwt_secret_change_in_production")

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// In release mode, refuse to run with default/placeholder JWT secret (use VOKZAL_TICKET_JWT_SECRET in production).
	if config.Server.Mode == "release" {
		s := strings.TrimSpace(config.JWT.Secret)
		if s == "" {
			return nil, fmt.Errorf("jwt.secret must not be empty or blank in release mode; set VOKZAL_TICKET_JWT_SECRET")
		}
		for _, bad := range insecureJWTSecrets {
			if s == bad {
				return nil, fmt.Errorf("jwt.secret must not be the default/placeholder in release mode; set VOKZAL_TICKET_JWT_SECRET")
			}
		}
	}

	return &config, nil
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
