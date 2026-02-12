// Package config загружает конфигурацию Document Service.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
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
	NATS     NATSConfig     `mapstructure:"nats"`
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
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

// MinIOConfig — настройки MinIO.
type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

// loadDotEnv loads .env from the current working directory (for local dev).
// Real credentials must not be committed; .env is in .gitignore.
func loadDotEnv() {
	for _, path := range []string{".", "./config"} {
		f, err := os.Open(path + "/.env")
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			idx := strings.Index(line, "=")
			if idx <= 0 {
				continue
			}
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			val = strings.Trim(val, "\"'")
			if key != "" {
				if setErr := os.Setenv(key, val); setErr != nil {
					continue
				}
			}
		}
		if closeErr := f.Close(); closeErr != nil {
			_ = closeErr // ignore close error on .env read
		}
		break
	}
}

var envExpandRe = regexp.MustCompile(`\$\{([^}:]+)(?::-([^}]*))?\}`)

// expandEnv replaces ${VAR} and ${VAR:-default} in s with os.Getenv("VAR") or default.
func expandEnv(s string) string {
	return envExpandRe.ReplaceAllStringFunc(s, func(match string) string {
		sub := envExpandRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		key := sub[1]
		val := os.Getenv(key)
		if val != "" {
			return val
		}
		if len(sub) >= 3 {
			return sub[2]
		}
		return ""
	})
}

// Load читает конфигурацию из файла и переменных окружения.
// For local dev, place a .env (see .env.example) in the service directory; do not commit .env.
func Load() (*Config, error) {
	loadDotEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vokzal/document")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VOKZAL_DOCUMENT")

	viper.SetDefault("server.port", "8089")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "admin")
	viper.SetDefault("database.password", "changeme")
	viper.SetDefault("database.dbname", "vokzal")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("nats.url", "nats://localhost:4222")
	viper.SetDefault("nats.user", "")
	viper.SetDefault("nats.password", "")
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("minio.endpoint", "localhost:9000")
	viper.SetDefault("minio.access_key", "")
	viper.SetDefault("minio.secret_key", "")
	viper.SetDefault("minio.bucket", "vokzal-documents")
	viper.SetDefault("minio.use_ssl", false)
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

	// Expand ${VAR} and ${VAR:-default} in string fields (secrets and URLs from env).
	config.Database.Host = expandEnv(config.Database.Host)
	config.Database.User = expandEnv(config.Database.User)
	config.Database.Password = expandEnv(config.Database.Password)
	config.Database.DBName = expandEnv(config.Database.DBName)
	config.Database.SSLMode = expandEnv(config.Database.SSLMode)
	config.NATS.URL = expandEnv(config.NATS.URL)
	config.NATS.User = expandEnv(config.NATS.User)
	config.NATS.Password = expandEnv(config.NATS.Password)
	config.MinIO.Endpoint = expandEnv(config.MinIO.Endpoint)
	config.MinIO.AccessKey = expandEnv(config.MinIO.AccessKey)
	config.MinIO.SecretKey = expandEnv(config.MinIO.SecretKey)
	config.MinIO.Bucket = expandEnv(config.MinIO.Bucket)
	config.JWT.Secret = expandEnv(config.JWT.Secret)

	if strings.TrimSpace(config.Database.User) == "" || strings.TrimSpace(config.Database.Password) == "" {
		return nil, fmt.Errorf("database.user and database.password must be set (e.g. via DB_USER, DB_PASSWORD or .env)")
	}

	// In release mode, refuse to run with default/placeholder JWT secret (use VOKZAL_DOCUMENT_JWT_SECRET in production).
	if config.Server.Mode == "release" {
		s := strings.TrimSpace(config.JWT.Secret)
		if s == "" {
			return nil, fmt.Errorf("jwt.secret must not be empty or blank in release mode; set VOKZAL_DOCUMENT_JWT_SECRET")
		}
		for _, bad := range insecureJWTSecrets {
			if s == bad {
				return nil, fmt.Errorf("jwt.secret must not be the default/placeholder in release mode; set VOKZAL_DOCUMENT_JWT_SECRET")
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
