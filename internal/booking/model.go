package booking

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	StatusScheduled  BookingStatus = "scheduled"
	StatusInProgress BookingStatus = "in_progress"
	StatusCompleted  BookingStatus = "completed"
	StatusCanceled   BookingStatus = "canceled"
)

type Booking struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID   string        `gorm:"type:uniqueidentifier" json:"vehicle_id"`
	CustomerID  string        `gorm:"type:uniqueidentifier" json:"customer_id"`
	MechanicID  *string       `gorm:"type:uniqueidentifier" json:"mechanic_id"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	DurationMin int           `json:"duration_min"`
	Status      BookingStatus `gorm:"type:varchar(30);default:'scheduled'" json:"status"`
	Notes       string        `json:"notes"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}
