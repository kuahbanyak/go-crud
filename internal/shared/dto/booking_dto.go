package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	VehicleID   uuid.UUID `json:"vehicle_id" validate:"required"`
	ServiceType string    `json:"service_type" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Notes       string    `json:"notes,omitempty"`
}

type UpdateBookingRequest struct {
	VehicleID   uuid.UUID `json:"vehicle_id,omitempty"`
	ServiceType string    `json:"service_type,omitempty"`
	StartDate   time.Time `json:"start_date,omitempty"`
	EndDate     time.Time `json:"end_date,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	Status      string    `json:"status,omitempty"`
}

type BookingResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	VehicleID   uuid.UUID `json:"vehicle_id"`
	ServiceType string    `json:"service_type"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Notes       string    `json:"notes"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookingListResponse struct {
	Bookings []BookingResponse `json:"bookings"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}
