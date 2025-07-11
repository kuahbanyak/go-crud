package middleware

import (
	"net/http"
	"strings"

	"go-crud/pkg/auth"
	"go-crud/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth provides JWT authentication middleware
func JWTAuth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required", "Missing authorization header")
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header", "Authorization header must be Bearer token")
			c.Abort()
			return
		}

		token := bearerToken[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token", err.Error())
			c.Abort()
			return
		}

		// Set user context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalAuth middleware that doesn't require authentication but extracts user info if present
func OptionalAuth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {
				token := bearerToken[1]
				claims, err := jwtService.ValidateToken(token)
				if err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("user_role", claims.Role)
				}
			}
		}
		c.Next()
	}
}

// AdminOnly middleware that requires admin role
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists || userRole != "admin" {
			response.Error(c, http.StatusForbidden, "Admin access required", "This endpoint requires admin privileges")
			c.Abort()
			return
		}
		c.Next()
	}
}
