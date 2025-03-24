package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	"github.com/nuhorizon/go-project-template/services/template/pkg/middlewares"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*domain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		authHeader       string
		mockUser         *domain.User
		mockError        error
		expectedStatus   int
		expectedUserIDIn bool
	}{
		{
			name:           "missing token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			mockUser:       nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:             "valid token",
			authHeader:       "Bearer valid-token",
			mockUser:         &domain.User{ID: "user-123"},
			mockError:        nil,
			expectedStatus:   http.StatusOK,
			expectedUserIDIn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWT := new(MockJWTService)
			if tt.authHeader != "" && tt.mockUser != nil || tt.mockError != nil {
				mockJWT.On("ValidateToken", strings.TrimPrefix(tt.authHeader, "Bearer ")).
					Return(tt.mockUser, tt.mockError)
			}

			// Target handler to check context propagation
			finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
				if tt.expectedUserIDIn {
					assert.True(t, ok)
					assert.Equal(t, "user-123", userID)
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()

			handler := middlewares.AuthMiddleware(mockJWT)(finalHandler)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockJWT.AssertExpectations(t)
		})
	}
}
