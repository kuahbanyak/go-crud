package dashboard

import (
	"time"

	"gorm.io/gorm"

	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
)

type VehicleHealthScore struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID uint            `json:"vehicle_id" gorm:"index"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	OverallScore    int       `json:"overall_score"` // 0-100
	EngineScore     int       `json:"engine_score"`
	BrakeScore      int       `json:"brake_score"`
	TireScore       int       `json:"tire_score"`
	BatteryScore    int       `json:"battery_score"`
	LastUpdated     time.Time `json:"last_updated"`
	Recommendations string    `json:"recommendations"`
}

type MaintenanceRecommendation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID uint            `json:"vehicle_id" gorm:"index"`
	Vehicle   vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	Title         string `json:"title"`
	Description   string `json:"description"`
	Priority      string `json:"priority"` // low, medium, high, critical
	EstimatedCost int    `json:"estimated_cost"`
	DueInDays     int    `json:"due_in_days"`
	Category      string `json:"category"`
	IsCompleted   bool   `json:"is_completed" gorm:"default:false"`
}

type CustomerBudget struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CustomerID uint      `json:"customer_id" gorm:"index"`
	Customer   user.User `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	VehicleID *uint            `json:"vehicle_id"`
	Vehicle   *vehicle.Vehicle `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	MonthlyBudget  int `json:"monthly_budget"`
	YearlyBudget   int `json:"yearly_budget"`
	SpentThisMonth int `json:"spent_this_month"`
	SpentThisYear  int `json:"spent_this_year"`
	AlertThreshold int `json:"alert_threshold"` // percentage
}

type ServiceTimeline struct {
	ID          uint      `json:"id"`
	Date        time.Time `json:"date"`
	ServiceType string    `json:"service_type"`
	Cost        int       `json:"cost"`
	Mileage     int       `json:"mileage"`
	Notes       string    `json:"notes"`
	Status      string    `json:"status"`
	Photos      []string  `json:"photos"`
}
