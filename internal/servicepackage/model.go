package servicepackage

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/booking"
	"github.com/kuahbanyak/go-crud/internal/scheduling"
	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
	"gorm.io/gorm"
)

type ServicePackage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
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

type PackageServiceType struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PackageID uint           `json:"package_id" gorm:"index"`
	Package   ServicePackage `gorm:"foreignKey:PackageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ServiceTypeID uint                   `json:"service_type_id" gorm:"index"`
	ServiceType   scheduling.ServiceType `gorm:"foreignKey:ServiceTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	Quantity int `json:"quantity" gorm:"default:1"`
}

type ServiceCategory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
}

type VehicleServiceHistory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID uint            `json:"vehicle_id" gorm:"index"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	BookingID uint            `json:"booking_id"`
	Booking   booking.Booking `gorm:"foreignKey:BookingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ServiceTypeID *uint                   `json:"service_type_id"`
	ServiceType   *scheduling.ServiceType `gorm:"foreignKey:ServiceTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	PackageID *uint           `json:"package_id"`
	Package   *ServicePackage `gorm:"foreignKey:PackageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	MechanicID uint      `json:"mechanic_id"`
	Mechanic   user.User `gorm:"foreignKey:MechanicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CompletedAt    time.Time `json:"completed_at"`
	Cost           int       `json:"cost"`
	Mileage        int       `json:"mileage"`
	Notes          string    `json:"notes"`
	Photos         string    `json:"photos"`         // JSON array of photo URLs
	QualityRating  int       `json:"quality_rating"` // 1-5 scale
	CustomerReview string    `json:"customer_review"`
}
