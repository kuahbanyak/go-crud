package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateJobRequest struct {
	CronJob string `json:"cron_job" validate:"required"`
	Task    string `json:"task" validate:"required"`
	Status  bool   `json:"status"`
}

type UpdateJobRequest struct {
	CronJob string `json:"cron_job" validate:"required"`
	Task    string `json:"task" validate:"required"`
	Status  bool   `json:"status"`
}

type UpdateJobStatusRequest struct {
	Status bool `json:"status"`
}

type JobResponse struct {
	ID        uuid.UUID `json:"id"`
	CronJob   string          `json:"cron_job"`
	Task      string          `json:"task"`
	Status    bool            `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type JobsListResponse struct {
	Jobs  []JobResponse `json:"jobs"`
	Total int           `json:"total"`
}
