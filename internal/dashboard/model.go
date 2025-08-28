package dashboard

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
)

type VehicleHealthScore struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID string          `gorm:"type:uniqueidentifier;index" json:"vehicle_id"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	OverallScore    int       `json:"overall_score"` // 0-100
	EngineScore     int       `json:"engine_score"`
	BrakeScore      int       `json:"brake_score"`
	TireScore       int       `json:"tire_score"`
	BatteryScore    int       `json:"battery_score"`
	LastUpdated     time.Time `json:"last_updated"`
	Recommendations string    `json:"recommendations"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (v *VehicleHealthScore) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	return
}

type MaintenanceRecommendation struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID string          `gorm:"type:uniqueidentifier;index" json:"vehicle_id"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	Title         string `json:"title"`
	Description   string `json:"description"`
	Priority      string `json:"priority"` // low, medium, high, critical
	EstimatedCost int    `json:"estimated_cost"`
	DueInDays     int    `json:"due_in_days"`
	Category      string `json:"category"`
	IsCompleted   bool   `json:"is_completed" gorm:"default:false"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *MaintenanceRecommendation) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

type CustomerBudget struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CustomerID string    `gorm:"type:uniqueidentifier;index" json:"customer_id"`
	Customer   user.User `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	VehicleID *string          `gorm:"type:uniqueidentifier" json:"vehicle_id"`
	Vehicle   *vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	MonthlyBudget  int `json:"monthly_budget"`
	YearlyBudget   int `json:"yearly_budget"`
	SpentThisMonth int `json:"spent_this_month"`
	SpentThisYear  int `json:"spent_this_year"`
	AlertThreshold int `json:"alert_threshold"` // percentage
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *CustomerBudget) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

type ServiceTimeline struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	ServiceType string    `json:"service_type"`
	Cost        int       `json:"cost"`
	Mileage     int       `json:"mileage"`
	Notes       string    `json:"notes"`
	Status      string    `json:"status"`
	Photos      []string  `json:"photos"`
}
