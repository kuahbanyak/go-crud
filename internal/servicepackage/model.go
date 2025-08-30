package servicepackage

import (
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/booking"
	"github.com/kuahbanyak/go-crud/internal/scheduling"
	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
	"gorm.io/gorm"
)

type ServicePackage struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Category    string `json:"category"`
	TotalPrice  int    `json:"total_price"`
	DiscountPct int    `json:"discount_pct"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}

func (s *ServicePackage) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

type PackageServiceType struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	PackageID string         `gorm:"type:uniqueidentifier;index" json:"package_id"`
	Package   ServicePackage `gorm:"foreignKey:PackageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ServiceTypeID string                 `gorm:"type:uniqueidentifier;index" json:"service_type_id"`
	ServiceType   scheduling.ServiceType `gorm:"foreignKey:ServiceTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	Quantity int `json:"quantity" gorm:"default:1"`
}

func (p *PackageServiceType) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

type ServiceCategory struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *ServiceCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

type VehicleServiceHistory struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID string          `gorm:"type:uniqueidentifier;index" json:"vehicle_id"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	BookingID string          `gorm:"type:uniqueidentifier" json:"booking_id"`
	Booking   booking.Booking `gorm:"foreignKey:BookingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ServiceTypeID *string                 `gorm:"type:uniqueidentifier" json:"service_type_id"`
	ServiceType   *scheduling.ServiceType `gorm:"foreignKey:ServiceTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	PackageID *string         `gorm:"type:uniqueidentifier" json:"package_id"`
	Package   *ServicePackage `gorm:"foreignKey:PackageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	MechanicID string    `gorm:"type:uniqueidentifier" json:"mechanic_id"`
	Mechanic   user.User `gorm:"foreignKey:MechanicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CompletedAt    time.Time `json:"completed_at"`
	Cost           int       `json:"cost"`
	Mileage        int       `json:"mileage"`
	Notes          string    `json:"notes"`
	Photos         string    `json:"photos"`         // JSON array of photo URLs
	QualityRating  int       `json:"quality_rating"` // 1-5 scale
	CustomerReview string    `json:"customer_review"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (v *VehicleServiceHistory) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	return
}
