package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nuhorizon/go-project-template/services/template/internal/delivery/handlers"
	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUseCase fully mocked
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) LoginOrRegister(ctx context.Context, firebaseToken string) (*domain.User, string, error) {
	args := m.Called(ctx, firebaseToken)
	user := args.Get(0)
	if user == nil {
		return nil, "", args.Error(2)
	}
	return user.(*domain.User), args.String(1), args.Error(2)
}

func (m *MockAuthUseCase) ResetPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        string
		mockUser       *domain.User
		mockToken      string
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "success",
			reqBody:        `{"firebase_token": "valid-firebase-token"}`,
			mockUser:       &domain.User{ID: "user-id"},
			mockToken:      "jwt-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json",
			reqBody:        `invalid-json`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "usecase error",
			reqBody:        `{"firebase_token": "invalid-token"}`,
			mockErr:        errors.New("firebase error"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(MockAuthUseCase)
			handler := handlers.NewAuthHandler(mockUC)

			if tt.mockUser != nil || tt.mockErr != nil {
				mockUC.On("LoginOrRegister", mock.Anything, mock.Anything).
					Return(tt.mockUser, tt.mockToken, tt.mockErr)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte(tt.reqBody)))
			rec := httptest.NewRecorder()

			handler.Login(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        string
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "success",
			reqBody:        `{"email": "user@example.com"}`,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid json",
			reqBody:        `invalid-json`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "usecase failure",
			reqBody:        `{"email": "fail@example.com"}`,
			mockErr:        errors.New("firebase error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(MockAuthUseCase)
			handler := handlers.NewAuthHandler(mockUC)

			if tt.mockErr != nil || tt.name == "success" {
				mockUC.On("ResetPassword", mock.Anything, mock.Anything).
					Return(tt.mockErr)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/reset-password", bytes.NewBuffer([]byte(tt.reqBody)))
			rec := httptest.NewRecorder()

			handler.ResetPassword(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
