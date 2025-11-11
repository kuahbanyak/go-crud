package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kuahbanyak/go-crud/pkg/response"
)

const MaxRequestSize = 10 << 20 // 10MB

func ValidateRequestSize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only validate for requests with body
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if r.Header.Get("Content-Type") != "application/json" {
				response.Error(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
				return
			}
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

			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		next.ServeHTTP(w, r)
	})
}
