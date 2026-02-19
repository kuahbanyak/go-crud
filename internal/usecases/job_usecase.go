package usecases

import (
	"context"
	"errors"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type JobUsecase struct {
	jobRepo repositories.JobRepository
}

func NewJobUsecase(jobRepo repositories.JobRepository) *JobUsecase {
	return &JobUsecase{
		jobRepo: jobRepo,
	}
}

func (u *JobUsecase) GetJob(ctx context.Context, id types.MSSQLUUID) (*entities.Job, error) {
	return u.jobRepo.GetByID(ctx, id)
}

func (u *JobUsecase) GetJobByTask(ctx context.Context, task string) (*entities.Job, error) {
	return u.jobRepo.GetByTask(ctx, task)
}

func (u *JobUsecase) GetAllJobs(ctx context.Context) ([]*entities.Job, error) {
	return u.jobRepo.GetAll(ctx)
}

func (u *JobUsecase) GetActiveJobs(ctx context.Context) ([]*entities.Job, error) {
	return u.jobRepo.GetActiveJobs(ctx)
}

func (u *JobUsecase) CreateJob(ctx context.Context, job *entities.Job) error {
	existing, err := u.jobRepo.GetByTask(ctx, job.Task)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("job with this task name already exists")
	}
	return u.jobRepo.Create(ctx, job)
}

func (u *JobUsecase) UpdateJob(ctx context.Context, id types.MSSQLUUID, cronJob, task string, status bool) error {
	job, err := u.jobRepo.GetByID(ctx, id)
	if err != nil || job == nil {
		return errors.New("job not found")
	}

	// Check if task name is being changed and if it conflicts
	if task != job.Task {
		existing, err := u.jobRepo.GetByTask(ctx, task)
		if err != nil {
			return err
		}
		if existing != nil && existing.ID != id {
			return errors.New("another job with this task name already exists")
		}
	}

	job.CronJob = cronJob
	job.Task = task
	job.Status = status

	return u.jobRepo.Update(ctx, job)
}

func (u *JobUsecase) UpdateJobStatus(ctx context.Context, id types.MSSQLUUID, status bool) error {
	job, err := u.jobRepo.GetByID(ctx, id)
	if err != nil || job == nil {
		return errors.New("job not found")
	}

	job.Status = status
	return u.jobRepo.Update(ctx, job)
}

func (u *JobUsecase) DeleteJob(ctx context.Context, id types.MSSQLUUID) error {
	job, err := u.jobRepo.GetByID(ctx, id)
	if err != nil || job == nil {
		return errors.New("job not found")
	}
	return u.jobRepo.Delete(ctx, id)
}

func (u *JobUsecase) HasAnyActiveJob(ctx context.Context) (bool, error) {
	return u.jobRepo.HasAnyActiveJob(ctx)
}

func (u *JobUsecase) IsJobEnabled(ctx context.Context, taskName string) bool {
	job, err := u.jobRepo.GetByTask(ctx, taskName)
	if err != nil || job == nil {
		return false
	}
	return job.Status
}

func (u *JobUsecase) GetJobSchedule(ctx context.Context, taskName string) string {
	job, err := u.jobRepo.GetByTask(ctx, taskName)
	if err != nil || job == nil {
		return ""
	}
	return job.CronJob
}
