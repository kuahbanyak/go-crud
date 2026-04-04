package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RoleID     uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	AssignedBy uuid.UUID `gorm:"type:uuid" json:"assigned_by"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
}

// TableName specifies the table name for UserRole
func (UserRole) TableName() string {
	return "user_roles"
}
