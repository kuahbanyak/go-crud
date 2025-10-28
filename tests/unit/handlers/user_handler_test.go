package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	handlers "github.com/kuahbanyak/go-crud/internal/adapters/handlers/http"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) Login(ctx context.Context, email, password string) (*entities.User, string, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*entities.User), args.String(1), args.Error(2)
}

func (m *MockUserUsecase) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUsecase) GetAll(ctx context.Context) ([]*entities.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserUsecase) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUsecase) RefreshToken(ctx context.Context, userID uint) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    dto.RegisterRequest
		mockSetup      func(*MockUserUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful registration",
			requestBody: dto.RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "1234567890",
			},
			mockSetup: func(m *MockUserUsecase) {
				m.On("Register", mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
					return u.Email == "test@example.com" && u.Name == "John Doe"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "User registered successfully", resp["message"])
			},
		},
		{
			name: "registration with existing email",
			requestBody: dto.RegisterRequest{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "Jane",
				LastName:  "Doe",
				Phone:     "9876543210",
			},
			mockSetup: func(m *MockUserUsecase) {
				m.On("Register", mock.Anything, mock.Anything).Return(errors.New("user with this email already exists"))
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Contains(t, resp["message"], "Registration failed")
			},
		},
		{
			name:           "invalid request body",
			requestBody:    dto.RegisterRequest{},
			mockSetup:      func(m *MockUserUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewUserHandler(&usecases.UserUsecase{})
			handler = &handlers.UserHandler{}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			handler.Register(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkResponse != nil {
				var response map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    dto.LoginRequest
		mockSetup      func(*MockUserUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful login",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockUserUsecase) {
				user := &entities.User{
					Email: "test@example.com",
					Name:  "John Doe",
					Role:  entities.RoleCustomer,
				}
				m.On("Login", mock.Anything, "test@example.com", "password123").Return(user, "mock-token", nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["token"])
			},
		},
		{
			name: "login with invalid credentials",
			requestBody: dto.LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockUserUsecase) {
				m.On("Login", mock.Anything, "wrong@example.com", "wrongpassword").Return(nil, "", errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Contains(t, resp["message"], "Login failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewUserHandler(&usecases.UserUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			handler.Login(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkResponse != nil {
				var response map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		mockSetup      func(*MockUserUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "get profile successfully",
			userID: 1,
			mockSetup: func(m *MockUserUsecase) {
				user := &entities.User{
					Email: "test@example.com",
					Name:  "John Doe",
					Phone: "1234567890",
				}
				m.On("GetByID", mock.Anything, uint(1)).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "test@example.com", data["email"])
			},
		},
		{
			name:   "user not found",
			userID: 999,
			mockSetup: func(m *MockUserUsecase) {
				m.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewUserHandler(&usecases.UserUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
			ctx := context.WithValue(req.Context(), "id", tt.userID)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.GetProfile(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkResponse != nil {
				var response map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockUserUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "get all users successfully",
			mockSetup: func(m *MockUserUsecase) {
				users := []*entities.User{
					{Email: "user1@example.com", Name: "User 1"},
					{Email: "user2@example.com", Name: "User 2"},
				}
				m.On("GetAll", mock.Anything).Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name: "database error",
			mockSetup: func(m *MockUserUsecase) {
				m.On("GetAll", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewUserHandler(&usecases.UserUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			rec := httptest.NewRecorder()

			handler.GetUsers(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkResponse != nil {
				var response map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    dto.UpdateUserRequest
		mockSetup      func(*MockUserUsecase)
		expectedStatus int
	}{
		{
			name:   "update user successfully",
			userID: "1",
			requestBody: dto.UpdateUserRequest{
				FirstName: "Updated",
				LastName:  "Name",
				Phone:     "9999999999",
			},
			mockSetup: func(m *MockUserUsecase) {
				existingUser := &entities.User{Email: "test@example.com"}
				m.On("GetByID", mock.Anything, uint(1)).Return(existingUser, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "invalid user ID",
			userID: "invalid",
			mockSetup: func(m *MockUserUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "user not found",
			userID: "999",
			requestBody: dto.UpdateUserRequest{
				FirstName: "Updated",
				LastName:  "Name",
			},
			mockSetup: func(m *MockUserUsecase) {
				m.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewUserHandler(&usecases.UserUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+tt.userID, bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})
			rec := httptest.NewRecorder()

			handler.UpdateUser(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserUsecase)
		expectedStatus int
	}{
		{
			name:   "delete user successfully",
			userID: "1",
			mockSetup: func(m *MockUserUsecase) {
				m.On("Delete", mock.Anything, uint(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "invalid user ID",
			userID: "invalid",
			mockSetup: func(m *MockUserUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "user not found",
			userID: "999",
			mockSetup: func(m *MockUserUsecase) {
				m.On("Delete", mock.Anything, uint(999)).Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewUserHandler(&usecases.UserUsecase{})

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+tt.userID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})
			rec := httptest.NewRecorder()

			handler.DeleteUser(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}
