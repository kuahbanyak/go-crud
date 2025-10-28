package dto
import (
	"time"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)
type CreateMaintenanceItemRequest struct {
	Category      string  `json:"category" validate:"required"` // e.g., "Engine", "Brakes"
	Name          string  `json:"name" validate:"required"`     // e.g., "Oil Change"
	Description   string  `json:"description"`
	EstimatedCost float64 `json:"estimated_cost"`
}
type AddDiscoveredItemRequest struct {
	WaitingListID    types.MSSQLUUID `json:"waiting_list_id" validate:"required"`
	Category         string          `json:"category" validate:"required"`
	Name             string          `json:"name" validate:"required"`
	Description      string          `json:"description" validate:"required"`
	Priority         string          `json:"priority" validate:"required,oneof=urgent high normal low"`
	EstimatedCost    float64         `json:"estimated_cost"`
	LaborHours       float64         `json:"labor_hours"`
	RequiresApproval bool            `json:"requires_approval"`
	ImageURL         string          `json:"image_url"`
	Notes            string          `json:"notes"`
}
type UpdateMaintenanceItemRequest struct {
	Status        string   `json:"status" validate:"omitempty,oneof=pending inspected approved rejected completed skipped"`
	Description   string   `json:"description"`
	EstimatedCost *float64 `json:"estimated_cost"`
	ActualCost    *float64 `json:"actual_cost"`
	LaborHours    *float64 `json:"labor_hours"`
	Priority      string   `json:"priority" validate:"omitempty,oneof=urgent high normal low"`
	Notes         string   `json:"notes"`
}
type ApproveMaintenanceItemRequest struct {
	ItemIDs []types.MSSQLUUID `json:"item_ids" validate:"required,min=1"`
	Approve bool              `json:"approve"` // true = approve, false = reject
	Notes   string            `json:"notes"`   // Customer feedback
}
type MaintenanceItemResponse struct {
	ID               types.MSSQLUUID  `json:"id"`
	WaitingListID    types.MSSQLUUID  `json:"waiting_list_id"`
	MechanicID       *types.MSSQLUUID `json:"mechanic_id,omitempty"`
	MechanicName     string           `json:"mechanic_name,omitempty"`
	ItemType         string           `json:"item_type"`
	Status           string           `json:"status"`
	Category         string           `json:"category"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Priority         string           `json:"priority"`
	EstimatedCost    float64          `json:"estimated_cost"`
	ActualCost       float64          `json:"actual_cost"`
	LaborHours       float64          `json:"labor_hours"`
	RequiresApproval bool             `json:"requires_approval"`
	ImageURL         string           `json:"image_url,omitempty"`
	Notes            string           `json:"notes"`
	InspectedAt      *time.Time       `json:"inspected_at,omitempty"`
	ApprovedAt       *time.Time       `json:"approved_at,omitempty"`
	CompletedAt      *time.Time       `json:"completed_at,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}
type MaintenanceItemListResponse struct {
	Items                []MaintenanceItemResponse `json:"items"`
	Total                int                       `json:"total"`
	TotalEstimatedCost   float64                   `json:"total_estimated_cost"`
	TotalActualCost      float64                   `json:"total_actual_cost"`
	PendingApprovalCount int                       `json:"pending_approval_count"`
	CompletedCount       int                       `json:"completed_count"`
}
type MaintenanceInspectionSummary struct {
	WaitingListID      types.MSSQLUUID           `json:"waiting_list_id"`
	QueueNumber        int                       `json:"queue_number"`
	VehicleBrand       string                    `json:"vehicle_brand"`
	VehicleModel       string                    `json:"vehicle_model"`
	LicensePlate       string                    `json:"license_plate"`
	InitialItems       []MaintenanceItemResponse `json:"initial_items"`
	DiscoveredItems    []MaintenanceItemResponse `json:"discovered_items"`
	TotalEstimatedCost float64                   `json:"total_estimated_cost"`
	RequiresApproval   bool                      `json:"requires_approval"`
	InspectedAt        time.Time                 `json:"inspected_at"`
}

