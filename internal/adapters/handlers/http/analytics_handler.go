package http

import (
	"net/http"

	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type AnalyticsHandler struct {
	usecase *usecases.AnalyticsUsecase
}

func NewAnalyticsHandler(usecase *usecases.AnalyticsUsecase) *AnalyticsHandler {
	return &AnalyticsHandler{usecase: usecase}
}

// GetOverview returns overall analytics overview
func (h *AnalyticsHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.usecase.GetOverview(r.Context())
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to get overview", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Analytics overview retrieved successfully", overview)
}

// GetRevenueStats returns revenue statistics
func (h *AnalyticsHandler) GetRevenueStats(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "daily"
	}

	stats, err := h.usecase.GetRevenueStats(r.Context(), period)
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to get revenue stats", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Revenue statistics retrieved successfully", stats)
}

// GetServiceStats returns service statistics
func (h *AnalyticsHandler) GetServiceStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.usecase.GetServiceStats(r.Context())
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to get service stats", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Service statistics retrieved successfully", stats)
}

// GetQueueStats returns queue statistics
func (h *AnalyticsHandler) GetQueueStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.usecase.GetQueueStats(r.Context())
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to get queue stats", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Queue statistics retrieved successfully", stats)
}

// GetMechanicPerformance returns mechanic performance metrics
func (h *AnalyticsHandler) GetMechanicPerformance(w http.ResponseWriter, r *http.Request) {
	performance, err := h.usecase.GetMechanicPerformance(r.Context())
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to get mechanic performance", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Mechanic performance retrieved successfully", performance)
}
