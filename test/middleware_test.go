package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-crud/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminOnlyMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		role           interface{}
		expectedStatus int
	}{
		{
			name:           "Admin role - access granted",
			role:           "admin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-admin role - access denied",
			role:           "user",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "No role - access denied",
			role:           nil,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder and a test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set the role in the context if provided
			if tt.role != nil {
				c.Set("role", tt.role)
			}

			// Create a dummy handler to test middleware
			handler := middleware.AdminOnly()
			handler(c)

			// Assert the response status
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
