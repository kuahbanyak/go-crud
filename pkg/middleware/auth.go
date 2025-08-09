package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(secret []byte) gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth"})
            return
        }
        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
            return
        }
        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return secret, nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        claims := token.Claims.(jwt.MapClaims)
        // normalize map[string]interface{} values for convenience
        m := map[string]interface{}{}
        for k, v := range claims {
            m[k] = v
        }
        c.Set("claims", m)
        c.Next()
    }
}
