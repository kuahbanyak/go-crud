package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleMechanic Role = "mechanic"
	RoleCustomer Role = "customer"
)

type User struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email    string `gorm:"not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     Role   `gorm:"type:varchar(20);default:'customer'" json:"role"`
	Address  string `json:"address"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
