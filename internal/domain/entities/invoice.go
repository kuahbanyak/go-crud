package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
)

type Invoice struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	WaitingListID *uuid.UUID     `gorm:"type:uuid" json:"waiting_list_id,omitempty"`
	CustomerID    uuid.UUID      `gorm:"type:uuid;not null" json:"customer_id"`
	Amount        int            `json:"amount"`
	TaxAmount     int            `json:"tax_amount"`
	TotalAmount   int            `json:"total_amount"`
	Status        InvoiceStatus  `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PDFURL        string         `json:"pdf_url,omitempty"`
	DueDate       *time.Time     `json:"due_date,omitempty"`
	PaidAt        *time.Time     `json:"paid_at,omitempty"`
	Notes         string         `json:"notes,omitempty"`

	// Relationships
	WaitingList *WaitingList `gorm:"foreignKey:WaitingListID" json:"waiting_list,omitempty"`
	Customer    *User        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
