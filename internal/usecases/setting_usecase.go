package usecases

import (
	"context"
	"errors"
	"strconv"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type SettingUsecase struct {
	settingRepo repositories.SettingRepository
}

func NewSettingUsecase(settingRepo repositories.SettingRepository) *SettingUsecase {
	return &SettingUsecase{
		settingRepo: settingRepo,
	}
}

func (u *SettingUsecase) GetSetting(ctx context.Context, key string) (*entities.Setting, error) {
	return u.settingRepo.GetByKey(ctx, key)
}

func (u *SettingUsecase) GetSettingsByCategory(ctx context.Context, category string) ([]*entities.Setting, error) {
	return u.settingRepo.GetByCategory(ctx, category)
}

func (u *SettingUsecase) GetAllSettings(ctx context.Context) ([]*entities.Setting, error) {
	return u.settingRepo.GetAll(ctx)
}

func (u *SettingUsecase) GetPublicSettings(ctx context.Context) ([]*entities.Setting, error) {
	return u.settingRepo.GetPublic(ctx)
}

func (u *SettingUsecase) UpdateSetting(ctx context.Context, key, value string) error {
	setting, err := u.settingRepo.GetByKey(ctx, key)
	if err != nil {
		return err
	}
	if setting == nil {
		return errors.New("setting not found")
	}

	if !setting.IsEditable {
		return errors.New("this setting cannot be edited")
	}

	setting.Value = value
	return u.settingRepo.Update(ctx, setting)
}

func (u *SettingUsecase) CreateSetting(ctx context.Context, setting *entities.Setting) error {
	// Check if key already exists
	existing, err := u.settingRepo.GetByKey(ctx, setting.Key)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("setting with this key already exists")
	}

	return u.settingRepo.Create(ctx, setting)
}

func (u *SettingUsecase) DeleteSetting(ctx context.Context, id types.MSSQLUUID) error {
	return u.settingRepo.Delete(ctx, id)
}

// Helper methods to get typed values

func (u *SettingUsecase) GetIntValue(ctx context.Context, key string, defaultValue int) int {
	setting, err := u.settingRepo.GetByKey(ctx, key)
	if err != nil || setting == nil {
		return defaultValue
	}

	val, err := strconv.Atoi(setting.Value)
	if err != nil {
		return defaultValue
	}
	return val
}

func (u *SettingUsecase) GetStringValue(ctx context.Context, key string, defaultValue string) string {
	setting, err := u.settingRepo.GetByKey(ctx, key)
	if err != nil || setting == nil {
		return defaultValue
	}
	return setting.Value
}

func (u *SettingUsecase) GetBoolValue(ctx context.Context, key string, defaultValue bool) bool {
	setting, err := u.settingRepo.GetByKey(ctx, key)
	if err != nil || setting == nil {
		return defaultValue
	}

	val, err := strconv.ParseBool(setting.Value)
	if err != nil {
		return defaultValue
	}
	return val
}

func (u *SettingUsecase) GetFloatValue(ctx context.Context, key string, defaultValue float64) float64 {
	setting, err := u.settingRepo.GetByKey(ctx, key)
	if err != nil || setting == nil {
		return defaultValue
	}

	val, err := strconv.ParseFloat(setting.Value, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// Specific waiting list settings helpers
func (u *SettingUsecase) GetMaxTicketsPerDay(ctx context.Context) int {
	return u.GetIntValue(ctx, "waiting_list.max_tickets_per_day", 10)
}

func (u *SettingUsecase) GetCleanupRetentionDays(ctx context.Context) int {
	return u.GetIntValue(ctx, "waiting_list.cleanup_retention_days", 7)
}

func (u *SettingUsecase) IsCleanupJobEnabled(ctx context.Context) bool {
	return u.GetBoolValue(ctx, "waiting_list.job_enabled", true)
}

func (u *SettingUsecase) GetJobSchedule(ctx context.Context) string {
	return u.GetStringValue(ctx, "waiting_list.job_schedule", "0 0 * * *")
}
