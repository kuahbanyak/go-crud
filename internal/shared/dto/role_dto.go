package dto

import "github.com/kuahbanyak/go-crud/internal/shared/types"

// RoleResponse represents role data in API responses
type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateRoleRequest represents request to create a new role
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	DisplayName string `json:"display_name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=255"`
	IsActive    bool   `json:"is_active"`
}

// UpdateRoleRequest represents request to update a role
type UpdateRoleRequest struct {
	DisplayName string `json:"display_name" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=255"`
	IsActive    *bool  `json:"is_active"`
}

// AssignRoleRequest represents request to assign role to user
type AssignRoleRequest struct {
	RoleID types.MSSQLUUID `json:"role_id" validate:"required"`
}

// RemoveRoleRequest represents request to remove role from user
type RemoveRoleRequest struct {
	RoleID types.MSSQLUUID `json:"role_id" validate:"required"`
}

