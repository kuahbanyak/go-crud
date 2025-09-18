package mocks

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id int, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, id, user)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) GetUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserService) AuthenticateUser(ctx context.Context, email, password string) (*entities.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserService) ValidateUser(user *entities.User) error {
	args := m.Called(user)
	return args.Error(0)
}
