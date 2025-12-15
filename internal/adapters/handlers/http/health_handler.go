package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/kuahbanyak/go-crud/pkg/response"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck - Simple liveness check
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	type HealthResponse struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}

	resp := HealthResponse{
		Status:    "ok",
		Message:   "Server is running",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Health check passed", resp)
}

// LivenessCheck - Basic liveness probe
func (h *HealthHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "alive",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadinessCheck - Check if app is ready to serve traffic
func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := make(map[string]interface{})
	allHealthy := true

	// Check database
	if err := h.db.PingContext(ctx); err != nil {
		checks["database"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
		allHealthy = false
	} else {
		checks["database"] = map[string]interface{}{
			"status": "healthy",
		}
	}

	// Check database stats
	stats := h.db.Stats()
	checks["database_stats"] = map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
	}

	status := "ready"
	statusCode := http.StatusOK
	if !allHealthy {
		status = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	resp := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks":    checks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// DetailedHealthCheck - Comprehensive health information (admin only)
func (h *HealthHandler) DetailedHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	checks := make(map[string]interface{})
	allHealthy := true

	// Database check
	dbStart := time.Now()
	if err := h.db.PingContext(ctx); err != nil {
		checks["database"] = map[string]interface{}{
			"status":        "unhealthy",
			"error":         err.Error(),
			"response_time": time.Since(dbStart).Milliseconds(),
		}
		allHealthy = false
	} else {
		stats := h.db.Stats()
		checks["database"] = map[string]interface{}{
			"status":              "healthy",
			"response_time":       time.Since(dbStart).Milliseconds(),
			"open_connections":    stats.OpenConnections,
			"in_use":              stats.InUse,
			"idle":                stats.Idle,
			"wait_count":          stats.WaitCount,
			"wait_duration":       stats.WaitDuration.Milliseconds(),
			"max_idle_closed":     stats.MaxIdleClosed,
			"max_lifetime_closed": stats.MaxLifetimeClosed,
		}
	}

	// System info
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	checks["system"] = map[string]interface{}{
		"goroutines":      runtime.NumGoroutine(),
		"memory_alloc_mb": memStats.Alloc / 1024 / 1024,
		"memory_total_mb": memStats.TotalAlloc / 1024 / 1024,
		"memory_sys_mb":   memStats.Sys / 1024 / 1024,
		"gc_runs":         memStats.NumGC,
		"cpu_count":       runtime.NumCPU(),
	}

	status := "healthy"
	statusCode := http.StatusOK
	if !allHealthy {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	resp := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
		"version":   getVersion(),
		"checks":    checks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

var startTime = time.Now()

func getVersion() map[string]string {
	return map[string]string{
		"api":   "v1",
		"go":    runtime.Version(),
		"build": "dev", // This can be set via build flags
	}
}
