// Package service содержит бизнес-логику маршрутов, расписаний и рейсов.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/vokzal-tech/schedule-service/internal/models"
	"github.com/vokzal-tech/schedule-service/internal/repository"
)

// ErrStationNotFound возвращается, когда station_id в запросе не ссылается на существующую станцию.
var ErrStationNotFound = errors.New("station not found")

// ErrInvalidCapacity возвращается, когда capacity меньше 1 (CreateBus/UpdateBus).
var ErrInvalidCapacity = errors.New("capacity must be at least 1")

// ErrTripNotFound возвращается, когда рейс не найден (UpdateTrip и др.).
var ErrTripNotFound = errors.New("trip not found")

// ErrBusNotFound возвращается, когда автобус не найден (UpdateBus, DeleteBus).
var ErrBusNotFound = errors.New("bus not found")

// ErrDriverNotFound возвращается, когда водитель не найден (UpdateDriver, DeleteDriver).
var ErrDriverNotFound = errors.New("driver not found")

// ScheduleService — интерфейс сервиса расписания (станции, маршруты, расписания, рейсы).
type ScheduleService interface {
	// Stations
	CreateStation(ctx context.Context, req *CreateStationRequest) (*models.Station, error)
	GetStation(ctx context.Context, id string) (*models.Station, error)
	ListStations(ctx context.Context, city string) ([]*models.Station, error)
	UpdateStation(ctx context.Context, id string, req *UpdateStationRequest) (*models.Station, error)
	DeleteStation(ctx context.Context, id string) error

	// Routes
	CreateRoute(ctx context.Context, req *CreateRouteRequest) (*models.Route, error)
	GetRoute(ctx context.Context, id string) (*models.Route, error)
	ListRoutes(ctx context.Context, activeOnly bool) ([]*models.Route, error)
	UpdateRoute(ctx context.Context, id string, req *UpdateRouteRequest) (*models.Route, error)
	DeleteRoute(ctx context.Context, id string) error

	// Schedules
	CreateSchedule(ctx context.Context, req *CreateScheduleRequest) (*models.Schedule, error)
	GetSchedule(ctx context.Context, id string) (*models.Schedule, error)
	ListSchedulesByRoute(ctx context.Context, routeID string) ([]*models.Schedule, error)
	UpdateSchedule(ctx context.Context, id string, req *UpdateScheduleRequest) (*models.Schedule, error)
	DeleteSchedule(ctx context.Context, id string) error

	// Trips
	CreateTrip(ctx context.Context, req *CreateTripRequest) (*models.Trip, error)
	GetTrip(ctx context.Context, id string) (*models.Trip, error)
	ListTripsByDate(ctx context.Context, date string) ([]*models.Trip, error)
	UpdateTripStatus(ctx context.Context, id string, status string, delayMinutes int) (*models.Trip, error)
	UpdateTrip(ctx context.Context, id string, req *UpdateTripRequest) (*models.Trip, error)
	GenerateTripsForSchedule(ctx context.Context, scheduleID string, fromDate, toDate time.Time) error
	GetDashboardStats(ctx context.Context, date string) (*DashboardStats, error)

	// Buses
	CreateBus(ctx context.Context, req *CreateBusRequest) (*models.Bus, error)
	GetBus(ctx context.Context, id string) (*models.Bus, error)
	ListBuses(ctx context.Context, stationID *string, status *string) ([]*models.Bus, error)
	UpdateBus(ctx context.Context, id string, req *UpdateBusRequest) (*models.Bus, error)
	DeleteBus(ctx context.Context, id string) error

	// Drivers
	CreateDriver(ctx context.Context, req *CreateDriverRequest) (*models.Driver, error)
	GetDriver(ctx context.Context, id string) (*models.Driver, error)
	ListDrivers(ctx context.Context, stationID *string) ([]*models.Driver, error)
	UpdateDriver(ctx context.Context, id string, req *UpdateDriverRequest) (*models.Driver, error)
	DeleteDriver(ctx context.Context, id string) error
}

type scheduleService struct {
	stationRepo  repository.StationRepository
	routeRepo    repository.RouteRepository
	scheduleRepo repository.ScheduleRepository
	tripRepo     repository.TripRepository
	busRepo      repository.BusRepository
	driverRepo   repository.DriverRepository
	natsConn     *nats.Conn
	logger       *zap.Logger
}

// CreateStationRequest — запрос на создание станции.
type CreateStationRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Address  string `json:"address"`
	Timezone string `json:"timezone"`
}

// UpdateStationRequest — запрос на обновление станции.
type UpdateStationRequest struct {
	Name     *string `json:"name"`
	Code     *string `json:"code"`
	Address  *string `json:"address"`
	Timezone *string `json:"timezone"`
}

// CreateRouteRequest — запрос на создание маршрута.
type CreateRouteRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Stops       []map[string]interface{} `json:"stops" binding:"required"`
	DistanceKm  float64                  `json:"distance_km"`
	DurationMin int                      `json:"duration_min"`
}

// UpdateRouteRequest — запрос на обновление маршрута.
type UpdateRouteRequest struct {
	Name        *string                  `json:"name"`
	DistanceKm  *float64                 `json:"distance_km"`
	DurationMin *int                     `json:"duration_min"`
	IsActive    *bool                    `json:"is_active"`
	Stops       []map[string]interface{} `json:"stops"`
}

// CreateScheduleRequest — запрос на создание расписания.
type CreateScheduleRequest struct {
	RouteID       string `json:"route_id" binding:"required"`
	DepartureTime string `json:"departure_time" binding:"required"`
	Platform      string `json:"platform"`
	DaysOfWeek    []int  `json:"days_of_week" binding:"required"`
}

// UpdateScheduleRequest — запрос на обновление расписания.
type UpdateScheduleRequest struct {
	DepartureTime *string `json:"departure_time"`
	DaysOfWeek    *[]int  `json:"days_of_week"`
	Platform      *string `json:"platform"`
	IsActive      *bool   `json:"is_active"`
}

// CreateTripRequest — запрос на создание рейса.
type CreateTripRequest struct {
	Platform   *string `json:"platform"`
	BusID      *string `json:"bus_id"`
	DriverID   *string `json:"driver_id"`
	ScheduleID string  `json:"schedule_id" binding:"required"`
	Date       string  `json:"date" binding:"required"`
}

// UpdateTripRequest — запрос на обновление рейса (перрон, автобус, водитель).
type UpdateTripRequest struct {
	Platform *string `json:"platform"`
	BusID    *string `json:"bus_id"`
	DriverID *string `json:"driver_id"`
}

// CreateBusRequest — запрос на создание автобуса.
// Поля упорядочены по убыванию размера (strings, затем int) для выравнивания и минимизации padding.
type CreateBusRequest struct {
	PlateNumber string `json:"plate_number" binding:"required"`
	Model       string `json:"model" binding:"required"`
	StationID   string `json:"station_id" binding:"required"`
	Status      string `json:"status"`
	Capacity    int    `json:"capacity" binding:"required"`
}

// UpdateBusRequest — запрос на обновление автобуса.
type UpdateBusRequest struct {
	PlateNumber *string `json:"plate_number"`
	Model       *string `json:"model"`
	Capacity    *int    `json:"capacity"`
	Status      *string `json:"status"`
}

// CreateDriverRequest — запрос на создание водителя.
type CreateDriverRequest struct {
	FullName        string  `json:"full_name" binding:"required"`
	LicenseNumber   string  `json:"license_number" binding:"required"`
	ExperienceYears *int    `json:"experience_years"`
	Phone           *string `json:"phone"`
	StationID       string  `json:"station_id" binding:"required"`
}

// UpdateDriverRequest — запрос на обновление водителя.
type UpdateDriverRequest struct {
	FullName        *string `json:"full_name"`
	LicenseNumber   *string `json:"license_number"`
	ExperienceYears *int    `json:"experience_years"`
	Phone           *string `json:"phone"`
}

// DashboardStats — статистика для дашборда (рейсы за дату).
type DashboardStats struct {
	TripsTotal     int `json:"trips_total"`
	TripsScheduled int `json:"trips_scheduled"`
	TripsBoarding  int `json:"trips_boarding"`
	TripsDeparted  int `json:"trips_departed"`
	TripsCancelled int `json:"trips_cancelled"` //nolint:misspell // British spelling; golangci-lint misspell (locale US) flags it
	TripsDelayed   int `json:"trips_delayed"`
	TripsArrived   int `json:"trips_arrived"`
	// TotalCapacity — сумма вместимостей автобусов по рейсам за дату (рейсы без автобуса считаются как 40 мест).
	TotalCapacity int `json:"total_capacity"`
}

// NewScheduleService создаёт сервис расписания.
func NewScheduleService(
	stationRepo repository.StationRepository,
	routeRepo repository.RouteRepository,
	scheduleRepo repository.ScheduleRepository,
	tripRepo repository.TripRepository,
	busRepo repository.BusRepository,
	driverRepo repository.DriverRepository,
	natsConn *nats.Conn,
	logger *zap.Logger,
) ScheduleService {
	return &scheduleService{
		stationRepo:  stationRepo,
		routeRepo:    routeRepo,
		scheduleRepo: scheduleRepo,
		tripRepo:     tripRepo,
		busRepo:      busRepo,
		driverRepo:   driverRepo,
		natsConn:     natsConn,
		logger:       logger,
	}
}

// CreateStation создаёт станцию.
func (s *scheduleService) CreateStation(ctx context.Context, req *CreateStationRequest) (*models.Station, error) {
	tz := req.Timezone
	if tz == "" {
		tz = "Europe/Moscow"
	}
	station := &models.Station{
		Name:     req.Name,
		Code:     req.Code,
		Address:  req.Address,
		Timezone: tz,
	}
	if err := s.stationRepo.Create(ctx, station); err != nil {
		return nil, err
	}
	s.logger.Info("Station created", zap.String("station_id", station.ID), zap.String("name", station.Name))
	return station, nil
}

func (s *scheduleService) GetStation(ctx context.Context, id string) (*models.Station, error) {
	return s.stationRepo.FindByID(ctx, id)
}

func (s *scheduleService) ListStations(ctx context.Context, city string) ([]*models.Station, error) {
	return s.stationRepo.FindAll(ctx, city, nil)
}

func (s *scheduleService) UpdateStation(ctx context.Context, id string, req *UpdateStationRequest) (*models.Station, error) {
	station, err := s.stationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		station.Name = *req.Name
	}
	if req.Code != nil {
		station.Code = *req.Code
	}
	if req.Address != nil {
		station.Address = *req.Address
	}
	if req.Timezone != nil {
		station.Timezone = *req.Timezone
	}
	if err := s.stationRepo.Update(ctx, station); err != nil {
		return nil, err
	}
	return station, nil
}

func (s *scheduleService) DeleteStation(ctx context.Context, id string) error {
	return s.stationRepo.Delete(ctx, id)
}

// CreateRoute создаёт маршрут.
func (s *scheduleService) CreateRoute(ctx context.Context, req *CreateRouteRequest) (*models.Route, error) {
	stopsJSON, err := json.Marshal(req.Stops)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stops: %w", err)
	}

	route := &models.Route{
		Name:        req.Name,
		Stops:       models.JSONB(stopsJSON),
		DistanceKm:  req.DistanceKm,
		DurationMin: req.DurationMin,
		IsActive:    true,
	}

	if err := s.routeRepo.Create(ctx, route); err != nil {
		return nil, err
	}

	s.logger.Info("Route created", zap.String("route_id", route.ID), zap.String("name", route.Name))
	return route, nil
}

func (s *scheduleService) GetRoute(ctx context.Context, id string) (*models.Route, error) {
	return s.routeRepo.FindByID(ctx, id)
}

func (s *scheduleService) ListRoutes(ctx context.Context, activeOnly bool) ([]*models.Route, error) {
	var isActive *bool
	if activeOnly {
		val := true
		isActive = &val
	}
	return s.routeRepo.FindAll(ctx, isActive)
}

func (s *scheduleService) UpdateRoute(ctx context.Context, id string, req *UpdateRouteRequest) (*models.Route, error) {
	route, err := s.routeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		route.Name = *req.Name
	}
	if req.Stops != nil {
		stopsJSON, err := json.Marshal(req.Stops)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal stops: %w", err)
		}
		route.Stops = models.JSONB(stopsJSON)
	}
	if req.DistanceKm != nil {
		route.DistanceKm = *req.DistanceKm
	}
	if req.DurationMin != nil {
		route.DurationMin = *req.DurationMin
	}
	if req.IsActive != nil {
		route.IsActive = *req.IsActive
	}

	if err := s.routeRepo.Update(ctx, route); err != nil {
		return nil, err
	}

	return route, nil
}

func (s *scheduleService) DeleteRoute(ctx context.Context, id string) error {
	return s.routeRepo.Delete(ctx, id)
}

// CreateSchedule создаёт расписание.
func (s *scheduleService) CreateSchedule(ctx context.Context, req *CreateScheduleRequest) (*models.Schedule, error) {
	daysJSON, err := json.Marshal(req.DaysOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal days: %w", err)
	}

	schedule := &models.Schedule{
		RouteID:       req.RouteID,
		DepartureTime: req.DepartureTime,
		DaysOfWeek:    models.JSONB(daysJSON),
		IsActive:      true,
	}

	if req.Platform != "" {
		schedule.Platform = &req.Platform
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	s.logger.Info("Schedule created", zap.String("schedule_id", schedule.ID))
	return schedule, nil
}

func (s *scheduleService) GetSchedule(ctx context.Context, id string) (*models.Schedule, error) {
	return s.scheduleRepo.FindByID(ctx, id)
}

func (s *scheduleService) ListSchedulesByRoute(ctx context.Context, routeID string) ([]*models.Schedule, error) {
	return s.scheduleRepo.FindByRouteID(ctx, routeID)
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id string, req *UpdateScheduleRequest) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.DepartureTime != nil {
		schedule.DepartureTime = *req.DepartureTime
	}
	if req.DaysOfWeek != nil {
		daysJSON, err := json.Marshal(*req.DaysOfWeek)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal days: %w", err)
		}
		schedule.DaysOfWeek = models.JSONB(daysJSON)
	}
	if req.Platform != nil {
		schedule.Platform = req.Platform
	}
	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, id string) error {
	return s.scheduleRepo.Delete(ctx, id)
}

// CreateTrip создаёт рейс.
func (s *scheduleService) CreateTrip(ctx context.Context, req *CreateTripRequest) (*models.Trip, error) {
	// Проверить, не существует ли уже рейс
	existing, err := s.tripRepo.FindByScheduleAndDate(ctx, req.ScheduleID, req.Date)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("trip already exists for this schedule and date")
	}

	trip := &models.Trip{
		ScheduleID: req.ScheduleID,
		Date:       req.Date,
		Status:     "scheduled",
		Platform:   req.Platform,
		BusID:      req.BusID,
		DriverID:   req.DriverID,
	}

	if err := s.tripRepo.Create(ctx, trip); err != nil {
		return nil, err
	}

	// Отправить событие в NATS
	s.publishTripEvent("trip.created", trip)

	s.logger.Info("Trip created", zap.String("trip_id", trip.ID), zap.String("date", trip.Date))
	return trip, nil
}

func (s *scheduleService) GetTrip(ctx context.Context, id string) (*models.Trip, error) {
	return s.tripRepo.FindByID(ctx, id)
}

func (s *scheduleService) ListTripsByDate(ctx context.Context, date string) ([]*models.Trip, error) {
	return s.tripRepo.FindByDate(ctx, date)
}

func (s *scheduleService) UpdateTripStatus(ctx context.Context, id, status string, delayMinutes int) (*models.Trip, error) {
	trip, err := s.tripRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	trip.Status = status
	trip.DelayMinutes = delayMinutes

	if status == "departed" && trip.DepartureActual == nil {
		now := time.Now()
		trip.DepartureActual = &now
	}

	if status == "arrived" && trip.ArrivalActual == nil {
		now := time.Now()
		trip.ArrivalActual = &now
	}

	if err := s.tripRepo.Update(ctx, trip); err != nil {
		return nil, err
	}

	// Отправить событие
	s.publishTripEvent("trip.status_changed", trip)

	s.logger.Info("Trip status updated",
		zap.String("trip_id", trip.ID),
		zap.String("status", status),
		zap.Int("delay", delayMinutes))

	return trip, nil
}

func (s *scheduleService) UpdateTrip(ctx context.Context, id string, req *UpdateTripRequest) (*models.Trip, error) {
	trip, err := s.tripRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTripNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, fmt.Errorf("find trip: %w", err)
	}
	if req.Platform != nil {
		trip.Platform = req.Platform
	}
	if req.BusID != nil {
		trip.BusID = req.BusID
	}
	if req.DriverID != nil {
		trip.DriverID = req.DriverID
	}
	if err := s.tripRepo.Update(ctx, trip); err != nil {
		return nil, fmt.Errorf("update trip: %w", err)
	}
	s.publishTripEvent("trip.updated", trip)
	return trip, nil
}

func (s *scheduleService) CreateBus(ctx context.Context, req *CreateBusRequest) (*models.Bus, error) {
	if req.Capacity < 1 {
		return nil, ErrInvalidCapacity
	}
	if _, err := s.stationRepo.FindByID(ctx, req.StationID); err != nil {
		if errors.Is(err, repository.ErrStationNotFound) {
			return nil, ErrStationNotFound
		}
		s.logger.Error("CreateBus: find station failed", zap.String("station_id", req.StationID), zap.Error(err))
		return nil, fmt.Errorf("find station: %w", err)
	}
	status := req.Status
	if status == "" {
		status = "active"
	}
	bus := &models.Bus{
		PlateNumber: req.PlateNumber,
		Model:       req.Model,
		Capacity:    req.Capacity,
		StationID:   req.StationID,
		Status:      status,
	}
	if err := s.busRepo.Create(ctx, bus); err != nil {
		s.logger.Error("CreateBus: create failed", zap.String("plate_number", req.PlateNumber), zap.String("station_id", req.StationID), zap.Error(err))
		return nil, fmt.Errorf("create bus: %w", err)
	}
	s.logger.Info("Bus created", zap.String("bus_id", bus.ID), zap.String("plate_number", bus.PlateNumber), zap.String("station_id", bus.StationID), zap.String("status", bus.Status))
	return bus, nil
}

func (s *scheduleService) GetBus(ctx context.Context, id string) (*models.Bus, error) {
	bus, err := s.busRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("GetBus failed", zap.String("bus_id", id), zap.Error(err))
		return nil, fmt.Errorf("find bus: %w", err)
	}
	return bus, nil
}

func (s *scheduleService) ListBuses(ctx context.Context, stationID, status *string) ([]*models.Bus, error) {
	buses, err := s.busRepo.FindAll(ctx, stationID, status)
	if err != nil {
		s.logger.Error("ListBuses failed", zap.Any("station_id", stationID), zap.Any("status", status), zap.Error(err))
		return nil, fmt.Errorf("list buses: %w", err)
	}
	return buses, nil
}

func (s *scheduleService) UpdateBus(ctx context.Context, id string, req *UpdateBusRequest) (*models.Bus, error) {
	s.logger.Debug("UpdateBus", zap.String("bus_id", id),
		zap.Any("plate_number", req.PlateNumber), zap.Any("model", req.Model),
		zap.Any("capacity", req.Capacity), zap.Any("status", req.Status))
	bus, err := s.busRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrBusNotFound) {
			return nil, ErrBusNotFound
		}
		s.logger.Error("UpdateBus: find bus failed", zap.String("bus_id", id), zap.Error(err))
		return nil, fmt.Errorf("find bus: %w", err)
	}
	if req.PlateNumber != nil {
		bus.PlateNumber = *req.PlateNumber
	}
	if req.Model != nil {
		bus.Model = *req.Model
	}
	if req.Capacity != nil {
		if *req.Capacity < 1 {
			return nil, ErrInvalidCapacity
		}
		bus.Capacity = *req.Capacity
	}
	if req.Status != nil {
		bus.Status = *req.Status
	}
	if err := s.busRepo.Update(ctx, bus); err != nil {
		s.logger.Error("UpdateBus: update failed", zap.String("bus_id", id), zap.Error(err))
		return nil, fmt.Errorf("update bus: %w", err)
	}
	s.logger.Info("Bus updated", zap.String("bus_id", id), zap.String("plate_number", bus.PlateNumber), zap.String("status", bus.Status))
	return bus, nil
}

func (s *scheduleService) DeleteBus(ctx context.Context, id string) error {
	if err := s.busRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrBusNotFound) {
			return ErrBusNotFound
		}
		s.logger.Error("DeleteBus failed", zap.String("bus_id", id), zap.Error(err))
		return fmt.Errorf("delete bus: %w", err)
	}
	s.logger.Info("Bus deleted", zap.String("bus_id", id))
	return nil
}

func (s *scheduleService) CreateDriver(ctx context.Context, req *CreateDriverRequest) (*models.Driver, error) {
	s.logger.Debug("CreateDriver", zap.String("station_id", req.StationID), zap.String("full_name", req.FullName), zap.String("license_number", req.LicenseNumber))
	if _, err := s.stationRepo.FindByID(ctx, req.StationID); err != nil {
		if errors.Is(err, repository.ErrStationNotFound) {
			return nil, ErrStationNotFound
		}
		s.logger.Error("CreateDriver failed: FindByID station", zap.String("station_id", req.StationID), zap.Error(err))
		return nil, fmt.Errorf("CreateDriver failed: FindByID station: %w", err)
	}
	driver := &models.Driver{
		FullName:        req.FullName,
		LicenseNumber:   req.LicenseNumber,
		ExperienceYears: req.ExperienceYears,
		Phone:           req.Phone,
		StationID:       req.StationID,
	}
	if err := s.driverRepo.Create(ctx, driver); err != nil {
		s.logger.Error("CreateDriver failed: Create", zap.String("station_id", req.StationID), zap.String("license_number", req.LicenseNumber), zap.Error(err))
		return nil, fmt.Errorf("CreateDriver failed: Create: %w", err)
	}
	s.logger.Info("Driver created", zap.String("driver_id", driver.ID), zap.String("full_name", driver.FullName), zap.String("station_id", driver.StationID))
	return driver, nil
}

func (s *scheduleService) GetDriver(ctx context.Context, id string) (*models.Driver, error) {
	s.logger.Debug("GetDriver", zap.String("driver_id", id))
	driver, err := s.driverRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("GetDriver failed: FindByID", zap.String("driver_id", id), zap.Error(err))
		return nil, fmt.Errorf("GetDriver failed: FindByID: %w", err)
	}
	return driver, nil
}

func (s *scheduleService) ListDrivers(ctx context.Context, stationID *string) ([]*models.Driver, error) {
	s.logger.Debug("ListDrivers", zap.Any("station_id", stationID))
	drivers, err := s.driverRepo.FindAll(ctx, stationID)
	if err != nil {
		s.logger.Error("ListDrivers failed: FindAll", zap.Any("station_id", stationID), zap.Error(err))
		return nil, fmt.Errorf("ListDrivers failed: FindAll: %w", err)
	}
	return drivers, nil
}

func (s *scheduleService) UpdateDriver(ctx context.Context, id string, req *UpdateDriverRequest) (*models.Driver, error) {
	s.logger.Debug("UpdateDriver", zap.String("driver_id", id),
		zap.Any("full_name", req.FullName), zap.Any("license_number", req.LicenseNumber),
		zap.Any("experience_years", req.ExperienceYears), zap.Any("phone", req.Phone))
	driver, err := s.driverRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			return nil, ErrDriverNotFound
		}
		s.logger.Error("UpdateDriver failed: FindByID", zap.String("driver_id", id), zap.Error(err))
		return nil, fmt.Errorf("UpdateDriver failed: FindByID: %w", err)
	}
	if req.FullName != nil {
		driver.FullName = *req.FullName
	}
	if req.LicenseNumber != nil {
		driver.LicenseNumber = *req.LicenseNumber
	}
	if req.ExperienceYears != nil {
		driver.ExperienceYears = req.ExperienceYears
	}
	if req.Phone != nil {
		driver.Phone = req.Phone
	}
	if err := s.driverRepo.Update(ctx, driver); err != nil {
		s.logger.Error("UpdateDriver failed: Update", zap.String("driver_id", id), zap.Error(err))
		return nil, fmt.Errorf("UpdateDriver failed: Update: %w", err)
	}
	s.logger.Info("Driver updated", zap.String("driver_id", id), zap.String("full_name", driver.FullName))
	return driver, nil
}

func (s *scheduleService) DeleteDriver(ctx context.Context, id string) error {
	s.logger.Debug("DeleteDriver", zap.String("driver_id", id))
	if err := s.driverRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			return ErrDriverNotFound
		}
		s.logger.Error("DeleteDriver failed: Delete", zap.String("driver_id", id), zap.Error(err))
		return fmt.Errorf("DeleteDriver failed: Delete: %w", err)
	}
	s.logger.Info("Driver deleted", zap.String("driver_id", id))
	return nil
}

const defaultCapacityPerTrip = 40

func (s *scheduleService) GetDashboardStats(ctx context.Context, date string) (*DashboardStats, error) {
	trips, err := s.tripRepo.FindByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("GetDashboardStats: find trips by date: %w", err)
	}
	stats := &DashboardStats{TripsTotal: len(trips)}
	var totalCapacity int
	for _, t := range trips {
		switch t.Status {
		case "scheduled":
			stats.TripsScheduled++
		case "boarding":
			stats.TripsBoarding++
		case "departed":
			stats.TripsDeparted++
		case "cancelled": //nolint:misspell // trip status; British spelling intentional
			stats.TripsCancelled++
		case "delayed":
			stats.TripsDelayed++
		case "arrived":
			stats.TripsArrived++
		}
		seats := defaultCapacityPerTrip
		if t.BusID != nil && *t.BusID != "" {
			bus, err := s.busRepo.FindByID(ctx, *t.BusID)
			if err == nil {
				seats = bus.Capacity
			}
		}
		totalCapacity += seats
	}
	stats.TotalCapacity = totalCapacity
	return stats, nil
}

func (s *scheduleService) GenerateTripsForSchedule(ctx context.Context, scheduleID string, fromDate, toDate time.Time) error {
	schedule, err := s.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	days, err := schedule.ParseDaysOfWeek()
	if err != nil {
		return fmt.Errorf("failed to parse days of week: %w", err)
	}

	for date := fromDate; !date.After(toDate); date = date.AddDate(0, 0, 1) {
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		// Проверить, входит ли день в расписание
		shouldCreate := false
		for _, day := range days {
			if day == weekday {
				shouldCreate = true
				break
			}
		}

		if shouldCreate {
			dateStr := date.Format("2006-01-02")
			existing, err := s.tripRepo.FindByScheduleAndDate(ctx, scheduleID, dateStr)
			if err != nil || existing != nil {
				continue
			}
			trip := &models.Trip{
				ScheduleID: scheduleID,
				Date:       dateStr,
				Status:     "scheduled",
				Platform:   schedule.Platform,
			}
			if err := s.tripRepo.Create(ctx, trip); err != nil {
				s.logger.Error("Failed to create trip", zap.Error(err), zap.String("date", dateStr))
				continue
			}
			s.publishTripEvent("trip.created", trip)
		}
	}

	return nil
}

func (s *scheduleService) publishTripEvent(subject string, trip *models.Trip) {
	data, err := json.Marshal(trip)
	if err != nil {
		s.logger.Error("Failed to marshal trip event", zap.Error(err))
		return
	}

	if err := s.natsConn.Publish(subject, data); err != nil {
		s.logger.Error("Failed to publish trip event", zap.Error(err), zap.String("subject", subject))
	}
}
