package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

type BookingService interface {
	CreateBooking(ctx context.Context, booking *entities.Booking) error
	GetBooking(ctx context.Context, id uuid.UUID) (*entities.Booking, error)
	UpdateBookingStatus(ctx context.Context, id uuid.UUID, status entities.BookingStatus) error
	AssignMechanic(ctx context.Context, bookingID, mechanicID uuid.UUID) error
	CancelBooking(ctx context.Context, id uuid.UUID) error
	GetCustomerBookings(ctx context.Context, customerID uuid.UUID) ([]*entities.Booking, error)
	GetMechanicBookings(ctx context.Context, mechanicID uuid.UUID) ([]*entities.Booking, error)
	GetBookingsByDate(ctx context.Context, date time.Time) ([]*entities.Booking, error)
}
