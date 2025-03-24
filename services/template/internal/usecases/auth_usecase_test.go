package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/services"
	"github.com/nuhorizon/go-project-template/services/template/internal/usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) UpsertByFirebaseUID(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

func (m *MockUserRepo) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	args := m.Called(ctx, firebaseUID)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

type MockFirebaseAuth struct{ mock.Mock }

func (m *MockFirebaseAuth) VerifyToken(ctx context.Context, firebaseToken string) (*services.FirebaseUser, error) {
	args := m.Called(ctx, firebaseToken)
	return args.Get(0).(*services.FirebaseUser), args.Error(1)
}

func (m *MockFirebaseAuth) SendPasswordReset(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

type MockJWTService struct{ mock.Mock }

func (m *MockJWTService) GenerateToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*domain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestAuthUseCase_LoginOrRegister(t *testing.T) {
	tests := []struct {
		name           string
		firebaseToken  string
		firebaseUser   *services.FirebaseUser
		verifyErr      error
		upsertResult   *domain.User
		upsertErr      error
		jwtToken       string
		jwtErr         error
		expectErr      bool
		expectedToken  string
		expectedUserID string
	}{
		{
			name: "success",
			firebaseUser: &services.FirebaseUser{
				UID:     "firebase-uid",
				Email:   "user@example.com",
				Name:    "Test User",
				Picture: "http://picture.url",
			},
			upsertResult:   &domain.User{ID: "user-id", FirebaseUID: "firebase-uid"},
			jwtToken:       "jwt-token",
			expectedToken:  "jwt-token",
			expectedUserID: "user-id",
		},
		{
			name:      "firebase verification fails",
			verifyErr: errors.New("invalid firebase token"),
			expectErr: true,
		},
		{
			name: "user repo fails",
			firebaseUser: &services.FirebaseUser{
				UID:   "firebase-uid",
				Email: "user@example.com",
			},
			upsertErr: errors.New("db error"),
			expectErr: true,
		},
		{
			name: "jwt generation fails",
			firebaseUser: &services.FirebaseUser{
				UID:   "firebase-uid",
				Email: "user@example.com",
			},
			upsertResult: &domain.User{ID: "user-id", FirebaseUID: "firebase-uid"},
			jwtErr:       errors.New("jwt fail"),
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockFirebase := new(MockFirebaseAuth)
			mockJWT := new(MockJWTService)
			authUC := usecases.NewAuthUseCase(mockRepo, mockFirebase, mockJWT)

			if tt.verifyErr != nil || tt.firebaseUser != nil {
				mockFirebase.On("VerifyToken", mock.Anything, tt.firebaseToken).
					Return(tt.firebaseUser, tt.verifyErr)
			}

			if tt.firebaseUser != nil {
				mockRepo.On("UpsertByFirebaseUID", mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(tt.upsertResult, tt.upsertErr)
			}

			if tt.upsertResult != nil {
				mockJWT.On("GenerateToken", tt.upsertResult).
					Return(tt.jwtToken, tt.jwtErr)
			}

			user, token, err := authUC.LoginOrRegister(context.Background(), tt.firebaseToken)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
				assert.Equal(t, tt.expectedUserID, user.ID)
			}
		})
	}
}

func TestAuthUseCase_ResetPassword(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		sendErr   error
		expectErr bool
	}{
		{
			name:  "success",
			email: "user@example.com",
		},
		{
			name:      "firebase fails",
			email:     "user@example.com",
			sendErr:   errors.New("firebase error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockFirebase := new(MockFirebaseAuth)
			mockJWT := new(MockJWTService)
			authUC := usecases.NewAuthUseCase(mockRepo, mockFirebase, mockJWT)

			mockFirebase.On("SendPasswordReset", mock.Anything, tt.email).
				Return(tt.sendErr)

			err := authUC.ResetPassword(context.Background(), tt.email)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
