// Package config загружает конфигурацию Payment Service.
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config — корневая конфигурация сервиса.
type Config struct {
	NATS     NATSConfig     `mapstructure:"nats"`
	Tinkoff  TinkoffConfig  `mapstructure:"tinkoff"`
	SBP      SBPConfig      `mapstructure:"sbp"`
	Server   ServerConfig   `mapstructure:"server"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
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

// TinkoffConfig — настройки Tinkoff Acquiring.
type TinkoffConfig struct {
	TerminalKey string `mapstructure:"terminal_key"`
	// Password — это API секрет терминала (не пароль пользователя), используется для подписи запросов
	Password string `mapstructure:"password"`
	APIUrl   string `mapstructure:"api_url"`
}

// SBPConfig — настройки СБП.
type SBPConfig struct {
	MerchantID string `mapstructure:"merchant_id"`
	APIUrl     string `mapstructure:"api_url"`
	APIKey     string `mapstructure:"api_key"`
}

// Load загружает конфигурацию из файла и переменных окружения.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vokzal/payment")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VOKZAL_PAYMENT")

	viper.SetDefault("server.port", "8085")
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
	viper.SetDefault("tinkoff.api_url", "https://securepay.tinkoff.ru/v2")
	viper.SetDefault("sbp.api_url", "https://api.sbp.nspk.ru")

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
