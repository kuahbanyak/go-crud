package mssql
import (
	"context"
	"time"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)
type MaintenanceItemRepositoryImpl struct {
	db *gorm.DB
}
func NewMaintenanceItemRepository(db *gorm.DB) *MaintenanceItemRepositoryImpl {
	return &MaintenanceItemRepositoryImpl{db: db}
}
func (r *MaintenanceItemRepositoryImpl) Create(ctx context.Context, item *entities.MaintenanceItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}
func (r *MaintenanceItemRepositoryImpl) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.MaintenanceItem, error) {
	var item entities.MaintenanceItem
	err := r.db.WithContext(ctx).
		Preload("Mechanic").
		Preload("WaitingList").
		First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
func (r *MaintenanceItemRepositoryImpl) Update(ctx context.Context, item *entities.MaintenanceItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}
func (r *MaintenanceItemRepositoryImpl) Delete(ctx context.Context, id types.MSSQLUUID) error {
	return r.db.WithContext(ctx).Delete(&entities.MaintenanceItem{}, "id = ?", id).Error
}
func (r *MaintenanceItemRepositoryImpl) GetByWaitingListID(ctx context.Context, waitingListID types.MSSQLUUID) ([]*entities.MaintenanceItem, error) {
	var items []*entities.MaintenanceItem
	err := r.db.WithContext(ctx).
		Preload("Mechanic").
		Where("waiting_list_id = ?", waitingListID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}
func (r *MaintenanceItemRepositoryImpl) GetByStatus(ctx context.Context, waitingListID types.MSSQLUUID, status entities.MaintenanceItemStatus) ([]*entities.MaintenanceItem, error) {
	var items []*entities.MaintenanceItem
	err := r.db.WithContext(ctx).
		Preload("Mechanic").
		Where("waiting_list_id = ? AND status = ?", waitingListID, status).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}
func (r *MaintenanceItemRepositoryImpl) GetByType(ctx context.Context, waitingListID types.MSSQLUUID, itemType entities.MaintenanceItemType) ([]*entities.MaintenanceItem, error) {
	var items []*entities.MaintenanceItem
	err := r.db.WithContext(ctx).
		Preload("Mechanic").
		Where("waiting_list_id = ? AND item_type = ?", waitingListID, itemType).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}
func (r *MaintenanceItemRepositoryImpl) GetPendingApproval(ctx context.Context, waitingListID types.MSSQLUUID) ([]*entities.MaintenanceItem, error) {
	var items []*entities.MaintenanceItem
	err := r.db.WithContext(ctx).
		Preload("Mechanic").
		Where("waiting_list_id = ? AND requires_approval = ? AND status = ?",
			waitingListID, true, entities.MaintenanceItemStatusInspected).
		Order("priority DESC, created_at ASC").
		Find(&items).Error
	return items, err
}
func (r *MaintenanceItemRepositoryImpl) GetInitialItems(ctx context.Context, waitingListID types.MSSQLUUID) ([]*entities.MaintenanceItem, error) {
	return r.GetByType(ctx, waitingListID, entities.MaintenanceItemTypeInitial)
}
func (r *MaintenanceItemRepositoryImpl) GetDiscoveredItems(ctx context.Context, waitingListID types.MSSQLUUID) ([]*entities.MaintenanceItem, error) {
	return r.GetByType(ctx, waitingListID, entities.MaintenanceItemTypeDiscovered)
}
func (r *MaintenanceItemRepositoryImpl) CreateMany(ctx context.Context, items []*entities.MaintenanceItem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}
func (r *MaintenanceItemRepositoryImpl) UpdateStatus(ctx context.Context, id types.MSSQLUUID, status entities.MaintenanceItemStatus) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}
	switch status {
	case entities.MaintenanceItemStatusInspected:
		updates["inspected_at"] = now
	case entities.MaintenanceItemStatusApproved:
		updates["approved_at"] = now
	case entities.MaintenanceItemStatusCompleted:
		updates["completed_at"] = now
	}
	return r.db.WithContext(ctx).
		Model(&entities.MaintenanceItem{}).
		Where("id = ?", id).
		Updates(updates).Error
}
func (r *MaintenanceItemRepositoryImpl) ApproveItems(ctx context.Context, ids []types.MSSQLUUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entities.MaintenanceItem{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":      entities.MaintenanceItemStatusApproved,
			"approved_at": now,
			"updated_at":  now,
		}).Error
}
func (r *MaintenanceItemRepositoryImpl) RejectItems(ctx context.Context, ids []types.MSSQLUUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entities.MaintenanceItem{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":     entities.MaintenanceItemStatusRejected,
			"updated_at": now,
		}).Error
}
func (r *MaintenanceItemRepositoryImpl) GetTotalCost(ctx context.Context, waitingListID types.MSSQLUUID) (estimated float64, actual float64, err error) {
	type Result struct {
		TotalEstimated float64
		TotalActual    float64
	}
	var result Result
	err = r.db.WithContext(ctx).
		Model(&entities.MaintenanceItem{}).
		Select("SUM(estimated_cost) as total_estimated, SUM(actual_cost) as total_actual").
		Where("waiting_list_id = ? AND status NOT IN ?",
			waitingListID,
			[]entities.MaintenanceItemStatus{entities.MaintenanceItemStatusRejected, entities.MaintenanceItemStatusSkipped}).
		Scan(&result).Error
	return result.TotalEstimated, result.TotalActual, err
}
func (r *MaintenanceItemRepositoryImpl) CountByStatus(ctx context.Context, waitingListID types.MSSQLUUID) (map[string]int, error) {
	type StatusCount struct {
		Status string
		Count  int
	}
	var results []StatusCount
	err := r.db.WithContext(ctx).
		Model(&entities.MaintenanceItem{}).
		Select("status, COUNT(*) as count").
		Where("waiting_list_id = ?", waitingListID).
		Group("status").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

