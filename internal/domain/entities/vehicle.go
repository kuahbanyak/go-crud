package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID        uuid.UUID      `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OwnerID      uuid.UUID `gorm:"type:uniqueidentifier;index" json:"owner_id"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	LicensePlate string    `json:"license_plate"`
	VIN          string    `json:"vin"`
	Mileage      int       `json:"mileage"`

	Owner    User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Bookings []Booking `gorm:"foreignKey:VehicleID" json:"bookings,omitempty"`
}

func (v *Vehicle) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}
