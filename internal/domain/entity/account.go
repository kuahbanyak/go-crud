package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uniqueidentifier;primary_key;default:NEWID()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string         `json:"username" gorm:"unique;not null;size:50" example:"johndoe"`
	Email     string         `json:"email" gorm:"unique;not null;size:100" example:"john@example.com"`
	Password  string         `json:"-" gorm:"not null"` // Never expose password in JSON
	FirstName string         `json:"first_name" gorm:"size:50" example:"John"`
	LastName  string         `json:"last_name" gorm:"size:50" example:"Doe"`
	IsActive  bool           `json:"is_active" gorm:"default:true" example:"true"`
	Role      string         `json:"role" gorm:"default:'user';size:20" example:"user"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Account) TableName() string {
	return "accounts"
}

type CreateAccountRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50" example:"johndoe"`
	Email     string `json:"email" validate:"required,email" example:"john@example.com"`
	Password  string `json:"password" validate:"required,min=6" example:"password123"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50" example:"John"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50" example:"Doe"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"johndoe"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// UpdateAccountRequest represents the request for updating an account
type UpdateAccountRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email" example:"newemail@example.com"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=50" example:"John"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=50" example:"Doe"`
	IsActive  *bool   `json:"is_active,omitempty" example:"false"`
}

// ChangePasswordRequest represents the request for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"oldpassword123"`
	NewPassword     string `json:"new_password" validate:"required,min=6" example:"newpassword123"`
}

// AccountResponse represents the response for account operations
type AccountResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	IsActive  bool      `json:"is_active" example:"true"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Account     *AccountResponse `json:"account"`
	AccessToken string           `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string           `json:"token_type" example:"Bearer"`
	ExpiresIn   int64            `json:"expires_in" example:"3600"`
}
