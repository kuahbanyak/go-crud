package entities
import (
	"time"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)
type SettingType string
const (
	SettingTypeInt    SettingType = "int"
	SettingTypeString SettingType = "string"
	SettingTypeBool   SettingType = "bool"
	SettingTypeFloat  SettingType = "float"
)
type Setting struct {
	ID        types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
	Key         string      `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value       string      `gorm:"type:varchar(500);not null" json:"value"`
	Type        SettingType `gorm:"type:varchar(20);not null" json:"type"`
	Description string      `gorm:"type:text" json:"description"`
	Category    string      `gorm:"type:varchar(50);index" json:"category"`
	IsEditable  bool        `gorm:"default:true" json:"is_editable"`
	IsPublic    bool        `gorm:"default:false" json:"is_public"`
}
func (s *Setting) BeforeCreate(_ *gorm.DB) error {
	if s.ID.String() == "00000000-0000-0000-0000-000000000000" {
		s.ID = types.NewMSSQLUUID()
	}
	return nil
}
var DefaultSettings = []Setting{
	{
		Key:         "waiting_list.max_tickets_per_day",
		Value:       "10",
		Type:        SettingTypeInt,
		Description: "Maximum number of waiting list tickets allowed per day",
		Category:    "waiting_list",
		IsEditable:  true,
		IsPublic:    true,
	},
	{
		Key:         "waiting_list.cleanup_retention_days",
		Value:       "7",
		Type:        SettingTypeInt,
		Description: "Number of days to keep completed/canceled entries before cleanup",
		Category:    "waiting_list",
		IsEditable:  true,
		IsPublic:    false,
	},
	{
		Key:         "waiting_list.job_enabled",
		Value:       "true",
		Type:        SettingTypeBool,
		Description: "Enable or disable the daily cleanup job",
		Category:    "waiting_list",
		IsEditable:  true,
		IsPublic:    false,
	},
	{
		Key:         "waiting_list.job_schedule",
		Value:       "0 0 * * *",
		Type:        SettingTypeString,
		Description: "Cron schedule for daily cleanup job (format: minute hour day month weekday)",
		Category:    "waiting_list",
		IsEditable:  true,
		IsPublic:    false,
	},
	{
		Key:         "waiting_list.allow_future_booking_days",
		Value:       "30",
		Type:        SettingTypeInt,
		Description: "How many days in advance customers can book tickets",
		Category:    "waiting_list",
		IsEditable:  true,
		IsPublic:    true,
	},
	{
		Key:         "business.shop_name",
		Value:       "Car Service Center",
		Type:        SettingTypeString,
		Description: "Name of the car service shop",
		Category:    "business",
		IsEditable:  true,
		IsPublic:    true,
	},
	{
		Key:         "business.opening_time",
		Value:       "08:00",
		Type:        SettingTypeString,
		Description: "Shop opening time (HH:MM format)",
		Category:    "business",
		IsEditable:  true,
		IsPublic:    true,
	},
	{
		Key:         "business.closing_time",
		Value:       "18:00",
		Type:        SettingTypeString,
		Description: "Shop closing time (HH:MM format)",
		Category:    "business",
		IsEditable:  true,
		IsPublic:    true,
	},
	{
		Key:         "business.working_days",
		Value:       "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday",
		Type:        SettingTypeString,
		Description: "Working days of the week (comma-separated)",
		Category:    "business",
		IsEditable:  true,
		IsPublic:    true,
	},
}

