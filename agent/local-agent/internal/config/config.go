package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config — конфигурация локального агента.
type Config struct {
	KKT     KKTConfig     `mapstructure:"kkt"`
	Printer PrinterConfig `mapstructure:"printer"`
	Scanner ScannerConfig `mapstructure:"scanner"`
	Server  ServerConfig  `mapstructure:"server"`
	Logger  LoggerConfig  `mapstructure:"logger"`
	TTS     TTSConfig     `mapstructure:"tts"`
}

// ServerConfig — настройки HTTP-сервера агента.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// KKTConfig — настройки ККТ (АТОЛ).
type KKTConfig struct {
	DevicePath string `mapstructure:"device_path"`
	ATOLDriver string `mapstructure:"atol_driver"`
	INN        string `mapstructure:"inn"`
	OFDUrl     string `mapstructure:"ofd_url"`
	Enabled    bool   `mapstructure:"enabled"`
}

// PrinterConfig — настройки принтера билетов.
type PrinterConfig struct {
	DevicePath string `mapstructure:"device_path"`
	Type       string `mapstructure:"type"`
	Name       string `mapstructure:"name"`
	Enabled    bool   `mapstructure:"enabled"`
}

// ScannerConfig — настройки сканера штрихкодов.
type ScannerConfig struct {
	DevicePath string `mapstructure:"device_path"`
	Mode       string `mapstructure:"mode"`
	Enabled    bool   `mapstructure:"enabled"`
}

// TTSConfig — настройки голосовых оповещений.
type TTSConfig struct {
	Engine  string `mapstructure:"engine"`
	Voice   string `mapstructure:"voice"`
	Volume  int    `mapstructure:"volume"`
	Enabled bool   `mapstructure:"enabled"`
}

// LoggerConfig — настройки логгера.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// Load загружает конфигурацию из файла и переменных окружения.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vokzal/agent")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VOKZAL_AGENT")

	viper.SetDefault("server.port", "8081")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("kkt.enabled", true)
	viper.SetDefault("kkt.device_path", "/dev/ttyUSB0")
	viper.SetDefault("kkt.atol_driver", "/opt/atol/driver")
	viper.SetDefault("printer.enabled", true)
	viper.SetDefault("printer.device_path", "/dev/usb/lp0")
	viper.SetDefault("printer.type", "escpos")
	viper.SetDefault("scanner.enabled", true)
	viper.SetDefault("scanner.mode", "hid")
	viper.SetDefault("tts.enabled", true)
	viper.SetDefault("tts.engine", "rhvoice")
	viper.SetDefault("tts.voice", "alena")
	viper.SetDefault("tts.volume", 80)
	viper.SetDefault("logger.level", "debug")

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
