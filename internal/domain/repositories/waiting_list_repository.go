package repositories

import (
	"context"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type WaitingListRepository interface {
	Create(ctx context.Context, waitingList *entities.WaitingList) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.WaitingList, error)
	GetByQueueNumber(ctx context.Context, queueNumber int, serviceDate time.Time) (*entities.WaitingList, error)
	GetByCustomerID(ctx context.Context, customerID types.MSSQLUUID) ([]*entities.WaitingList, error)
	GetByServiceDate(ctx context.Context, serviceDate time.Time) ([]*entities.WaitingList, error)
	GetByStatus(ctx context.Context, status entities.WaitingListStatus, serviceDate time.Time) ([]*entities.WaitingList, error)
	GetNextQueueNumber(ctx context.Context, serviceDate time.Time) (int, error)
	Update(ctx context.Context, waitingList *entities.WaitingList) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.WaitingList, error)
}
