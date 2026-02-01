package helpers

import (
	"net/http"

	apperrors "github.com/kuahbanyak/go-crud/pkg/errors"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

// RespondSuccess sends a success response with context
func RespondSuccess(w http.ResponseWriter, r *http.Request, status int, message string, data interface{}) {
	response.SuccessWithContext(r.Context(), w, status, message, data)
}

// RespondError sends an error response based on error type
func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	// Default to internal server error for unknown errors
	response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "An error occurred", err.Error())
}

// RespondBadRequest sends a bad request error response
func RespondBadRequest(w http.ResponseWriter, r *http.Request, message string, err error) {
	var errDetail interface{}
	if err != nil {
		errDetail = err.Error()
	}
	response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, message, errDetail)
}

// RespondNotFound sends a not found error response
func RespondNotFound(w http.ResponseWriter, r *http.Request, message string) {
	response.ErrorWithContext(r.Context(), w, http.StatusNotFound, message, nil)
}

// RespondUnauthorized sends an unauthorized error response
func RespondUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	response.ErrorWithContext(r.Context(), w, http.StatusUnauthorized, message, nil)
}

// RespondForbidden sends a forbidden error response
func RespondForbidden(w http.ResponseWriter, r *http.Request, message string) {
	response.ErrorWithContext(r.Context(), w, http.StatusForbidden, message, nil)
}

// RespondInternalError sends an internal server error response
func RespondInternalError(w http.ResponseWriter, r *http.Request, message string, err error) {
	var errDetail interface{}
	if err != nil {
		errDetail = err.Error()
	}
	response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, message, errDetail)
}

// RespondValidationError sends a validation error response
func RespondValidationError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	appErr := apperrors.NewValidationError("Validation failed", errors)
	response.ErrorFromAppError(r.Context(), w, appErr)
}

// RespondConflict sends a conflict error response
func RespondConflict(w http.ResponseWriter, r *http.Request, message string) {
	response.ErrorWithContext(r.Context(), w, http.StatusConflict, message, nil)
}
