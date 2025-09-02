package entities

import (
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleMechanic Role = "mechanic"
	RoleCustomer Role = "customer"
)

type User struct {
	ID        types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     Role   `gorm:"type:varchar(20);default:'customer'" json:"role"`
	Address  string `json:"address"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID.String() == "00000000-0000-0000-0000-000000000000" {
		u.ID = types.NewMSSQLUUID()
	}
	return nil
}
