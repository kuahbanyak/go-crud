package repositories

import (
	"context"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *entities.Booking) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Booking, error)
	GetByCustomerID(ctx context.Context, customerID types.MSSQLUUID) ([]*entities.Booking, error)
	GetByMechanicID(ctx context.Context, mechanicID types.MSSQLUUID) ([]*entities.Booking, error)
	GetByVehicleID(ctx context.Context, vehicleID types.MSSQLUUID) ([]*entities.Booking, error)
	GetByStatus(ctx context.Context, status entities.BookingStatus) ([]*entities.Booking, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.Booking, error)
	Update(ctx context.Context, booking *entities.Booking) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Booking, error)
}
