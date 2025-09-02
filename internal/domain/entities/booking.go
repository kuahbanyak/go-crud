package entities

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
	ID        uuid.UUID      `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID   uuid.UUID     `gorm:"type:uniqueidentifier" json:"vehicle_id"`
	CustomerID  uuid.UUID     `gorm:"type:uniqueidentifier" json:"customer_id"`
	MechanicID  *uuid.UUID    `gorm:"type:uniqueidentifier" json:"mechanic_id"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	DurationMin int           `json:"duration_min"`
	Status      BookingStatus `gorm:"type:varchar(30);default:'scheduled'" json:"status"`
	Notes       string        `json:"notes"`

	Vehicle  Vehicle   `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Customer User      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Mechanic *User     `gorm:"foreignKey:MechanicID" json:"mechanic,omitempty"`
	Invoices []Invoice `gorm:"foreignKey:BookingID" json:"invoices,omitempty"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
