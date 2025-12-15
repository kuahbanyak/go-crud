package http

import (
	"net/http"
	"runtime"
	"time"

	"github.com/kuahbanyak/go-crud/pkg/response"
)

type VersionHandler struct{}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

func (h *VersionHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	versionInfo := map[string]interface{}{
		"api_version": "v1",
		"go_version":  runtime.Version(),
		"build_time":  buildTime,
		"build_mode":  getBuildMode(),
		"endpoints": map[string]interface{}{
			"base_url":      "/api/v1",
			"health":        "/health",
			"health_live":   "/health/live",
			"health_ready":  "/health/ready",
			"health_detail": "/health/detail",
			"version":       "/api/version",
		},
		"features": []string{
			"authentication",
			"user_management",
			"product_management",
			"vehicle_management",
			"waiting_list",
			"maintenance_tracking",
			"settings_management",
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "API version information", versionInfo)
}

var buildTime = time.Now().Format(time.RFC3339) // This would be set via build flags

func getBuildMode() string {
	// Simple check, could be enhanced with build tags
	if runtime.GOOS == "windows" {
		return "development"
	}
	return "production"
}
