package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceItem struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"type:varchar(100);not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	Category        string         `gorm:"type:varchar(50);not null" json:"category"` // Engine, Brakes, Tires, etc.
	EstimatedTime   int            `json:"estimated_time"`                            // in minutes
	EstimatedCost   float64        `json:"estimated_cost"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	DisplayOrder    int            `gorm:"default:0" json:"display_order"`
	RequiresBooking bool           `gorm:"default:true" json:"requires_booking"`
}

func (s *ServiceItem) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (ServiceItem) TableName() string {
	return "service_items"
}
