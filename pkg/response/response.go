package response

import (
	"encoding/json"
	"net/http"
	"os"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, status int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Sanitize error in production - don't expose internal details
	var errorDetail interface{}
	if isProduction() {
		// In production, only include error for 4xx client errors, hide 5xx details
		if status >= 400 && status < 500 {
			if err != nil {
				errorDetail = message // Use the message instead of raw error
			}
		} else {
			errorDetail = "An internal error occurred" // Generic message for server errors
		}
	} else {
		// In development, include full error details
		errorDetail = err
	}

	response := APIResponse{
		Success: false,
		Message: message,
		Error:   errorDetail,
	}

	json.NewEncoder(w).Encode(response)
}

// isProduction checks if the application is running in production mode
func isProduction() bool {
	env := getEnv("GIN_MODE", "debug")
	return env == "release" || getEnv("RAILWAY_ENVIRONMENT", "") != ""
}

// getEnv retrieves an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
