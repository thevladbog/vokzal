// Package config загружает конфигурацию Notify Service.
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config — корневая конфигурация сервиса.
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	NATS        NATSConfig        `mapstructure:"nats"`
	Logger      LoggerConfig      `mapstructure:"logger"`
	SMS         SMSConfig         `mapstructure:"sms"`
	Email       EmailConfig       `mapstructure:"email"`
	Telegram    TelegramConfig    `mapstructure:"telegram"`
	LocalAgent  LocalAgentConfig  `mapstructure:"local_agent"`
}

// ServerConfig — настройки HTTP-сервера.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig — настройки БД.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
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

// SMSConfig — настройки SMS (SMS.ru).
type SMSConfig struct {
	APIID string `mapstructure:"api_id"`
	URL   string `mapstructure:"url"`
}

// EmailConfig — настройки SMTP для email.
type EmailConfig struct {
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// TelegramConfig — настройки Telegram-бота.
type TelegramConfig struct {
	BotToken string `mapstructure:"bot_token"`
	WebhookURL string `mapstructure:"webhook_url"`
}

// LocalAgentConfig — настройки локального агента (TTS и т.п.).
type LocalAgentConfig struct {
	URL string `mapstructure:"url"`
}

// Load загружает конфигурацию из файла и переменных окружения.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vokzal/notify")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VOKZAL_NOTIFY")

	viper.SetDefault("server.port", "8087")
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
	viper.SetDefault("sms.url", "https://sms.ru/sms/send")
	viper.SetDefault("email.smtp_port", 587)
	viper.SetDefault("email.from", "noreply@vokzal.tech")
	viper.SetDefault("local_agent.url", "http://localhost:8081")

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
