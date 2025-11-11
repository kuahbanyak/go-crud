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
	ServiceDate    time.Time         `gorm:"uniqueIndex:idx_queue_date;not null" json:"service_date"`
	ServiceType    string            `gorm:"type:varchar(100);not null" json:"service_type"`
	EstimatedTime  int               `json:"estimated_time"` // in minutes
	Status         WaitingListStatus `gorm:"type:varchar(30);default:'waiting'" json:"status"`
	CalledAt       *time.Time        `json:"called_at,omitempty"`
	ServiceStartAt *time.Time        `json:"service_start_at,omitempty"`
	ServiceEndAt   *time.Time        `json:"service_end_at,omitempty"`
	Notes          string            `gorm:"type:text" json:"notes"`
	Vehicle        Vehicle           `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Customer       User              `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (w *WaitingList) BeforeCreate(_ *gorm.DB) error {
	if w.ID.String() == "00000000-0000-0000-0000-000000000000" {
		w.ID = types.NewMSSQLUUID()
	}
	return nil
}
