package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

type JobRepository interface {
	Create(ctx context.Context, job *entities.Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Job, error)
	GetByTask(ctx context.Context, task string) (*entities.Job, error)
	GetAll(ctx context.Context) ([]*entities.Job, error)
	GetActiveJobs(ctx context.Context) ([]*entities.Job, error)
	Update(ctx context.Context, job *entities.Job) error
	Delete(ctx context.Context, id uuid.UUID) error
	SeedDefaults(ctx context.Context) error
	HasAnyActiveJob(ctx context.Context) (bool, error)
}
