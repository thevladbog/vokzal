package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	// ErrInvalidToken возвращается при невалидном или подписанном неверным ключом токене.
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken возвращается при истёкшем токене.
	ErrExpiredToken = errors.New("token expired")
)

// Claims содержит JWT claims для Вокзал.ТЕХ.
type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	StationID string `json:"station_id"`
	jwt.RegisteredClaims
}

// Config — конфигурация JWT.
type Config struct {
	Secret     string
	Issuer     string
	Expiration time.Duration
}

// Manager управляет JWT токенами.
type Manager struct {
	config Config
}

// NewManager создаёт новый JWT manager.
func NewManager(config Config) *Manager {
	return &Manager{config: config}
}

// Generate генерирует новый JWT токен.
func (m *Manager) Generate(userID, username, role, stationID string) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		StationID: stationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.Expiration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Verify проверяет и парсит JWT токен.
func (m *Manager) Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
