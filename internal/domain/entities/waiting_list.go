package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type WaitingListStatus string

const (
	WaitingListStatusWaiting   WaitingListStatus = "waiting"
	WaitingListStatusCalled    WaitingListStatus = "called"
	WaitingListStatusInService WaitingListStatus = "in_service"
	WaitingListStatusCompleted WaitingListStatus = "completed"
	WaitingListStatusCanceled  WaitingListStatus = "canceled"
	WaitingListStatusNoShow    WaitingListStatus = "no_show"
)

type WaitingList struct {
	ID             types.MSSQLUUID   `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      gorm.DeletedAt    `gorm:"index" json:"-"`
	QueueNumber    int               `gorm:"uniqueIndex:idx_queue_date;not null" json:"queue_number"`
	VehicleID      types.MSSQLUUID   `gorm:"type:uniqueidentifier;not null" json:"vehicle_id"`
	CustomerID     types.MSSQLUUID   `gorm:"type:uniqueidentifier;not null" json:"customer_id"`
	MechanicID     *types.MSSQLUUID  `gorm:"type:uniqueidentifier" json:"mechanic_id,omitempty"`     // Mechanic assigned to service
	ServiceItemID  *types.MSSQLUUID  `gorm:"type:uniqueidentifier" json:"service_item_id,omitempty"` // Selected service item
	ServiceDate    time.Time         `gorm:"uniqueIndex:idx_queue_date;not null" json:"service_date"`
	ServiceType    string            `gorm:"type:varchar(100);not null" json:"service_type"`
	EstimatedTime  int               `json:"estimated_time"` // in minutes, default 0, updated by mechanic
	EstimatedCost  float64           `json:"estimated_cost"` // estimated cost from service item
	Status         WaitingListStatus `gorm:"type:varchar(30);default:'waiting'" json:"status"`
	CalledAt       *time.Time        `json:"called_at,omitempty"`
	ServiceStartAt *time.Time        `json:"service_start_at,omitempty"`
	ServiceEndAt   *time.Time        `json:"service_end_at,omitempty"`
	Notes          string            `gorm:"type:text" json:"notes"`          // Customer notes
	MechanicNotes  string            `gorm:"type:text" json:"mechanic_notes"` // Mechanic notes after inspection
	Vehicle        Vehicle           `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Customer       User              `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Mechanic       *User             `gorm:"foreignKey:MechanicID" json:"mechanic,omitempty"`        // Mechanic who serviced
	ServiceItem    *ServiceItem      `gorm:"foreignKey:ServiceItemID" json:"service_item,omitempty"` // Selected service
}

func (w *WaitingList) BeforeCreate(_ *gorm.DB) error {
	if w.ID.String() == "00000000-0000-0000-0000-000000000000" {
		w.ID = types.NewMSSQLUUID()
	}
	return nil
}
