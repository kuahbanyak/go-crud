package vehicle

import (
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OwnerID      uint   `gorm:"index" json:"owner_id"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string ` json:"license_plate"`
	VIN          string ` json:"vin"`
	Mileage      int    `json:"mileage"`
}
