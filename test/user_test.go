package test

import (
	"testing"
	"time"

	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepo is a mock implementation of user.Repository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) FindByID(id uint) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) Update(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepo) FindAll() ([]*user.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func TestUserModel_Validation(t *testing.T) {
	tests := []struct {
		name     string
		user     user.User
		expected bool
	}{
		{
			name: "Valid user",
			user: user.User{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "John Doe",
				Phone:    "1234567890",
				Role:     user.RoleCustomer,
				Address:  "123 Main St",
			},
			expected: true,
		},
		{
			name: "Valid admin user",
			user: user.User{
				Email:    "admin@example.com",
				Password: "admin123",
				Name:     "Admin User",
				Role:     user.RoleAdmin,
			},
			expected: true,
		},
		{
			name: "Valid mechanic user",
			user: user.User{
				Email:    "mechanic@example.com",
				Password: "mechanic123",
				Name:     "Mechanic User",
				Role:     user.RoleMechanic,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that user struct can be created with valid data
			assert.NotEmpty(t, tt.user.Email)
			assert.NotEmpty(t, tt.user.Password)
			assert.NotEmpty(t, tt.user.Name)
			assert.Contains(t, []user.Role{user.RoleAdmin, user.RoleMechanic, user.RoleCustomer}, tt.user.Role)
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	mockRepo := new(MockUserRepo)

	testUser := &user.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "John Doe",
		Phone:    "1234567890",
		Role:     user.RoleCustomer,
		Address:  "123 Main St",
	}

	mockRepo.On("Create", testUser).Return(nil)

	err := mockRepo.Create(testUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	mockRepo := new(MockUserRepo)

	expectedUser := &user.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "John Doe",
		Role:      user.RoleCustomer,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(expectedUser, nil)
	mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, gorm.ErrRecordNotFound)

	// Test successful find
	foundUser, err := mockRepo.FindByEmail("test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, expectedUser.Email, foundUser.Email)
	assert.Equal(t, expectedUser.Name, foundUser.Name)

	// Test user not found
	notFoundUser, err := mockRepo.FindByEmail("notfound@example.com")
	assert.Error(t, err)
	assert.Nil(t, notFoundUser)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestUserRepository_FindByID(t *testing.T) {
	mockRepo := new(MockUserRepo)

	expectedUser := &user.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "John Doe",
		Role:      user.RoleCustomer,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedUser, nil)
	mockRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// Test successful find
	foundUser, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, uint(1), foundUser.ID)

	// Test user not found
	notFoundUser, err := mockRepo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, notFoundUser)

	mockRepo.AssertExpectations(t)
}

func TestUserRepository_Update(t *testing.T) {
	mockRepo := new(MockUserRepo)

	testUser := &user.User{
		ID:    1,
		Email: "updated@example.com",
		Name:  "Updated Name",
		Role:  user.RoleCustomer,
	}

	mockRepo.On("Update", testUser).Return(nil)

	err := mockRepo.Update(testUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_FindAll(t *testing.T) {
	mockRepo := new(MockUserRepo)

	expectedUsers := []*user.User{
		{
			ID:    1,
			Email: "user1@example.com",
			Name:  "User 1",
			Role:  user.RoleCustomer,
		},
		{
			ID:    2,
			Email: "user2@example.com",
			Name:  "User 2",
			Role:  user.RoleMechanic,
		},
	}

	mockRepo.On("FindAll").Return(expectedUsers, nil)

	users, err := mockRepo.FindAll()

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].Email, users[0].Email)
	assert.Equal(t, expectedUsers[1].Email, users[1].Email)

	mockRepo.AssertExpectations(t)
}

func TestUserRoles(t *testing.T) {
	tests := []struct {
		name string
		role user.Role
	}{
		{"Admin role", user.RoleAdmin},
		{"Mechanic role", user.RoleMechanic},
		{"Customer role", user.RoleCustomer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, string(tt.role))
			assert.Contains(t, []user.Role{user.RoleAdmin, user.RoleMechanic, user.RoleCustomer}, tt.role)
		})
	}
}
