// Package repository содержит слой доступа к данным маршрутов, расписаний и рейсов.
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/vokzal-tech/schedule-service/internal/models"
)

var (
	// ErrStationNotFound возвращается, когда станция не найдена.
	ErrStationNotFound = errors.New("station not found")
	// ErrRouteNotFound возвращается, когда маршрут не найден.
	ErrRouteNotFound = errors.New("route not found")
	// ErrScheduleNotFound возвращается, когда расписание не найдено.
	ErrScheduleNotFound = errors.New("schedule not found")
	// ErrTripNotFound возвращается, когда рейс не найден.
	ErrTripNotFound = errors.New("trip not found")
	// ErrBusNotFound возвращается, когда автобус не найден.
	ErrBusNotFound = errors.New("bus not found")
	// ErrDriverNotFound возвращается, когда водитель не найден.
	ErrDriverNotFound = errors.New("driver not found")
)

// BusRepository — интерфейс репозитория автобусов.
type BusRepository interface {
	Create(ctx context.Context, bus *models.Bus) error
	FindByID(ctx context.Context, id string) (*models.Bus, error)
	FindAll(ctx context.Context, stationID *string, status *string) ([]*models.Bus, error)
	Update(ctx context.Context, bus *models.Bus) error
	Delete(ctx context.Context, id string) error
}

// DriverRepository — интерфейс репозитория водителей.
type DriverRepository interface {
	Create(ctx context.Context, driver *models.Driver) error
	FindByID(ctx context.Context, id string) (*models.Driver, error)
	FindAll(ctx context.Context, stationID *string) ([]*models.Driver, error)
	Update(ctx context.Context, driver *models.Driver) error
	Delete(ctx context.Context, id string) error
}

// StationRepository — интерфейс репозитория станций.
type StationRepository interface {
	Create(ctx context.Context, station *models.Station) error
	FindByID(ctx context.Context, id string) (*models.Station, error)
	FindAll(ctx context.Context, city string, active *bool) ([]*models.Station, error)
	Update(ctx context.Context, station *models.Station) error
	Delete(ctx context.Context, id string) error
}

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

type stationRepository struct {
	db *gorm.DB
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

type busRepository struct {
	db *gorm.DB
}

type driverRepository struct {
	db *gorm.DB
}

// NewBusRepository создаёт репозиторий автобусов.
func NewBusRepository(db *gorm.DB) BusRepository {
	return &busRepository{db: db}
}

// NewDriverRepository создаёт репозиторий водителей.
func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepository{db: db}
}

// NewStationRepository создаёт репозиторий станций.
func NewStationRepository(db *gorm.DB) StationRepository {
	return &stationRepository{db: db}
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

// Station repository implementation.
func (r *stationRepository) Create(ctx context.Context, station *models.Station) error {
	return r.db.WithContext(ctx).Create(station).Error
}

func (r *stationRepository) FindByID(ctx context.Context, id string) (*models.Station, error) {
	var station models.Station
	if err := r.db.WithContext(ctx).First(&station, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStationNotFound
		}
		return nil, err
	}
	return &station, nil
}

func (r *stationRepository) FindAll(ctx context.Context, city string, _ *bool) ([]*models.Station, error) {
	var stations []*models.Station
	query := r.db.WithContext(ctx)
	if city != "" {
		query = query.Where("name ILIKE ? OR address ILIKE ?", "%"+city+"%", "%"+city+"%")
	}
	if err := query.Order("name ASC").Find(&stations).Error; err != nil {
		return nil, err
	}
	return stations, nil
}

func (r *stationRepository) Update(ctx context.Context, station *models.Station) error {
	return r.db.WithContext(ctx).Save(station).Error
}

func (r *stationRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Station{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStationNotFound
	}
	return nil
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

// Bus repository implementation.
func (r *busRepository) Create(ctx context.Context, bus *models.Bus) error {
	return r.db.WithContext(ctx).Create(bus).Error
}

func (r *busRepository) FindByID(ctx context.Context, id string) (*models.Bus, error) {
	var bus models.Bus
	if err := r.db.WithContext(ctx).First(&bus, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBusNotFound
		}
		return nil, err
	}
	return &bus, nil
}

func (r *busRepository) FindAll(ctx context.Context, stationID, status *string) ([]*models.Bus, error) {
	var buses []*models.Bus
	query := r.db.WithContext(ctx)
	if stationID != nil && *stationID != "" {
		query = query.Where("station_id = ?", *stationID)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if err := query.Order("plate_number ASC").Find(&buses).Error; err != nil {
		return nil, err
	}
	return buses, nil
}

func (r *busRepository) Update(ctx context.Context, bus *models.Bus) error {
	return r.db.WithContext(ctx).Save(bus).Error
}

func (r *busRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Bus{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrBusNotFound
	}
	return nil
}

// Driver repository implementation.
func (r *driverRepository) Create(ctx context.Context, driver *models.Driver) error {
	return r.db.WithContext(ctx).Create(driver).Error
}

func (r *driverRepository) FindByID(ctx context.Context, id string) (*models.Driver, error) {
	var driver models.Driver
	if err := r.db.WithContext(ctx).First(&driver, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDriverNotFound
		}
		return nil, err
	}
	return &driver, nil
}

func (r *driverRepository) FindAll(ctx context.Context, stationID *string) ([]*models.Driver, error) {
	var drivers []*models.Driver
	query := r.db.WithContext(ctx)
	if stationID != nil && *stationID != "" {
		query = query.Where("station_id = ?", *stationID)
	}
	if err := query.Order("full_name ASC").Find(&drivers).Error; err != nil {
		return nil, err
	}
	return drivers, nil
}

func (r *driverRepository) Update(ctx context.Context, driver *models.Driver) error {
	return r.db.WithContext(ctx).Save(driver).Error
}

func (r *driverRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Driver{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDriverNotFound
	}
	return nil
}
