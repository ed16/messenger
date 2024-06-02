package middleware

import (
	"log"
	"net/http"
)

// Custom ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// Middleware to log each request and its response status code
func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the request and record the response status code
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		// Log the request and response status code
		log.Printf("%s %s - Response status code: %d", r.Method, r.URL.Path, rw.status)
	})
}

// Override WriteHeader to capture the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
