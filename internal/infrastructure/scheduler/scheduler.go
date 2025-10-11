package scheduler

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/logger"
)

type Scheduler struct {
	scheduler gocron.Scheduler
	jobs      []Job
}

type Job interface {
	Name() string
	Run(ctx context.Context) error
	Schedule() string // Cron expression
}

func NewScheduler() (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Scheduler{
		scheduler: s,
		jobs:      make([]Job, 0),
	}, nil
}

func (s *Scheduler) RegisterJob(job Job) error {
	logger.Info(fmt.Sprintf("Registering job: %s with schedule: %s", job.Name(), job.Schedule()))

	_, err := s.scheduler.NewJob(
		gocron.CronJob(job.Schedule(), false),
		gocron.NewTask(func() {
			ctx := context.Background()
			logger.Info(fmt.Sprintf("Starting job: %s", job.Name()))

			if err := job.Run(ctx); err != nil {
				logger.Error(fmt.Sprintf("Job %s failed: %v", job.Name(), err))
			} else {
				logger.Info(fmt.Sprintf("Job %s completed successfully", job.Name()))
			}
		}),
	)

	if err != nil {
		return fmt.Errorf("failed to register job %s: %w", job.Name(), err)
	}

	s.jobs = append(s.jobs, job)
	return nil
}

func (s *Scheduler) Start() {
	logger.Info(fmt.Sprintf("Starting scheduler with %d jobs", len(s.jobs)))
	s.scheduler.Start()
}

func (s *Scheduler) Stop() error {
	logger.Info("Stopping scheduler")
	return s.scheduler.Shutdown()
}

func (s *Scheduler) RunJobNow(jobName string) error {
	ctx := context.Background()
	for _, job := range s.jobs {
		if job.Name() == jobName {
			logger.Info(fmt.Sprintf("Manually running job: %s", jobName))
			return job.Run(ctx)
		}
	}
	return fmt.Errorf("job not found: %s", jobName)
}

func (s *Scheduler) ListJobs() []string {
	jobs := make([]string, len(s.jobs))
	for i, job := range s.jobs {
		jobs[i] = fmt.Sprintf("%s (Schedule: %s)", job.Name(), job.Schedule())
	}
	return jobs
}
