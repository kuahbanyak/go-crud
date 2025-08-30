package servicehistory_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/servicehistory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServiceHistoryRepo struct {
	mock.Mock
}

func (m *MockServiceHistoryRepo) Create(s *servicehistory.ServiceRecord) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockServiceHistoryRepo) ListByVehicle(vehicle string) ([]servicehistory.ServiceRecord, error) {
	args := m.Called(vehicle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]servicehistory.ServiceRecord), args.Error(1)
}

func TestServiceRecordModel_Validation(t *testing.T) {
	tests := []struct {
		name   string
		record servicehistory.ServiceRecord
		valid  bool
	}{
		{
			name: "Valid service record",
			record: servicehistory.ServiceRecord{
				BookingID:   uuid.New().String(),
				VehicleID:   uuid.New().String(),
				Description: "Oil change and filter replacement",
				Cost:        150000,
				ReceiptURL:  "https://example.com/receipt.pdf",
			},
			valid: true,
		},
		{
			name: "Valid service record without receipt",
			record: servicehistory.ServiceRecord{
				BookingID:   uuid.New().String(),
				VehicleID:   uuid.New().String(),
				Description: "Brake pad replacement",
				Cost:        200000,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.record.BookingID)
			assert.NotEmpty(t, tt.record.VehicleID)
			assert.NotEmpty(t, tt.record.Description)
			assert.Greater(t, tt.record.Cost, 0)
		})
	}
}

func TestServiceHistoryRepository_Create(t *testing.T) {
	mockRepo := new(MockServiceHistoryRepo)

	testRecord := &servicehistory.ServiceRecord{
		BookingID:   uuid.New().String(),
		VehicleID:   uuid.New().String(),
		Description: "Oil change and filter replacement",
		Cost:        150000,
		ReceiptURL:  "https://example.com/receipt.pdf",
	}

	mockRepo.On("Create", testRecord).Return(nil)

	err := mockRepo.Create(testRecord)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestServiceHistoryRepository_ListByVehicle(t *testing.T) {
	mockRepo := new(MockServiceHistoryRepo)
	vehicleID := uuid.New().String()

	expectedRecords := []servicehistory.ServiceRecord{
		{
			ID:          uuid.New().String(),
			BookingID:   uuid.New().String(),
			VehicleID:   vehicleID,
			Description: "Oil change",
			Cost:        100000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			BookingID:   uuid.New().String(),
			VehicleID:   vehicleID,
			Description: "Brake service",
			Cost:        250000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("ListByVehicle", vehicleID).Return(expectedRecords, nil)
	mockRepo.On("ListByVehicle", "non-existent-id").Return([]servicehistory.ServiceRecord{}, nil)

	// Test successful list
	records, err := mockRepo.ListByVehicle(vehicleID)
	assert.NoError(t, err)
	assert.NotNil(t, records)
	assert.Len(t, records, 2)
	assert.Equal(t, vehicleID, records[0].VehicleID)
	assert.Equal(t, vehicleID, records[1].VehicleID)

	// Test empty list
	emptyRecords, err := mockRepo.ListByVehicle("non-existent-id")
	assert.NoError(t, err)
	assert.Len(t, emptyRecords, 0)

	mockRepo.AssertExpectations(t)
}
