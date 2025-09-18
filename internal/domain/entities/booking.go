package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
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
	ID        types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	VehicleID   types.MSSQLUUID  `gorm:"type:uniqueidentifier" json:"vehicle_id"`
	CustomerID  types.MSSQLUUID  `gorm:"type:uniqueidentifier" json:"customer_id"`
	MechanicID  *types.MSSQLUUID `gorm:"type:uniqueidentifier" json:"mechanic_id"`
	ScheduledAt time.Time        `json:"scheduled_at"`
	DurationMin int              `json:"duration_min"`
	Status      BookingStatus    `gorm:"type:varchar(30);default:'scheduled'" json:"status"`
	Notes       string           `json:"notes"`

	Vehicle  Vehicle   `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Customer User      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Mechanic *User     `gorm:"foreignKey:MechanicID" json:"mechanic,omitempty"`
	Invoices []Invoice `gorm:"foreignKey:BookingID" json:"invoices,omitempty"`
}

func (i *Booking) BeforeCreate(_ *gorm.DB) error {
	if i.ID.String() == "00000000-0000-0000-0000-000000000000" {
		i.ID = types.NewMSSQLUUID()
	}
	return nil
}
