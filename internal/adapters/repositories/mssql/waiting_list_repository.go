package mssql
import (
	"context"
	"time"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)
type waitingListRepository struct {
	db *gorm.DB
}
func NewWaitingListRepository(db *gorm.DB) repositories.WaitingListRepository {
	return &waitingListRepository{db: db}
}
func (r *waitingListRepository) Create(ctx context.Context, waitingList *entities.WaitingList) error {
	return r.db.WithContext(ctx).Create(waitingList).Error
}
func (r *waitingListRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.WaitingList, error) {
	var waitingList entities.WaitingList
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Where("id = ?", id).First(&waitingList).Error
	if err != nil {
		return nil, err
	}
	return &waitingList, nil
}
func (r *waitingListRepository) GetByQueueNumber(ctx context.Context, queueNumber int, serviceDate time.Time) (*entities.WaitingList, error) {
	var waitingList entities.WaitingList
	startOfDay := time.Date(serviceDate.Year(), serviceDate.Month(), serviceDate.Day(), 0, 0, 0, 0, serviceDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Where("queue_number = ? AND service_date >= ? AND service_date < ?", queueNumber, startOfDay, endOfDay).
		First(&waitingList).Error
	if err != nil {
		return nil, err
	}
	return &waitingList, nil
}
func (r *waitingListRepository) GetByCustomerID(ctx context.Context, customerID types.MSSQLUUID) ([]*entities.WaitingList, error) {
	var waitingLists []*entities.WaitingList
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Where("customer_id = ?", customerID.String()).
		Order("service_date DESC, queue_number ASC").
		Find(&waitingLists).Error
	return waitingLists, err
}
func (r *waitingListRepository) GetByServiceDate(ctx context.Context, serviceDate time.Time) ([]*entities.WaitingList, error) {
	var waitingLists []*entities.WaitingList
	startOfDay := time.Date(serviceDate.Year(), serviceDate.Month(), serviceDate.Day(), 0, 0, 0, 0, serviceDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Where("service_date >= ? AND service_date < ?", startOfDay, endOfDay).
		Order("queue_number ASC").
		Find(&waitingLists).Error
	return waitingLists, err
}
func (r *waitingListRepository) GetByStatus(ctx context.Context, status entities.WaitingListStatus, serviceDate time.Time) ([]*entities.WaitingList, error) {
	var waitingLists []*entities.WaitingList
	startOfDay := time.Date(serviceDate.Year(), serviceDate.Month(), serviceDate.Day(), 0, 0, 0, 0, serviceDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Where("status = ? AND service_date >= ? AND service_date < ?", status, startOfDay, endOfDay).
		Order("queue_number ASC").
		Find(&waitingLists).Error
	return waitingLists, err
}
func (r *waitingListRepository) GetNextQueueNumber(ctx context.Context, serviceDate time.Time) (int, error) {
	var maxQueue int
	startOfDay := time.Date(serviceDate.Year(), serviceDate.Month(), serviceDate.Day(), 0, 0, 0, 0, serviceDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Model(&entities.WaitingList{}).
		Where("service_date >= ? AND service_date < ?", startOfDay, endOfDay).
		Select("COALESCE(MAX(queue_number), 0)").
		Scan(&maxQueue).Error
	return maxQueue + 1, err
}
func (r *waitingListRepository) Update(ctx context.Context, waitingList *entities.WaitingList) error {
	return r.db.WithContext(ctx).Save(waitingList).Error
}
func (r *waitingListRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.WaitingList{}).Error
}
func (r *waitingListRepository) List(ctx context.Context, limit, offset int) ([]*entities.WaitingList, error) {
	var waitingLists []*entities.WaitingList
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Limit(limit).Offset(offset).
		Order("service_date DESC, queue_number ASC").
		Find(&waitingLists).Error
	return waitingLists, err
}

