package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/logger"
)

type JobStatus struct {
	Name          string    `json:"name"`
	Schedule      string    `json:"schedule"`
	IsRunning     bool      `json:"is_running"`
	LastRunTime   time.Time `json:"last_run_time,omitempty"`
	LastRunStatus string    `json:"last_run_status,omitempty"`
	NextRunTime   time.Time `json:"next_run_time,omitempty"`
	RunCount      int       `json:"run_count"`
	ErrorCount    int       `json:"error_count"`
}

type Scheduler struct {
	scheduler      gocron.Scheduler
	jobs           []Job
	jobStatus      map[string]*JobStatus
	statusLock     sync.RWMutex
	autoRunEnabled bool
	autoRunLock    sync.RWMutex
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
		scheduler:      s,
		jobs:           make([]Job, 0),
		jobStatus:      make(map[string]*JobStatus),
		statusLock:     sync.RWMutex{},
		autoRunEnabled: false, // Default to disabled
		autoRunLock:    sync.RWMutex{},
	}, nil
}

func (s *Scheduler) RegisterJob(job Job) error {
	logger.Info(fmt.Sprintf("Registering job: %s with schedule: %s", job.Name(), job.Schedule()))

	// Initialize job status
	s.statusLock.Lock()
	s.jobStatus[job.Name()] = &JobStatus{
		Name:          job.Name(),
		Schedule:      job.Schedule(),
		IsRunning:     false,
		LastRunStatus: "Not started yet",
		RunCount:      0,
		ErrorCount:    0,
	}
	s.statusLock.Unlock()

	_, err := s.scheduler.NewJob(
		gocron.CronJob(job.Schedule(), false),
		gocron.NewTask(func() {
			// Check if auto-run is enabled before executing
			s.autoRunLock.RLock()
			enabled := s.autoRunEnabled
			s.autoRunLock.RUnlock()

			if !enabled {
				logger.Info(fmt.Sprintf("Job %s skipped - automatic execution is disabled", job.Name()))
				return
			}

			ctx := context.Background()

			// Mark job as running
			s.updateJobStatus(job.Name(), func(status *JobStatus) {
				status.IsRunning = true
				status.LastRunTime = time.Now()
			})

			logger.Info(fmt.Sprintf("Starting job: %s", job.Name()))

			err := job.Run(ctx)

			// Update job status after completion
			s.updateJobStatus(job.Name(), func(status *JobStatus) {
				status.IsRunning = false
				status.RunCount++
				if err != nil {
					status.ErrorCount++
					status.LastRunStatus = fmt.Sprintf("Failed: %v", err)
					logger.Error(fmt.Sprintf("Job %s failed: %v", job.Name(), err))
				} else {
					status.LastRunStatus = "Success"
					logger.Info(fmt.Sprintf("Job %s completed successfully", job.Name()))
				}
			})
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to register job %s: %w", job.Name(), err)
	}
	s.jobs = append(s.jobs, job)
	return nil
}

func (s *Scheduler) updateJobStatus(jobName string, updateFn func(*JobStatus)) {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()
	if status, exists := s.jobStatus[jobName]; exists {
		updateFn(status)
	}
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

			// Mark job as running
			s.updateJobStatus(job.Name(), func(status *JobStatus) {
				status.IsRunning = true
				status.LastRunTime = time.Now()
			})

			// Run the job
			err := job.Run(ctx)

			// Update job status after completion
			s.updateJobStatus(job.Name(), func(status *JobStatus) {
				status.IsRunning = false
				status.RunCount++
				if err != nil {
					status.ErrorCount++
					status.LastRunStatus = fmt.Sprintf("Failed: %v", err)
					logger.Error(fmt.Sprintf("Manually triggered job %s failed: %v", job.Name(), err))
				} else {
					status.LastRunStatus = "Success"
					logger.Info(fmt.Sprintf("Manually triggered job %s completed successfully", job.Name()))
				}
			})

			return err
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

func (s *Scheduler) GetJobsStatus() []*JobStatus {
	s.statusLock.RLock()
	defer s.statusLock.RUnlock()

	statuses := make([]*JobStatus, 0, len(s.jobStatus))
	for _, status := range s.jobStatus {
		// Create a copy to avoid race conditions
		statusCopy := *status
		statuses = append(statuses, &statusCopy)
	}
	return statuses
}

func (s *Scheduler) GetJobStatus(jobName string) (*JobStatus, bool) {
	s.statusLock.RLock()
	defer s.statusLock.RUnlock()

	if status, exists := s.jobStatus[jobName]; exists {
		// Create a copy to avoid race conditions
		statusCopy := *status
		return &statusCopy, true
	}
	return nil, false
}

func (s *Scheduler) IsRunning() bool {
	return len(s.jobs) > 0
}

// IsAutoRunEnabled returns whether automatic job execution is enabled
func (s *Scheduler) IsAutoRunEnabled() bool {
	s.autoRunLock.RLock()
	defer s.autoRunLock.RUnlock()
	return s.autoRunEnabled
}

// EnableAutoRun enables automatic job execution
func (s *Scheduler) EnableAutoRun() {
	s.autoRunLock.Lock()
	defer s.autoRunLock.Unlock()
	s.autoRunEnabled = true
	logger.Info("Automatic job execution enabled")
}

// DisableAutoRun disables automatic job execution
func (s *Scheduler) DisableAutoRun() {
	s.autoRunLock.Lock()
	defer s.autoRunLock.Unlock()
	s.autoRunEnabled = false
	logger.Info("Automatic job execution disabled")
}

// SetAutoRun sets the automatic job execution state
func (s *Scheduler) SetAutoRun(enabled bool) {
	if enabled {
		s.EnableAutoRun()
	} else {
		s.DisableAutoRun()
	}
}
