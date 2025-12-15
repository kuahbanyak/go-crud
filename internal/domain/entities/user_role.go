package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID         types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:NEWID()" json:"id"`
	UserID     types.MSSQLUUID `gorm:"type:uniqueidentifier;not null" json:"user_id"`
	RoleID     types.MSSQLUUID `gorm:"type:uniqueidentifier;not null" json:"role_id"`
	AssignedBy types.MSSQLUUID `gorm:"type:uniqueidentifier" json:"assigned_by"`
	AssignedAt time.Time       `gorm:"autoCreateTime" json:"assigned_at"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
}

// TableName specifies the table name for UserRole
func (UserRole) TableName() string {
	return "user_roles"
}
