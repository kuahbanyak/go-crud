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

		// Reduce logging in production to prevent Railway rate limit
		ginMode := os.Getenv("GIN_MODE")
		isProduction := ginMode == "release" || os.Getenv("RAILWAY_ENVIRONMENT") != ""

		if isProduction {
			// Only log errors and slow requests in production
			if wrapped.statusCode >= 400 || duration > 1*time.Second {
				log.Printf("ERROR/SLOW: %s %s %d %v",
					r.Method, r.RequestURI, wrapped.statusCode, duration)
			}
		} else {
			// Full logging in development
			log.Printf("%s %s %d %v %s",
				r.Method,
				r.RequestURI,
				wrapped.statusCode,
				duration,
				r.RemoteAddr,
			)
		}
	})
}
