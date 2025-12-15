package pagination

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Params holds pagination parameters
type Params struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by,omitempty"`
	SortDir  string `json:"sort_dir,omitempty"`
}

// Response holds paginated response data
type Response struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	TotalItems int64       `json:"total_items"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// ParseParams extracts pagination parameters from request
func ParseParams(r *http.Request) Params {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = DefaultPage
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortDir := r.URL.Query().Get("sort_dir")
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}
}

// Apply applies pagination to GORM query
func (p Params) Apply(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.PageSize

	query := db.Limit(p.PageSize).Offset(offset)

	if p.SortBy != "" {
		orderClause := fmt.Sprintf("%s %s", p.SortBy, p.SortDir)
		query = query.Order(orderClause)
	}

	return query
}

// GetOffset returns the offset for current page
func (p Params) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the limit (page size)
func (p Params) GetLimit() int {
	return p.PageSize
}

// BuildResponse creates a paginated response
func BuildResponse(data interface{}, totalItems int64, params Params) Response {
	totalPages := int(math.Ceil(float64(totalItems) / float64(params.PageSize)))

	return Response{
		Data:       data,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}

// FilterParams holds common filter parameters
type FilterParams struct {
	Search    string `json:"search,omitempty"`
	Status    string `json:"status,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// ParseFilterParams extracts filter parameters from request
func ParseFilterParams(r *http.Request) FilterParams {
	return FilterParams{
		Search:    r.URL.Query().Get("search"),
		Status:    r.URL.Query().Get("status"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}
}

// ApplySearch applies search filter to GORM query
func ApplySearch(db *gorm.DB, search string, fields ...string) *gorm.DB {
	if search == "" || len(fields) == 0 {
		return db
	}

	var conditions []interface{}
	var values []interface{}

	for _, field := range fields {
		conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
		values = append(values, "%"+search+"%")
	}

	// Build OR condition
	query := ""
	for i, condition := range conditions {
		if i > 0 {
			query += " OR "
		}
		query += condition.(string)
	}

	return db.Where(query, values...)
}

// ApplyStatusFilter applies status filter to GORM query
func ApplyStatusFilter(db *gorm.DB, status string) *gorm.DB {
	if status != "" {
		return db.Where("status = ?", status)
	}
	return db
}

// ApplyDateRangeFilter applies date range filter to GORM query
func ApplyDateRangeFilter(db *gorm.DB, field, startDate, endDate string) *gorm.DB {
	if startDate != "" {
		db = db.Where(fmt.Sprintf("%s >= ?", field), startDate)
	}
	if endDate != "" {
		db = db.Where(fmt.Sprintf("%s <= ?", field), endDate)
	}
	return db
}
