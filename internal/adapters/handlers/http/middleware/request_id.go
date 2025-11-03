package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// RequestIDKey is the context key for the request ID
type contextKey string

const RequestIDKey contextKey = "request_id"

// RequestID middleware adds a unique request ID to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists in header
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate new UUID for request ID
			requestID = uuid.New().String()
		}

		// Add to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Add to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Pass to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
