package mssql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
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
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12)
	`

	now := time.Now()
	invoice.CreatedAt = now
	invoice.UpdatedAt = now

	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", invoice.ID),
		sql.Named("p2", invoice.WaitingListID),
		sql.Named("p3", invoice.CustomerID),
		sql.Named("p4", invoice.Amount),
		sql.Named("p5", invoice.TaxAmount),
		sql.Named("p6", invoice.TotalAmount),
		sql.Named("p7", invoice.Status),
		sql.Named("p8", invoice.PDFURL),
		sql.Named("p9", invoice.DueDate),
		sql.Named("p10", invoice.Notes),
		sql.Named("p11", invoice.CreatedAt),
		sql.Named("p12", invoice.UpdatedAt),
	)

	return err
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Invoice, error) {
	query := `
		SELECT id, waiting_list_id, customer_id, amount, tax_amount, total_amount, status, pdf_url, due_date, paid_at, notes, created_at, updated_at
		FROM invoices
		WHERE id = @p1 AND deleted_at IS NULL
	`

	invoice := &entities.Invoice{}
	var waitingListID, pdfURL, notes sql.NullString
	var dueDate, paidAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, sql.Named("p1", id)).Scan(
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
		WHERE waiting_list_id = @p1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", bookingID))
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
		WHERE status = @p1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", status))
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
		SET waiting_list_id = @p1, customer_id = @p2, amount = @p3, tax_amount = @p4, 
		    total_amount = @p5, status = @p6, pdf_url = @p7, due_date = @p8, 
		    paid_at = @p9, notes = @p10, updated_at = @p11
		WHERE id = @p12 AND deleted_at IS NULL
	`

	invoice.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", invoice.WaitingListID),
		sql.Named("p2", invoice.CustomerID),
		sql.Named("p3", invoice.Amount),
		sql.Named("p4", invoice.TaxAmount),
		sql.Named("p5", invoice.TotalAmount),
		sql.Named("p6", invoice.Status),
		sql.Named("p7", invoice.PDFURL),
		sql.Named("p8", invoice.DueDate),
		sql.Named("p9", invoice.PaidAt),
		sql.Named("p10", invoice.Notes),
		sql.Named("p11", invoice.UpdatedAt),
		sql.Named("p12", invoice.ID),
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
		SET deleted_at = @p1
		WHERE id = @p2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", time.Now()),
		sql.Named("p2", id),
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
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", offset),
		sql.Named("p2", limit),
	)
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
