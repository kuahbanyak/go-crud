package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type ServiceItem struct {
	ID              types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
	Name            string          `gorm:"type:varchar(100);not null" json:"name"`
	Description     string          `gorm:"type:text" json:"description"`
	Category        string          `gorm:"type:varchar(50);not null" json:"category"` // Engine, Brakes, Tires, etc.
	EstimatedTime   int             `json:"estimated_time"`                            // in minutes
	EstimatedCost   float64         `json:"estimated_cost"`
	IsActive        bool            `gorm:"default:true" json:"is_active"`
	DisplayOrder    int             `gorm:"default:0" json:"display_order"`
	RequiresBooking bool            `gorm:"default:true" json:"requires_booking"`
}

func (s *ServiceItem) BeforeCreate(_ *gorm.DB) error {
	if s.ID.String() == "00000000-0000-0000-0000-000000000000" {
		s.ID = types.NewMSSQLUUID()
	}
	return nil
}

func (ServiceItem) TableName() string {
	return "service_items"
}
