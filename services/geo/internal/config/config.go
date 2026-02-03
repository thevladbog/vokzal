// Package config загружает конфигурацию Geo Service.
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config — корневая конфигурация сервиса.
//
//nolint:govet // fieldalignment: keep field order for mapstructure/config clarity
type Config struct {
	Logger     LoggerConfig     `mapstructure:"logger"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Server     ServerConfig     `mapstructure:"server"`
	YandexMaps YandexMapsConfig `mapstructure:"yandex_maps"`
}

// ServerConfig — настройки HTTP-сервера.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// RedisConfig — настройки Redis.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
}

// LoggerConfig — настройки логгера.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// YandexMapsConfig — настройки Yandex Geocoder API.
type YandexMapsConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

// Load читает конфигурацию из файла и переменных окружения.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vokzal/geo")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VOKZAL_GEO")

	viper.SetDefault("server.port", "8090")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "vokzal_redis_2026")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("yandex_maps.base_url", "https://geocode-maps.yandex.ru/1.x/")

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

// Address возвращает адрес Redis.
func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
