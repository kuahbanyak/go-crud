package dto

import (
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Address   string `json:"address,omitempty" validate:"omitempty,min=1,max=100"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
}
type UserResponse struct {
	ID    types.MSSQLUUID `json:"id"`
	Email string          `json:"email"`
	Name  string          `json:"name"`
	Phone string          `json:"phone"`
	Roles []RoleResponse  `json:"roles,omitempty"` // RBAC roles from roles table
}
type LoginResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
	ExpiresIn   int64        `json:"expires_in"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
