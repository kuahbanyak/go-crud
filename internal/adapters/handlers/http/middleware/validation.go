package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kuahbanyak/go-crud/pkg/response"
)

// MaxRequestSize defines the maximum allowed request body size (10MB)
const MaxRequestSize = 10 << 20 // 10MB

// ValidateRequestSize middleware limits the size of incoming requests
func ValidateRequestSize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply to requests with body
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateJSON middleware ensures the request body is valid JSON
func ValidateJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only validate for requests with body
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if r.Header.Get("Content-Type") != "application/json" {
				response.Error(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
				return
			}

			// Try to decode to validate JSON structure
			body, err := io.ReadAll(r.Body)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Failed to read request body", err)
				return
			}
			defer r.Body.Close()

			var js json.RawMessage
			if err := json.Unmarshal(body, &js); err != nil {
				response.Error(w, http.StatusBadRequest, "Invalid JSON format", err)
				return
			}

			// Create a new reader with the body for the next handler
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		next.ServeHTTP(w, r)
	})
}
