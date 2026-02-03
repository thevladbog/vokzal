// Package service содержит бизнес-логику маршрутов, расписаний и рейсов.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/vokzal-tech/schedule-service/internal/models"
	"github.com/vokzal-tech/schedule-service/internal/repository"
)

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
	GenerateTripsForSchedule(ctx context.Context, scheduleID string, fromDate, toDate time.Time) error
}

type scheduleService struct {
	stationRepo  repository.StationRepository
	routeRepo    repository.RouteRepository
	scheduleRepo repository.ScheduleRepository
	tripRepo     repository.TripRepository
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

// NewScheduleService создаёт сервис расписания.
func NewScheduleService(
	stationRepo repository.StationRepository,
	routeRepo repository.RouteRepository,
	scheduleRepo repository.ScheduleRepository,
	tripRepo repository.TripRepository,
	natsConn *nats.Conn,
	logger *zap.Logger,
) ScheduleService {
	return &scheduleService{
		stationRepo:  stationRepo,
		routeRepo:    routeRepo,
		scheduleRepo: scheduleRepo,
		tripRepo:     tripRepo,
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
