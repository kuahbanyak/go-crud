package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type Job struct {
	ID        types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
	CronJob   string          `gorm:"type:varchar(100);not null" json:"cron_job"`
	Status    bool            `gorm:"default:false" json:"status"`
	Task      string          `gorm:"type:varchar(200);not null;unique" json:"task"`
}

func (j *Job) BeforeCreate(_ *gorm.DB) error {
	if j.ID.String() == "00000000-0000-0000-0000-000000000000" {
		j.ID = types.NewMSSQLUUID()
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
