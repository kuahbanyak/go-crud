package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVehicleRepository is a mock for VehicleRepository
type MockVehicleRepository struct {
	mock.Mock
}

func (m *MockVehicleRepository) Create(ctx context.Context, vehicle *entities.Vehicle) error {
	args := m.Called(ctx, vehicle)
	return args.Error(0)
}

func (m *MockVehicleRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Vehicle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Vehicle), args.Error(1)
}

func (m *MockVehicleRepository) GetByOwnerID(ctx context.Context, ownerID types.MSSQLUUID) ([]*entities.Vehicle, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Vehicle), args.Error(1)
}

func (m *MockVehicleRepository) Update(ctx context.Context, vehicle *entities.Vehicle) error {
	args := m.Called(ctx, vehicle)
	return args.Error(0)
}

func (m *MockVehicleRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVehicleRepository) List(ctx context.Context, limit, offset int) ([]*entities.Vehicle, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Vehicle), args.Error(1)
}

func (m *MockVehicleRepository) ListPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.Vehicle, int64, error) {
	args := m.Called(ctx, pagParams, filterParams)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Vehicle), int64(args.Int(1)), args.Error(2)
}

func TestVehicleUseCase_CreateVehicle_Success(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	req := &dto.CreateVehicleRequest{
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC123",
		VIN:          "1234567890",
		Mileage:      1000,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Vehicle")).Return(nil)

	result, err := usecase.CreateVehicle(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Toyota", result.Brand)
	assert.Equal(t, "Camry", result.Model)
	assert.Equal(t, 2023, result.Year)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_CreateVehicle_Error(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	req := &dto.CreateVehicleRequest{Brand: "Toyota"}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Vehicle")).Return(errors.New("database error"))

	result, err := usecase.CreateVehicle(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetMyVehicles_Success(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicles := []*entities.Vehicle{
		{
			ID:           types.MSSQLUUID{},
			OwnerID:      userID,
			Brand:        "Toyota",
			Model:        "Camry",
			Year:         2023,
			LicensePlate: "ABC123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           types.MSSQLUUID{UUID: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
			OwnerID:      userID,
			Brand:        "Honda",
			Model:        "Accord",
			Year:         2022,
			LicensePlate: "XYZ789",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	mockRepo.On("GetByOwnerID", ctx, userID).Return(vehicles, nil)

	result, err := usecase.GetMyVehicles(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Toyota", result[0].Brand)
	assert.Equal(t, "Honda", result[1].Brand)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetMyVehicles_Error(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}

	mockRepo.On("GetByOwnerID", ctx, userID).Return(nil, errors.New("database error"))

	result, err := usecase.GetMyVehicles(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetVehicleByID_Success(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:           vehicleID,
		OwnerID:      userID,
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)

	result, err := usecase.GetVehicleByID(ctx, userID, vehicleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Toyota", result.Brand)
	assert.Equal(t, "Camry", result.Model)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetVehicleByID_NotFound(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}

	mockRepo.On("GetByID", ctx, vehicleID).Return(nil, nil)

	result, err := usecase.GetVehicleByID(ctx, userID, vehicleID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "vehicle not found")
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetVehicleByID_Unauthorized(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	differentUserID := types.MSSQLUUID{UUID: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:      vehicleID,
		OwnerID: differentUserID,
		Brand:   "Toyota",
	}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)

	result, err := usecase.GetVehicleByID(ctx, userID, vehicleID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetVehicleByID_RepoError(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}

	mockRepo.On("GetByID", ctx, vehicleID).Return(nil, errors.New("database error"))

	result, err := usecase.GetVehicleByID(ctx, userID, vehicleID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_UpdateVehicle_Success(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:           vehicleID,
		OwnerID:      userID,
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2020,
		LicensePlate: "OLD123",
		Mileage:      10000,
	}

	req := &dto.UpdateVehicleRequest{
		Brand:        "Honda",
		Model:        "Accord",
		Year:         2024,
		LicensePlate: "NEW456",
		VIN:          "VIN123456",
		Mileage:      5000,
	}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entities.Vehicle")).Return(nil)

	result, err := usecase.UpdateVehicle(ctx, userID, vehicleID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Honda", result.Brand)
	assert.Equal(t, "Accord", result.Model)
	assert.Equal(t, 2024, result.Year)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_UpdateVehicle_PartialUpdate(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:      vehicleID,
		OwnerID: userID,
		Brand:   "Toyota",
		Model:   "Camry",
	}

	req := &dto.UpdateVehicleRequest{
		Brand: "Honda",
	}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entities.Vehicle")).Return(nil)

	result, err := usecase.UpdateVehicle(ctx, userID, vehicleID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Honda", result.Brand)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_UpdateVehicle_Unauthorized(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	differentUserID := types.MSSQLUUID{UUID: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:      vehicleID,
		OwnerID: differentUserID,
	}

	req := &dto.UpdateVehicleRequest{Brand: "Honda"}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)

	result, err := usecase.UpdateVehicle(ctx, userID, vehicleID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_UpdateVehicle_UpdateError(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:      vehicleID,
		OwnerID: userID,
		Brand:   "Toyota",
	}

	req := &dto.UpdateVehicleRequest{Brand: "Honda"}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entities.Vehicle")).Return(errors.New("update failed"))

	result, err := usecase.UpdateVehicle(ctx, userID, vehicleID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_DeleteVehicle_Success(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:      vehicleID,
		OwnerID: userID,
		Brand:   "Toyota",
	}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)
	mockRepo.On("Delete", ctx, vehicleID).Return(nil)

	err := usecase.DeleteVehicle(ctx, userID, vehicleID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_DeleteVehicle_Unauthorized(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	differentUserID := types.MSSQLUUID{UUID: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:      vehicleID,
		OwnerID: differentUserID,
	}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)

	err := usecase.DeleteVehicle(ctx, userID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_DeleteVehicle_NotFound(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}

	mockRepo.On("GetByID", ctx, vehicleID).Return(nil, nil)

	err := usecase.DeleteVehicle(ctx, userID, vehicleID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_DeleteVehicle_DeleteError(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	vehicleID := types.MSSQLUUID{}
	vehicle := &entities.Vehicle{
		ID:      vehicleID,
		OwnerID: userID,
	}

	mockRepo.On("GetByID", ctx, vehicleID).Return(vehicle, nil)
	mockRepo.On("Delete", ctx, vehicleID).Return(errors.New("delete failed"))

	err := usecase.DeleteVehicle(ctx, userID, vehicleID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetAllVehicles_Success(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	vehicles := []*entities.Vehicle{
		{Brand: "Toyota", Model: "Camry", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Brand: "Honda", Model: "Accord", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockRepo.On("List", ctx, 10, 0).Return(vehicles, nil)

	result, err := usecase.GetAllVehicles(ctx, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Toyota", result[0].Brand)
	assert.Equal(t, "Honda", result[1].Brand)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetAllVehicles_Error(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	mockRepo.On("List", ctx, 10, 0).Return(nil, errors.New("database error"))

	result, err := usecase.GetAllVehicles(ctx, 10, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetAllVehiclesPaginated_Success(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	vehicles := []*entities.Vehicle{
		{Brand: "Toyota", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Brand: "Honda", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	pagParams := pagination.Params{Page: 1, PageSize: 10}
	filterParams := pagination.FilterParams{}

	mockRepo.On("ListPaginated", ctx, pagParams, filterParams).Return(vehicles, 2, nil)

	result, total, err := usecase.GetAllVehiclesPaginated(ctx, pagParams, filterParams)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestVehicleUseCase_GetAllVehiclesPaginated_Error(t *testing.T) {
	mockRepo := new(MockVehicleRepository)
	usecase := usecases.NewVehicleUseCase(mockRepo)
	ctx := context.Background()

	pagParams := pagination.Params{Page: 1, PageSize: 10}
	filterParams := pagination.FilterParams{}

	mockRepo.On("ListPaginated", ctx, pagParams, filterParams).Return(nil, 0, errors.New("database error"))

	result, total, err := usecase.GetAllVehiclesPaginated(ctx, pagParams, filterParams)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}
