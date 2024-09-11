package entity

import (
	"github.com/google/uuid"
)

// Account represents a user management.
type Account struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	RepeatPassword string    `json:"repeat_password"`
	CreatedAt      uuid.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      uuid.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
}
