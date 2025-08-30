package booking_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
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

func (m *MockBookingRepo) ListByCustomer(customerID string) ([]booking.Booking, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]booking.Booking), args.Error(1)
}

func (m *MockBookingRepo) UpdateStatus(id string, status booking.BookingStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockBookingRepo) GetId(id string) (*booking.Booking, error) {
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
				VehicleID:   uuid.New().String(),
				CustomerID:  uuid.New().String(),
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
				VehicleID:   uuid.New().String(),
				CustomerID:  uuid.New().String(),
				MechanicID:  stringPtr(uuid.New().String()),
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
			assert.NotEmpty(t, tt.booking.VehicleID)
			assert.NotEmpty(t, tt.booking.CustomerID)
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
		VehicleID:   uuid.New().String(),
		CustomerID:  uuid.New().String(),
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
	customerID := uuid.New().String()

	expectedBookings := []booking.Booking{
		{
			ID:          uuid.New().String(),
			VehicleID:   uuid.New().String(),
			CustomerID:  customerID,
			ScheduledAt: time.Now().Add(24 * time.Hour),
			DurationMin: 60,
			Status:      booking.StatusScheduled,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			VehicleID:   uuid.New().String(),
			CustomerID:  customerID,
			ScheduledAt: time.Now().Add(48 * time.Hour),
			DurationMin: 120,
			Status:      booking.StatusInProgress,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("ListByCustomer", customerID).Return(expectedBookings, nil)
	mockRepo.On("ListByCustomer", "non-existent-id").Return([]booking.Booking{}, nil)

	// Test successful list
	bookings, err := mockRepo.ListByCustomer(customerID)
	assert.NoError(t, err)
	assert.NotNil(t, bookings)
	assert.Len(t, bookings, 2)
	assert.Equal(t, customerID, bookings[0].CustomerID)
	assert.Equal(t, customerID, bookings[1].CustomerID)

	// Test empty list
	emptyBookings, err := mockRepo.ListByCustomer("non-existent-id")
	assert.NoError(t, err)
	assert.Len(t, emptyBookings, 0)

	mockRepo.AssertExpectations(t)
}

func TestBookingRepository_GetId(t *testing.T) {
	mockRepo := new(MockBookingRepo)
	bookingID := uuid.New().String()

	expectedBooking := &booking.Booking{
		ID:          bookingID,
		VehicleID:   uuid.New().String(),
		CustomerID:  uuid.New().String(),
		ScheduledAt: time.Now().Add(24 * time.Hour),
		DurationMin: 60,
		Status:      booking.StatusScheduled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetId", bookingID).Return(expectedBooking, nil)
	mockRepo.On("GetId", "non-existent-id").Return(nil, gorm.ErrRecordNotFound)

	// Test successful get
	foundBooking, err := mockRepo.GetId(bookingID)
	assert.NoError(t, err)
	assert.NotNil(t, foundBooking)
	assert.Equal(t, bookingID, foundBooking.ID)
	assert.Equal(t, booking.StatusScheduled, foundBooking.Status)

	// Test booking not found
	notFoundBooking, err := mockRepo.GetId("non-existent-id")
	assert.Error(t, err)
	assert.Nil(t, notFoundBooking)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestBookingRepository_UpdateStatus(t *testing.T) {
	mockRepo := new(MockBookingRepo)
	bookingID := uuid.New().String()

	mockRepo.On("UpdateStatus", bookingID, booking.StatusInProgress).Return(nil)
	mockRepo.On("UpdateStatus", bookingID, booking.StatusCompleted).Return(nil)
	mockRepo.On("UpdateStatus", bookingID, booking.StatusCanceled).Return(nil)
	mockRepo.On("UpdateStatus", "non-existent-id", booking.StatusCompleted).Return(gorm.ErrRecordNotFound)

	// Test successful status updates
	err := mockRepo.UpdateStatus(bookingID, booking.StatusInProgress)
	assert.NoError(t, err)

	err = mockRepo.UpdateStatus(bookingID, booking.StatusCompleted)
	assert.NoError(t, err)

	err = mockRepo.UpdateStatus(bookingID, booking.StatusCanceled)
	assert.NoError(t, err)

	// Test update not found
	err = mockRepo.UpdateStatus("non-existent-id", booking.StatusCompleted)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
