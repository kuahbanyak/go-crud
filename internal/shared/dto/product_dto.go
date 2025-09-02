package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	Category    string  `json:"category" validate:"required,min=1,max=100"`
	SKU         string  `json:"sku,omitempty" validate:"max=50"`
	IsActive    bool    `json:"is_active"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description string  `json:"description,omitempty" validate:"max=1000"`
	Price       float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Stock       int     `json:"stock,omitempty" validate:"omitempty,min=0"`
	Category    string  `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	SKU         string  `json:"sku,omitempty" validate:"max=50"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	SKU         string    `json:"sku"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock" validate:"required,min=0"`
}
