package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/shared/utils"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock for ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll(ctx context.Context, filter *entities.ProductFilter) ([]*entities.Product, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, id types.MSSQLUUID, product *entities.Product) (*entities.Product, error) {
	args := m.Called(ctx, id, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id types.MSSQLUUID, stock int) error {
	args := m.Called(ctx, id, stock)
	return args.Error(0)
}

func (m *MockProductRepository) GetByCategory(ctx context.Context, category string) ([]*entities.Product, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Product), args.Error(1)
}

func (m *MockProductRepository) Count(ctx context.Context, filter *entities.ProductFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func TestProductUsecase_CreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	product := &entities.Product{
		Name:     "Test Product",
		Category: "Electronics",
		Price:    100.0,
		Stock:    10,
	}

	mockRepo.On("GetBySKU", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Product")).Return(product, nil)

	result, err := usecase.CreateProduct(ctx, product)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Product", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_CreateProduct_WithSKU(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	product := &entities.Product{
		Name:     "Test Product",
		Category: "Electronics",
		SKU:      "TEST-SKU-001",
		Price:    100.0,
		Stock:    10,
	}

	mockRepo.On("GetBySKU", ctx, "TEST-SKU-001").Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Product")).Return(product, nil)

	result, err := usecase.CreateProduct(ctx, product)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_CreateProduct_DuplicateSKU(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	product := &entities.Product{
		Name:     "Test Product",
		SKU:      "EXISTING-SKU",
		Category: "Electronics",
		Price:    100.0,
		Stock:    10,
	}

	existingProduct := &entities.Product{SKU: "EXISTING-SKU"}
	mockRepo.On("GetBySKU", ctx, "EXISTING-SKU").Return(existingProduct, nil)

	result, err := usecase.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_CreateProduct_ValidationError_EmptyName(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	product := &entities.Product{
		Name:  "",
		Price: 100.0,
	}

	result, err := usecase.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name is required")
}

func TestProductUsecase_CreateProduct_ValidationError_NegativePrice(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	product := &entities.Product{
		Name:  "Test Product",
		Price: -10.0,
	}

	result, err := usecase.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "price cannot be negative")
}

func TestProductUsecase_CreateProduct_ValidationError_NegativeStock(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	product := &entities.Product{
		Name:  "Test Product",
		Price: 100.0,
		Stock: -5,
	}

	result, err := usecase.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "stock cannot be negative")
}

func TestProductUsecase_CreateProduct_ValidationError_LongName(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	longName := ""
	for range make([]byte, 300) {
		longName = "a" + longName
	}

	product := &entities.Product{
		Name:  longName,
		Price: 100.0,
	}

	result, err := usecase.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestProductUsecase_GetProductByID_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}
	product := &entities.Product{
		ID:    productID,
		Name:  "Test Product",
		Price: 100.0,
	}

	mockRepo.On("GetByID", ctx, productID).Return(product, nil)

	result, err := usecase.GetProductByID(ctx, productID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Product", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_GetProductByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}

	mockRepo.On("GetByID", ctx, productID).Return(nil, nil)

	result, err := usecase.GetProductByID(ctx, productID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "product not found")
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_GetProductByID_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}

	mockRepo.On("GetByID", ctx, productID).Return(nil, errors.New("database error"))

	result, err := usecase.GetProductByID(ctx, productID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_GetProducts_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	filter := &entities.ProductFilter{
		Name:     "Test",
		Category: "Electronics",
		Limit:    10,
		Offset:   0,
	}

	products := []*entities.Product{
		{Name: "Product 1", Price: 100.0},
		{Name: "Product 2", Price: 200.0},
	}

	mockRepo.On("GetAll", ctx, filter).Return(products, nil)

	result, err := usecase.GetProducts(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_GetProducts_NilFilter(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	products := []*entities.Product{
		{Name: "Product 1"},
	}

	mockRepo.On("GetAll", ctx, mock.AnythingOfType("*entities.ProductFilter")).Return(products, nil)

	result, err := usecase.GetProducts(ctx, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_GetProducts_ExceedsMaxLimit(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	filter := &entities.ProductFilter{
		Limit: 200, // Exceeds max of 100
	}

	products := []*entities.Product{}
	mockRepo.On("GetAll", ctx, mock.AnythingOfType("*entities.ProductFilter")).Return(products, nil)

	result, err := usecase.GetProducts(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}
	existingProduct := &entities.Product{
		ID:    productID,
		Name:  "Old Product",
		Price: 50.0,
	}

	updatedProduct := &entities.Product{
		Name:  "New Product",
		Price: 100.0,
		Stock: 20,
	}

	mockRepo.On("GetByID", ctx, productID).Return(existingProduct, nil)
	mockRepo.On("Update", ctx, productID, updatedProduct).Return(updatedProduct, nil)

	result, err := usecase.UpdateProduct(ctx, productID, updatedProduct)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}
	product := &entities.Product{Name: "Test", Price: 100}

	mockRepo.On("GetByID", ctx, productID).Return(nil, nil)

	result, err := usecase.UpdateProduct(ctx, productID, product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "product not found")
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProduct_ValidationError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}
	existingProduct := &entities.Product{ID: productID}

	invalidProduct := &entities.Product{
		Name:  "",
		Price: -10,
	}

	mockRepo.On("GetByID", ctx, productID).Return(existingProduct, nil)

	result, err := usecase.UpdateProduct(ctx, productID, invalidProduct)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_DeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}
	product := &entities.Product{ID: productID}

	mockRepo.On("GetByID", ctx, productID).Return(product, nil)
	mockRepo.On("Delete", ctx, productID).Return(nil)

	err := usecase.DeleteProduct(ctx, productID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_DeleteProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}

	mockRepo.On("GetByID", ctx, productID).Return(nil, nil)

	err := usecase.DeleteProduct(ctx, productID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_DeleteProduct_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}
	product := &entities.Product{ID: productID}

	mockRepo.On("GetByID", ctx, productID).Return(product, nil)
	mockRepo.On("Delete", ctx, productID).Return(errors.New("delete failed"))

	err := usecase.DeleteProduct(ctx, productID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProductStock_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}
	product := &entities.Product{ID: productID, Stock: 10}

	mockRepo.On("GetByID", ctx, productID).Return(product, nil)
	mockRepo.On("UpdateStock", ctx, productID, 50).Return(nil)

	err := usecase.UpdateProductStock(ctx, productID, 50)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProductStock_NegativeStock(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}

	err := usecase.UpdateProductStock(ctx, productID, -10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stock cannot be negative")
}

func TestProductUsecase_UpdateProductStock_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}

	mockRepo.On("GetByID", ctx, productID).Return(nil, nil)

	err := usecase.UpdateProductStock(ctx, productID, 50)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProductStock_UpdateError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	validator := utils.NewValidator()
	usecase := usecases.NewProductUsecase(mockRepo, validator)
	ctx := context.Background()

	productID := types.MSSQLUUID{}
	product := &entities.Product{ID: productID}

	mockRepo.On("GetByID", ctx, productID).Return(product, nil)
	mockRepo.On("UpdateStock", ctx, productID, 50).Return(errors.New("update failed"))

	err := usecase.UpdateProductStock(ctx, productID, 50)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
