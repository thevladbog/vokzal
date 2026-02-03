// Package handlers — HTTP-обработчики Geo Service.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/geo-service/internal/service"
	"go.uber.org/zap"
)

// GeoHandler обрабатывает HTTP-запросы к API геокодирования.
type GeoHandler struct {
	service service.GeoService
	logger  *zap.Logger
}

// NewGeoHandler создаёт новый GeoHandler.
func NewGeoHandler(service service.GeoService, logger *zap.Logger) *GeoHandler {
	return &GeoHandler{
		service: service,
		logger:  logger,
	}
}

// Geocode возвращает координаты по адресу.
func (h *GeoHandler) Geocode(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	result, err := h.service.Geocode(c.Request.Context(), address)
	if err != nil {
		h.logger.Error("Geocoding failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to geocode address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ReverseGeocode возвращает адрес по координатам.
func (h *GeoHandler) ReverseGeocode(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon parameters are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat parameter"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lon parameter"})
		return
	}

	result, err := h.service.ReverseGeocode(c.Request.Context(), lat, lon)
	if err != nil {
		h.logger.Error("Reverse geocoding failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reverse geocode"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetDistance возвращает расстояние между двумя точками (км).
func (h *GeoHandler) GetDistance(c *gin.Context) {
	lat1Str := c.Query("lat1")
	lon1Str := c.Query("lon1")
	lat2Str := c.Query("lat2")
	lon2Str := c.Query("lon2")

	if lat1Str == "" || lon1Str == "" || lat2Str == "" || lon2Str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat1, lon1, lat2, lon2 parameters are required"})
		return
	}

	lat1, _ := strconv.ParseFloat(lat1Str, 64)
	lon1, _ := strconv.ParseFloat(lon1Str, 64)
	lat2, _ := strconv.ParseFloat(lat2Str, 64)
	lon2, _ := strconv.ParseFloat(lon2Str, 64)

	distance := h.service.GetDistance(c.Request.Context(), lat1, lon1, lat2, lon2)

	c.JSON(http.StatusOK, gin.H{
		"distance_km": distance,
	})
}
