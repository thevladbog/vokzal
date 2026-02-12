// Package middleware — HTTP middleware для Payment Service.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Claims — JWT claims (must match auth-service token shape).
type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	StationID string `json:"station_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware returns a gin.HandlerFunc that validates the JWT in the Authorization header,
// verifies signature/claims, sets user_id, username, role, station_id in context, and aborts with 401 if invalid.
func AuthMiddleware(jwtSecret string, logger *zap.Logger) gin.HandlerFunc {
	if jwtSecret == "" {
		logger.Fatal("AuthMiddleware: jwtSecret must not be empty; set JWT secret in config (e.g. VOKZAL_PAYMENT_JWT_SECRET)")
	}
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			logger.Warn("Invalid or expired JWT", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		if claims.StationID != "" {
			c.Set("station_id", claims.StationID)
		}
		c.Next()
	}
}
