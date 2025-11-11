package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *entities.Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Invoice, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*entities.Invoice, error)
	GetByStatus(ctx context.Context, status entities.InvoiceStatus) ([]*entities.Invoice, error)
	Update(ctx context.Context, invoice *entities.Invoice) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Invoice, error)
}
type PartRepository interface {
	Create(ctx context.Context, part *entities.Part) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Part, error)
	GetBySKU(ctx context.Context, sku string) (*entities.Part, error)
	Update(ctx context.Context, part *entities.Part) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Part, error)
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}
