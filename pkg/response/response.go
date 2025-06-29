package response

import (
	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:"Error description"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success" example:"true"`
	Message    string      `json:"message" example:"Data retrieved successfully"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Error      string      `json:"error,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Total  int64 `json:"total" example:"100"`
	Limit  int   `json:"limit" example:"10"`
	Offset int   `json:"offset" example:"0"`
	Page   int   `json:"page" example:"1"`
	Pages  int   `json:"pages" example:"10"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, statusCode int, message string, data interface{}, total int64, limit, offset int) {
	page := (offset / limit) + 1
	pages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
			Page:   page,
			Pages:  pages,
		},
	})
}
