package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		// Get request ID from context
		requestID := GetRequestID(r.Context())

		ginMode := os.Getenv("GIN_MODE")
		isProduction := ginMode == "release" || os.Getenv("RAILWAY_ENVIRONMENT") != ""

		if isProduction {
			// Only log errors and slow requests in production
			if wrapped.statusCode >= 400 || duration > 1*time.Second {
				log.Printf("[%s] ERROR/SLOW: %s %s %d %v",
					requestID, r.Method, r.RequestURI, wrapped.statusCode, duration)
			}
		} else {
			// Log all requests in development with more details
			log.Printf("[%s] %s %s %d %v %s",
				requestID,
				r.Method,
				r.RequestURI,
				wrapped.statusCode,
				duration,
				r.RemoteAddr,
			)
		}
	})
}
