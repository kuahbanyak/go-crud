package dto

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type TakeQueueRequest struct {
	VehicleID     types.MSSQLUUID `json:"vehicle_id" validate:"required"`
	ServiceType   string          `json:"service_type" validate:"required"`
	ServiceDate   string          `json:"service_date" validate:"required"` // Changed to string for date-only format (YYYY-MM-DD)
	EstimatedTime int             `json:"estimated_time"`                   // in minutes
	Notes         string          `json:"notes,omitempty"`
}

type UpdateWaitingListRequest struct {
	ServiceType   string `json:"service_type,omitempty"`
	EstimatedTime int    `json:"estimated_time,omitempty"`
	Notes         string `json:"notes,omitempty"`
	Status        string `json:"status,omitempty"`
}

type WaitingListResponse struct {
	ID             types.MSSQLUUID `json:"id"`
	QueueNumber    int             `json:"queue_number"`
	VehicleID      types.MSSQLUUID `json:"vehicle_id"`
	CustomerID     types.MSSQLUUID `json:"customer_id"`
	ServiceDate    time.Time       `json:"service_date"`
	ServiceType    string          `json:"service_type"`
	EstimatedTime  int             `json:"estimated_time"`
	Status         string          `json:"status"`
	CalledAt       *time.Time      `json:"called_at,omitempty"`
	ServiceStartAt *time.Time      `json:"service_start_at,omitempty"`
	ServiceEndAt   *time.Time      `json:"service_end_at,omitempty"`
	Notes          string          `json:"notes"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type WaitingListWithDetailsResponse struct {
	ID             types.MSSQLUUID `json:"id"`
	QueueNumber    int             `json:"queue_number"`
	VehicleID      types.MSSQLUUID `json:"vehicle_id"`
	VehicleBrand   string          `json:"vehicle_brand,omitempty"`
	VehicleModel   string          `json:"vehicle_model,omitempty"`
	LicensePlate   string          `json:"license_plate,omitempty"`
	CustomerID     types.MSSQLUUID `json:"customer_id"`
	CustomerName   string          `json:"customer_name,omitempty"`
	CustomerPhone  string          `json:"customer_phone,omitempty"`
	ServiceDate    time.Time       `json:"service_date"`
	ServiceType    string          `json:"service_type"`
	EstimatedTime  int             `json:"estimated_time"`
	Status         string          `json:"status"`
	CalledAt       *time.Time      `json:"called_at,omitempty"`
	ServiceStartAt *time.Time      `json:"service_start_at,omitempty"`
	ServiceEndAt   *time.Time      `json:"service_end_at,omitempty"`
	Notes          string          `json:"notes"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type WaitingListListResponse struct {
	WaitingLists []WaitingListWithDetailsResponse `json:"waiting_lists"`
	Total        int                              `json:"total"`
	Date         string                           `json:"date"`
}

type QueueStatusResponse struct {
	QueueNumber      int       `json:"queue_number"`
	Status           string    `json:"status"`
	CurrentlyServing int       `json:"currently_serving"`
	WaitingAhead     int       `json:"waiting_ahead"`
	EstimatedWaitMin int       `json:"estimated_wait_minutes"`
	ServiceDate      time.Time `json:"service_date"`
}

// ServiceProgressResponse provides detailed progress information for customers
type ServiceProgressResponse struct {
	ID            types.MSSQLUUID `json:"id"`
	QueueNumber   int             `json:"queue_number"`
	Status        string          `json:"status"`
	StatusMessage string          `json:"status_message"`
	VehicleBrand  string          `json:"vehicle_brand"`
	VehicleModel  string          `json:"vehicle_model"`
	LicensePlate  string          `json:"license_plate"`
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
