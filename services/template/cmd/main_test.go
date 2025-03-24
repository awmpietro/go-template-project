package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/services"
	middlewares "github.com/nuhorizon/go-project-template/services/template/pkg/middlewares"
)

type MockSQLConnector struct {
	mock.Mock
}

func (m *MockSQLConnector) InitDB() error { return nil }
func (m *MockSQLConnector) GetDB() *sql.DB {
	db, _, _ := sqlmock.New()
	return db
}
func (m *MockSQLConnector) CloseDB()           {}
func (m *MockSQLConnector) Stats() sql.DBStats { return sql.DBStats{} }

type MockFirebaseAuthService struct {
	mock.Mock
}

func (m *MockFirebaseAuthService) VerifyToken(ctx context.Context, token string) (*services.FirebaseUser, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*services.FirebaseUser), args.Error(1)
}

func (m *MockFirebaseAuthService) SendPasswordReset(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

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

func TestSetupHandlersAndRoutes(t *testing.T) {
	tests := []struct {
		method string
		route  string
	}{
		{method: http.MethodPost, route: "/auth/login"},
		{method: http.MethodPost, route: "/auth/register"},
		{method: http.MethodPost, route: "/auth/reset-password"},
		{method: http.MethodGet, route: "/cats"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.route, func(t *testing.T) {
			db, _, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			dbMock := &MockSQLConnector{}
			dbMock.On("GetDB").Return(db)

			firebaseMock := new(MockFirebaseAuthService)
			jwtMock := new(MockJWTService)

			mux := chi.NewMux()
			setupHandlers(mux, dbMock, firebaseMock, jwtMock)

			req := httptest.NewRequest(tt.method, tt.route, nil)
			// 👇 Add userID if route needs auth context (like /cats)
			if tt.route == "/cats" {
				ctx := context.WithValue(req.Context(), middlewares.UserIDKey, "test-user-id")
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			assert.NotEqual(t, http.StatusNotFound, rec.Code)
		})
	}

}

func TestInitializeMux(t *testing.T) {
	mux := chi.NewMux()
	initializeMux(mux)
	assert.NotNil(t, mux)
}
