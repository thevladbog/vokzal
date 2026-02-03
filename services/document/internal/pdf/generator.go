// Package pdf — генерация PDF-документов (билеты, ПД-2).
package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

// Generator — генератор PDF-документов.
type Generator struct {
	logger *zap.Logger
}

// NewGenerator создаёт новый Generator.
func NewGenerator(logger *zap.Logger) *Generator {
	return &Generator{logger: logger}
}

// TicketData — данные для билета.
type TicketData struct {
	TicketID     string
	PassengerFIO string
	PassengerDoc string
	Route        string
	Date         string
	Time         string
	Platform     string
	Seat         string
	Price        float64
	QRCode       string
	BarCode      string
}

// GenerateTicket генерирует билет в PDF.
func (g *Generator) GenerateTicket(data *TicketData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("DejaVu", "", "DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 14)

	// Заголовок
	pdf.CellFormat(190, 10, "Вокзал.ТЕХ - Электронный билет", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Данные билета
	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(60, 6, "Номер билета:", "", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, data.TicketID, "", 1, "L", false, 0, "")

	pdf.CellFormat(60, 6, "Пассажир:", "", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, data.PassengerFIO, "", 1, "L", false, 0, "")

	pdf.CellFormat(60, 6, "Документ:", "", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, data.PassengerDoc, "", 1, "L", false, 0, "")

	pdf.Ln(3)
	pdf.SetFont("DejaVu", "", 12)
	pdf.CellFormat(60, 8, "Маршрут:", "", 0, "L", false, 0, "")
	pdf.CellFormat(130, 8, data.Route, "", 1, "L", false, 0, "")

	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(60, 6, "Дата:", "", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, data.Date, "", 1, "L", false, 0, "")

	pdf.CellFormat(60, 6, "Время отправления:", "", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, data.Time, "", 1, "L", false, 0, "")

	pdf.CellFormat(60, 6, "Перрон:", "", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, data.Platform, "", 1, "L", false, 0, "")

	if data.Seat != "" {
		pdf.CellFormat(60, 6, "Место:", "", 0, "L", false, 0, "")
		pdf.CellFormat(130, 6, data.Seat, "", 1, "L", false, 0, "")
	}

	pdf.CellFormat(60, 6, "Стоимость:", "", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, fmt.Sprintf("%.2f руб.", data.Price), "", 1, "L", false, 0, "")

	pdf.Ln(5)

	// QR код
	if data.QRCode != "" {
		qr, err := qrcode.Encode(data.QRCode, qrcode.Medium, 256)
		if err == nil {
			pdf.RegisterImageReader("qr", "png", bytes.NewReader(qr))
			pdf.Image("qr", 80, pdf.GetY(), 50, 50, false, "", 0, "")
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// PD2Data — данные для ПД-2 (проездной документ).
type PD2Data struct {
	Number        string
	Series        string
	PassengerFIO  string
	PassengerDoc  string
	RouteFrom     string
	RouteTo       string
	Date          string
	Price         float64
	IssueDate     string
	IssuerName    string
	BusNumber     string
}

// GeneratePD2 генерирует форму ПД-2.
func (g *Generator) GeneratePD2(data *PD2Data) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("DejaVu", "", "DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 14)

	// Заголовок
	pdf.CellFormat(190, 10, "ПРОЕЗДНОЙ ДОКУМЕНТ (ПД-2)", "", 1, "C", false, 0, "")
	pdf.Ln(3)

	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(40, 6, "Серия:", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.Series, "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "Номер:", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.Number, "", 1, "L", false, 0, "")

	pdf.Ln(2)
	pdf.CellFormat(40, 6, "Пассажир:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.PassengerFIO, "", 1, "L", false, 0, "")

	pdf.CellFormat(40, 6, "Документ:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.PassengerDoc, "", 1, "L", false, 0, "")

	pdf.Ln(2)
	pdf.CellFormat(40, 6, "Откуда:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.RouteFrom, "", 1, "L", false, 0, "")

	pdf.CellFormat(40, 6, "Куда:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.RouteTo, "", 1, "L", false, 0, "")

	pdf.Ln(2)
	pdf.CellFormat(40, 6, "Дата отправления:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.Date, "", 1, "L", false, 0, "")

	pdf.CellFormat(40, 6, "Автобус №:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.BusNumber, "", 1, "L", false, 0, "")

	pdf.Ln(2)
	pdf.CellFormat(40, 6, "Стоимость:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, fmt.Sprintf("%.2f руб.", data.Price), "", 1, "L", false, 0, "")

	pdf.Ln(5)
	pdf.CellFormat(40, 6, "Дата выдачи:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.IssueDate, "", 1, "L", false, 0, "")

	pdf.CellFormat(40, 6, "Кассир:", "", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.IssuerName, "", 1, "L", false, 0, "")

	pdf.Ln(10)
	pdf.SetFont("DejaVu", "", 8)
	pdf.CellFormat(190, 4, "Документ действителен только при предъявлении паспорта", "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 4, "© "+time.Now().Format("2006")+" Вокзал.ТЕХ", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}

	return buf.Bytes(), nil
}
