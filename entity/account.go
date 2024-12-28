package entity

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	Id             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id" `
	Username       string    `json:"username" `
	Password       string    `json:"password" `
	RepeatPassword string    `json:"repeat_password" `
	CreatedAt      time.Time `json:"created_at" `
	UpdateAt       time.Time `json:"updated_at" `
}
