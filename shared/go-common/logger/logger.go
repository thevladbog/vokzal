package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger создаёт новый логгер для Вокзал.ТЕХ
func NewLogger(serviceName string, debug bool) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	
	if debug {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	
	config.InitialFields = map[string]interface{}{
		"service": serviceName,
		"project": "vokzal-tech",
	}
	
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	
	return logger, nil
}

// Sugar возвращает sugared logger для удобного логирования
func Sugar(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}
