package repositories

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type JobRepository interface {
	Create(ctx context.Context, job *entities.Job) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Job, error)
	GetByTask(ctx context.Context, task string) (*entities.Job, error)
	GetAll(ctx context.Context) ([]*entities.Job, error)
	GetActiveJobs(ctx context.Context) ([]*entities.Job, error)
	Update(ctx context.Context, job *entities.Job) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
	SeedDefaults(ctx context.Context) error
	HasAnyActiveJob(ctx context.Context) (bool, error)
}
