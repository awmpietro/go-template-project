package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nuhorizon/go-project-template/services/template/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuthMiddleware(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("SWAGGER_USER_AUTH", "testuser")
	os.Setenv("SWAGGER_PASSWORD_AUTH", "testpass")

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Authorized",
			authHeader:     "Basic dGVzdHVzZXI6dGVzdHBhc3M=", // base64 of testuser:testpass
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized - wrong credentials",
			authHeader:     "Basic d3Jvbmd1c2VyOndyb25ncGFzcw==", // base64 of wronguser:wrongpass
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Unauthorized - missing header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false

			// Dummy handler to verify if it was called
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			middleware := middlewares.BasicAuthMiddleware(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.True(t, handlerCalled, "Expected next handler to be called")
			} else {
				assert.False(t, handlerCalled, "Expected next handler NOT to be called")
			}
		})
	}
}
