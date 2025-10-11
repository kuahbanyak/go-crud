package dto

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type CreateSettingRequest struct {
	Key         string `json:"key" validate:"required"`
	Value       string `json:"value" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=int string bool float"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"required"`
	IsEditable  bool   `json:"is_editable"`
	IsPublic    bool   `json:"is_public"`
}

type UpdateSettingRequest struct {
	Value string `json:"value" validate:"required"`
}

type SettingResponse struct {
	ID          types.MSSQLUUID `json:"id"`
	Key         string          `json:"key"`
	Value       string          `json:"value"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	IsEditable  bool            `json:"is_editable"`
	IsPublic    bool            `json:"is_public"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type SettingsListResponse struct {
	Settings []SettingResponse `json:"settings"`
	Total    int               `json:"total"`
}

type SettingsByCategoryResponse struct {
	Category string            `json:"category"`
	Settings []SettingResponse `json:"settings"`
}
