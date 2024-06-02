package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ed16/messenger/internal/middleware"
	"github.com/stretchr/testify/assert"
)

// testHandler is a simple HTTP handler for testing purposes.
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

// TestLogRequests tests the LogRequests middleware.
func TestLogRequests(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a test handler wrapped with the LogRequests middleware
	handler := middleware.LogRequests(http.HandlerFunc(testHandler))

	// Create a logger with a buffer to capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr) // Restore original logger output

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	assert.Equal(t, "Hello, World!", rr.Body.String())

	// Check the log output
	logOutput := buf.String()
	expectedLog := "GET /test - Response status code: 200"
	assert.Contains(t, logOutput, expectedLog)
}
