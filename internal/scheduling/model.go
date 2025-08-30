package scheduling

import (
	"time"

	"github.com/google/uuid"
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
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	MechanicID string    `gorm:"type:uniqueidentifier;index" json:"mechanic_id"`
	Mechanic   user.User `gorm:"foreignKey:MechanicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	Date      time.Time          `json:"date" gorm:"type:date"`
	StartTime time.Time          `json:"start_time"`
	EndTime   time.Time          `json:"end_time"`
	Status    AvailabilityStatus `json:"status" gorm:"type:varchar(20);default:'available'"`
	Notes     string             `json:"notes"`
}

func (m *MechanicAvailability) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

type ServiceType struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name             string `json:"name" gorm:"not null"`
	Description      string `json:"description"`
	EstimatedMinutes int    `json:"estimated_minutes"`
	BasePrice        int    `json:"base_price"`
	Category         string `json:"category"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *ServiceType) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

type MaintenanceReminder struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID string          `gorm:"type:uniqueidentifier;index" json:"vehicle_id"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ServiceTypeID string      `gorm:"type:uniqueidentifier" json:"service_type_id"`
	ServiceType   ServiceType `gorm:"foreignKey:ServiceTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	DueDate     time.Time  `json:"due_date"`
	DueMileage  int        `json:"due_mileage"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *MaintenanceReminder) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

type BookingWaitlist struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CustomerID string    `gorm:"type:uniqueidentifier;index" json:"customer_id"`
	Customer   user.User `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	VehicleID string          `gorm:"type:uniqueidentifier" json:"vehicle_id"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ServiceTypeID string      `gorm:"type:uniqueidentifier" json:"service_type_id"`
	ServiceType   ServiceType `gorm:"foreignKey:ServiceTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	PreferredDate time.Time `json:"preferred_date"`
	Notes         string    `json:"notes"`
	IsNotified    bool      `json:"is_notified" gorm:"default:false"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (b *BookingWaitlist) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}
