package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	handlers "github.com/kuahbanyak/go-crud/internal/adapters/handlers/http"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVehicleUsecase struct {
	mock.Mock
}

func (m *MockVehicleUsecase) Create(ctx context.Context, vehicle *entities.Vehicle) error {
	args := m.Called(ctx, vehicle)
	return args.Error(0)
}

func (m *MockVehicleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entities.Vehicle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Vehicle), args.Error(1)
}

func (m *MockVehicleUsecase) GetByUserID(ctx context.Context, userID uint) ([]*entities.Vehicle, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Vehicle), args.Error(1)
}

func (m *MockVehicleUsecase) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Vehicle), args.Error(1)
}

func (m *MockVehicleUsecase) Update(ctx context.Context, vehicle *entities.Vehicle) error {
	args := m.Called(ctx, vehicle)
	return args.Error(0)
}

func (m *MockVehicleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestVehicleHandler_CreateVehicle(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		requestBody    dto.CreateVehicleRequest
		mockSetup      func(*MockVehicleUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "create vehicle successfully",
			userID: 1,
			requestBody: dto.CreateVehicleRequest{
				LicensePlate: "ABC123",
				Brand:        "Toyota",
				Model:        "Camry",
				Year:         2023,
				Color:        "Black",
			},
			mockSetup: func(m *MockVehicleUsecase) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(v *entities.Vehicle) bool {
					return v.LicensePlate == "ABC123" && v.UserID == 1
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "Vehicle created successfully", resp["message"])
			},
		},
		{
			name:   "invalid request body",
			userID: 1,
			requestBody: dto.CreateVehicleRequest{
				LicensePlate: "",
			},
			mockSetup:      func(m *MockVehicleUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "usecase error",
			userID: 1,
			requestBody: dto.CreateVehicleRequest{
				LicensePlate: "ABC123",
				Brand:        "Toyota",
				Model:        "Camry",
				Year:         2023,
			},
			mockSetup: func(m *MockVehicleUsecase) {
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockVehicleUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewVehicleHandler(&usecases.VehicleUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/vehicles", bytes.NewBuffer(body))
			ctx := context.WithValue(req.Context(), "id", tt.userID)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.CreateVehicle(rec, req)

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

func TestVehicleHandler_GetVehicle(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		vehicleID      string
		mockSetup      func(*MockVehicleUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "get vehicle successfully",
			vehicleID: validID.String(),
			mockSetup: func(m *MockVehicleUsecase) {
				vehicle := &entities.Vehicle{
					ID:           validID,
					LicensePlate: "ABC123",
					Brand:        "Toyota",
					Model:        "Camry",
				}
				m.On("GetByID", mock.Anything, validID).Return(vehicle, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "ABC123", data["license_plate"])
			},
		},
		{
			name:           "invalid vehicle ID",
			vehicleID:      "invalid-uuid",
			mockSetup:      func(m *MockVehicleUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "vehicle not found",
			vehicleID: validID.String(),
			mockSetup: func(m *MockVehicleUsecase) {
				m.On("GetByID", mock.Anything, validID).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockVehicleUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewVehicleHandler(&usecases.VehicleUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/vehicles/"+tt.vehicleID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.vehicleID})
			rec := httptest.NewRecorder()

			handler.GetVehicle(rec, req)

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

func TestVehicleHandler_GetMyVehicles(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		mockSetup      func(*MockVehicleUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "get my vehicles successfully",
			userID: 1,
			mockSetup: func(m *MockVehicleUsecase) {
				vehicles := []*entities.Vehicle{
					{LicensePlate: "ABC123", Brand: "Toyota"},
					{LicensePlate: "XYZ789", Brand: "Honda"},
				}
				m.On("GetByUserID", mock.Anything, uint(1)).Return(vehicles, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:   "no vehicles found",
			userID: 999,
			mockSetup: func(m *MockVehicleUsecase) {
				m.On("GetByUserID", mock.Anything, uint(999)).Return([]*entities.Vehicle{}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockVehicleUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewVehicleHandler(&usecases.VehicleUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/vehicles", nil)
			ctx := context.WithValue(req.Context(), "id", tt.userID)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.GetMyVehicles(rec, req)

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

func TestVehicleHandler_GetAllVehicles(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockVehicleUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "get all vehicles successfully",
			mockSetup: func(m *MockVehicleUsecase) {
				vehicles := []*entities.Vehicle{
					{LicensePlate: "ABC123", Brand: "Toyota"},
					{LicensePlate: "XYZ789", Brand: "Honda"},
					{LicensePlate: "DEF456", Brand: "Ford"},
				}
				m.On("GetAll", mock.Anything).Return(vehicles, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 3)
			},
		},
		{
			name: "database error",
			mockSetup: func(m *MockVehicleUsecase) {
				m.On("GetAll", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockVehicleUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewVehicleHandler(&usecases.VehicleUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/vehicles", nil)
			rec := httptest.NewRecorder()

			handler.GetAllVehicles(rec, req)

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

func TestVehicleHandler_UpdateVehicle(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		vehicleID      string
		requestBody    dto.UpdateVehicleRequest
		mockSetup      func(*MockVehicleUsecase)
		expectedStatus int
	}{
		{
			name:      "update vehicle successfully",
			vehicleID: validID.String(),
			requestBody: dto.UpdateVehicleRequest{
				LicensePlate: "ABC123",
				Brand:        "Toyota",
				Model:        "Camry Updated",
			},
			mockSetup: func(m *MockVehicleUsecase) {
				existingVehicle := &entities.Vehicle{ID: validID}
				m.On("GetByID", mock.Anything, validID).Return(existingVehicle, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid vehicle ID",
			vehicleID:      "invalid-uuid",
			mockSetup:      func(m *MockVehicleUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "vehicle not found",
			vehicleID: validID.String(),
			requestBody: dto.UpdateVehicleRequest{
				Brand: "Toyota",
			},
			mockSetup: func(m *MockVehicleUsecase) {
				m.On("GetByID", mock.Anything, validID).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockVehicleUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewVehicleHandler(&usecases.VehicleUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/vehicles/"+tt.vehicleID, bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.vehicleID})
			rec := httptest.NewRecorder()

			handler.UpdateVehicle(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestVehicleHandler_DeleteVehicle(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		vehicleID      string
		mockSetup      func(*MockVehicleUsecase)
		expectedStatus int
	}{
		{
			name:      "delete vehicle successfully",
			vehicleID: validID.String(),
			mockSetup: func(m *MockVehicleUsecase) {
				m.On("Delete", mock.Anything, validID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid vehicle ID",
			vehicleID:      "invalid-uuid",
			mockSetup:      func(m *MockVehicleUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "vehicle not found",
			vehicleID: validID.String(),
			mockSetup: func(m *MockVehicleUsecase) {
				m.On("Delete", mock.Anything, validID).Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockVehicleUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewVehicleHandler(&usecases.VehicleUsecase{})

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/vehicles/"+tt.vehicleID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.vehicleID})
			rec := httptest.NewRecorder()

			handler.DeleteVehicle(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}
