package test

import (
	"testing"
	"time"

	"github.com/kuahbanyak/go-crud/internal/servicehistory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServiceHistoryRepo is a mock implementation of servicehistory.Repository
type MockServiceHistoryRepo struct {
	mock.Mock
}

func (m *MockServiceHistoryRepo) Create(s *servicehistory.ServiceRecord) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockServiceHistoryRepo) ListByVehicle(vehicle uint) ([]servicehistory.ServiceRecord, error) {
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
				BookingID:   1,
				VehicleID:   1,
				Description: "Oil change and filter replacement",
				Cost:        150000,
				ReceiptURL:  "https://example.com/receipt.pdf",
			},
			valid: true,
		},
		{
			name: "Valid service record without receipt",
			record: servicehistory.ServiceRecord{
				BookingID:   2,
				VehicleID:   2,
				Description: "Brake pad replacement",
				Cost:        200000,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotZero(t, tt.record.BookingID)
			assert.NotZero(t, tt.record.VehicleID)
			assert.NotEmpty(t, tt.record.Description)
			assert.GreaterOrEqual(t, tt.record.Cost, 0)
		})
	}
}

func TestServiceHistoryRepository_Create(t *testing.T) {
	mockRepo := new(MockServiceHistoryRepo)

	testRecord := &servicehistory.ServiceRecord{
		BookingID:   1,
		VehicleID:   1,
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

	expectedRecords := []servicehistory.ServiceRecord{
		{
			ID:          1,
			BookingID:   1,
			VehicleID:   1,
			Description: "Oil change",
			Cost:        150000,
			ReceiptURL:  "https://example.com/receipt1.pdf",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			BookingID:   3,
			VehicleID:   1,
			Description: "Brake service",
			Cost:        200000,
			ReceiptURL:  "https://example.com/receipt2.pdf",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("ListByVehicle", uint(1)).Return(expectedRecords, nil)
	mockRepo.On("ListByVehicle", uint(999)).Return([]servicehistory.ServiceRecord{}, nil)

	// Test successful list
	records, err := mockRepo.ListByVehicle(1)
	assert.NoError(t, err)
	assert.NotNil(t, records)
	assert.Len(t, records, 2)
	assert.Equal(t, uint(1), records[0].VehicleID)
	assert.Equal(t, uint(1), records[1].VehicleID)

	// Test empty list
	emptyRecords, err := mockRepo.ListByVehicle(999)
	assert.NoError(t, err)
	assert.Len(t, emptyRecords, 0)

	mockRepo.AssertExpectations(t)
}

func TestServiceRecord_Cost(t *testing.T) {
	tests := []struct {
		name string
		cost int
	}{
		{"Low cost service", 50000},
		{"Medium cost service", 150000},
		{"High cost service", 500000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.GreaterOrEqual(t, tt.cost, 0)
			assert.LessOrEqual(t, tt.cost, 10000000) // Reasonable upper limit
		})
	}
}
