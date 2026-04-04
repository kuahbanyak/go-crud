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
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductUsecase struct {
	mock.Mock
}

func (m *MockProductUsecase) CreateProduct(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductUsecase) GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductUsecase) GetAllProducts(ctx context.Context, filter *entities.ProductFilter) ([]*entities.Product, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Product), args.Error(1)
}

func (m *MockProductUsecase) UpdateProduct(ctx context.Context, id uuid.UUID, product *entities.Product) (*entities.Product, error) {
	args := m.Called(ctx, id, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductUsecase) UpdateProductStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockProductUsecase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    entities.Product
		mockSetup      func(*MockProductUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "create product successfully",
			requestBody: entities.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       99.99,
				Stock:       100,
				Category:    "Electronics",
			},
			mockSetup: func(m *MockProductUsecase) {
				product := &entities.Product{
					ID:          uuid.New(),
					Name:        "Test Product",
					Description: "Test Description",
					Price:       99.99,
					Stock:       100,
					Category:    "Electronics",
				}
				m.On("CreateProduct", mock.Anything, mock.AnythingOfType("*entities.Product")).Return(product, nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "Product created successfully", resp["message"])
			},
		},
		{
			name:        "invalid json body",
			requestBody: entities.Product{},
			mockSetup: func(m *MockProductUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			requestBody: entities.Product{
				Name:  "Test Product",
				Price: 99.99,
			},
			mockSetup: func(m *MockProductUsecase) {
				m.On("CreateProduct", mock.Anything, mock.AnythingOfType("*entities.Product")).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockProductUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewProductHandler(&usecases.ProductUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			handler.CreateProduct(rec, req)

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

func TestProductHandler_GetProduct(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*MockProductUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "get product successfully",
			productID: validID.String(),
			mockSetup: func(m *MockProductUsecase) {
				product := &entities.Product{
					ID:          validID,
					Name:        "Test Product",
					Description: "Test Description",
					Price:       99.99,
					Stock:       100,
				}
				m.On("GetProductByID", mock.Anything, validID).Return(product, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "Test Product", data["name"])
			},
		},
		{
			name:      "invalid product ID",
			productID: "invalid-uuid",
			mockSetup: func(m *MockProductUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "product not found",
			productID: validID.String(),
			mockSetup: func(m *MockProductUsecase) {
				m.On("GetProductByID", mock.Anything, validID).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockProductUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewProductHandler(&usecases.ProductUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+tt.productID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			rec := httptest.NewRecorder()

			handler.GetProduct(rec, req)

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

func TestProductHandler_GetAllProducts(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockProductUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "get all products successfully",
			queryParams: "",
			mockSetup: func(m *MockProductUsecase) {
				products := []*entities.Product{
					{Name: "Product 1", Price: 10.99},
					{Name: "Product 2", Price: 20.99},
				}
				m.On("GetAllProducts", mock.Anything, mock.AnythingOfType("*entities.ProductFilter")).Return(products, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:        "get products with filter",
			queryParams: "?name=Product&category=Electronics&min_price=10&max_price=100",
			mockSetup: func(m *MockProductUsecase) {
				products := []*entities.Product{
					{Name: "Product 1", Category: "Electronics", Price: 50.99},
				}
				m.On("GetAllProducts", mock.Anything, mock.AnythingOfType("*entities.ProductFilter")).Return(products, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
			},
		},
		{
			name:        "database error",
			queryParams: "",
			mockSetup: func(m *MockProductUsecase) {
				m.On("GetAllProducts", mock.Anything, mock.AnythingOfType("*entities.ProductFilter")).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockProductUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewProductHandler(&usecases.ProductUsecase{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/products"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			handler.GetAllProducts(rec, req)

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

func TestProductHandler_UpdateProduct(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		requestBody    entities.Product
		mockSetup      func(*MockProductUsecase)
		expectedStatus int
	}{
		{
			name:      "update product successfully",
			productID: validID.String(),
			requestBody: entities.Product{
				Name:  "Updated Product",
				Price: 199.99,
			},
			mockSetup: func(m *MockProductUsecase) {
				updatedProduct := &entities.Product{
					ID:    validID,
					Name:  "Updated Product",
					Price: 199.99,
				}
				m.On("UpdateProduct", mock.Anything, validID, mock.AnythingOfType("*entities.Product")).Return(updatedProduct, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "invalid product ID",
			productID: "invalid-uuid",
			mockSetup: func(m *MockProductUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "product not found",
			productID: validID.String(),
			requestBody: entities.Product{
				Name: "Updated Product",
			},
			mockSetup: func(m *MockProductUsecase) {
				m.On("UpdateProduct", mock.Anything, validID, mock.AnythingOfType("*entities.Product")).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockProductUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewProductHandler(&usecases.ProductUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/products/"+tt.productID, bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			rec := httptest.NewRecorder()

			handler.UpdateProduct(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_UpdateProductStock(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		requestBody    map[string]int
		mockSetup      func(*MockProductUsecase)
		expectedStatus int
	}{
		{
			name:      "update stock successfully",
			productID: validID.String(),
			requestBody: map[string]int{
				"quantity": 50,
			},
			mockSetup: func(m *MockProductUsecase) {
				m.On("UpdateProductStock", mock.Anything, validID, 50).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "invalid product ID",
			productID: "invalid-uuid",
			mockSetup: func(m *MockProductUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "product not found",
			productID: validID.String(),
			requestBody: map[string]int{
				"quantity": 50,
			},
			mockSetup: func(m *MockProductUsecase) {
				m.On("UpdateProductStock", mock.Anything, validID, 50).Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockProductUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewProductHandler(&usecases.ProductUsecase{})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/products/"+tt.productID+"/stock", bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			rec := httptest.NewRecorder()

			handler.UpdateProductStock(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*MockProductUsecase)
		expectedStatus int
	}{
		{
			name:      "delete product successfully",
			productID: validID.String(),
			mockSetup: func(m *MockProductUsecase) {
				m.On("DeleteProduct", mock.Anything, validID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "invalid product ID",
			productID: "invalid-uuid",
			mockSetup: func(m *MockProductUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "product not found",
			productID: validID.String(),
			mockSetup: func(m *MockProductUsecase) {
				m.On("DeleteProduct", mock.Anything, validID).Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockProductUsecase)
			tt.mockSetup(mockUsecase)

			handler := handlers.NewProductHandler(&usecases.ProductUsecase{})

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/"+tt.productID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			rec := httptest.NewRecorder()

			handler.DeleteProduct(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}
