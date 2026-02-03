// Package repository содержит слой доступа к данным маршрутов, расписаний и рейсов.
package repository

import (
	"context"
	"errors"

	"github.com/vokzal-tech/schedule-service/internal/models"
	"gorm.io/gorm"
)

var (
	// ErrRouteNotFound возвращается, когда маршрут не найден.
	ErrRouteNotFound = errors.New("route not found")
	// ErrScheduleNotFound возвращается, когда расписание не найдено.
	ErrScheduleNotFound = errors.New("schedule not found")
	// ErrTripNotFound возвращается, когда рейс не найден.
	ErrTripNotFound = errors.New("trip not found")
)

// RouteRepository — интерфейс репозитория маршрутов.
type RouteRepository interface {
	Create(ctx context.Context, route *models.Route) error
	FindByID(ctx context.Context, id string) (*models.Route, error)
	FindAll(ctx context.Context, isActive *bool) ([]*models.Route, error)
	Update(ctx context.Context, route *models.Route) error
	Delete(ctx context.Context, id string) error
}

// ScheduleRepository — интерфейс репозитория расписаний.
type ScheduleRepository interface {
	Create(ctx context.Context, schedule *models.Schedule) error
	FindByID(ctx context.Context, id string) (*models.Schedule, error)
	FindByRouteID(ctx context.Context, routeID string) ([]*models.Schedule, error)
	Update(ctx context.Context, schedule *models.Schedule) error
	Delete(ctx context.Context, id string) error
}

// TripRepository — интерфейс репозитория рейсов.
type TripRepository interface {
	Create(ctx context.Context, trip *models.Trip) error
	FindByID(ctx context.Context, id string) (*models.Trip, error)
	FindByDate(ctx context.Context, date string) ([]*models.Trip, error)
	FindByScheduleAndDate(ctx context.Context, scheduleID, date string) (*models.Trip, error)
	Update(ctx context.Context, trip *models.Trip) error
	Delete(ctx context.Context, id string) error
}

type routeRepository struct {
	db *gorm.DB
}

type scheduleRepository struct {
	db *gorm.DB
}

type tripRepository struct {
	db *gorm.DB
}

// NewRouteRepository создаёт репозиторий маршрутов.
func NewRouteRepository(db *gorm.DB) RouteRepository {
	return &routeRepository{db: db}
}

// NewScheduleRepository создаёт репозиторий расписаний.
func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

// NewTripRepository создаёт репозиторий рейсов.
func NewTripRepository(db *gorm.DB) TripRepository {
	return &tripRepository{db: db}
}

// Create создаёт маршрут.
func (r *routeRepository) Create(ctx context.Context, route *models.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *routeRepository) FindByID(ctx context.Context, id string) (*models.Route, error) {
	var route models.Route
	if err := r.db.WithContext(ctx).First(&route, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}
	return &route, nil
}

func (r *routeRepository) FindAll(ctx context.Context, isActive *bool) ([]*models.Route, error) {
	var routes []*models.Route
	query := r.db.WithContext(ctx)
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	if err := query.Find(&routes).Error; err != nil {
		return nil, err
	}
	return routes, nil
}

func (r *routeRepository) Update(ctx context.Context, route *models.Route) error {
	return r.db.WithContext(ctx).Save(route).Error
}

func (r *routeRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Route{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRouteNotFound
	}
	return nil
}

// Create создаёт расписание.
func (r *scheduleRepository) Create(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *scheduleRepository) FindByID(ctx context.Context, id string) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.WithContext(ctx).Preload("Route").First(&schedule, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) FindByRouteID(ctx context.Context, routeID string) ([]*models.Schedule, error) {
	var schedules []*models.Schedule
	if err := r.db.WithContext(ctx).Preload("Route").Where("route_id = ? AND is_active = ?", routeID, true).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *scheduleRepository) Update(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

func (r *scheduleRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Schedule{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrScheduleNotFound
	}
	return nil
}

// Create создаёт рейс.
func (r *tripRepository) Create(ctx context.Context, trip *models.Trip) error {
	return r.db.WithContext(ctx).Create(trip).Error
}

func (r *tripRepository) FindByID(ctx context.Context, id string) (*models.Trip, error) {
	var trip models.Trip
	if err := r.db.WithContext(ctx).Preload("Schedule.Route").First(&trip, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) FindByDate(ctx context.Context, date string) ([]*models.Trip, error) {
	var trips []*models.Trip
	if err := r.db.WithContext(ctx).Preload("Schedule.Route").Where("date = ?", date).Order("date ASC").Find(&trips).Error; err != nil {
		return nil, err
	}
	return trips, nil
}

func (r *tripRepository) FindByScheduleAndDate(ctx context.Context, scheduleID, date string) (*models.Trip, error) {
	var trip models.Trip
	if err := r.db.WithContext(ctx).Where("schedule_id = ? AND date = ?", scheduleID, date).First(&trip).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) Update(ctx context.Context, trip *models.Trip) error {
	return r.db.WithContext(ctx).Save(trip).Error
}

func (r *tripRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Trip{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTripNotFound
	}
	return nil
}
