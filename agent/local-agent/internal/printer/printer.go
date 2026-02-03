package printer

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// PrinterClient — клиент для печати билетов.
//
//nolint:revive // exported: имя PrinterClient намеренно (Client слишком общее в пакете printer).
type PrinterClient struct {
	logger      *zap.Logger
	devicePath  string
	printerType string
	enabled     bool
}

// TicketData — данные билета для печати.
type TicketData struct {
	TicketID     string
	Route        string
	Date         string
	Time         string
	Platform     string
	Seat         string
	PassengerFIO string
	QRCode       string
	BarCode      string
	Price        float64
}

// NewPrinterClient создаёт клиент принтера.
func NewPrinterClient(devicePath, printerType string, enabled bool, logger *zap.Logger) *PrinterClient {
	return &PrinterClient{
		devicePath:  devicePath,
		printerType: printerType,
		enabled:     enabled,
		logger:      logger,
	}
}

// PrintTicket печатает билет.
func (p *PrinterClient) PrintTicket(data *TicketData) error {
	if !p.enabled {
		p.logger.Info("Printer disabled, simulating ticket print")
		return nil
	}

	p.logger.Info("Printing ticket", zap.String("ticket_id", data.TicketID))

	// ESC/POS команды для термопринтера
	commands := p.generateESCPOSCommands(data)

	// Отправить команды на принтер
	file, err := os.OpenFile(p.devicePath, os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open printer: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			p.logger.Warn("failed to close printer device", zap.Error(err))
		}
	}()

	if _, err := file.Write(commands); err != nil {
		return fmt.Errorf("failed to write to printer: %w", err)
	}

	p.logger.Info("Ticket printed successfully")
	return nil
}

func (p *PrinterClient) generateESCPOSCommands(data *TicketData) []byte {
	var buf []byte

	// ESC @ — Reset; выравнивание по центру; жирный шрифт.
	buf = append(buf, 0x1B, 0x40, 0x1B, 0x61, 0x01, 0x1B, 0x45, 0x01)
	buf = append(buf, []byte("Вокзал.ТЕХ\n")...)
	buf = append(buf, 0x1B, 0x45, 0x00)

	// Обычный шрифт
	buf = append(buf, []byte("Электронный билет\n\n")...)

	// Выравнивание слева
	buf = append(buf, 0x1B, 0x61, 0x00)

	buf = append(buf, []byte(fmt.Sprintf("Билет: %s\n", data.TicketID))...)
	buf = append(buf, []byte(fmt.Sprintf("Маршрут: %s\n", data.Route))...)
	buf = append(buf, []byte(fmt.Sprintf("Дата: %s, %s\n", data.Date, data.Time))...)
	buf = append(buf, []byte(fmt.Sprintf("Перрон: %s\n", data.Platform))...)

	if data.Seat != "" {
		buf = append(buf, []byte(fmt.Sprintf("Место: %s\n", data.Seat))...)
	}

	if data.PassengerFIO != "" {
		buf = append(buf, []byte(fmt.Sprintf("Пассажир: %s\n", data.PassengerFIO))...)
	}

	buf = append(buf, []byte(fmt.Sprintf("Цена: %.2f руб.\n\n", data.Price))...)

	// QR код (если поддерживается)
	if data.QRCode != "" {
		buf = append(buf, 0x1B, 0x61, 0x01) // Center
		buf = append(buf, []byte("[QR CODE]\n")...)
		buf = append(buf, []byte(data.QRCode)...)
		buf = append(buf, []byte("\n\n")...)
	}

	// Barcode
	if data.BarCode != "" {
		buf = append(buf, []byte(fmt.Sprintf("ШК: %s\n", data.BarCode))...)
	}

	buf = append(buf, []byte("\n--------------------------------\n")...)
	buf = append(buf, []byte("© 2025 Вокзал.ТЕХ\n")...)
	buf = append(buf, []byte("Билет действителен при предъявлении паспорта\n\n")...)

	// Отрезать чек
	buf = append(buf, 0x1D, 0x56, 0x00)

	return buf
}

// GetStatus получает статус принтера.
func (p *PrinterClient) GetStatus() (map[string]interface{}, error) {
	if !p.enabled {
		return map[string]interface{}{
			"status": "simulated",
			"online": true,
			"paper":  true,
		}, nil
	}

	// Проверить доступность устройства
	if _, err := os.Stat(p.devicePath); os.IsNotExist(err) {
		return map[string]interface{}{
			"status": "offline",
			"online": false,
			"error":  "device not found",
		}, nil
	}

	return map[string]interface{}{
		"status": "online",
		"online": true,
		"paper":  true,
	}, nil
}
