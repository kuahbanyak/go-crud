package mssql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"gorm.io/gorm"
)

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) repositories.BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *entities.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Preload("Mechanic").
		Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Mechanic").
		Where("customer_id = ?", customerID).
		Order("scheduled_at DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetByMechanicID(ctx context.Context, mechanicID uuid.UUID) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Where("mechanic_id = ?", mechanicID).
		Order("scheduled_at ASC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Mechanic").
		Where("vehicle_id = ?", vehicleID).
		Order("scheduled_at DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetByStatus(ctx context.Context, status entities.BookingStatus) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Preload("Mechanic").
		Where("status = ?", status).
		Order("scheduled_at ASC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Preload("Mechanic").
		Where("scheduled_at BETWEEN ? AND ?", startDate, endDate).
		Order("scheduled_at ASC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) Update(ctx context.Context, booking *entities.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}

func (r *bookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Booking{}, id).Error
}

func (r *bookingRepository) List(ctx context.Context, limit, offset int) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Preload("Mechanic").
		Limit(limit).Offset(offset).
		Order("scheduled_at DESC").
		Find(&bookings).Error
	return bookings, err
}
