package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vokzal-tech/auth-service/internal/config"
	"github.com/vokzal-tech/auth-service/internal/models"
	"github.com/vokzal-tech/auth-service/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository - мок репозитория пользователей
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSessionRepository - мок репозитория сессий
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *models.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) FindByToken(ctx context.Context, tokenHash string) (*models.Session, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockSessionRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Session), args.Error(1)
}

func (m *MockSessionRepository) Delete(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestAuthService_Login(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	// Создаем bcrypt hash для "password123"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	tests := []struct {
		name          string
		username      string
		password      string
		stationID     string
		setupMock     func(*MockUserRepository, *MockSessionRepository)
		expectError   bool
	}{
		{
			name:      "успешная авторизация",
			username:  "admin",
			password:  "password123",
			stationID: "",
			setupMock: func(userRepo *MockUserRepository, sessionRepo *MockSessionRepository) {
				user := &models.User{
					ID:           "user-123",
					Username:     "admin",
					PasswordHash: string(passwordHash),
					Role:         "admin",
					IsActive:     true,
				}
				userRepo.On("FindByUsername", mock.Anything, "admin").Return(user, nil)
				sessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Session")).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "пользователь не найден",
			username:  "nonexistent",
			password:  "password123",
			stationID: "",
			setupMock: func(userRepo *MockUserRepository, sessionRepo *MockSessionRepository) {
				userRepo.On("FindByUsername", mock.Anything, "nonexistent").Return(nil, repository.ErrUserNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(MockUserRepository)
			mockSessionRepo := new(MockSessionRepository)
			tt.setupMock(mockUserRepo, mockSessionRepo)

			jwtConfig := config.JWTConfig{
				Secret:            "test-secret",
				Issuer:            "test",
				AccessExpiration:  15 * time.Minute,
				RefreshExpiration: 7 * 24 * time.Hour,
			}

			service := NewAuthService(mockUserRepo, mockSessionRepo, jwtConfig, logger)

			ctx := context.Background()

			// Act
			result, err := service.Login(ctx, tt.username, tt.password, tt.stationID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			mockUserRepo.AssertExpectations(t)
			mockSessionRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("успешный выход", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockSessionRepo := new(MockSessionRepository)

		mockSessionRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)

		jwtConfig := config.JWTConfig{
			Secret:            "test-secret",
			Issuer:            "test",
			AccessExpiration:  15 * time.Minute,
			RefreshExpiration: 7 * 24 * time.Hour,
		}

		service := NewAuthService(mockUserRepo, mockSessionRepo, jwtConfig, logger)

		ctx := context.Background()

		// Act
		err := service.Logout(ctx, "some-token")

		// Assert
		assert.NoError(t, err)
		mockSessionRepo.AssertExpectations(t)
	})
}
