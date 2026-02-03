// Package repository — слой доступа к данным Auth Service.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vokzal-tech/auth-service/internal/models"

	"gorm.io/gorm"
)

var (
	// ErrUserNotFound возвращается, когда пользователь не найден.
	ErrUserNotFound = errors.New("user not found")
	// ErrUsernameExists возвращается, когда имя пользователя уже занято.
	ErrUsernameExists = errors.New("username already exists")
	// ErrSessionNotFound возвращается, когда сессия не найдена.
	ErrSessionNotFound = errors.New("session not found")
)

// UserRepository — интерфейс репозитория пользователей.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

// SessionRepository — интерфейс репозитория сессий.
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	FindByToken(ctx context.Context, tokenHash string) (*models.Session, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.Session, error)
	Delete(ctx context.Context, tokenHash string) error
	DeleteExpired(ctx context.Context) error
}

type userRepository struct {
	db *gorm.DB
}

type sessionRepository struct {
	db *gorm.DB
}

// NewUserRepository создаёт новый UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// NewSessionRepository создаёт новый SessionRepository.
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

// User Repository Implementation

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrUsernameExists
		}
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

// findUserBy находит пользователя по полю и значению (устраняет dupl между FindByID и FindByUsername).
func (r *userRepository) findUserBy(ctx context.Context, field, value string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, field+" = ?", value)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	return r.findUserBy(ctx, "id", id)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	return r.findUserBy(ctx, "username", username)
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// Session Repository Implementation

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	result := r.db.WithContext(ctx).Create(session)
	if result.Error != nil {
		return fmt.Errorf("failed to create session: %w", result.Error)
	}
	return nil
}

func (r *sessionRepository) FindByToken(ctx context.Context, tokenHash string) (*models.Session, error) {
	var session models.Session
	result := r.db.WithContext(ctx).
		Preload("User").
		First(&session, "token_hash = ?", tokenHash)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to find session: %w", result.Error)
	}

	// Проверить истечение срока
	if session.ExpiresAt.Before(time.Now()) {
		return nil, ErrSessionNotFound
	}

	return &session, nil
}

func (r *sessionRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Session, error) {
	var sessions []*models.Session
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Find(&sessions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find sessions: %w", result.Error)
	}
	return sessions, nil
}

func (r *sessionRepository) Delete(ctx context.Context, tokenHash string) error {
	result := r.db.WithContext(ctx).Delete(&models.Session{}, "token_hash = ?", tokenHash)
	if result.Error != nil {
		return fmt.Errorf("failed to delete session: %w", result.Error)
	}
	return nil
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.Session{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", result.Error)
	}
	return nil
}
