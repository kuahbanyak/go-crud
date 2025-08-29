package test

import (
	"testing"
	"time"

	"github.com/kuahbanyak/go-crud/internal/booking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockBookingRepo is a mock implementation of booking.Repository
type MockBookingRepo struct {
	mock.Mock
}

func (m *MockBookingRepo) Create(b *booking.Booking) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *MockBookingRepo) ListByCustomer(customer uint) ([]booking.Booking, error) {
	args := m.Called(customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]booking.Booking), args.Error(1)
}

func (m *MockBookingRepo) UpdateStatus(id uint, status booking.BookingStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockBookingRepo) GetId(id uint) (*booking.Booking, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*booking.Booking), args.Error(1)
}

func TestBookingModel_Validation(t *testing.T) {
	scheduledTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name    string
		booking booking.Booking
		valid   bool
	}{
		{
			name: "Valid booking",
			booking: booking.Booking{
				VehicleID:   1,
				CustomerID:  1,
				MechanicID:  nil,
				ScheduledAt: scheduledTime,
				DurationMin: 60,
				Status:      booking.StatusScheduled,
				Notes:       "Regular maintenance",
			},
			valid: true,
		},
		{
			name: "Valid booking with mechanic",
			booking: booking.Booking{
				VehicleID:   2,
				CustomerID:  1,
				MechanicID:  uintPtr(2),
				ScheduledAt: scheduledTime,
				DurationMin: 120,
				Status:      booking.StatusInProgress,
				Notes:       "Engine repair",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotZero(t, tt.booking.VehicleID)
			assert.NotZero(t, tt.booking.CustomerID)
			assert.NotZero(t, tt.booking.ScheduledAt)
			assert.Greater(t, tt.booking.DurationMin, 0)
			assert.Contains(t, []booking.BookingStatus{
				booking.StatusScheduled,
				booking.StatusInProgress,
				booking.StatusCompleted,
				booking.StatusCanceled,
			}, tt.booking.Status)
		})
	}
}

func TestBookingStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status booking.BookingStatus
	}{
		{"Scheduled status", booking.StatusScheduled},
		{"In progress status", booking.StatusInProgress},
		{"Completed status", booking.StatusCompleted},
		{"Canceled status", booking.StatusCanceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, string(tt.status))
			assert.Contains(t, []booking.BookingStatus{
				booking.StatusScheduled,
				booking.StatusInProgress,
				booking.StatusCompleted,
				booking.StatusCanceled,
			}, tt.status)
		})
	}
}

func TestBookingRepository_Create(t *testing.T) {
	mockRepo := new(MockBookingRepo)

	testBooking := &booking.Booking{
		VehicleID:   1,
		CustomerID:  1,
		ScheduledAt: time.Now().Add(24 * time.Hour),
		DurationMin: 60,
		Status:      booking.StatusScheduled,
		Notes:       "Regular maintenance",
	}

	mockRepo.On("Create", testBooking).Return(nil)

	err := mockRepo.Create(testBooking)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookingRepository_ListByCustomer(t *testing.T) {
	mockRepo := new(MockBookingRepo)

	expectedBookings := []booking.Booking{
		{
			ID:          1,
			VehicleID:   1,
			CustomerID:  1,
			ScheduledAt: time.Now().Add(24 * time.Hour),
			DurationMin: 60,
			Status:      booking.StatusScheduled,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			VehicleID:   2,
			CustomerID:  1,
			ScheduledAt: time.Now().Add(48 * time.Hour),
			DurationMin: 120,
			Status:      booking.StatusInProgress,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("ListByCustomer", uint(1)).Return(expectedBookings, nil)
	mockRepo.On("ListByCustomer", uint(999)).Return([]booking.Booking{}, nil)

	// Test successful list
	bookings, err := mockRepo.ListByCustomer(1)
	assert.NoError(t, err)
	assert.NotNil(t, bookings)
	assert.Len(t, bookings, 2)
	assert.Equal(t, uint(1), bookings[0].CustomerID)
	assert.Equal(t, uint(1), bookings[1].CustomerID)

	// Test empty list
	emptyBookings, err := mockRepo.ListByCustomer(999)
	assert.NoError(t, err)
	assert.Len(t, emptyBookings, 0)

	mockRepo.AssertExpectations(t)
}

func TestBookingRepository_GetId(t *testing.T) {
	mockRepo := new(MockBookingRepo)

	expectedBooking := &booking.Booking{
		ID:          1,
		VehicleID:   1,
		CustomerID:  1,
		ScheduledAt: time.Now().Add(24 * time.Hour),
		DurationMin: 60,
		Status:      booking.StatusScheduled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetId", uint(1)).Return(expectedBooking, nil)
	mockRepo.On("GetId", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// Test successful get
	foundBooking, err := mockRepo.GetId(1)
	assert.NoError(t, err)
	assert.NotNil(t, foundBooking)
	assert.Equal(t, uint(1), foundBooking.ID)
	assert.Equal(t, booking.StatusScheduled, foundBooking.Status)

	// Test booking not found
	notFoundBooking, err := mockRepo.GetId(999)
	assert.Error(t, err)
	assert.Nil(t, notFoundBooking)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestBookingRepository_UpdateStatus(t *testing.T) {
	mockRepo := new(MockBookingRepo)

	mockRepo.On("UpdateStatus", uint(1), booking.StatusInProgress).Return(nil)
	mockRepo.On("UpdateStatus", uint(1), booking.StatusCompleted).Return(nil)
	mockRepo.On("UpdateStatus", uint(1), booking.StatusCanceled).Return(nil)
	mockRepo.On("UpdateStatus", uint(999), booking.StatusCompleted).Return(gorm.ErrRecordNotFound)

	// Test successful status updates
	err := mockRepo.UpdateStatus(1, booking.StatusInProgress)
	assert.NoError(t, err)

	err = mockRepo.UpdateStatus(1, booking.StatusCompleted)
	assert.NoError(t, err)

	err = mockRepo.UpdateStatus(1, booking.StatusCanceled)
	assert.NoError(t, err)

	// Test update not found
	err = mockRepo.UpdateStatus(999, booking.StatusCompleted)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

// Helper function to create uint pointer
func uintPtr(u uint) *uint {
	return &u
}
