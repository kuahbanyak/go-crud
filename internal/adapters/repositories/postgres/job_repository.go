package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"gorm.io/gorm"
)

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) repositories.JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *entities.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Job, error) {
	var job entities.Job
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetByTask(ctx context.Context, task string) (*entities.Job, error) {
	var job entities.Job
	err := r.db.WithContext(ctx).Where("task = ?", task).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetAll(ctx context.Context) ([]*entities.Job, error) {
	var jobs []*entities.Job
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) GetActiveJobs(ctx context.Context) ([]*entities.Job, error) {
	var jobs []*entities.Job
	err := r.db.WithContext(ctx).
		Where("status = ?", true).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) Update(ctx context.Context, job *entities.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *jobRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Job{}).Error
}

func (r *jobRepository) SeedDefaults(ctx context.Context) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.Job{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	for _, job := range entities.DefaultJobs {
		if err := r.Create(ctx, &job); err != nil {
			return err
		}
	}
	return nil
}

func (r *jobRepository) HasAnyActiveJob(ctx context.Context) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Job{}).Where("status = ?", true).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
