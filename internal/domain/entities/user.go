package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type User struct {
	ID        types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
	Email     string          `gorm:"not null;unique" json:"email"`
	Password  string          `gorm:"not null" json:"-"`
	Name      string          `json:"name"`
	Phone     string          `json:"phone"`
	Address   string          `json:"address"`

	// Many-to-many relationship with roles table (RBAC system)
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID.String() == "00000000-0000-0000-0000-000000000000" {
		u.ID = types.NewMSSQLUUID()
	}
	return nil
}
