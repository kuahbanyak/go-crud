package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

// CreateInvoiceRequest represents a request to create an invoice
type CreateInvoiceRequest struct {
	WaitingListID *uuid.UUID `json:"waiting_list_id,omitempty"`
	CustomerID    uuid.UUID  `json:"customer_id" validate:"required"`
	Amount        int        `json:"amount" validate:"required,gt=0"`
	TaxAmount     int        `json:"tax_amount" validate:"gte=0"`
	Notes         string     `json:"notes,omitempty"`
	DueDays       int        `json:"due_days" validate:"gte=0"`
}

// UpdateInvoiceRequest represents a request to update an invoice
type UpdateInvoiceRequest struct {
	Amount    *int                    `json:"amount,omitempty" validate:"omitempty,gt=0"`
	TaxAmount *int                    `json:"tax_amount,omitempty" validate:"omitempty,gte=0"`
	Status    *entities.InvoiceStatus `json:"status,omitempty"`
	Notes     *string                 `json:"notes,omitempty"`
	DueDate   *time.Time              `json:"due_date,omitempty"`
}

// PayInvoiceRequest represents a request to pay an invoice
type PayInvoiceRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required"`
	PaymentRef    string `json:"payment_ref,omitempty"`
}

// InvoiceResponse represents an invoice response
type InvoiceResponse struct {
	ID            uuid.UUID              `json:"id"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	WaitingListID *uuid.UUID             `json:"waiting_list_id,omitempty"`
	CustomerID    uuid.UUID              `json:"customer_id"`
	CustomerName  string                 `json:"customer_name,omitempty"`
	Amount        int                    `json:"amount"`
	TaxAmount     int                    `json:"tax_amount"`
	TotalAmount   int                    `json:"total_amount"`
	Status        entities.InvoiceStatus `json:"status"`
	PDFURL        string                 `json:"pdf_url,omitempty"`
	DueDate       *time.Time             `json:"due_date,omitempty"`
	PaidAt        *time.Time             `json:"paid_at,omitempty"`
	Notes         string                 `json:"notes,omitempty"`
}

// InvoiceListResponse represents a paginated list of invoices
type InvoiceListResponse struct {
	Invoices   []InvoiceResponse `json:"invoices"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
}

// ToInvoiceResponse converts an Invoice entity to a response DTO
func ToInvoiceResponse(invoice *entities.Invoice) *InvoiceResponse {
	return &InvoiceResponse{
		ID:            invoice.ID,
		CreatedAt:     invoice.CreatedAt,
		UpdatedAt:     invoice.UpdatedAt,
		WaitingListID: invoice.WaitingListID,
		CustomerID:    invoice.CustomerID,
		Amount:        invoice.Amount,
		TaxAmount:     invoice.TaxAmount,
		TotalAmount:   invoice.TotalAmount,
		Status:        invoice.Status,
		PDFURL:        invoice.PDFURL,
		DueDate:       invoice.DueDate,
		PaidAt:        invoice.PaidAt,
		Notes:         invoice.Notes,
	}
}
