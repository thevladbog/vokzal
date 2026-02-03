package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	YandexMaps YandexMapsConfig `mapstructure:"yandex_maps"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type YandexMapsConfig struct {
	APIKey string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

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
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
