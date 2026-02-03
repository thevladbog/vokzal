package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/schedule-service/internal/service"
	"go.uber.org/zap"
)

type ScheduleHandler struct {
	service service.ScheduleService
	logger  *zap.Logger
}

func NewScheduleHandler(service service.ScheduleService, logger *zap.Logger) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
		logger:  logger,
	}
}

// Routes
func (h *ScheduleHandler) CreateRoute(c *gin.Context) {
	var req service.CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.service.CreateRoute(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create route", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create route"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": route})
}

func (h *ScheduleHandler) GetRoute(c *gin.Context) {
	id := c.Param("id")
	route, err := h.service.GetRoute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": route})
}

func (h *ScheduleHandler) ListRoutes(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "false") == "true"
	routes, err := h.service.ListRoutes(c.Request.Context(), activeOnly)
	if err != nil {
		h.logger.Error("Failed to list routes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list routes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": routes})
}

func (h *ScheduleHandler) UpdateRoute(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.service.UpdateRoute(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update route", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update route"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": route})
}

func (h *ScheduleHandler) DeleteRoute(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteRoute(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete route", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete route"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Route deleted"})
}

// Schedules
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req service.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.service.CreateSchedule(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create schedule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": schedule})
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	id := c.Param("id")
	schedule, err := h.service.GetSchedule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schedule})
}

func (h *ScheduleHandler) ListSchedulesByRoute(c *gin.Context) {
	routeID := c.Query("route_id")
	if routeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "route_id is required"})
		return
	}

	schedules, err := h.service.ListSchedulesByRoute(c.Request.Context(), routeID)
	if err != nil {
		h.logger.Error("Failed to list schedules", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list schedules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schedules})
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.service.UpdateSchedule(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update schedule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schedule})
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteSchedule(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete schedule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}

// Trips
func (h *ScheduleHandler) CreateTrip(c *gin.Context) {
	var req service.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip, err := h.service.CreateTrip(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create trip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": trip})
}

func (h *ScheduleHandler) GetTrip(c *gin.Context) {
	id := c.Param("id")
	trip, err := h.service.GetTrip(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trip})
}

func (h *ScheduleHandler) ListTripsByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required (YYYY-MM-DD)"})
		return
	}

	trips, err := h.service.ListTripsByDate(c.Request.Context(), date)
	if err != nil {
		h.logger.Error("Failed to list trips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list trips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trips})
}

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

	trip, err := h.service.UpdateTripStatus(c.Request.Context(), id, req.Status, req.DelayMinutes)
	if err != nil {
		h.logger.Error("Failed to update trip status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trip})
}

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

	if err := h.service.GenerateTripsForSchedule(c.Request.Context(), req.ScheduleID, from, to); err != nil {
		h.logger.Error("Failed to generate trips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate trips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trips generated successfully"})
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
