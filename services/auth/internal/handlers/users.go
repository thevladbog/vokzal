// Package handlers — HTTP-обработчики Auth Service (Users CRUD).
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/vokzal-tech/auth-service/internal/models"
	"github.com/vokzal-tech/auth-service/internal/repository"
	"github.com/vokzal-tech/auth-service/internal/service"

	"go.uber.org/zap"
)

// ListUsers — GET /v1/users (admin only).
func (h *AuthHandler) ListUsers(c *gin.Context) {
	role := c.Query("role")
	stationID := c.Query("station_id")
	page, errPage := strconv.Atoi(c.DefaultQuery("page", "1"))
	if errPage != nil {
		page = 1
	}
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if errLimit != nil {
		limit = 20
	}

	var rolePtr, stationIDPtr *string
	if role != "" {
		rolePtr = &role
	}
	if stationID != "" {
		stationIDPtr = &stationID
	}

	result, err := h.authService.ListUsers(c.Request.Context(), rolePtr, stationIDPtr, page, limit)
	if err != nil {
		h.logger.Error("ListUsers failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to list users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"users": result.Users,
			"total": result.Total,
			"page":  result.Page,
			"limit": result.Limit,
		},
	})
}

// CreateUserRequest — тело запроса создания пользователя.
//
//nolint:govet // fieldalignment: binding tags order
type CreateUserRequest struct {
	Username  string  `json:"username" binding:"required,min=3,max=50"`
	Password  string  `json:"password" binding:"required,min=8"`
	FullName  string  `json:"full_name" binding:"required,max=100"`
	Role      string  `json:"role" binding:"required,oneof=admin dispatcher cashier controller"`
	StationID *string `json:"station_id"`
}

// CreateUser — POST /v1/users (admin only).
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	in := &service.CreateUserInput{
		Username:  req.Username,
		Password:  req.Password,
		FullName:  req.FullName,
		Role:      req.Role,
		StationID: req.StationID,
	}
	user, err := h.authService.CreateUser(c.Request.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "Username already exists",
			})
			return
		}
		h.logger.Error("CreateUser failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    userToResponse(user),
	})
}

// userToResponse builds API response from User (excludes password_hash).
func userToResponse(user *models.User) gin.H {
	res := gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"full_name":  user.FullName,
		"role":       user.Role,
		"station_id": user.StationID,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	return res
}

// GetUser — GET /v1/users/:id (admin only).
func (h *AuthHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "User ID required"})
		return
	}

	user, err := h.authService.GetUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "User not found",
			})
			return
		}
		h.logger.Error("GetUser failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userToResponse(user),
	})
}

// UpdateUserRequest — тело запроса обновления пользователя.
type UpdateUserRequest struct {
	FullName  *string `json:"full_name"`
	Password  *string `json:"password"`
	Role      *string `json:"role"`
	StationID *string `json:"station_id"`
	IsActive  *bool   `json:"is_active"`
}

// UpdateUser — PUT /v1/users/:id (admin only).
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "User ID required"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	in := &service.UpdateUserInput{
		FullName:  req.FullName,
		Password:  req.Password,
		Role:      req.Role,
		StationID: req.StationID,
		IsActive:  req.IsActive,
	}
	if req.Role != nil && *req.Role != "" {
		valid := false
		for _, r := range []string{"admin", "dispatcher", "cashier", "controller"} {
			if *req.Role == r {
				valid = true
				break
			}
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid role"})
			return
		}
	}
	if req.Password != nil && *req.Password != "" && len(*req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Password must be at least 8 characters"})
		return
	}

	user, err := h.authService.UpdateUser(c.Request.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
			return
		}
		h.logger.Error("UpdateUser failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userToResponse(user),
	})
}

// DeleteUser — DELETE /v1/users/:id (admin only).
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "User ID required"})
		return
	}

	err := h.authService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
			return
		}
		h.logger.Error("DeleteUser failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}
