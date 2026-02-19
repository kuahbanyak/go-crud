package dto

type ServiceItemResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	EstimatedTime   int     `json:"estimated_time"`
	EstimatedCost   float64 `json:"estimated_cost"`
	DisplayOrder    int     `json:"display_order"`
	RequiresBooking bool    `json:"requires_booking"`
}

type CreateServiceItemRequest struct {
	Name            string  `json:"name" validate:"required"`
	Description     string  `json:"description"`
	Category        string  `json:"category" validate:"required"`
	EstimatedTime   int     `json:"estimated_time"`
	EstimatedCost   float64 `json:"estimated_cost"`
	DisplayOrder    int     `json:"display_order"`
	RequiresBooking bool    `json:"requires_booking"`
}

type UpdateServiceItemRequest struct {
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	Category        *string  `json:"category"`
	EstimatedTime   *int     `json:"estimated_time"`
	EstimatedCost   *float64 `json:"estimated_cost"`
	IsActive        *bool    `json:"is_active"`
	DisplayOrder    *int     `json:"display_order"`
	RequiresBooking *bool    `json:"requires_booking"`
}

type ServiceItemsByCategoryResponse struct {
	Category string                `json:"category"`
	Items    []ServiceItemResponse `json:"items"`
}
