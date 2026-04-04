package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/utils"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(ctx context.Context, invoice *entities.Invoice) error {
	query := `
		INSERT INTO invoices (id, waiting_list_id, customer_id, amount, tax_amount, total_amount, status, pdf_url, due_date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	now := utils.NowWIB()
	invoice.CreatedAt = now
	invoice.UpdatedAt = now

	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query,
		invoice.ID,
		invoice.WaitingListID,
		invoice.CustomerID,
		invoice.Amount,
		invoice.TaxAmount,
		invoice.TotalAmount,
		invoice.Status,
		invoice.PDFURL,
		invoice.DueDate,
		invoice.Notes,
		invoice.CreatedAt,
		invoice.UpdatedAt,
	)

	return err
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Invoice, error) {
	query := `
		SELECT id, waiting_list_id, customer_id, amount, tax_amount, total_amount, status, pdf_url, due_date, paid_at, notes, created_at, updated_at
		FROM invoices
		WHERE id = $1 AND deleted_at IS NULL
	`

	invoice := &entities.Invoice{}
	var waitingListID, pdfURL, notes sql.NullString
	var dueDate, paidAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invoice.ID,
		&waitingListID,
		&invoice.CustomerID,
		&invoice.Amount,
		&invoice.TaxAmount,
		&invoice.TotalAmount,
		&invoice.Status,
		&pdfURL,
		&dueDate,
		&paidAt,
		&notes,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	if waitingListID.Valid {
		wlID, _ := uuid.Parse(waitingListID.String)
		invoice.WaitingListID = &wlID
	}
	if pdfURL.Valid {
		invoice.PDFURL = pdfURL.String
	}
	if dueDate.Valid {
		invoice.DueDate = &dueDate.Time
	}
	if paidAt.Valid {
		invoice.PaidAt = &paidAt.Time
	}
	if notes.Valid {
		invoice.Notes = notes.String
	}

	return invoice, nil
}

func (r *InvoiceRepository) GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*entities.Invoice, error) {
	query := `
		SELECT id, waiting_list_id, customer_id, amount, tax_amount, total_amount, status, pdf_url, due_date, paid_at, notes, created_at, updated_at
		FROM invoices
		WHERE waiting_list_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*entities.Invoice
	for rows.Next() {
		invoice := &entities.Invoice{}
		var waitingListID, pdfURL, notes sql.NullString
		var dueDate, paidAt sql.NullTime

		err := rows.Scan(
			&invoice.ID,
			&waitingListID,
			&invoice.CustomerID,
			&invoice.Amount,
			&invoice.TaxAmount,
			&invoice.TotalAmount,
			&invoice.Status,
			&pdfURL,
			&dueDate,
			&paidAt,
			&notes,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if waitingListID.Valid {
			wlID, _ := uuid.Parse(waitingListID.String)
			invoice.WaitingListID = &wlID
		}
		if pdfURL.Valid {
			invoice.PDFURL = pdfURL.String
		}
		if dueDate.Valid {
			invoice.DueDate = &dueDate.Time
		}
		if paidAt.Valid {
			invoice.PaidAt = &paidAt.Time
		}
		if notes.Valid {
			invoice.Notes = notes.String
		}

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (r *InvoiceRepository) GetByStatus(ctx context.Context, status entities.InvoiceStatus) ([]*entities.Invoice, error) {
	query := `
		SELECT id, waiting_list_id, customer_id, amount, tax_amount, total_amount, status, pdf_url, due_date, paid_at, notes, created_at, updated_at
		FROM invoices
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*entities.Invoice
	for rows.Next() {
		invoice := &entities.Invoice{}
		var waitingListID, pdfURL, notes sql.NullString
		var dueDate, paidAt sql.NullTime

		err := rows.Scan(
			&invoice.ID,
			&waitingListID,
			&invoice.CustomerID,
			&invoice.Amount,
			&invoice.TaxAmount,
			&invoice.TotalAmount,
			&invoice.Status,
			&pdfURL,
			&dueDate,
			&paidAt,
			&notes,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if waitingListID.Valid {
			wlID, _ := uuid.Parse(waitingListID.String)
			invoice.WaitingListID = &wlID
		}
		if pdfURL.Valid {
			invoice.PDFURL = pdfURL.String
		}
		if dueDate.Valid {
			invoice.DueDate = &dueDate.Time
		}
		if paidAt.Valid {
			invoice.PaidAt = &paidAt.Time
		}
		if notes.Valid {
			invoice.Notes = notes.String
		}

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (r *InvoiceRepository) Update(ctx context.Context, invoice *entities.Invoice) error {
	query := `
		UPDATE invoices
		SET waiting_list_id = $1, customer_id = $2, amount = $3, tax_amount = $4, 
		    total_amount = $5, status = $6, pdf_url = $7, due_date = $8, 
		    paid_at = $9, notes = $10, updated_at = $11
		WHERE id = $12 AND deleted_at IS NULL
	`

	invoice.UpdatedAt = utils.NowWIB()

	result, err := r.db.ExecContext(ctx, query,
		invoice.WaitingListID,
		invoice.CustomerID,
		invoice.Amount,
		invoice.TaxAmount,
		invoice.TotalAmount,
		invoice.Status,
		invoice.PDFURL,
		invoice.DueDate,
		invoice.PaidAt,
		invoice.Notes,
		invoice.UpdatedAt,
		invoice.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("invoice not found")
	}

	return nil
}

func (r *InvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE invoices
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		utils.NowWIB(),
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("invoice not found")
	}

	return nil
}

func (r *InvoiceRepository) List(ctx context.Context, limit, offset int) ([]*entities.Invoice, error) {
	query := `
		SELECT id, waiting_list_id, customer_id, amount, tax_amount, total_amount, status, pdf_url, due_date, paid_at, notes, created_at, updated_at
		FROM invoices
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*entities.Invoice
	for rows.Next() {
		invoice := &entities.Invoice{}
		var waitingListID, pdfURL, notes sql.NullString
		var dueDate, paidAt sql.NullTime

		err := rows.Scan(
			&invoice.ID,
			&waitingListID,
			&invoice.CustomerID,
			&invoice.Amount,
			&invoice.TaxAmount,
			&invoice.TotalAmount,
			&invoice.Status,
			&pdfURL,
			&dueDate,
			&paidAt,
			&notes,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if waitingListID.Valid {
			wlID, _ := uuid.Parse(waitingListID.String)
			invoice.WaitingListID = &wlID
		}
		if pdfURL.Valid {
			invoice.PDFURL = pdfURL.String
		}
		if dueDate.Valid {
			invoice.DueDate = &dueDate.Time
		}
		if paidAt.Valid {
			invoice.PaidAt = &paidAt.Time
		}
		if notes.Valid {
			invoice.Notes = notes.String
		}

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (r *InvoiceRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM invoices WHERE deleted_at IS NULL`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
