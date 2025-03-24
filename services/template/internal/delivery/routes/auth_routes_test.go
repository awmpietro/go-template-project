package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/nuhorizon/go-project-template/services/template/internal/delivery/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthHandler to verify that each route is correctly wired
type MockAuthHandler struct {
	mock.Mock
}

func (m *MockAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func (m *MockAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusCreated)
}

func (m *MockAuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusAccepted)
}

func (m *MockAuthHandler) ExchangeToken(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusNoContent)
}

func TestRegisterAuthRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		expectCode int
		mockMethod string
	}{
		{"Login route", http.MethodPost, "/auth/login", http.StatusOK, "Login"},
		{"Register route", http.MethodPost, "/auth/register", http.StatusCreated, "Register"},
		{"ResetPassword route", http.MethodPost, "/auth/reset-password", http.StatusAccepted, "ResetPassword"},
		{"ExchangeToken route", http.MethodPost, "/auth/exchange-token", http.StatusNoContent, "ExchangeToken"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			mockHandler := new(MockAuthHandler)

			// Expect the correct method to be called
			mockHandler.On(tt.mockMethod, mock.Anything, mock.Anything).Once()

			// Register routes with the mock
			routes.RegisterAuthRoutes(r, mockHandler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectCode, rec.Code)
			mockHandler.AssertExpectations(t)
		})
	}
}
