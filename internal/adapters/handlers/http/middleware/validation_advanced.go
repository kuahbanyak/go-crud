package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	apperrors "github.com/kuahbanyak/go-crud/pkg/errors"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

// Validator interface for custom validation
type Validator interface {
	Validate() error
}

// ValidateRequest validates request body against struct tags
func ValidateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip validation for GET, DELETE, HEAD, OPTIONS
		if r.Method == "GET" || r.Method == "DELETE" ||
			r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Read and validate JSON structure
		body, err := io.ReadAll(r.Body)
		if err != nil {
			appErr := apperrors.NewBadRequestError("Failed to read request body")
			response.ErrorWithContext(r.Context(), w, appErr.StatusCode, appErr.Message, nil)
			return
		}
		defer r.Body.Close()

		// Validate JSON syntax
		var js json.RawMessage
		if err := json.Unmarshal(body, &js); err != nil {
			appErr := apperrors.NewValidationError("Invalid JSON format", map[string]string{
				"error": "Request body must be valid JSON",
			})
			response.ErrorWithContext(r.Context(), w, appErr.StatusCode, appErr.Message, appErr.Details)
			return
		}

		// Restore body for handlers to read
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		next.ServeHTTP(w, r)
	})
}

// ValidateStruct performs validation on struct based on tags
func ValidateStruct(v interface{}) *apperrors.AppError {
	errors := make(apperrors.ValidationErrors)

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Get validation tags
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		fieldName := getJSONFieldName(field)
		rules := strings.Split(validateTag, ",")

		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if err := validateField(fieldName, value, rule); err != nil {
				errors.Add(fieldName, err.Error())
			}
		}
	}

	if errors.HasErrors() {
		return apperrors.NewValidationErrors(errors)
	}

	// Check if struct implements Validator interface
	if validator, ok := v.(Validator); ok {
		if err := validator.Validate(); err != nil {
			return apperrors.NewValidationError("Custom validation failed", err.Error())
		}
	}

	return nil
}

func validateField(fieldName string, value reflect.Value, rule string) error {
	parts := strings.Split(rule, "=")
	ruleName := parts[0]
	var ruleValue string
	if len(parts) > 1 {
		ruleValue = parts[1]
	}

	switch ruleName {
	case "required":
		if isZero(value) {
			return fmt.Errorf("%s is required", fieldName)
		}

	case "min":
		if value.Kind() == reflect.String {
			if len(value.String()) < parseInt(ruleValue) {
				return fmt.Errorf("%s must be at least %s characters", fieldName, ruleValue)
			}
		} else if isNumeric(value) {
			if getNumericValue(value) < float64(parseInt(ruleValue)) {
				return fmt.Errorf("%s must be at least %s", fieldName, ruleValue)
			}
		}

	case "max":
		if value.Kind() == reflect.String {
			if len(value.String()) > parseInt(ruleValue) {
				return fmt.Errorf("%s must be at most %s characters", fieldName, ruleValue)
			}
		} else if isNumeric(value) {
			if getNumericValue(value) > float64(parseInt(ruleValue)) {
				return fmt.Errorf("%s must be at most %s", fieldName, ruleValue)
			}
		}

	case "email":
		if value.Kind() == reflect.String && value.String() != "" {
			if !isValidEmail(value.String()) {
				return fmt.Errorf("%s must be a valid email address", fieldName)
			}
		}

	case "phone":
		if value.Kind() == reflect.String && value.String() != "" {
			if !isValidPhone(value.String()) {
				return fmt.Errorf("%s must be a valid phone number", fieldName)
			}
		}

	case "uuid":
		if value.Kind() == reflect.String && value.String() != "" {
			if !isValidUUID(value.String()) {
				return fmt.Errorf("%s must be a valid UUID", fieldName)
			}
		}

	case "oneof":
		allowed := strings.Split(ruleValue, " ")
		found := false
		strVal := fmt.Sprintf("%v", value.Interface())
		for _, a := range allowed {
			if strVal == a {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s must be one of: %s", fieldName, strings.Join(allowed, ", "))
		}
	}

	return nil
}

// ValidatePathParams validates URL path parameters
func ValidatePathParams(params map[string]string) *apperrors.AppError {
	errors := make(apperrors.ValidationErrors)

	for key, value := range params {
		if value == "" {
			errors.Add(key, fmt.Sprintf("%s is required", key))
		}
	}

	if errors.HasErrors() {
		return apperrors.NewValidationErrors(errors)
	}
	return nil
}

// GetValidatedPathParam extracts and validates a path parameter
func GetValidatedPathParam(r *http.Request, paramName string) (string, *apperrors.AppError) {
	vars := mux.Vars(r)
	value, exists := vars[paramName]
	if !exists || value == "" {
		return "", apperrors.NewBadRequestError(fmt.Sprintf("%s is required", paramName))
	}
	return value, nil
}

// Helper functions
func getJSONFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return field.Name
	}
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}

func isNumeric(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func getNumericValue(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	return 0
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func isValidPhone(phone string) bool {
	pattern := `^[\d\s\-\+\(\)]{10,20}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func isValidUUID(uuid string) bool {
	pattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	matched, _ := regexp.MatchString(pattern, uuid)
	return matched
}

// Context keys for validated data
type validationContextKey string

const (
	ValidatedBodyKey validationContextKey = "validated_body"
)

// SetValidatedBody stores validated body in context
func SetValidatedBody(ctx context.Context, body interface{}) context.Context {
	return context.WithValue(ctx, ValidatedBodyKey, body)
}

// GetValidatedBody retrieves validated body from context
func GetValidatedBody(ctx context.Context) interface{} {
	return ctx.Value(ValidatedBodyKey)
}
