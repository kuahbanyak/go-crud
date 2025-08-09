package booking

import (
    "time"
    "gorm.io/gorm"
)

type BookingStatus string

const (
    StatusScheduled BookingStatus = "scheduled"
    StatusInProgress BookingStatus = "in_progress"
    StatusCompleted BookingStatus = "completed"
    StatusCanceled BookingStatus = "canceled"
)

type Booking struct {
    ID uint `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    VehicleID uint `json:"vehicle_id"`
    CustomerID uint `json:"customer_id"`
    MechanicID *uint `json:"mechanic_id"`
    ScheduledAt time.Time `json:"scheduled_at"`
    DurationMin int `json:"duration_min"`
    Status BookingStatus `gorm:"type:varchar(30);default:'scheduled'" json:"status"`
    Notes string `json:"notes"`
}
