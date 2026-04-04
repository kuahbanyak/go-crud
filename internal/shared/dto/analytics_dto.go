package dto

import "time"

// AnalyticsOverviewResponse represents overall analytics data
type AnalyticsOverviewResponse struct {
	TodayRevenue      int       `json:"today_revenue"`
	TotalRevenue      int       `json:"total_revenue"`
	TotalCustomers    int       `json:"total_customers"`
	ActiveServices    int       `json:"active_services"`
	CompletedServices int       `json:"completed_services"`
	PendingInvoices   int       `json:"pending_invoices"`
	TodayQueue        int       `json:"today_queue"`
	AverageWaitTime   float64   `json:"average_wait_time_minutes"`
	Timestamp         time.Time `json:"timestamp"`
}

// RevenueStatsResponse represents revenue statistics
type RevenueStatsResponse struct {
	Period     string             `json:"period"`
	TotalCount int                `json:"total_count"`
	Data       []RevenueDataPoint `json:"data"`
}

type RevenueDataPoint struct {
	Date   string `json:"date"`
	Amount int    `json:"amount"`
	Count  int    `json:"count"`
}

// ServiceStatsResponse represents service statistics
type ServiceStatsResponse struct {
	TotalServices     int                `json:"total_services"`
	StatusBreakdown   map[string]int     `json:"status_breakdown"`
	ServicesByType    []ServiceTypeCount `json:"services_by_type"`
	AverageCompletion float64            `json:"average_completion_hours"`
}

type ServiceTypeCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// QueueStatsResponse represents queue statistics
type QueueStatsResponse struct {
	TodayTotal       int     `json:"today_total"`
	CurrentWaiting   int     `json:"current_waiting"`
	CurrentInService int     `json:"current_in_service"`
	Completed        int     `json:"completed"`
	NoShow           int     `json:"no_show"`
	Cancelled        int     `json:"cancelled"`
	AverageWaitTime  float64 `json:"average_wait_time_minutes"`
	PeakHour         int     `json:"peak_hour"`
}

// MechanicPerformanceResponse represents mechanic performance data
type MechanicPerformanceResponse struct {
	MechanicID        string  `json:"mechanic_id"`
	MechanicName      string  `json:"mechanic_name"`
	TotalServices     int     `json:"total_services"`
	CompletedServices int     `json:"completed_services"`
	AverageCompletion float64 `json:"average_completion_hours"`
	CustomerRating    float64 `json:"customer_rating"`
	Efficiency        float64 `json:"efficiency_percentage"`
}

// VisitorStatsResponse represents visitor/customer statistics by time period
type VisitorStatsResponse struct {
	Period             string                        `json:"period"`              // "7days", "30days", "3months"
	TotalVisitors      int                           `json:"total_visitors"`      // Total unique visitors
	TotalTickets       int                           `json:"total_tickets"`       // Total tickets created
	NewCustomers       int                           `json:"new_customers"`       // New customers registered
	ReturningCustomers int                           `json:"returning_customers"` // Customers with multiple visits
	ByServiceCategory  []ServiceCategoryVisitorStats `json:"by_service_category"`
	ByDate             []VisitorDataPoint            `json:"by_date"`
	Timestamp          time.Time                     `json:"timestamp"`
}

// ServiceCategoryVisitorStats represents visitor statistics grouped by service category
type ServiceCategoryVisitorStats struct {
	Category     string  `json:"category"`      // Engine, Brakes, Tires, etc.
	VisitorCount int     `json:"visitor_count"` // Unique visitors for this category
	TicketCount  int     `json:"ticket_count"`  // Total tickets for this category
	Percentage   float64 `json:"percentage"`    // Percentage of total
}

// VisitorDataPoint represents visitor count for a specific date
type VisitorDataPoint struct {
	Date         string `json:"date"`          // YYYY-MM-DD
	VisitorCount int    `json:"visitor_count"` // Unique visitors
	TicketCount  int    `json:"ticket_count"`  // Tickets created
	NewCustomers int    `json:"new_customers"` // New registrations
}
