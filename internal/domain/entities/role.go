package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

// Role represents a user role in the system
type Role struct {
	ID          types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:NEWID()" json:"id"`
	Name        string          `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	DisplayName string          `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string          `gorm:"type:varchar(255)" json:"description"`
	IsActive    bool            `gorm:"default:1" json:"is_active"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}
