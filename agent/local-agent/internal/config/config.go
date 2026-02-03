package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	KKT     KKTConfig     `mapstructure:"kkt"`
	Printer PrinterConfig `mapstructure:"printer"`
	Scanner ScannerConfig `mapstructure:"scanner"`
	TTS     TTSConfig     `mapstructure:"tts"`
	Logger  LoggerConfig  `mapstructure:"logger"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type KKTConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	DevicePath string `mapstructure:"device_path"`
	ATOLDriver string `mapstructure:"atol_driver"` // path to ATOL driver
	INN        string `mapstructure:"inn"`
	OFDUrl     string `mapstructure:"ofd_url"`
}

type PrinterConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	DevicePath string `mapstructure:"device_path"`
	Type       string `mapstructure:"type"` // escpos, cups
	Name       string `mapstructure:"name"`
}

type ScannerConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	DevicePath string `mapstructure:"device_path"`
	Mode       string `mapstructure:"mode"` // hid, serial
}

type TTSConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Engine  string `mapstructure:"engine"` // rhvoice, espeak
	Voice   string `mapstructure:"voice"`  // alena, david
	Volume  int    `mapstructure:"volume"` // 0-100
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

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
