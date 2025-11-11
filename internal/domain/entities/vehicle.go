package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID           types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
	OwnerID      types.MSSQLUUID `gorm:"type:uniqueidentifier;index" json:"owner_id"`
	Brand        string          `json:"brand"`
	Model        string          `json:"model"`
	Year         int             `json:"year"`
	LicensePlate string          `json:"license_plate"`
	VIN          string          `json:"vin"`
	Mileage      int             `json:"mileage"`
	Owner        User            `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	WaitingLists []WaitingList   `gorm:"foreignKey:VehicleID" json:"waiting_lists,omitempty"`
}

func (v *Vehicle) BeforeCreate(_ *gorm.DB) error {
	if v.ID.String() == "00000000-0000-0000-0000-000000000000" {
		v.ID = types.NewMSSQLUUID()
	}
	return nil
}
