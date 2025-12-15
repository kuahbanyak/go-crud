package response

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	apperrors "github.com/kuahbanyak/go-crud/pkg/errors"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Get request ID from context if available
	requestID := ""
	if r, ok := w.(interface{ Context() context.Context }); ok {
		requestID = getRequestIDFromContext(r.Context())
	}

	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: requestID,
	}

	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, status int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Get request ID from context if available
	requestID := ""
	if r, ok := w.(interface{ Context() context.Context }); ok {
		requestID = getRequestIDFromContext(r.Context())
	}

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
		Success:   false,
		Message:   message,
		Error:     errorDetail,
		RequestID: requestID,
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

// getRequestIDFromContext retrieves request ID from context
func getRequestIDFromContext(ctx context.Context) string {
	type contextKey string
	const requestIDKey contextKey = "request_id"

	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// SuccessWithContext sends success response with request ID from context
func SuccessWithContext(ctx context.Context, w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: getRequestIDFromContext(ctx),
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorWithContext sends error response with request ID from context
func ErrorWithContext(ctx context.Context, w http.ResponseWriter, status int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Sanitize error in production
	var errorDetail interface{}
	if isProduction() {
		if status >= 400 && status < 500 {
			if err != nil {
				errorDetail = message
			}
		} else {
			errorDetail = "An internal error occurred"
		}
	} else {
		errorDetail = err
	}

	response := APIResponse{
		Success:   false,
		Message:   message,
		Error:     errorDetail,
		RequestID: getRequestIDFromContext(ctx),
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorFromAppError sends error response from AppError type
func ErrorFromAppError(ctx context.Context, w http.ResponseWriter, appErr *apperrors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)

	// Sanitize error in production
	var errorDetail interface{}
	if isProduction() {
		if appErr.StatusCode >= 400 && appErr.StatusCode < 500 {
			errorDetail = appErr.Details
		} else {
			errorDetail = "An internal error occurred"
		}
	} else {
		if appErr.Details != nil {
			errorDetail = appErr.Details
		} else if appErr.Internal != nil {
			errorDetail = appErr.Internal.Error()
		}
	}

	response := APIResponse{
		Success:   false,
		Message:   appErr.Message,
		Error:     errorDetail,
		RequestID: getRequestIDFromContext(ctx),
	}

	json.NewEncoder(w).Encode(response)
}
