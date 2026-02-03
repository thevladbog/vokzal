package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONB тип для PostgreSQL
type JSONB []byte

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

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

// Route модель маршрута
type Route struct {
	ID         string    `gorm:"type:uuid;primary_key" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Stops      JSONB     `gorm:"type:jsonb;not null" json:"stops"`
	DistanceKm float64   `gorm:"type:decimal(8,2)" json:"distance_km"`
	DurationMin int      `gorm:"type:integer" json:"duration_min"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Schedule модель расписания
type Schedule struct {
	ID            string    `gorm:"type:uuid;primary_key" json:"id"`
	RouteID       string    `gorm:"type:uuid;not null;index" json:"route_id"`
	DepartureTime string    `gorm:"type:time;not null" json:"departure_time"`
	DaysOfWeek    JSONB     `gorm:"type:jsonb;not null" json:"days_of_week"`
	Platform      *string   `gorm:"type:varchar(10)" json:"platform,omitempty"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	Route Route `gorm:"foreignKey:RouteID" json:"route,omitempty"`
}

// Trip модель рейса
type Trip struct {
	ID              string     `gorm:"type:uuid;primary_key" json:"id"`
	ScheduleID      string     `gorm:"type:uuid;not null;index" json:"schedule_id"`
	Date            string     `gorm:"type:date;not null;index" json:"date"`
	Status          string     `gorm:"type:varchar(20);not null;default:'scheduled'" json:"status"`
	DelayMinutes    int        `gorm:"default:0" json:"delay_minutes"`
	Platform        *string    `gorm:"type:varchar(10)" json:"platform,omitempty"`
	DepartureActual *time.Time `json:"departure_actual,omitempty"`
	ArrivalActual   *time.Time `json:"arrival_actual,omitempty"`
	BusID           *string    `gorm:"type:uuid" json:"bus_id,omitempty"`
	DriverID        *string    `gorm:"type:uuid" json:"driver_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	
	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
}

// Stop информация об остановке
type Stop struct {
	StationID        string `json:"station_id"`
	Order            int    `json:"order"`
	ArrivalOffsetMin int    `json:"arrival_offset_min"`
}

func (Route) TableName() string {
	return "routes"
}

func (Schedule) TableName() string {
	return "schedules"
}

func (Trip) TableName() string {
	return "trips"
}

func (r *Route) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// ParseStops парсит JSONB stops в []Stop
func (r *Route) ParseStops() ([]Stop, error) {
	var stops []Stop
	if err := json.Unmarshal(r.Stops, &stops); err != nil {
		return nil, err
	}
	return stops, nil
}

// ParseDaysOfWeek парсит JSONB days_of_week в []int
func (s *Schedule) ParseDaysOfWeek() ([]int, error) {
	var days []int
	if err := json.Unmarshal(s.DaysOfWeek, &days); err != nil {
		return nil, err
	}
	return days, nil
}
