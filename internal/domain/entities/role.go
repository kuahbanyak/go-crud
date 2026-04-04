package entities

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user role in the system
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	DisplayName string    `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}
