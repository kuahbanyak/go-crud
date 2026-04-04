package dto

import (
	"time"

	"github.com/google/uuid"
)

type TakeQueueRequest struct {
	VehicleID     *uuid.UUID      `json:"vehicle_id,omitempty"`
	NewVehicle    *CreateVehicleRequest `json:"new_vehicle,omitempty"`
	ServiceItemID *uuid.UUID      `json:"service_item_id,omitempty"` // Selected service item from list
	ServiceType   string                `json:"service_type" validate:"required"`
	ServiceDate   string                `json:"service_date" validate:"required"` // Changed to string for date-only format (YYYY-MM-DD)
	EstimatedTime int                   `json:"estimated_time"`                   // in minutes, default 0
	Notes         string                `json:"notes,omitempty"`
}

type UpdateWaitingListRequest struct {
	ServiceType   string `json:"service_type,omitempty"`
	EstimatedTime int    `json:"estimated_time,omitempty"` // Estimated time in minutes from mechanic
	Notes         string `json:"notes,omitempty"`
	MechanicNotes string `json:"mechanic_notes,omitempty"` // Notes from mechanic after inspection
	Status        string `json:"status,omitempty"`
}

type AssignMechanicRequest struct {
	QueueID uuid.UUID `json:"queue_id" validate:"required"`
}

type WaitingListResponse struct {
	ID             uuid.UUID  `json:"id"`
	QueueNumber    int              `json:"queue_number"`
	VehicleID      uuid.UUID  `json:"vehicle_id"`
	CustomerID     uuid.UUID  `json:"customer_id"`
	MechanicID     *uuid.UUID `json:"mechanic_id,omitempty"`
	ServiceDate    time.Time        `json:"service_date"`
	ServiceType    string           `json:"service_type"`
	EstimatedTime  int              `json:"estimated_time"`
	Status         string           `json:"status"`
	CalledAt       *time.Time       `json:"called_at,omitempty"`
	ServiceStartAt *time.Time       `json:"service_start_at,omitempty"`
	ServiceEndAt   *time.Time       `json:"service_end_at,omitempty"`
	Notes          string           `json:"notes"`
	MechanicNotes  string           `json:"mechanic_notes"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

type WaitingListWithDetailsResponse struct {
	ID             uuid.UUID  `json:"id"`
	QueueNumber    int              `json:"queue_number"`
	VehicleID      uuid.UUID  `json:"vehicle_id"`
	VehicleBrand   string           `json:"vehicle_brand,omitempty"`
	VehicleModel   string           `json:"vehicle_model,omitempty"`
	LicensePlate   string           `json:"license_plate,omitempty"`
	CustomerID     uuid.UUID  `json:"customer_id"`
	CustomerName   string           `json:"customer_name,omitempty"`
	CustomerPhone  string           `json:"customer_phone,omitempty"`
	MechanicID     *uuid.UUID `json:"mechanic_id,omitempty"`
	MechanicName   string           `json:"mechanic_name,omitempty"`
	ServiceDate    time.Time        `json:"service_date"`
	ServiceType    string           `json:"service_type"`
	EstimatedTime  int              `json:"estimated_time"`
	Status         string           `json:"status"`
	CalledAt       *time.Time       `json:"called_at,omitempty"`
	ServiceStartAt *time.Time       `json:"service_start_at,omitempty"`
	ServiceEndAt   *time.Time       `json:"service_end_at,omitempty"`
	Notes          string           `json:"notes"`
	MechanicNotes  string           `json:"mechanic_notes"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

type WaitingListListResponse struct {
	WaitingLists []WaitingListWithDetailsResponse `json:"waiting_lists"`
	Total        int                              `json:"total"`
	Date         string                           `json:"date"`
}

type QueueStatusResponse struct {
	CurrentQueue     int       `json:"current_queue"`
	TotalToday       int       `json:"total_today"`
	AverageWaitTime  int       `json:"average_wait_time_minutes"`
	EstimatedWaitMin int       `json:"estimated_wait_minutes"`
	ServiceDate      time.Time `json:"service_date"`
}

type QueueAvailabilityResponse struct {
	Date              string `json:"date"`
	IsAvailable       bool   `json:"is_available"`
	CurrentQueueCount int    `json:"current_queue_count"`
	MaxDailyLimit     int    `json:"max_daily_limit"`
	RemainingSlots    int    `json:"remaining_slots"`
	Message           string `json:"message"`
}

type ServiceProgressResponse struct {
	ID            uuid.UUID `json:"id"`
	QueueNumber   int             `json:"queue_number"`
	Status        string          `json:"status"`
	StatusMessage string          `json:"status_message"`
	VehicleBrand  string          `json:"vehicle_brand,omitempty"`
	VehicleModel  string          `json:"vehicle_model,omitempty"`
	LicensePlate  string          `json:"license_plate,omitempty"`
	CustomerName  string          `json:"customer_name,omitempty"`
	CustomerPhone string          `json:"customer_phone,omitempty"`
	MechanicName  string          `json:"mechanic_name,omitempty"`
	ServiceType   string          `json:"service_type"`
	ServiceDate   time.Time       `json:"service_date"`
	EstimatedTime int             `json:"estimated_time_minutes"`
	QueuePosition int             `json:"queue_position"`
	PeopleAhead   int             `json:"people_ahead"`
	EstimatedWait int             `json:"estimated_wait_minutes"`
	Timeline      Timeline        `json:"timeline"`
	Notes         string          `json:"notes,omitempty"`
}

type Timeline struct {
	QueueTakenAt   time.Time  `json:"queue_taken_at"`
	CalledAt       *time.Time `json:"called_at,omitempty"`
	ServiceStartAt *time.Time `json:"service_start_at,omitempty"`
	ServiceEndAt   *time.Time `json:"service_end_at,omitempty"`
}

type TicketCountResponse struct {
	TotalTickets         int64  `json:"total_tickets"`
	ActiveTickets        int64  `json:"active_tickets"`
	CompletedTickets     int64  `json:"completed_tickets"`
	CompletedOnlyTickets int64  `json:"completed_only_tickets,omitempty"`
	CanceledTickets      int64  `json:"canceled_tickets,omitempty"`
	NoShowTickets        int64  `json:"no_show_tickets,omitempty"`
	DeletedTickets       int64  `json:"deleted_tickets,omitempty"`
	SystemActive         bool   `json:"system_active,omitempty"`
	AcceptingBookings    bool   `json:"accepting_bookings,omitempty"`
	Available            bool   `json:"available,omitempty"`
	RemainingTickets     int    `json:"remaining_tickets,omitempty"`
	MaxTicketsPerWeek    int    `json:"max_tickets_per_week,omitempty"`
	WeekStart            string `json:"week_start,omitempty"`
	WeekEnd              string `json:"week_end,omitempty"`
	CurrentDate          string `json:"date,omitempty"`
	Message              string `json:"message"`
}
