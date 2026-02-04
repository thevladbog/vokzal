package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response — стандартная структура ответа API.
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Success bool        `json:"success"`
}

// ErrorInfo — детали ошибки.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Success отправляет успешный ответ.
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// Error отправляет ответ с ошибкой.
func Error(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// ErrorWithDetails отправляет ответ с ошибкой и деталями.
func ErrorWithDetails(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest — ответ 400.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "bad_request", message)
}

// Unauthorized — ответ 401.
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "unauthorized", message)
}

// Forbidden — ответ 403.
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "forbidden", message)
}

// NotFound — ответ 404.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "not_found", message)
}

// InternalError — ответ 500.
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "internal_error", message)
}
