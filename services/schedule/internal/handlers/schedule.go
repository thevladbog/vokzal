// Package handlers содержит HTTP-обработчики API расписания.
//
//nolint:dupl // CRUD-обработчики для routes/schedules/trips намеренно однотипны.
package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vokzal-tech/schedule-service/internal/service"
)

// allowedBusStatuses — допустимые значения статуса автобуса (CreateBusRequest, UpdateBusRequest).
var allowedBusStatuses = map[string]bool{
	"active":         true,
	"maintenance":    true,
	"out_of_service": true,
}

func isValidBusStatus(status string) bool {
	return allowedBusStatuses[strings.TrimSpace(status)]
}

// ScheduleHandler — обработчик HTTP-запросов для маршрутов, расписаний и рейсов.
type ScheduleHandler struct {
	svc    service.ScheduleService
	logger *zap.Logger
}

// NewScheduleHandler создаёт обработчик расписания.
func NewScheduleHandler(svc service.ScheduleService, logger *zap.Logger) *ScheduleHandler {
	return &ScheduleHandler{
		svc:    svc,
		logger: logger,
	}
}

// CreateStation создаёт станцию.
func (h *ScheduleHandler) CreateStation(c *gin.Context) {
	var req service.CreateStationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	station, err := h.svc.CreateStation(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create station", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create station"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": station})
}

// GetStation возвращает станцию по ID.
func (h *ScheduleHandler) GetStation(c *gin.Context) {
	id := c.Param("id")
	station, err := h.svc.GetStation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": station})
}

// ListStations возвращает список станций.
func (h *ScheduleHandler) ListStations(c *gin.Context) {
	city := c.Query("city")
	stations, err := h.svc.ListStations(c.Request.Context(), city)
	if err != nil {
		h.logger.Error("Failed to list stations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list stations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stations})
}

// UpdateStation обновляет станцию.
func (h *ScheduleHandler) UpdateStation(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateStationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	station, err := h.svc.UpdateStation(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update station", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update station"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": station})
}

// DeleteStation удаляет станцию.
func (h *ScheduleHandler) DeleteStation(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteStation(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete station", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete station"})
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateRoute создаёт маршрут.
func (h *ScheduleHandler) CreateRoute(c *gin.Context) {
	var req service.CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.svc.CreateRoute(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create route", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create route"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": route})
}

// GetRoute возвращает маршрут по ID.
func (h *ScheduleHandler) GetRoute(c *gin.Context) {
	id := c.Param("id")
	route, err := h.svc.GetRoute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": route})
}

// ListRoutes возвращает список маршрутов.
func (h *ScheduleHandler) ListRoutes(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "false") == "true"
	routes, err := h.svc.ListRoutes(c.Request.Context(), activeOnly)
	if err != nil {
		h.logger.Error("Failed to list routes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list routes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": routes})
}

// UpdateRoute обновляет маршрут.
func (h *ScheduleHandler) UpdateRoute(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.svc.UpdateRoute(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update route", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update route"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": route})
}

// DeleteRoute удаляет маршрут.
func (h *ScheduleHandler) DeleteRoute(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRoute(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete route", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete route"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Route deleted"})
}

// CreateSchedule создаёт расписание.
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req service.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.svc.CreateSchedule(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create schedule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": schedule})
}

// GetSchedule возвращает расписание по ID.
func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	id := c.Param("id")
	schedule, err := h.svc.GetSchedule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schedule})
}

// ListSchedulesByRoute возвращает расписания по маршруту.
func (h *ScheduleHandler) ListSchedulesByRoute(c *gin.Context) {
	routeID := c.Query("route_id")
	if routeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "route_id is required"})
		return
	}

	schedules, err := h.svc.ListSchedulesByRoute(c.Request.Context(), routeID)
	if err != nil {
		h.logger.Error("Failed to list schedules", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list schedules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schedules})
}

// UpdateSchedule обновляет расписание.
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.svc.UpdateSchedule(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update schedule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schedule})
}

// DeleteSchedule удаляет расписание.
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteSchedule(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete schedule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}

// CreateTrip создаёт рейс.
func (h *ScheduleHandler) CreateTrip(c *gin.Context) {
	var req service.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip, err := h.svc.CreateTrip(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create trip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": trip})
}

// GetTrip возвращает рейс по ID.
func (h *ScheduleHandler) GetTrip(c *gin.Context) {
	id := c.Param("id")
	trip, err := h.svc.GetTrip(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trip})
}

// ListTripsByDate возвращает рейсы на дату.
func (h *ScheduleHandler) ListTripsByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required (YYYY-MM-DD)"})
		return
	}

	trips, err := h.svc.ListTripsByDate(c.Request.Context(), date)
	if err != nil {
		h.logger.Error("Failed to list trips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list trips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trips})
}

// UpdateTripStatus обновляет статус рейса.
func (h *ScheduleHandler) UpdateTripStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status       string `json:"status" binding:"required"`
		DelayMinutes int    `json:"delay_minutes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip, err := h.svc.UpdateTripStatus(c.Request.Context(), id, req.Status, req.DelayMinutes)
	if err != nil {
		h.logger.Error("Failed to update trip status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trip})
}

// GenerateTrips генерирует рейсы по расписанию на период.
func (h *ScheduleHandler) GenerateTrips(c *gin.Context) {
	var req struct {
		ScheduleID string `json:"schedule_id" binding:"required"`
		FromDate   string `json:"from_date" binding:"required"`
		ToDate     string `json:"to_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	from, err := parseDate(req.FromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_date format"})
		return
	}

	to, err := parseDate(req.ToDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_date format"})
		return
	}

	if err := h.svc.GenerateTripsForSchedule(c.Request.Context(), req.ScheduleID, from, to); err != nil {
		h.logger.Error("Failed to generate trips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate trips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trips generated successfully"})
}

// UpdateTrip обновляет рейс (перрон, автобус, водитель).
func (h *ScheduleHandler) UpdateTrip(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	trip, err := h.svc.UpdateTrip(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrTripNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
			return
		}
		h.logger.Error("Failed to update trip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": trip})
}

// CreateBus создаёт автобус.
func (h *ScheduleHandler) CreateBus(c *gin.Context) {
	var req service.CreateBusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Status != "" && !isValidBusStatus(req.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status: must be one of active, maintenance, out_of_service",
		})
		return
	}
	if req.Capacity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "capacity: must be at least 1"})
		return
	}
	bus, err := h.svc.CreateBus(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrStationNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "station_id: station not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidCapacity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "capacity: must be at least 1"})
			return
		}
		h.logger.Error("Failed to create bus", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bus"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": bus})
}

// GetBus возвращает автобус по ID.
func (h *ScheduleHandler) GetBus(c *gin.Context) {
	id := c.Param("id")
	bus, err := h.svc.GetBus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bus})
}

// ListBuses возвращает список автобусов.
func (h *ScheduleHandler) ListBuses(c *gin.Context) {
	stationID := c.Query("station_id")
	status := c.Query("status")
	if status != "" && !isValidBusStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status: must be one of active, maintenance, out_of_service",
		})
		return
	}
	var sid, st *string
	if stationID != "" {
		sid = &stationID
	}
	if status != "" {
		st = &status
	}
	buses, err := h.svc.ListBuses(c.Request.Context(), sid, st)
	if err != nil {
		h.logger.Error("Failed to list buses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list buses"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": buses})
}

// UpdateBus обновляет автобус.
func (h *ScheduleHandler) UpdateBus(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateBusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Status != nil && !isValidBusStatus(*req.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status: must be one of active, maintenance, out_of_service",
		})
		return
	}
	if req.Capacity != nil && *req.Capacity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "capacity: must be at least 1"})
		return
	}
	bus, err := h.svc.UpdateBus(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrBusNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidCapacity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "capacity: must be at least 1"})
			return
		}
		h.logger.Error("Failed to update bus", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bus})
}

// DeleteBus удаляет автобус.
func (h *ScheduleHandler) DeleteBus(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteBus(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrBusNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bus not found"})
			return
		}
		h.logger.Error("Failed to delete bus", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bus"})
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateDriver создаёт водителя.
func (h *ScheduleHandler) CreateDriver(c *gin.Context) {
	var req service.CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	driver, err := h.svc.CreateDriver(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrStationNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "station_id: station not found"})
			return
		}
		h.logger.Error("Failed to create driver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": driver})
}

// GetDriver возвращает водителя по ID.
func (h *ScheduleHandler) GetDriver(c *gin.Context) {
	id := c.Param("id")
	driver, err := h.svc.GetDriver(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": driver})
}

// ListDrivers возвращает список водителей.
func (h *ScheduleHandler) ListDrivers(c *gin.Context) {
	stationID := c.Query("station_id")
	var sid *string
	if stationID != "" {
		sid = &stationID
	}
	drivers, err := h.svc.ListDrivers(c.Request.Context(), sid)
	if err != nil {
		h.logger.Error("Failed to list drivers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list drivers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": drivers})
}

// UpdateDriver обновляет водителя.
func (h *ScheduleHandler) UpdateDriver(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	driver, err := h.svc.UpdateDriver(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrDriverNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
			return
		}
		h.logger.Error("Failed to update driver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": driver})
}

// DeleteDriver удаляет водителя.
func (h *ScheduleHandler) DeleteDriver(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteDriver(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrDriverNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
			return
		}
		h.logger.Error("Failed to delete driver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete driver"})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetDashboardStats возвращает статистику рейсов за дату для дашборда.
func (h *ScheduleHandler) GetDashboardStats(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	parsed, err := parseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (use YYYY-MM-DD)"})
		return
	}
	date := parsed.Format("2006-01-02")
	stats, err := h.svc.GetDashboardStats(c.Request.Context(), date)
	if err != nil {
		h.logger.Error("Failed to get dashboard stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
