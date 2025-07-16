package entity

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uint           `json:"id" gorm:"primarykey" example:"1"`
	Name        string         `json:"name" gorm:"not null;size:255" validate:"required,min=3,max=255" example:"Product Name"`
	Description string         `json:"description" gorm:"type:text" example:"Product description"`
	Price       float64        `json:"price" gorm:"not null" validate:"required,gt=0" example:"99.99"`
	Quantity    int            `json:"quantity" gorm:"not null;default:0" validate:"min=0" example:"100"`
	Category    string         `json:"category" gorm:"size:100" example:"Electronics"`
	SKU         string         `json:"sku" gorm:"unique;not null;size:100" validate:"required" example:"SKU001"`
	IsActive    bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

// CreateProductRequest represents the request for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255" example:"Product Name"`
	Description string  `json:"description" example:"Product description"`
	Price       float64 `json:"price" validate:"required,gt=0" example:"99.99"`
	Quantity    int     `json:"quantity" validate:"min=0" example:"100"`
	Category    string  `json:"category" example:"Electronics"`
	SKU         string  `json:"sku" validate:"required" example:"SKU001"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255" example:"Updated Product Name"`
	Description *string  `json:"description,omitempty" example:"Updated description"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0" example:"89.99"`
	Quantity    *int     `json:"quantity,omitempty" validate:"omitempty,min=0" example:"150"`
	Category    *string  `json:"category,omitempty" example:"Updated Category"`
	IsActive    *bool    `json:"is_active,omitempty" example:"false"`
}

type ProductResponse struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"Product Name"`
	Description string    `json:"description" example:"Product description"`
	Price       float64   `json:"price" example:"99.99"`
	Quantity    int       `json:"quantity" example:"100"`
	Category    string    `json:"category" example:"Electronics"`
	SKU         string    `json:"sku" example:"SKU001"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}
