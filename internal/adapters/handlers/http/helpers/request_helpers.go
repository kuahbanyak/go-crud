package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	apperrors "github.com/kuahbanyak/go-crud/pkg/errors"
)

// DecodeJSONBody decodes JSON request body into the provided destination
func DecodeJSONBody(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	return nil
}

// GetPathParam retrieves a path parameter by key
func GetPathParam(r *http.Request, key string) (string, bool) {
	vars := mux.Vars(r)
	value, exists := vars[key]
	return value, exists
}

// GetUUIDFromPath retrieves and parses a UUID from path parameters
func GetUUIDFromPath(r *http.Request, key string) (uuid.UUID, error) {
	value, exists := GetPathParam(r, key)
	if !exists {
		return uuid.UUID{}, fmt.Errorf("%s is required", key)
	}

	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s: %w", key, err)
	}

	return parsedUUID, nil
}

// GetAuthenticatedUserID retrieves the authenticated user ID from context
func GetAuthenticatedUserID(r *http.Request) (uuid.UUID, error) {
	userID, ok := r.Context().Value("id").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, apperrors.NewUnauthorizedError("User not authenticated")
	}
	return userID, nil
}

// GetAuthenticatedUserRole retrieves the authenticated user role from context
func GetAuthenticatedUserRole(r *http.Request) (string, error) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		return "", apperrors.NewUnauthorizedError("User role not found")
	}
	return role, nil
}

// ParsePaginationParams parses limit and offset from query parameters
func ParsePaginationParams(r *http.Request, defaultLimit, maxLimit int) (limit, offset int) {
	limit = defaultLimit
	offset = 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > maxLimit {
				limit = maxLimit
			}
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return limit, offset
}

// ParseIntQueryParam parses an integer query parameter
func ParseIntQueryParam(r *http.Request, key string, defaultValue int) int {
	if valueStr := r.URL.Query().Get(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// ParseFloatQueryParam parses a float query parameter
func ParseFloatQueryParam(r *http.Request, key string, defaultValue float64) float64 {
	if valueStr := r.URL.Query().Get(key); valueStr != "" {
		if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
			return value
		}
	}
	return defaultValue
}

// ParseBoolQueryParam parses a boolean query parameter
func ParseBoolQueryParam(r *http.Request, key string) *bool {
	if valueStr := r.URL.Query().Get(key); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return &value
		}
	}
	return nil
}

// GetQueryParam retrieves a query parameter value
func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
