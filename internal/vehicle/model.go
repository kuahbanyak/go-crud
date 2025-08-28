package vehicle

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OwnerID      string `gorm:"type:uniqueidentifier;index" json:"owner_id"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string ` json:"license_plate"`
	VIN          string ` json:"vin"`
	Mileage      int    `json:"mileage"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (v *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	return
}
