package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

type WaitingListRepository interface {
	Create(ctx context.Context, waitingList *entities.WaitingList) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.WaitingList, error)
	GetByQueueNumber(ctx context.Context, queueNumber int, serviceDate time.Time) (*entities.WaitingList, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.WaitingList, error)
	GetByServiceDate(ctx context.Context, serviceDate time.Time) ([]*entities.WaitingList, error)
	GetByWeekRange(ctx context.Context, weekStart, weekEnd time.Time) ([]*entities.WaitingList, error)
	GetByStatus(ctx context.Context, status entities.WaitingListStatus, serviceDate time.Time) ([]*entities.WaitingList, error)
	GetNextQueueNumber(ctx context.Context, serviceDate time.Time) (int, error)
	Update(ctx context.Context, waitingList *entities.WaitingList) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.WaitingList, error)
	CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error)
	CountActiveByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error)
	CountAll(ctx context.Context) (int64, error)
	CountActiveAll(ctx context.Context) (int64, error)
	CountByCustomerIDIncludingDeleted(ctx context.Context, customerID uuid.UUID) (int64, error)
	CountAllIncludingDeleted(ctx context.Context) (int64, error)
	CountByStatusAll(ctx context.Context, status entities.WaitingListStatus) (int64, error)
	CountByStatusCustomer(ctx context.Context, customerID uuid.UUID, status entities.WaitingListStatus) (int64, error)
}
