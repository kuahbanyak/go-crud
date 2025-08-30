package vehicle_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
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

func (m *MockVehicleRepo) ListByOwner(owner string) ([]vehicle.Vehicle, error) {
	args := m.Called(owner)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]vehicle.Vehicle), args.Error(1)
}

func (m *MockVehicleRepo) Get(id string) (*vehicle.Vehicle, error) {
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

func (m *MockVehicleRepo) Delete(id string) error {
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
				OwnerID:      uuid.New().String(),
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
				OwnerID: uuid.New().String(),
				Brand:   "Honda",
				Model:   "Civic",
				Year:    2019,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.vehicle.OwnerID)
			assert.NotEmpty(t, tt.vehicle.Brand)
			assert.NotEmpty(t, tt.vehicle.Model)
			assert.Greater(t, tt.vehicle.Year, 1900)
		})
	}
}

func TestVehicleRepository_Create(t *testing.T) {
	mockRepo := new(MockVehicleRepo)

	testVehicle := &vehicle.Vehicle{
		OwnerID:      uuid.New().String(),
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
	ownerID := uuid.New().String()

	expectedVehicles := []vehicle.Vehicle{
		{
			ID:           uuid.New().String(),
			OwnerID:      ownerID,
			Brand:        "Toyota",
			Model:        "Camry",
			Year:         2020,
			LicensePlate: "ABC-123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New().String(),
			OwnerID:      ownerID,
			Brand:        "Honda",
			Model:        "Civic",
			Year:         2019,
			LicensePlate: "XYZ-789",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	mockRepo.On("ListByOwner", ownerID).Return(expectedVehicles, nil)
	mockRepo.On("ListByOwner", "non-existent-id").Return([]vehicle.Vehicle{}, nil)

	// Test successful list
	vehicles, err := mockRepo.ListByOwner(ownerID)
	assert.NoError(t, err)
	assert.NotNil(t, vehicles)
	assert.Len(t, vehicles, 2)
	assert.Equal(t, ownerID, vehicles[0].OwnerID)
	assert.Equal(t, ownerID, vehicles[1].OwnerID)

	// Test empty list
	emptyVehicles, err := mockRepo.ListByOwner("non-existent-id")
	assert.NoError(t, err)
	assert.Len(t, emptyVehicles, 0)

	mockRepo.AssertExpectations(t)
}

func TestVehicleRepository_Get(t *testing.T) {
	mockRepo := new(MockVehicleRepo)
	vehicleID := uuid.New().String()

	expectedVehicle := &vehicle.Vehicle{
		ID:           vehicleID,
		OwnerID:      uuid.New().String(),
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2020,
		LicensePlate: "ABC-123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("Get", vehicleID).Return(expectedVehicle, nil)
	mockRepo.On("Get", "non-existent-id").Return(nil, gorm.ErrRecordNotFound)

	// Test successful get
	foundVehicle, err := mockRepo.Get(vehicleID)
	assert.NoError(t, err)
	assert.NotNil(t, foundVehicle)
	assert.Equal(t, vehicleID, foundVehicle.ID)
	assert.Equal(t, "Toyota", foundVehicle.Brand)

	// Test vehicle not found
	notFoundVehicle, err := mockRepo.Get("non-existent-id")
	assert.Error(t, err)
	assert.Nil(t, notFoundVehicle)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestVehicleRepository_Update(t *testing.T) {
	mockRepo := new(MockVehicleRepo)

	testVehicle := &vehicle.Vehicle{
		ID:           uuid.New().String(),
		OwnerID:      uuid.New().String(),
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
	vehicleID := uuid.New().String()

	mockRepo.On("Delete", vehicleID).Return(nil)
	mockRepo.On("Delete", "non-existent-id").Return(gorm.ErrRecordNotFound)

	// Test successful delete
	err := mockRepo.Delete(vehicleID)
	assert.NoError(t, err)

	// Test delete not found
	err = mockRepo.Delete("non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}
