// Package repository содержит слой доступа к данным билетов и посадки.
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/vokzal-tech/ticket-service/internal/models"
	"gorm.io/gorm"
)

var (
	// ErrTicketNotFound возвращается, когда билет не найден.
	ErrTicketNotFound = errors.New("ticket not found")
	// ErrSeatAlreadyTaken возвращается, когда место уже занято.
	ErrSeatAlreadyTaken = errors.New("seat already taken")
	// ErrBoardingAlreadyStarted возвращается, когда посадка уже начата.
	ErrBoardingAlreadyStarted = errors.New("boarding already started")
	// ErrBoardingNotStarted возвращается, когда посадка ещё не начата.
	ErrBoardingNotStarted = errors.New("boarding not started")
)

// TicketRepository — интерфейс репозитория билетов.
type TicketRepository interface {
	Create(ctx context.Context, ticket *models.Ticket) error
	FindByID(ctx context.Context, id string) (*models.Ticket, error)
	FindByQRCode(ctx context.Context, qrCode string) (*models.Ticket, error)
	FindByTripID(ctx context.Context, tripID string) ([]*models.Ticket, error)
	CheckSeatAvailability(ctx context.Context, tripID, seatID string) (bool, error)
	Update(ctx context.Context, ticket *models.Ticket) error
	Delete(ctx context.Context, id string) error
	GetTripDepartureTime(ctx context.Context, tripID string) (*time.Time, error)
}

// BoardingRepository — интерфейс репозитория событий и отметок посадки.
type BoardingRepository interface {
	CreateEvent(ctx context.Context, event *models.BoardingEvent) error
	FindEventByTripID(ctx context.Context, tripID string) (*models.BoardingEvent, error)
	CreateMark(ctx context.Context, mark *models.BoardingMark) error
	FindMarksByTripID(ctx context.Context, tripID string) ([]*models.BoardingMark, error)
	CheckIfMarked(ctx context.Context, ticketID string) (bool, error)
}

type ticketRepository struct {
	db *gorm.DB
}

type boardingRepository struct {
	db *gorm.DB
}

// NewTicketRepository создаёт репозиторий билетов.
func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

// NewBoardingRepository создаёт репозиторий посадки.
func NewBoardingRepository(db *gorm.DB) BoardingRepository {
	return &boardingRepository{db: db}
}

// Create создаёт билет.
func (r *ticketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) FindByID(ctx context.Context, id string) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.WithContext(ctx).First(&ticket, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByQRCode(ctx context.Context, qrCode string) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.WithContext(ctx).First(&ticket, "qr_code = ?", qrCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByTripID(ctx context.Context, tripID string) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	if err := r.db.WithContext(ctx).Where("trip_id = ?", tripID).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) CheckSeatAvailability(ctx context.Context, tripID, seatID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Ticket{}).
		Where("trip_id = ? AND seat_id = ? AND status = ?", tripID, seatID, "active").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *ticketRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Ticket{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTicketNotFound
	}
	return nil
}

// GetTripDepartureTime возвращает время отправления рейса (date + schedule.departure_time из БД).
func (r *ticketRepository) GetTripDepartureTime(ctx context.Context, tripID string) (*time.Time, error) {
	var dep time.Time
	err := r.db.WithContext(ctx).Raw(
		"SELECT (t.date + s.departure_time) AS dep FROM trips t JOIN schedules s ON s.id = t.schedule_id WHERE t.id = ?",
		tripID,
	).Scan(&dep).Error
	if err != nil {
		return nil, err
	}
	if dep.IsZero() {
		return nil, nil
	}
	return &dep, nil
}

// CreateEvent создаёт событие начала посадки.
func (r *boardingRepository) CreateEvent(ctx context.Context, event *models.BoardingEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *boardingRepository) FindEventByTripID(ctx context.Context, tripID string) (*models.BoardingEvent, error) {
	var event models.BoardingEvent
	if err := r.db.WithContext(ctx).First(&event, "trip_id = ?", tripID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *boardingRepository) CreateMark(ctx context.Context, mark *models.BoardingMark) error {
	return r.db.WithContext(ctx).Create(mark).Error
}

func (r *boardingRepository) FindMarksByTripID(ctx context.Context, tripID string) ([]*models.BoardingMark, error) {
	var marks []*models.BoardingMark
	err := r.db.WithContext(ctx).
		Joins("JOIN tickets ON tickets.id = boarding_marks.ticket_id").
		Where("tickets.trip_id = ?", tripID).
		Find(&marks).Error
	if err != nil {
		return nil, err
	}
	return marks, nil
}

func (r *boardingRepository) CheckIfMarked(ctx context.Context, ticketID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.BoardingMark{}).
		Where("ticket_id = ?", ticketID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
