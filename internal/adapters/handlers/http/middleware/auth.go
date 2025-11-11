package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kuahbanyak/go-crud/internal/domain/services"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

var authService services.AuthService

func SetAuthService(service services.AuthService) {
	authService = service
}
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "Authorization header required", nil)
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "Invalid authorization header format", nil)
			return
		}
		token := tokenParts[1]
		if authService == nil {
			response.Error(w, http.StatusInternalServerError, "Auth service not initialized", nil)
			return
		}
		userID, role, err := authService.ValidateToken(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid token", err)
			return
		}
		ctx := context.WithValue(r.Context(), "id", userID)
		ctx = context.WithValue(ctx, "role", string(role))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("role").(string)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "User role not found in context", nil)
				return
			}
			hasRole := false
			for _, role := range roles {
				if userRole == role {
					hasRole = true
					break
				}
			}
			if !hasRole {
				response.Error(w, http.StatusForbidden, "Insufficient permissions", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
