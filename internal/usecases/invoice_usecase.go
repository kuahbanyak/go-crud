package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type InvoiceUsecase struct {
	invoiceRepo     repositories.InvoiceRepository
	waitingListRepo repositories.WaitingListRepository
	userRepo        repositories.UserRepository
}

func NewInvoiceUsecase(
	invoiceRepo repositories.InvoiceRepository,
	waitingListRepo repositories.WaitingListRepository,
	userRepo repositories.UserRepository,
) *InvoiceUsecase {
	return &InvoiceUsecase{
		invoiceRepo:     invoiceRepo,
		waitingListRepo: waitingListRepo,
		userRepo:        userRepo,
	}
}

func (u *InvoiceUsecase) CreateInvoice(ctx context.Context, req *dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error) {
	// Verify customer exists
	_, err := u.userRepo.GetByID(ctx, types.FromUUID(req.CustomerID))
	if err != nil {
		return nil, errors.New("customer not found")
	}

	// If waiting list ID provided, verify it exists
	if req.WaitingListID != nil {
		_, err := u.waitingListRepo.GetByID(ctx, types.FromUUID(*req.WaitingListID))
		if err != nil {
			return nil, errors.New("waiting list not found")
		}
	}

	// Calculate total amount
	totalAmount := req.Amount + req.TaxAmount

	// Calculate due date
	var dueDate *time.Time
	if req.DueDays > 0 {
		due := time.Now().AddDate(0, 0, req.DueDays)
		dueDate = &due
	}

	invoice := &entities.Invoice{
		WaitingListID: req.WaitingListID,
		CustomerID:    req.CustomerID,
		Amount:        req.Amount,
		TaxAmount:     req.TaxAmount,
		TotalAmount:   totalAmount,
		Status:        entities.InvoiceStatusPending,
		Notes:         req.Notes,
		DueDate:       dueDate,
	}

	if err := u.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	return dto.ToInvoiceResponse(invoice), nil
}

func (u *InvoiceUsecase) GetInvoice(ctx context.Context, id uuid.UUID) (*dto.InvoiceResponse, error) {
	invoice, err := u.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := dto.ToInvoiceResponse(invoice)

	// Get customer name
	customer, err := u.userRepo.GetByID(ctx, types.FromUUID(invoice.CustomerID))
	if err == nil {
		response.CustomerName = customer.Name
	}

	return response, nil
}

func (u *InvoiceUsecase) GetInvoicesByWaitingList(ctx context.Context, waitingListID uuid.UUID) ([]dto.InvoiceResponse, error) {
	invoices, err := u.invoiceRepo.GetByBookingID(ctx, waitingListID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.InvoiceResponse, len(invoices))
	for i, invoice := range invoices {
		responses[i] = *dto.ToInvoiceResponse(invoice)
	}

	return responses, nil
}

func (u *InvoiceUsecase) ListInvoices(ctx context.Context, page, pageSize int, status string) (*dto.InvoiceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var invoices []*entities.Invoice
	var err error

	if status != "" {
		invoices, err = u.invoiceRepo.GetByStatus(ctx, entities.InvoiceStatus(status))
	} else {
		invoices, err = u.invoiceRepo.List(ctx, pageSize, offset)
	}

	if err != nil {
		return nil, err
	}

	responses := make([]dto.InvoiceResponse, len(invoices))
	for i, invoice := range invoices {
		resp := dto.ToInvoiceResponse(invoice)

		// Get customer name
		customer, err := u.userRepo.GetByID(ctx, types.FromUUID(invoice.CustomerID))
		if err == nil {
			resp.CustomerName = customer.Name
		}

		responses[i] = *resp
	}

	totalCount, _ := u.invoiceRepo.Count(ctx)

	return &dto.InvoiceListResponse{
		Invoices:   responses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (u *InvoiceUsecase) UpdateInvoice(ctx context.Context, id uuid.UUID, req *dto.UpdateInvoiceRequest) (*dto.InvoiceResponse, error) {
	invoice, err := u.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Amount != nil {
		invoice.Amount = *req.Amount
	}
	if req.TaxAmount != nil {
		invoice.TaxAmount = *req.TaxAmount
	}
	if req.Amount != nil || req.TaxAmount != nil {
		invoice.TotalAmount = invoice.Amount + invoice.TaxAmount
	}
	if req.Status != nil {
		invoice.Status = *req.Status
	}
	if req.Notes != nil {
		invoice.Notes = *req.Notes
	}
	if req.DueDate != nil {
		invoice.DueDate = req.DueDate
	}

	if err := u.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return dto.ToInvoiceResponse(invoice), nil
}

func (u *InvoiceUsecase) PayInvoice(ctx context.Context, id uuid.UUID, req *dto.PayInvoiceRequest) (*dto.InvoiceResponse, error) {
	invoice, err := u.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if invoice.Status == entities.InvoiceStatusPaid {
		return nil, errors.New("invoice already paid")
	}

	if invoice.Status == entities.InvoiceStatusCancelled {
		return nil, errors.New("cannot pay cancelled invoice")
	}

	// Update invoice status
	invoice.Status = entities.InvoiceStatusPaid
	now := time.Now()
	invoice.PaidAt = &now

	if err := u.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return dto.ToInvoiceResponse(invoice), nil
}

func (u *InvoiceUsecase) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	invoice, err := u.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if invoice.Status == entities.InvoiceStatusPaid {
		return errors.New("cannot delete paid invoice")
	}

	return u.invoiceRepo.Delete(ctx, id)
}
