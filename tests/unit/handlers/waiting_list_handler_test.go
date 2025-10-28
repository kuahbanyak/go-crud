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

type MockWaitingListUsecase struct {
	mock.Mock
}

func (m *MockWaitingListUsecase) TakeQueue(ctx context.Context, userID uint, vehicleID uuid.UUID) (*entities.WaitingList, error) {
	args := m.Called(ctx, userID, vehicleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.WaitingList), args.Error(1)
}

func (m *MockWaitingListUsecase) GetMyQueue(ctx context.Context, userID uint) ([]*entities.WaitingList, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.WaitingList), args.Error(1)
}

func (m *MockWaitingListUsecase) GetTodayQueue(ctx context.Context) ([]*entities.WaitingList, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.WaitingList), args.Error(1)
}

func (m *MockWaitingListUsecase) GetQueueByDate(ctx context.Context, date string) ([]*entities.WaitingList, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.WaitingList), args.Error(1)
}

func (m *MockWaitingListUsecase) GetQueueByNumber(ctx context.Context, number int) (*entities.WaitingList, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.WaitingList), args.Error(1)
}

func (m *MockWaitingListUsecase) CheckAvailability(ctx context.Context, date string) (*dto.QueueAvailabilityResponse, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.QueueAvailabilityResponse), args.Error(1)
}

func (m *MockWaitingListUsecase) CancelQueue(ctx context.Context, id uuid.UUID, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockWaitingListUsecase) CallCustomer(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWaitingListUsecase) StartService(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWaitingListUsecase) CompleteService(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWaitingListUsecase) MarkNoShow(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWaitingListUsecase) GetServiceProgress(ctx context.Context, id uuid.UUID) (*dto.ServiceProgressResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ServiceProgressResponse), args.Error(1)
}

func TestWaitingListHandler_TakeQueueNumber(t *testing.T) {
	vehicleID := uuid.New()

	tests := []struct {
		name           string
		userID         uint
		requestBody    dto.TakeQueueRequest
		mockSetup      func(*MockWaitingListUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "take queue successfully",
			userID: 1,
			requestBody: dto.TakeQueueRequest{
				VehicleID: vehicleID,
			},
			mockSetup: func(m *MockWaitingListUsecase) {
				queue := &entities.WaitingList{
					ID:          uuid.New(),
					UserID:      1,
					VehicleID:   vehicleID,
					QueueNumber: 10,
					Status:      entities.WaitingListStatusWaiting,
				}
				m.On("TakeQueue", mock.Anything, uint(1), vehicleID).Return(queue, nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, float64(10), data["queue_number"])
			},
		},
		{
			name:   "queue full",
			userID: 1,
			requestBody: dto.TakeQueueRequest{
				VehicleID: vehicleID,
			},
			mockSetup: func(m *MockWaitingListUsecase) {
				m.On("TakeQueue", mock.Anything, uint(1), vehicleID).Return(nil, errors.New("queue is full"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockWaitingListUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewWaitingListHandler(&usecases.WaitingListUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/waiting-list/take", bytes.NewBuffer(body))
			ctx := context.WithValue(req.Context(), "id", tt.userID)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.TakeQueueNumber(rec, req)

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

func TestWaitingListHandler_GetMyQueue(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		mockSetup      func(*MockWaitingListUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "get my queue successfully",
			userID: 1,
			mockSetup: func(m *MockWaitingListUsecase) {
				queues := []*entities.WaitingList{
					{QueueNumber: 10, Status: entities.WaitingListStatusWaiting},
					{QueueNumber: 5, Status: entities.WaitingListStatusCompleted},
				}
				m.On("GetMyQueue", mock.Anything, uint(1)).Return(queues, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:   "no queues found",
			userID: 999,
			mockSetup: func(m *MockWaitingListUsecase) {
				m.On("GetMyQueue", mock.Anything, uint(999)).Return([]*entities.WaitingList{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockWaitingListUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewWaitingListHandler(&usecases.WaitingListUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/waiting-list/my-queue", nil)
			ctx := context.WithValue(req.Context(), "id", tt.userID)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.GetMyQueue(rec, req)

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

func TestWaitingListHandler_GetTodayQueue(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockWaitingListUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "get today queue successfully",
			mockSetup: func(m *MockWaitingListUsecase) {
				queues := []*entities.WaitingList{
					{QueueNumber: 1, Status: entities.WaitingListStatusWaiting},
					{QueueNumber: 2, Status: entities.WaitingListStatusInProgress},
				}
				m.On("GetTodayQueue", mock.Anything).Return(queues, nil)
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
			mockSetup: func(m *MockWaitingListUsecase) {
				m.On("GetTodayQueue", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockWaitingListUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewWaitingListHandler(&usecases.WaitingListUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/waiting-list/today", nil)
			rec := httptest.NewRecorder()

			handler.GetTodayQueue(rec, req)

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

func TestWaitingListHandler_CheckAvailability(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		mockSetup      func(*MockWaitingListUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "check availability successfully",
			queryParam: "?date=2025-10-28",
			mockSetup: func(m *MockWaitingListUsecase) {
				availability := &dto.QueueAvailabilityResponse{
					Date:              "2025-10-28",
					TotalCapacity:     50,
					CurrentCount:      30,
					AvailableSlots:    20,
					IsAvailable:       true,
					EstimatedWaitTime: 45,
				}
				m.On("CheckAvailability", mock.Anything, "2025-10-28").Return(availability, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.True(t, data["is_available"].(bool))
				assert.Equal(t, float64(20), data["available_slots"])
			},
		},
		{
			name:           "missing date parameter",
			queryParam:     "",
			mockSetup:      func(m *MockWaitingListUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockWaitingListUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewWaitingListHandler(&usecases.WaitingListUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/waiting-list/availability"+tt.queryParam, nil)
			rec := httptest.NewRecorder()

			handler.CheckAvailability(rec, req)

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

func TestWaitingListHandler_CancelQueue(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		queueID        string
		userID         uint
		mockSetup      func(*MockWaitingListUsecase)
		expectedStatus int
	}{
		{
			name:    "cancel queue successfully",
			queueID: validID.String(),
			userID:  1,
			mockSetup: func(m *MockWaitingListUsecase) {
				m.On("CancelQueue", mock.Anything, validID, uint(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid queue ID",
			queueID:        "invalid-uuid",
			userID:         1,
			mockSetup:      func(m *MockWaitingListUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "queue not found",
			queueID: validID.String(),
			userID:  1,
			mockSetup: func(m *MockWaitingListUsecase) {
				m.On("CancelQueue", mock.Anything, validID, uint(1)).Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockWaitingListUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewWaitingListHandler(&usecases.WaitingListUsecase{})

			req := httptest.NewRequest(http.MethodPut, "/api/v1/waiting-list/"+tt.queueID+"/cancel", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.queueID})
			ctx := context.WithValue(req.Context(), "id", tt.userID)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.CancelQueue(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestWaitingListHandler_CallCustomer(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		queueID        string
		mockSetup      func(*MockWaitingListUsecase)
		expectedStatus int
	}{
		{
			name:    "call customer successfully",
			queueID: validID.String(),
			mockSetup: func(m *MockWaitingListUsecase) {
				m.On("CallCustomer", mock.Anything, validID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid queue ID",
			queueID:        "invalid-uuid",
			mockSetup:      func(m *MockWaitingListUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockWaitingListUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewWaitingListHandler(&usecases.WaitingListUsecase{})

			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/waiting-list/"+tt.queueID+"/call", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.queueID})
			rec := httptest.NewRecorder()

			handler.CallCustomer(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestWaitingListHandler_StartService(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		queueID        string
		mockSetup      func(*MockWaitingListUsecase)
		expectedStatus int
	}{
		{
			name:    "start service successfully",
			queueID: validID.String(),
			mockSetup: func(m *MockWaitingListUsecase) {
				m.On("StartService", mock.Anything, validID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid queue ID",
			queueID:        "invalid-uuid",
			mockSetup:      func(m *MockWaitingListUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockWaitingListUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewWaitingListHandler(&usecases.WaitingListUsecase{})

			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/waiting-list/"+tt.queueID+"/start", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.queueID})
			rec := httptest.NewRecorder()

			handler.StartService(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}
