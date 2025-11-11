package events

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

const (
	EventQueueNumberAssigned  = "event.queue.assigned"
	EventQueuePositionChanged = "event.queue.position_changed"
	EventServiceCalled        = "event.service.called"
	EventServiceStarted       = "event.service.started"
	EventServiceCompleted     = "event.service.completed"
	EventIssueDiscovered      = "event.approval.issue_discovered"
	EventApprovalNeeded       = "event.approval.needed"
	EventNotificationEmail    = "notification.email"
	EventNotificationSMS      = "notification.sms"
	EventNotificationPush     = "notification.push"
	EventPaymentRequired      = "payment.process"
	EventAuditLog             = "audit.log"
)

type BaseEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

type QueueNumberAssignedEvent struct {
	BaseEvent
	WaitingListID types.MSSQLUUID `json:"waiting_list_id"`
	CustomerID    types.MSSQLUUID `json:"customer_id"`
	CustomerEmail string          `json:"customer_email"`
	CustomerName  string          `json:"customer_name"`
	CustomerPhone string          `json:"customer_phone"`
	QueueNumber   int             `json:"queue_number"`
	ServiceDate   time.Time       `json:"service_date"`
	ServiceType   string          `json:"service_type"`
	VehicleBrand  string          `json:"vehicle_brand"`
	VehicleModel  string          `json:"vehicle_model"`
	LicensePlate  string          `json:"license_plate"`
}

type ServiceStartedEvent struct {
	BaseEvent
	WaitingListID types.MSSQLUUID `json:"waiting_list_id"`
	CustomerID    types.MSSQLUUID `json:"customer_id"`
	CustomerEmail string          `json:"customer_email"`
	CustomerName  string          `json:"customer_name"`
	CustomerPhone string          `json:"customer_phone"`
	QueueNumber   int             `json:"queue_number"`
	ServiceType   string          `json:"service_type"`
	VehicleBrand  string          `json:"vehicle_brand"`
	VehicleModel  string          `json:"vehicle_model"`
	StartedAt     time.Time       `json:"started_at"`
}

type IssueDiscoveredEvent struct {
	BaseEvent
	WaitingListID    types.MSSQLUUID `json:"waiting_list_id"`
	ItemID           types.MSSQLUUID `json:"item_id"`
	CustomerID       types.MSSQLUUID `json:"customer_id"`
	CustomerEmail    string          `json:"customer_email"`
	CustomerName     string          `json:"customer_name"`
	CustomerPhone    string          `json:"customer_phone"`
	MechanicID       types.MSSQLUUID `json:"mechanic_id"`
	MechanicName     string          `json:"mechanic_name"`
	Category         string          `json:"category"`
	ItemName         string          `json:"item_name"`
	Description      string          `json:"description"`
	Priority         string          `json:"priority"`
	EstimatedCost    float64         `json:"estimated_cost"`
	ImageURL         string          `json:"image_url,omitempty"`
	RequiresApproval bool            `json:"requires_approval"`
}

type ApprovalNeededEvent struct {
	BaseEvent
	WaitingListID types.MSSQLUUID `json:"waiting_list_id"`
	CustomerID    types.MSSQLUUID `json:"customer_id"`
	CustomerEmail string          `json:"customer_email"`
	CustomerName  string          `json:"customer_name"`
	CustomerPhone string          `json:"customer_phone"`
	ItemCount     int             `json:"item_count"`
	TotalCost     float64         `json:"total_cost"`
	ApprovalURL   string          `json:"approval_url"`
}

type ServiceCompletedEvent struct {
	BaseEvent
	WaitingListID  types.MSSQLUUID `json:"waiting_list_id"`
	CustomerID     types.MSSQLUUID `json:"customer_id"`
	CustomerEmail  string          `json:"customer_email"`
	CustomerName   string          `json:"customer_name"`
	CustomerPhone  string          `json:"customer_phone"`
	QueueNumber    int             `json:"queue_number"`
	ServiceType    string          `json:"service_type"`
	CompletedAt    time.Time       `json:"completed_at"`
	TotalCost      float64         `json:"total_cost"`
	ItemsCompleted int             `json:"items_completed"`
}

type EmailNotificationEvent struct {
	BaseEvent
	To           string                 `json:"to"`
	Subject      string                 `json:"subject"`
	Body         string                 `json:"body"`
	Template     string                 `json:"template,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	Priority     string                 `json:"priority"`
}

type SMSNotificationEvent struct {
	BaseEvent
	To       string `json:"to"`
	Message  string `json:"message"`
	Priority string `json:"priority"`
}

type PushNotificationEvent struct {
	BaseEvent
	UserID   types.MSSQLUUID        `json:"user_id"`
	Title    string                 `json:"title"`
	Body     string                 `json:"body"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Priority string                 `json:"priority"`
}

type AuditLogEvent struct {
	BaseEvent
	UserID     types.MSSQLUUID        `json:"user_id"`
	UserRole   string                 `json:"user_role"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id"`
	Changes    map[string]interface{} `json:"changes,omitempty"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
}
