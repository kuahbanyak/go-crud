package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CronJob   string         `gorm:"type:varchar(100);not null" json:"cron_job"`
	Status    bool           `gorm:"default:false" json:"status"`
	Task      string         `gorm:"type:varchar(200);not null;unique" json:"task"`
}

func (j *Job) BeforeCreate(_ *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}

// Default jobs for seeding
var DefaultJobs = []Job{
	{
		Task:    "DailyWaitingListCleanup",
		CronJob: "0 0 * * *",
		Status:  false,
	},
}
