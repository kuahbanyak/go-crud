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

type MockSettingUsecase struct {
	mock.Mock
}

func (m *MockSettingUsecase) GetAllSettings(ctx context.Context) ([]*entities.Setting, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Setting), args.Error(1)
}

func (m *MockSettingUsecase) GetSettingByKey(ctx context.Context, key string) (*entities.Setting, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Setting), args.Error(1)
}

func (m *MockSettingUsecase) GetSettingsByCategory(ctx context.Context, category string) ([]*entities.Setting, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Setting), args.Error(1)
}

func (m *MockSettingUsecase) GetPublicSettings(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockSettingUsecase) CreateSetting(ctx context.Context, setting *entities.Setting) error {
	args := m.Called(ctx, setting)
	return args.Error(0)
}

func (m *MockSettingUsecase) UpdateSetting(ctx context.Context, key string, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockSettingUsecase) DeleteSetting(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSettingHandler_GetAllSettings(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockSettingUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "get all settings successfully",
			mockSetup: func(m *MockSettingUsecase) {
				settings := []*entities.Setting{
					{Key: "max_queue_per_day", Value: "50", Category: "queue"},
					{Key: "business_hours_start", Value: "08:00", Category: "business"},
				}
				m.On("GetAllSettings", mock.Anything).Return(settings, nil)
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
			mockSetup: func(m *MockSettingUsecase) {
				m.On("GetAllSettings", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockSettingUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewSettingHandler(&usecases.SettingUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/settings", nil)
			rec := httptest.NewRecorder()

			handler.GetAllSettings(rec, req)

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

func TestSettingHandler_GetPublicSettings(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockSettingUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "get public settings successfully",
			mockSetup: func(m *MockSettingUsecase) {
				publicSettings := map[string]interface{}{
					"max_queue_per_day":     50,
					"business_hours_start":  "08:00",
					"business_hours_end":    "17:00",
					"estimated_service_min": 30,
				}
				m.On("GetPublicSettings", mock.Anything).Return(publicSettings, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, float64(50), data["max_queue_per_day"])
			},
		},
		{
			name: "database error",
			mockSetup: func(m *MockSettingUsecase) {
				m.On("GetPublicSettings", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockSettingUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewSettingHandler(&usecases.SettingUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/settings/public", nil)
			rec := httptest.NewRecorder()

			handler.GetPublicSettings(rec, req)

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

func TestSettingHandler_GetSetting(t *testing.T) {
	tests := []struct {
		name           string
		settingKey     string
		mockSetup      func(*MockSettingUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "get setting successfully",
			settingKey: "max_queue_per_day",
			mockSetup: func(m *MockSettingUsecase) {
				setting := &entities.Setting{
					Key:      "max_queue_per_day",
					Value:    "50",
					Category: "queue",
				}
				m.On("GetSettingByKey", mock.Anything, "max_queue_per_day").Return(setting, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "max_queue_per_day", data["key"])
			},
		},
		{
			name:       "setting not found",
			settingKey: "non_existent",
			mockSetup: func(m *MockSettingUsecase) {
				m.On("GetSettingByKey", mock.Anything, "non_existent").Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockSettingUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewSettingHandler(&usecases.SettingUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/settings/key/"+tt.settingKey, nil)
			req = mux.SetURLVars(req, map[string]string{"key": tt.settingKey})
			rec := httptest.NewRecorder()

			handler.GetSetting(rec, req)

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

func TestSettingHandler_CreateSetting(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    dto.CreateSettingRequest
		mockSetup      func(*MockSettingUsecase)
		expectedStatus int
	}{
		{
			name: "create setting successfully",
			requestBody: dto.CreateSettingRequest{
				Key:      "new_setting",
				Value:    "100",
				Category: "queue",
			},
			mockSetup: func(m *MockSettingUsecase) {
				m.On("CreateSetting", mock.Anything, mock.MatchedBy(func(s *entities.Setting) bool {
					return s.Key == "new_setting" && s.Value == "100"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "duplicate key",
			requestBody: dto.CreateSettingRequest{
				Key:   "existing_setting",
				Value: "50",
			},
			mockSetup: func(m *MockSettingUsecase) {
				m.On("CreateSetting", mock.Anything, mock.Anything).Return(errors.New("setting already exists"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockSettingUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewSettingHandler(&usecases.SettingUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/settings", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			handler.CreateSetting(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSettingHandler_UpdateSetting(t *testing.T) {
	tests := []struct {
		name           string
		settingKey     string
		requestBody    dto.UpdateSettingRequest
		mockSetup      func(*MockSettingUsecase)
		expectedStatus int
	}{
		{
			name:       "update setting successfully",
			settingKey: "max_queue_per_day",
			requestBody: dto.UpdateSettingRequest{
				Value: "75",
			},
			mockSetup: func(m *MockSettingUsecase) {
				m.On("UpdateSetting", mock.Anything, "max_queue_per_day", "75").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "setting not found",
			settingKey: "non_existent",
			requestBody: dto.UpdateSettingRequest{
				Value: "100",
			},
			mockSetup: func(m *MockSettingUsecase) {
				m.On("UpdateSetting", mock.Anything, "non_existent", "100").Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockSettingUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewSettingHandler(&usecases.SettingUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/settings/key/"+tt.settingKey, bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"key": tt.settingKey})
			rec := httptest.NewRecorder()

			handler.UpdateSetting(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSettingHandler_DeleteSetting(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		settingID      string
		mockSetup      func(*MockSettingUsecase)
		expectedStatus int
	}{
		{
			name:      "delete setting successfully",
			settingID: validID.String(),
			mockSetup: func(m *MockSettingUsecase) {
				m.On("DeleteSetting", mock.Anything, validID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid setting ID",
			settingID:      "invalid-uuid",
			mockSetup:      func(m *MockSettingUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "setting not found",
			settingID: validID.String(),
			mockSetup: func(m *MockSettingUsecase) {
				m.On("DeleteSetting", mock.Anything, validID).Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockSettingUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewSettingHandler(&usecases.SettingUsecase{})

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/settings/"+tt.settingID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.settingID})
			rec := httptest.NewRecorder()

			handler.DeleteSetting(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSettingHandler_GetSettingsByCategory(t *testing.T) {
	tests := []struct {
		name           string
		category       string
		mockSetup      func(*MockSettingUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:     "get settings by category successfully",
			category: "queue",
			mockSetup: func(m *MockSettingUsecase) {
				settings := []*entities.Setting{
					{Key: "max_queue_per_day", Value: "50", Category: "queue"},
					{Key: "queue_timeout_min", Value: "15", Category: "queue"},
				}
				m.On("GetSettingsByCategory", mock.Anything, "queue").Return(settings, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:     "category not found",
			category: "non_existent",
			mockSetup: func(m *MockSettingUsecase) {
				m.On("GetSettingsByCategory", mock.Anything, "non_existent").Return([]*entities.Setting{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockSettingUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewSettingHandler(&usecases.SettingUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/settings/category/"+tt.category, nil)
			req = mux.SetURLVars(req, map[string]string{"category": tt.category})
			rec := httptest.NewRecorder()

			handler.GetSettingsByCategory(rec, req)

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
