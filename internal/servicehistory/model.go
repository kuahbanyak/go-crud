package servicehistory

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRecord struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	BookingID   string `gorm:"type:uniqueidentifier" json:"booking_id"`
	VehicleID   string `gorm:"type:uniqueidentifier" json:"vehicle_id"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	ReceiptURL  string `json:"receipt_url"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *ServiceRecord) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}
