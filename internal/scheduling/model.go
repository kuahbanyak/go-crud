package scheduling

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
	"gorm.io/gorm"
)

type AvailabilityStatus string

const (
	StatusAvailable   AvailabilityStatus = "available"
	StatusBooked      AvailabilityStatus = "booked"
	StatusUnavailable AvailabilityStatus = "unavailable"
)

type MechanicAvailability struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	MechanicID uint      `json:"mechanic_id" gorm:"index"`
	Mechanic   user.User `gorm:"foreignKey:MechanicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	Date      time.Time          `json:"date" gorm:"type:date"`
	StartTime time.Time          `json:"start_time"`
	EndTime   time.Time          `json:"end_time"`
	Status    AvailabilityStatus `json:"status" gorm:"type:varchar(20);default:'available'"`
	Notes     string             `json:"notes"`
}

type ServiceType struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name             string `json:"name" gorm:"not null"`
	Description      string `json:"description"`
	EstimatedMinutes int    `json:"estimated_minutes"`
	BasePrice        int    `json:"base_price"`
	Category         string `json:"category"`
}

type MaintenanceReminder struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID uint            `json:"vehicle_id" gorm:"index"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ServiceTypeID uint        `json:"service_type_id"`
	ServiceType   ServiceType `gorm:"foreignKey:ServiceTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	DueDate     time.Time  `json:"due_date"`
	DueMileage  int        `json:"due_mileage"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at"`
}

type BookingWaitlist struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CustomerID uint      `json:"customer_id" gorm:"index"`
	Customer   user.User `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	VehicleID uint            `json:"vehicle_id"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ServiceTypeID uint        `json:"service_type_id"`
	ServiceType   ServiceType `gorm:"foreignKey:ServiceTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	PreferredDate time.Time `json:"preferred_date"`
	Notes         string    `json:"notes"`
	IsNotified    bool      `json:"is_notified" gorm:"default:false"`
}
