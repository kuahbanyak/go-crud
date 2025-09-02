package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCreation(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockUserRepository)
	mockAuthService := &MockAuthService{}
	userUsecase := usecases.NewUserUsecase(mockRepo, mockAuthService)

	// Create test user
	user := &entities.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockRepo.On("Create", mock.Anything, user).Return(user, nil)

	// Execute
	result, err := userUsecase.CreateUser(context.Background(), user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Email, result.Email)

	mockRepo.AssertExpectations(t)
}

// MockAuthService for testing
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateToken(userID uuid.UUID, role entities.Role) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (uuid.UUID, entities.Role, error) {
	args := m.Called(token)
	return args.Get(0).(uuid.UUID), args.Get(1).(entities.Role), args.Error(2)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) CheckPassword(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}
