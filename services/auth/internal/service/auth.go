// Package service — бизнес-логика Auth Service.
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/vokzal-tech/auth-service/internal/config"
	"github.com/vokzal-tech/auth-service/internal/models"
	"github.com/vokzal-tech/auth-service/internal/repository"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials возвращается при неверном логине или пароле.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserInactive возвращается, когда пользователь деактивирован.
	ErrUserInactive = errors.New("user is inactive")
	// ErrInvalidToken возвращается при невалидном или истёкшем токене.
	ErrInvalidToken = errors.New("invalid token")
)

// ListUsersResult — результат списка пользователей.
type ListUsersResult struct {
	Users []*models.User `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// CreateUserInput — входные данные для создания пользователя.
//
//nolint:govet // fieldalignment: JSON tags order
type CreateUserInput struct {
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	FullName  string  `json:"full_name"`
	Role      string  `json:"role"`
	StationID *string `json:"station_id,omitempty"`
}

// UpdateUserInput — входные данные для обновления пользователя.
type UpdateUserInput struct {
	FullName  *string `json:"full_name,omitempty"`
	Password  *string `json:"password,omitempty"`
	Role      *string `json:"role,omitempty"`
	StationID *string `json:"station_id,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// AuthService — интерфейс сервиса аутентификации.
type AuthService interface {
	Login(ctx context.Context, username, password, stationID string) (*LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*LoginResponse, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*models.User, error)
	// User CRUD (admin only, enforced by middleware)
	ListUsers(ctx context.Context, role, stationID *string, page, limit int) (*ListUsersResult, error)
	CreateUser(ctx context.Context, in *CreateUserInput) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, in *UpdateUserInput) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

//
//nolint:govet // fieldalignment: порядок полей
type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtConfig   config.JWTConfig
	logger      *zap.Logger
}

// LoginResponse — ответ при успешном входе.
//
//nolint:govet // fieldalignment: порядок полей для JSON
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

// UserResponse — данные пользователя в ответе API.
//
//nolint:govet // fieldalignment: порядок полей для JSON
type UserResponse struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	FullName  string  `json:"full_name"`
	Role      string  `json:"role"`
	StationID *string `json:"station_id,omitempty"`
}

// NewAuthService создаёт новый AuthService.
func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwtConfig config.JWTConfig,
	logger *zap.Logger,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtConfig:   jwtConfig,
		logger:      logger,
	}
}

func (s *authService) Login(ctx context.Context, username, password, stationID string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Проверить пароль
	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); compareErr != nil {
		s.logger.Warn("Invalid password attempt",
			zap.String("username", username),
			zap.String("user_id", user.ID))
		return nil, ErrInvalidCredentials
	}

	// Проверить активность
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Проверить станцию
	if stationID != "" && user.StationID != nil && *user.StationID != stationID {
		return nil, fmt.Errorf("user station mismatch")
	}

	// Генерировать токены
	accessToken, err := s.generateToken(user, s.jwtConfig.AccessExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(user, s.jwtConfig.RefreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Сохранить сессию
	tokenHash := hashToken(refreshToken)
	session := &models.Session{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.jwtConfig.RefreshExpiration),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info("User logged in",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("role", user.Role))

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtConfig.AccessExpiration.Seconds()),
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      user.Role,
			StationID: user.StationID,
		},
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	tokenHash := hashToken(refreshToken)

	session, err := s.sessionRepo.FindByToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	// Проверить токен
	claims, err := s.verifyToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims.UserID != session.UserID {
		return nil, ErrInvalidToken
	}

	// Получить пользователя
	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Генерировать новые токены
	newAccessToken, err := s.generateToken(user, s.jwtConfig.AccessExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.generateToken(user, s.jwtConfig.RefreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Удалить старую сессию
	if err := s.sessionRepo.Delete(ctx, tokenHash); err != nil {
		s.logger.Error("Failed to delete old session", zap.Error(err))
	}

	// Создать новую сессию
	newTokenHash := hashToken(newRefreshToken)
	newSession := &models.Session{
		UserID:    user.ID,
		TokenHash: newTokenHash,
		ExpiresAt: time.Now().Add(s.jwtConfig.RefreshExpiration),
	}

	if err := s.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.jwtConfig.AccessExpiration.Seconds()),
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      user.Role,
			StationID: user.StationID,
		},
	}, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	tokenHash := hashToken(token)
	return s.sessionRepo.Delete(ctx, tokenHash)
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	claims, err := s.verifyToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	return user, nil
}

// Claims — JWT claims для access-токена.
type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	StationID string `json:"station_id"`
	jwt.RegisteredClaims
}

func (s *authService) generateToken(user *models.User, expiration time.Duration) (string, error) {
	now := time.Now()

	stationID := ""
	if user.StationID != nil {
		stationID = *user.StationID
	}

	claims := Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		StationID: stationID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.jwtConfig.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.Secret))
}

func (s *authService) verifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *authService) ListUsers(ctx context.Context, role, stationID *string, page, limit int) (*ListUsersResult, error) {
	filter := repository.ListUsersFilter{
		Role:      role,
		StationID: stationID,
		Page:      page,
		Limit:     limit,
	}
	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return &ListUsersResult{
		Users: users,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *authService) CreateUser(ctx context.Context, in *CreateUserInput) (*models.User, error) {
	existing, err := s.userRepo.FindByUsername(ctx, in.Username)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("check username: %w", err)
	}
	if existing != nil {
		return nil, repository.ErrUsernameExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user := &models.User{
		Username:     in.Username,
		PasswordHash: string(hash),
		FullName:     in.FullName,
		Role:         in.Role,
		StationID:    in.StationID,
		IsActive:     true,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	s.logger.Info("User created",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("role", user.Role))
	return user, nil
}

func (s *authService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *authService) UpdateUser(ctx context.Context, id string, in *UpdateUserInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.FullName != nil {
		user.FullName = *in.FullName
	}
	if in.Password != nil && *in.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}
	if in.Role != nil {
		user.Role = *in.Role
	}
	if in.StationID != nil {
		user.StationID = in.StationID
	}
	if in.IsActive != nil {
		user.IsActive = *in.IsActive
	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	s.logger.Info("User updated",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username))
	return user, nil
}

func (s *authService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
