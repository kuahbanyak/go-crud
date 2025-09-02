package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *entities.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Booking, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.Booking, error)
	GetByMechanicID(ctx context.Context, mechanicID uuid.UUID) ([]*entities.Booking, error)
	GetByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]*entities.Booking, error)
	GetByStatus(ctx context.Context, status entities.BookingStatus) ([]*entities.Booking, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.Booking, error)
	Update(ctx context.Context, booking *entities.Booking) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Booking, error)
}
