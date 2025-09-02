package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type BookingUsecase struct {
	bookingRepo repositories.BookingRepository
	vehicleRepo repositories.VehicleRepository
	userRepo    repositories.UserRepository
}

func NewBookingUsecase(
	bookingRepo repositories.BookingRepository,
	vehicleRepo repositories.VehicleRepository,
	userRepo repositories.UserRepository,
) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo: bookingRepo,
		vehicleRepo: vehicleRepo,
		userRepo:    userRepo,
	}
}

func (b *BookingUsecase) CreateBooking(ctx context.Context, booking *entities.Booking) error {
	if b.vehicleRepo != nil {
		_, err := b.vehicleRepo.GetByID(ctx, booking.VehicleID)
		if err != nil {
			return errors.New("vehicle not found")
		}
	}

	_, err := b.userRepo.GetByID(ctx, types.FromUUID(booking.CustomerID))
	if err != nil {
		return errors.New("customer not found")
	}

	// Set default status
	booking.Status = entities.StatusScheduled

	return b.bookingRepo.Create(ctx, booking)
}

func (b *BookingUsecase) GetBooking(ctx context.Context, id uuid.UUID) (*entities.Booking, error) {
	return b.bookingRepo.GetByID(ctx, id)
}

func (b *BookingUsecase) GetBookingsByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*entities.Booking, error) {
	allBookings, err := b.bookingRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	// Apply pagination manually since repository doesn't support it
	start := offset
	if start >= len(allBookings) {
		return []*entities.Booking{}, nil
	}

	end := start + limit
	if end > len(allBookings) {
		end = len(allBookings)
	}

	return allBookings[start:end], nil
}

func (b *BookingUsecase) UpdateBooking(ctx context.Context, id uuid.UUID, booking *entities.Booking) error {
	existing, err := b.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("booking not found")
	}

	// Validate vehicle exists if VehicleID is being updated
	if booking.VehicleID != uuid.Nil && b.vehicleRepo != nil {
		_, err := b.vehicleRepo.GetByID(ctx, booking.VehicleID)
		if err != nil {
			return errors.New("vehicle not found")
		}
	}

	// Set the ID to ensure we're updating the right booking
	booking.ID = id
	return b.bookingRepo.Update(ctx, booking)
}

func (b *BookingUsecase) DeleteBooking(ctx context.Context, id uuid.UUID) error {
	existing, err := b.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("booking not found")
	}

	return b.bookingRepo.Delete(ctx, id)
}

func (b *BookingUsecase) AssignMechanic(ctx context.Context, bookingID, mechanicID uuid.UUID) error {
	booking, err := b.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	// Validate mechanic exists
	_, err = b.userRepo.GetByID(ctx, types.FromUUID(mechanicID))
	if err != nil {
		return errors.New("mechanic not found")
	}

	booking.MechanicID = &mechanicID
	return b.bookingRepo.Update(ctx, booking)
}

func (b *BookingUsecase) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.BookingStatus) error {
	booking, err := b.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	booking.Status = status
	return b.bookingRepo.Update(ctx, booking)
}
