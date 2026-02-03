package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vokzal-tech/auth-service/internal/service"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(authService service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	StationID string `json:"station_id"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login обработчик входа
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req.Username, req.Password, req.StationID)
	if err != nil {
		h.logger.Error("Login failed",
			zap.String("username", req.Username),
			zap.Error(err))

		statusCode := http.StatusUnauthorized
		message := "Invalid credentials"

		if err == service.ErrUserInactive {
			statusCode = http.StatusForbidden
			message = "User is inactive"
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// Refresh обработчик обновления токена
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	resp, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error("Token refresh failed", zap.Error(err))

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid or expired refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// Logout обработчик выхода
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Refresh token required",
		})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		h.logger.Error("Logout failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// Me получить информацию о текущем пользователе
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user_id":    userID,
			"username":   c.GetString("username"),
			"role":       c.GetString("role"),
			"station_id": c.GetString("station_id"),
		},
	})
}
