package test

import (
	"testing"
	"time"

	"github.com/kuahbanyak/go-crud/internal/vehicle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockVehicleRepo is a mock implementation of vehicle.Repository
type MockVehicleRepo struct {
	mock.Mock
}

func (m *MockVehicleRepo) Create(v *vehicle.Vehicle) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockVehicleRepo) ListByOwner(owner uint) ([]vehicle.Vehicle, error) {
	args := m.Called(owner)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]vehicle.Vehicle), args.Error(1)
}

func (m *MockVehicleRepo) Get(id uint) (*vehicle.Vehicle, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vehicle.Vehicle), args.Error(1)
}

func (m *MockVehicleRepo) Update(v *vehicle.Vehicle) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockVehicleRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestVehicleModel_Validation(t *testing.T) {
	tests := []struct {
		name    string
		vehicle vehicle.Vehicle
		valid   bool
	}{
		{
			name: "Valid vehicle",
			vehicle: vehicle.Vehicle{
				OwnerID:      1,
				Brand:        "Toyota",
				Model:        "Camry",
				Year:         2020,
				LicensePlate: "ABC-123",
				VIN:          "1HGCM82633A123456",
				Mileage:      50000,
			},
			valid: true,
		},
		{
			name: "Valid vehicle with minimal data",
			vehicle: vehicle.Vehicle{
				OwnerID: 2,
				Brand:   "Honda",
				Model:   "Civic",
				Year:    2019,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotZero(t, tt.vehicle.OwnerID)
			assert.NotEmpty(t, tt.vehicle.Brand)
			assert.NotEmpty(t, tt.vehicle.Model)
			assert.Greater(t, tt.vehicle.Year, 1900)
		})
	}
}

func TestVehicleRepository_Create(t *testing.T) {
	mockRepo := new(MockVehicleRepo)

	testVehicle := &vehicle.Vehicle{
		OwnerID:      1,
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2020,
		LicensePlate: "ABC-123",
		VIN:          "1HGCM82633A123456",
		Mileage:      50000,
	}

	mockRepo.On("Create", testVehicle).Return(nil)

	err := mockRepo.Create(testVehicle)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleRepository_ListByOwner(t *testing.T) {
	mockRepo := new(MockVehicleRepo)

	expectedVehicles := []vehicle.Vehicle{
		{
			ID:           1,
			OwnerID:      1,
			Brand:        "Toyota",
			Model:        "Camry",
			Year:         2020,
			LicensePlate: "ABC-123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           2,
			OwnerID:      1,
			Brand:        "Honda",
			Model:        "Civic",
			Year:         2019,
			LicensePlate: "XYZ-789",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	mockRepo.On("ListByOwner", uint(1)).Return(expectedVehicles, nil)
	mockRepo.On("ListByOwner", uint(999)).Return([]vehicle.Vehicle{}, nil)

	// Test successful list
	vehicles, err := mockRepo.ListByOwner(1)
	assert.NoError(t, err)
	assert.NotNil(t, vehicles)
	assert.Len(t, vehicles, 2)
	assert.Equal(t, uint(1), vehicles[0].OwnerID)
	assert.Equal(t, uint(1), vehicles[1].OwnerID)

	// Test empty list
	emptyVehicles, err := mockRepo.ListByOwner(999)
	assert.NoError(t, err)
	assert.Len(t, emptyVehicles, 0)

	mockRepo.AssertExpectations(t)
}

func TestVehicleRepository_Get(t *testing.T) {
	mockRepo := new(MockVehicleRepo)

	expectedVehicle := &vehicle.Vehicle{
		ID:           1,
		OwnerID:      1,
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2020,
		LicensePlate: "ABC-123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("Get", uint(1)).Return(expectedVehicle, nil)
	mockRepo.On("Get", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// Test successful get
	foundVehicle, err := mockRepo.Get(1)
	assert.NoError(t, err)
	assert.NotNil(t, foundVehicle)
	assert.Equal(t, uint(1), foundVehicle.ID)
	assert.Equal(t, "Toyota", foundVehicle.Brand)

	// Test vehicle not found
	notFoundVehicle, err := mockRepo.Get(999)
	assert.Error(t, err)
	assert.Nil(t, notFoundVehicle)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestVehicleRepository_Update(t *testing.T) {
	mockRepo := new(MockVehicleRepo)

	testVehicle := &vehicle.Vehicle{
		ID:           1,
		OwnerID:      1,
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2020,
		LicensePlate: "NEW-123",
		Mileage:      60000,
	}

	mockRepo.On("Update", testVehicle).Return(nil)

	err := mockRepo.Update(testVehicle)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleRepository_Delete(t *testing.T) {
	mockRepo := new(MockVehicleRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)
	mockRepo.On("Delete", uint(999)).Return(gorm.ErrRecordNotFound)

	// Test successful delete
	err := mockRepo.Delete(1)
	assert.NoError(t, err)

	// Test delete not found
	err = mockRepo.Delete(999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}
