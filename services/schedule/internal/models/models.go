// Package models содержит модели данных сервиса расписания.
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONB — тип для PostgreSQL JSONB.
type JSONB []byte

// Value реализует driver.Valuer для JSONB.
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan реализует sql.Scanner для JSONB.
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONB value")
	}
	*j = s
	return nil
}

// Station — модель станции (автовокзала).
//
//nolint:govet // fieldalignment: GORM/JSON order
type Station struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Code      string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Address   string    `gorm:"type:text" json:"address,omitempty"`
	Timezone  string    `gorm:"type:varchar(50);default:'Europe/Moscow'" json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName возвращает имя таблицы для Station.
func (Station) TableName() string {
	return "stations"
}

// BeforeCreate генерирует UUID для Station.
func (s *Station) BeforeCreate(_ *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// Route — модель маршрута.
type Route struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ID          string    `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Stops       JSONB     `gorm:"type:jsonb;not null" json:"stops"`
	DistanceKm  float64   `gorm:"type:decimal(8,2)" json:"distance_km"`
	DurationMin int       `gorm:"type:integer" json:"duration_min"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}

// Schedule — модель расписания.
type Schedule struct {
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Platform      *string   `gorm:"type:varchar(10)" json:"platform,omitempty"`
	ID            string    `gorm:"type:uuid;primary_key" json:"id"`
	RouteID       string    `gorm:"type:uuid;not null;index" json:"route_id"`
	DepartureTime string    `gorm:"type:time;not null" json:"departure_time"`
	DaysOfWeek    JSONB     `gorm:"type:jsonb;not null" json:"days_of_week"`
	Route         Route     `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
}

// Trip — модель рейса.
type Trip struct {
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedAt       time.Time  `json:"created_at"`
	ArrivalActual   *time.Time `json:"arrival_actual,omitempty"`
	DriverID        *string    `gorm:"type:uuid" json:"driver_id,omitempty"`
	Platform        *string    `gorm:"type:varchar(10)" json:"platform,omitempty"`
	DepartureActual *time.Time `json:"departure_actual,omitempty"`
	BusID           *string    `gorm:"type:uuid" json:"bus_id,omitempty"`
	ID              string     `gorm:"type:uuid;primary_key" json:"id"`
	Status          string     `gorm:"type:varchar(20);not null;default:'scheduled'" json:"status"`
	Date            string     `gorm:"type:date;not null;index" json:"date"`
	ScheduleID      string     `gorm:"type:uuid;not null;index" json:"schedule_id"`
	Schedule        Schedule   `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	DelayMinutes    int        `gorm:"default:0" json:"delay_minutes"`
}

// Stop — информация об остановке.
type Stop struct {
	StationID        string `json:"station_id"`
	Order            int    `json:"order"`
	ArrivalOffsetMin int    `json:"arrival_offset_min"`
}

// Bus — модель автобуса.
//
//nolint:govet // fieldalignment: field order kept for GORM/JSON readability
type Bus struct {
	Capacity    int       `gorm:"type:integer;not null" json:"capacity"`
	ID          string    `gorm:"type:uuid;primary_key" json:"id"`
	PlateNumber string    `gorm:"type:varchar(12);uniqueIndex;not null" json:"plate_number"`
	Model       string    `gorm:"type:varchar(50);not null" json:"model"`
	Status      string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	StationID   string    `gorm:"type:uuid;not null;index" json:"station_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Driver — модель водителя.
//
//nolint:govet // fieldalignment: field order kept for GORM/JSON readability
type Driver struct {
	ExperienceYears *int      `gorm:"type:integer" json:"experience_years,omitempty"`
	Phone           *string   `gorm:"type:varchar(15)" json:"phone,omitempty"`
	ID              string    `gorm:"type:uuid;primary_key" json:"id"`
	FullName        string    `gorm:"type:varchar(100);not null" json:"full_name"`
	LicenseNumber   string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"license_number"`
	StationID       string    `gorm:"type:uuid;not null;index" json:"station_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName возвращает имя таблицы для GORM (Route).
func (Route) TableName() string {
	return "routes"
}

// TableName возвращает имя таблицы для GORM (Bus).
func (Bus) TableName() string {
	return "buses"
}

// TableName возвращает имя таблицы для GORM (Driver).
func (Driver) TableName() string {
	return "drivers"
}

// BeforeCreate генерирует UUID для Bus.
func (b *Bus) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate генерирует UUID для Driver.
func (d *Driver) BeforeCreate(_ *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

// TableName возвращает имя таблицы для GORM (Schedule).
func (Schedule) TableName() string {
	return "schedules"
}

// TableName возвращает имя таблицы для GORM (Trip).
func (Trip) TableName() string {
	return "trips"
}

// BeforeCreate генерирует UUID для новой записи (Route).
func (r *Route) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate генерирует UUID для новой записи (Schedule).
func (s *Schedule) BeforeCreate(_ *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate генерирует UUID для новой записи (Trip).
func (t *Trip) BeforeCreate(_ *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// ParseStops парсит JSONB stops в []Stop.
func (r *Route) ParseStops() ([]Stop, error) {
	var stops []Stop
	if err := json.Unmarshal(r.Stops, &stops); err != nil {
		return nil, err
	}
	return stops, nil
}

// ParseDaysOfWeek парсит JSONB days_of_week в []int.
func (s *Schedule) ParseDaysOfWeek() ([]int, error) {
	var days []int
	if err := json.Unmarshal(s.DaysOfWeek, &days); err != nil {
		return nil, err
	}
	return days, nil
}
