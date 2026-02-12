// Package config загружает конфигурацию Fiscal Service.
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

// Config — корневая конфигурация сервиса (поля по убыванию размера для выравнивания).
type Config struct {
	NATS       NATSConfig       `mapstructure:"nats"`
	ATOL       ATOLConfig       `mapstructure:"atol"`
	Server     ServerConfig     `mapstructure:"server"`
	LocalAgent LocalAgentConfig `mapstructure:"local_agent"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Database   DatabaseConfig   `mapstructure:"database"`
}

// JWTConfig — настройки JWT (тот же секрет, что у auth-service, для проверки токенов).
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// ServerConfig — настройки HTTP-сервера.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig — настройки PostgreSQL.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Port     int    `mapstructure:"port"`
}

// NATSConfig — настройки NATS.
type NATSConfig struct {
	URL      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// LoggerConfig — настройки логгера.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// ATOLConfig — настройки АТОЛ ККТ.
type ATOLConfig struct {
	CompanyINN  string `mapstructure:"company_inn"`
	CompanyName string `mapstructure:"company_name"`
	TaxSystem   string `mapstructure:"tax_system"`
}

// LocalAgentConfig — настройки локального агента ККТ.
type LocalAgentConfig struct {
	URL string `mapstructure:"url"`
}

// Load читает конфигурацию из файла и переменных окружения.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vokzal/fiscal")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VOKZAL_FISCAL")

	viper.SetDefault("server.port", "8084")
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
	viper.SetDefault("atol.company_inn", "1234567890")
	viper.SetDefault("atol.company_name", "ООО «Вокзал.ТЕХ»")
	viper.SetDefault("atol.tax_system", "osn")
	viper.SetDefault("local_agent.url", "http://localhost:8081")
	viper.SetDefault("jwt.secret", "vokzal-tech-jwt-secret-change-me") // override in production; must match auth-service

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

	// In release mode, refuse to run with default/placeholder JWT secret (use VOKZAL_FISCAL_JWT_SECRET in production).
	if config.Server.Mode == "release" {
		s := strings.TrimSpace(config.JWT.Secret)
		if s == "" {
			return nil, fmt.Errorf("jwt.secret must not be empty or blank in release mode; set VOKZAL_FISCAL_JWT_SECRET")
		}
		for _, bad := range insecureJWTSecrets {
			if s == bad {
				return nil, fmt.Errorf("jwt.secret must not be the default/placeholder in release mode; set VOKZAL_FISCAL_JWT_SECRET")
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
