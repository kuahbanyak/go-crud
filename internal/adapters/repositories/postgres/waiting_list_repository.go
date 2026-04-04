package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
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
func (r *waitingListRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.WaitingList, error) {
	var waitingList entities.WaitingList
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Preload("Mechanic").
		Preload("ServiceItem").
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
		Preload("Mechanic").
		Preload("ServiceItem").
		Where("queue_number = ? AND service_date >= ? AND service_date < ?", queueNumber, startOfDay, endOfDay).
		First(&waitingList).Error
	if err != nil {
		return nil, err
	}
	return &waitingList, nil
}
func (r *waitingListRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.WaitingList, error) {
	var waitingLists []*entities.WaitingList
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Preload("Mechanic").
		Preload("ServiceItem").
		Where("customer_id = ?", customerID).
		Order("CASE WHEN status IN ('waiting', 'called', 'in_service') THEN 0 ELSE 1 END, service_date DESC, queue_number ASC").
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
		Preload("Mechanic").
		Preload("ServiceItem").
		Where("service_date >= ? AND service_date < ?", startOfDay, endOfDay).
		Order("CASE WHEN status IN ('waiting', 'called', 'in_service') THEN 0 ELSE 1 END, queue_number ASC").
		Find(&waitingLists).Error
	return waitingLists, err
}

func (r *waitingListRepository) GetByWeekRange(ctx context.Context, weekStart, weekEnd time.Time) ([]*entities.WaitingList, error) {
	var waitingLists []*entities.WaitingList
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Preload("Mechanic").
		Preload("ServiceItem").
		Where("service_date >= ? AND service_date <= ?", weekStart, weekEnd).
		Order("CASE WHEN status IN ('waiting', 'called', 'in_service') THEN 0 ELSE 1 END, service_date ASC, queue_number ASC").
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
		Preload("Mechanic").
		Preload("ServiceItem").
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
func (r *waitingListRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.WaitingList{}).Error
}
func (r *waitingListRepository) List(ctx context.Context, limit, offset int) ([]*entities.WaitingList, error) {
	var waitingLists []*entities.WaitingList
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Preload("ServiceItem").
		Limit(limit).Offset(offset).
		Order("CASE WHEN status IN ('waiting', 'called', 'in_service') THEN 0 ELSE 1 END, service_date DESC, queue_number ASC").
		Find(&waitingLists).Error
	return waitingLists, err
}

func (r *waitingListRepository) CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.WaitingList{}).
		Where("customer_id = ?", customerID).
		Count(&count).Error
	return count, err
}

func (r *waitingListRepository) CountActiveByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.WaitingList{}).
		Where("customer_id = ? AND status IN ('waiting', 'called', 'in_service')", customerID).
		Count(&count).Error
	return count, err
}

func (r *waitingListRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.WaitingList{}).
		Count(&count).Error
	return count, err
}

func (r *waitingListRepository) CountActiveAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.WaitingList{}).
		Where("status IN ('waiting', 'called', 'in_service')").
		Count(&count).Error
	return count, err
}

func (r *waitingListRepository) CountByCustomerIDIncludingDeleted(ctx context.Context, customerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Unscoped().
		Model(&entities.WaitingList{}).
		Where("customer_id = ?", customerID).
		Count(&count).Error
	return count, err
}

func (r *waitingListRepository) CountAllIncludingDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Unscoped().
		Model(&entities.WaitingList{}).
		Count(&count).Error
	return count, err
}

func (r *waitingListRepository) CountByStatusAll(ctx context.Context, status entities.WaitingListStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.WaitingList{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *waitingListRepository) CountByStatusCustomer(ctx context.Context, customerID uuid.UUID, status entities.WaitingListStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.WaitingList{}).
		Where("customer_id = ? AND status = ?", customerID, status).
		Count(&count).Error
	return count, err
}
