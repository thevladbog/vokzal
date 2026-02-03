// Package config загружает конфигурацию Auth Service из YAML и переменных окружения.
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config — корневая конфигурация сервиса.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig — настройки HTTP-сервера.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release
}

// DatabaseConfig — настройки подключения к PostgreSQL.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig — настройки JWT.
type JWTConfig struct {
	Secret            string        `mapstructure:"secret"`
	Issuer            string        `mapstructure:"issuer"`
	AccessExpiration  time.Duration `mapstructure:"access_expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
}

// LoggerConfig — настройки логгера.
type LoggerConfig struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
}

// Load читает конфигурацию из файла и переменных окружения.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vokzal/auth")

	// Переменные окружения
	viper.AutomaticEnv()
	viper.SetEnvPrefix("VOKZAL_AUTH")

	// Значения по умолчанию
	viper.SetDefault("server.port", "8081")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "admin")
	viper.SetDefault("database.password", "vokzal_secret_2026")
	viper.SetDefault("database.dbname", "vokzal")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.secret", "vokzal-tech-jwt-secret-change-me")
	viper.SetDefault("jwt.issuer", "vokzal.tech")
	viper.SetDefault("jwt.access_expiration", "24h")
	viper.SetDefault("jwt.refresh_expiration", "168h") // 7 дней
	viper.SetDefault("logger.level", "debug")

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Файл не найден - используем значения по умолчанию
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
