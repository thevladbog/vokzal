package scanner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ScannerClient — клиент для работы со сканером.
//
//nolint:revive // exported: имя ScannerClient намеренно (Client слишком общее в пакете scanner).
type ScannerClient struct {
	lastTime   time.Time
	logger     *zap.Logger
	devicePath string
	mode       string
	lastScan   string
	enabled    bool
}

// NewScannerClient создаёт клиент сканера.
func NewScannerClient(devicePath, mode string, enabled bool, logger *zap.Logger) *ScannerClient {
	return &ScannerClient{
		devicePath: devicePath,
		mode:       mode,
		enabled:    enabled,
		logger:     logger,
	}
}

// ReadBarcode читает штрихкод/QR код.
func (s *ScannerClient) ReadBarcode() (string, error) {
	if !s.enabled {
		s.logger.Debug("Scanner disabled")
		return "", fmt.Errorf("scanner disabled")
	}

	// Открыть устройство HID сканера
	file, err := os.Open(s.devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to open scanner device: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.Warn("failed to close scanner device", zap.Error(err))
		}
	}()

	s.logger.Debug("Waiting for barcode scan")

	// Читать данные
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	barcode := strings.TrimSpace(scanner.Text())

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read barcode: %w", err)
	}

	s.lastScan = barcode
	s.lastTime = time.Now()

	s.logger.Info("Barcode scanned", zap.String("barcode", barcode))
	return barcode, nil
}

// GetLastScan возвращает последний отсканированный код.
func (s *ScannerClient) GetLastScan() (string, time.Time) {
	return s.lastScan, s.lastTime
}

// GetStatus получает статус сканера.
func (s *ScannerClient) GetStatus() map[string]interface{} {
	if !s.enabled {
		return map[string]interface{}{
			"status": "disabled",
			"online": false,
		}
	}

	// Проверить доступность устройства
	if _, err := os.Stat(s.devicePath); os.IsNotExist(err) {
		return map[string]interface{}{
			"status": "offline",
			"online": false,
			"error":  "device not found",
		}
	}

	return map[string]interface{}{
		"status":    "online",
		"online":    true,
		"last_scan": s.lastScan,
		"last_time": s.lastTime,
	}
}
