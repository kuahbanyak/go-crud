package repositories

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type MaintenanceItemRepository interface {
	Create(ctx context.Context, item *entities.MaintenanceItem) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.MaintenanceItem, error)
	Update(ctx context.Context, item *entities.MaintenanceItem) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
	GetByWaitingListID(ctx context.Context, waitingListID types.MSSQLUUID) ([]*entities.MaintenanceItem, error)
	GetByStatus(ctx context.Context, waitingListID types.MSSQLUUID, status entities.MaintenanceItemStatus) ([]*entities.MaintenanceItem, error)
	GetByType(ctx context.Context, waitingListID types.MSSQLUUID, itemType entities.MaintenanceItemType) ([]*entities.MaintenanceItem, error)
	GetPendingApproval(ctx context.Context, waitingListID types.MSSQLUUID) ([]*entities.MaintenanceItem, error)
	GetInitialItems(ctx context.Context, waitingListID types.MSSQLUUID) ([]*entities.MaintenanceItem, error)
	GetDiscoveredItems(ctx context.Context, waitingListID types.MSSQLUUID) ([]*entities.MaintenanceItem, error)
	CreateMany(ctx context.Context, items []*entities.MaintenanceItem) error
	UpdateStatus(ctx context.Context, id types.MSSQLUUID, status entities.MaintenanceItemStatus) error
	ApproveItems(ctx context.Context, ids []types.MSSQLUUID) error
	RejectItems(ctx context.Context, ids []types.MSSQLUUID) error
	GetTotalCost(ctx context.Context, waitingListID types.MSSQLUUID) (estimated float64, actual float64, err error)
	CountByStatus(ctx context.Context, waitingListID types.MSSQLUUID) (map[string]int, error)
}
