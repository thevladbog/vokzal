// Package config загружает конфигурацию Ticket Service.
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config — корневая конфигурация сервиса.
type Config struct {
	NATS     NATSConfig     `mapstructure:"nats"`
	Server   ServerConfig   `mapstructure:"server"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
	Business BusinessConfig `mapstructure:"business"`
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

	return &config, nil
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
