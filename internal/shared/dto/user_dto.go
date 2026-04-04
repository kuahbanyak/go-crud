package dto

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Address  string `json:"address,omitempty" validate:"omitempty,min=1,max=100"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type UpdateUserRequest struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Address  string `json:"address,omitempty" validate:"omitempty,min=1,max=100"`
	Password string `json:"password,omitempty" validate:"omitempty,min=10,max=20"`
}
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string          `json:"email"`
	Name  string          `json:"name"`
	Phone string          `json:"phone"`
	Roles []RoleResponse  `json:"roles,omitempty"`
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
