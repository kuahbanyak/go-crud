package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a structured application error
type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	StatusCode int         `json:"-"`
	Details    interface{} `json:"details,omitempty"`
	Internal   error       `json:"-"` // Internal error (not exposed to client)
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Error constructors
func NewBadRequestError(message string, details ...interface{}) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details:    getDetails(details...),
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:       "FORBIDDEN",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewValidationError(message string, details interface{}) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
}

func NewInternalError(message string, internal error) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Internal:   internal,
	}
}

// Business logic errors
func NewBusinessError(code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusUnprocessableEntity,
	}
}

// Specific business errors
var (
	ErrDailyLimitReached = &AppError{
		Code:       "DAILY_LIMIT_REACHED",
		Message:    "Daily ticket limit has been reached",
		StatusCode: http.StatusUnprocessableEntity,
	}

	ErrServiceNotStarted = &AppError{
		Code:       "SERVICE_NOT_STARTED",
		Message:    "Service must be started before performing this action",
		StatusCode: http.StatusUnprocessableEntity,
	}

	ErrInvalidStatus = &AppError{
		Code:       "INVALID_STATUS",
		Message:    "Invalid status transition",
		StatusCode: http.StatusUnprocessableEntity,
	}

	ErrUnauthorizedVehicle = &AppError{
		Code:       "UNAUTHORIZED_VEHICLE",
		Message:    "You don't have permission to access this vehicle",
		StatusCode: http.StatusForbidden,
	}

	ErrApprovalRequired = &AppError{
		Code:       "APPROVAL_REQUIRED",
		Message:    "Customer approval is required before proceeding",
		StatusCode: http.StatusUnprocessableEntity,
	}
)

func getDetails(details ...interface{}) interface{} {
	if len(details) == 0 {
		return nil
	}
	if len(details) == 1 {
		return details[0]
	}
	return details
}

// Wrap wraps an internal error with a message
func Wrap(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		appErr.Internal = err
		return appErr
	}
	return NewInternalError(message, err)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors map[string][]string

func (v ValidationErrors) Add(field, message string) {
	v[field] = append(v[field], message)
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

func (v ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %d error(s)", len(v))
}

// NewValidationErrors creates validation error with field-specific messages
func NewValidationErrors(errors ValidationErrors) *AppError {
	return NewValidationError("Validation failed", errors)
}
