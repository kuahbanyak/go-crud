package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/scheduler"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type SchedulerHandler struct {
	scheduler  *scheduler.Scheduler
	jobUsecase *usecases.JobUsecase
}

func NewSchedulerHandler(scheduler *scheduler.Scheduler, jobUsecase *usecases.JobUsecase) *SchedulerHandler {
	return &SchedulerHandler{
		scheduler:  scheduler,
		jobUsecase: jobUsecase,
	}
}

// GetJobsStatus returns the status of all registered jobs
func (h *SchedulerHandler) GetJobsStatus(w http.ResponseWriter, r *http.Request) {
	statuses := h.scheduler.GetJobsStatus()

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Jobs status retrieved successfully", map[string]interface{}{
		"scheduler_running": h.scheduler.IsRunning(),
		"total_jobs":        len(statuses),
		"jobs":              statuses,
	})
}

// GetJobStatus returns the status of a specific job
func (h *SchedulerHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["name"]

	status, exists := h.scheduler.GetJobStatus(jobName)
	if !exists {
		response.ErrorWithContext(r.Context(), w, http.StatusNotFound, "Job not found", nil)
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Job status retrieved successfully", status)
}

// RunJobManually triggers a job to run immediately
func (h *SchedulerHandler) RunJobManually(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["name"]

	if err := h.scheduler.RunJobNow(jobName); err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Job triggered successfully", map[string]string{
		"job_name": jobName,
		"message":  "Job is running in the background",
	})
}

// ListJobs returns a simple list of all registered jobs
func (h *SchedulerHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs := h.scheduler.ListJobs()

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Jobs list retrieved successfully", map[string]interface{}{
		"scheduler_running": h.scheduler.IsRunning(),
		"auto_run_enabled":  h.scheduler.IsAutoRunEnabled(),
		"total_jobs":        len(jobs),
		"jobs":              jobs,
	})
}

// StartAutoRun enables automatic job execution
func (h *SchedulerHandler) StartAutoRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Enable in scheduler
	h.scheduler.EnableAutoRun()

	// Enable all jobs in database
	if h.jobUsecase != nil {
		jobs, err := h.jobUsecase.GetAllJobs(ctx)
		if err == nil {
			for _, job := range jobs {
				// Enable each job
				_ = h.jobUsecase.UpdateJobStatus(ctx, job.ID, true)
			}
		}
	}

	response.SuccessWithContext(ctx, w, http.StatusOK, "Automatic job execution enabled", map[string]interface{}{
		"auto_run_enabled": true,
		"message":          "Jobs will now run automatically according to their schedule",
	})
}

// StopAutoRun disables automatic job execution
func (h *SchedulerHandler) StopAutoRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Disable in scheduler
	h.scheduler.DisableAutoRun()

	// Disable all jobs in database
	if h.jobUsecase != nil {
		jobs, err := h.jobUsecase.GetAllJobs(ctx)
		if err == nil {
			for _, job := range jobs {
				// Disable each job
				_ = h.jobUsecase.UpdateJobStatus(ctx, job.ID, false)
			}
		}
	}

	response.SuccessWithContext(ctx, w, http.StatusOK, "Automatic job execution disabled", map[string]interface{}{
		"auto_run_enabled": false,
		"message":          "Jobs will only run when manually triggered by admin",
	})
}

// GetAutoRunStatus returns the current auto-run status
func (h *SchedulerHandler) GetAutoRunStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	autoRunEnabled := h.scheduler.IsAutoRunEnabled()

	response.SuccessWithContext(ctx, w, http.StatusOK, "Auto-run status retrieved", map[string]interface{}{
		"auto_run_enabled": autoRunEnabled,
		"message":          getAutoRunMessage(autoRunEnabled),
	})
}

func getAutoRunMessage(enabled bool) string {
	if enabled {
		return "Jobs are running automatically according to schedule"
	}
	return "Jobs will only run when manually triggered by admin"
}
