package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type MaintenanceItemStatus string
type MaintenanceItemType string

const (
	// Status of each maintenance item
	MaintenanceItemStatusPending   MaintenanceItemStatus = "pending"   // Initial request by customer
	MaintenanceItemStatusInspected MaintenanceItemStatus = "inspected" // Mechanic found this issue
	MaintenanceItemStatusApproved  MaintenanceItemStatus = "approved"  // Customer approved
	MaintenanceItemStatusRejected  MaintenanceItemStatus = "rejected"  // Customer rejected
	MaintenanceItemStatusCompleted MaintenanceItemStatus = "completed" // Work completed
	MaintenanceItemStatusSkipped   MaintenanceItemStatus = "skipped"   // Skipped for some reason

	// Type of maintenance item
	MaintenanceItemTypeInitial    MaintenanceItemType = "initial"    // Selected by customer when booking
	MaintenanceItemTypeDiscovered MaintenanceItemType = "discovered" // Found by mechanic during inspection
)

// MaintenanceItem represents individual maintenance/repair items for a service
type MaintenanceItem struct {
	ID        types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relations
	WaitingListID types.MSSQLUUID  `gorm:"type:uniqueidentifier;not null;index" json:"waiting_list_id"`
	MechanicID    *types.MSSQLUUID `gorm:"type:uniqueidentifier" json:"mechanic_id,omitempty"` // Who found/performed the work

	// Item Details
	ItemType    MaintenanceItemType   `gorm:"type:varchar(20);not null;default:'initial'" json:"item_type"`
	Status      MaintenanceItemStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Category    string                `gorm:"type:varchar(100);not null" json:"category"`        // e.g., "Engine", "Brakes", "Suspension"
	Name        string                `gorm:"type:varchar(200);not null" json:"name"`            // e.g., "Oil Change", "Brake Pad Replacement"
	Description string                `gorm:"type:text" json:"description"`                      // Detailed description/notes
	Priority    string                `gorm:"type:varchar(20);default:'normal'" json:"priority"` // urgent, high, normal, low

	// Pricing
	EstimatedCost float64 `gorm:"type:decimal(10,2);default:0" json:"estimated_cost"`
	ActualCost    float64 `gorm:"type:decimal(10,2);default:0" json:"actual_cost"`
	LaborHours    float64 `gorm:"type:decimal(5,2);default:0" json:"labor_hours"`

	// Timestamps for tracking
	InspectedAt *time.Time `json:"inspected_at,omitempty"` // When mechanic found it
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`  // When customer approved
	CompletedAt *time.Time `json:"completed_at,omitempty"` // When work was completed

	// Additional info
	RequiresApproval bool   `gorm:"default:false" json:"requires_approval"`       // Does this need customer approval?
	ImageURL         string `gorm:"type:varchar(500)" json:"image_url,omitempty"` // Photo of the issue
	Notes            string `gorm:"type:text" json:"notes"`                       // Mechanic notes

	// Relations
	WaitingList *WaitingList `gorm:"foreignKey:WaitingListID" json:"waiting_list,omitempty"`
	Mechanic    *User        `gorm:"foreignKey:MechanicID" json:"mechanic,omitempty"`
}

func (m *MaintenanceItem) BeforeCreate(_ *gorm.DB) error {
	if m.ID.String() == "00000000-0000-0000-0000-000000000000" {
		m.ID = types.NewMSSQLUUID()
	}
	return nil
}

// TableName specifies the table name for MaintenanceItem
func (MaintenanceItem) TableName() string {
	return "maintenance_items"
}
