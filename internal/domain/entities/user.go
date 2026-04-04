package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"not null;unique" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `json:"name"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`

	// Many-to-many relationship with roles table (RBAC system)
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
