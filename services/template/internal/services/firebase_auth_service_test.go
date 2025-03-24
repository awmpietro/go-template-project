package services_test

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/v4/auth"
	portservices "github.com/nuhorizon/go-project-template/services/template/internal/ports/services"
	internalservices "github.com/nuhorizon/go-project-template/services/template/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFirebaseClient mocks *auth.Client
type MockFirebaseClient struct {
	mock.Mock
}

func (m *MockFirebaseClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	args := m.Called(ctx, idToken)
	token, _ := args.Get(0).(*auth.Token)
	return token, args.Error(1)
}

func (m *MockFirebaseClient) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	args := m.Called(ctx, uid)
	user, _ := args.Get(0).(*auth.UserRecord)
	return user, args.Error(1)
}

func (m *MockFirebaseClient) PasswordResetLinkWithSettings(ctx context.Context, email string, settings *auth.ActionCodeSettings) (string, error) {
	args := m.Called(ctx, email, settings)
	return args.String(0), args.Error(1)
}

func TestFirebaseAuthService_VerifyToken(t *testing.T) {
	tests := []struct {
		name           string
		mockVerifyErr  error
		mockUserErr    error
		expectedErr    bool
		expectedResult *portservices.FirebaseUser
	}{
		{
			name: "success",
			expectedResult: &portservices.FirebaseUser{
				UID:     "user123",
				Email:   "test@example.com",
				Name:    "Test User",
				Picture: "http://pic.url",
			},
		},
		{
			name:          "invalid token",
			mockVerifyErr: errors.New("invalid token"),
			expectedErr:   true,
		},
		{
			name:        "get user failure",
			mockUserErr: errors.New("user not found"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockFirebaseClient)
			ctx := context.Background()

			if tt.mockVerifyErr == nil {
				mockClient.On("VerifyIDToken", ctx, "valid-token").
					Return(&auth.Token{UID: "user123"}, nil)
			} else {
				mockClient.On("VerifyIDToken", ctx, "valid-token").
					Return((*auth.Token)(nil), tt.mockVerifyErr)
			}

			if tt.mockVerifyErr == nil {
				if tt.mockUserErr == nil {
					mockClient.On("GetUser", ctx, "user123").Return(&auth.UserRecord{
						UserInfo: &auth.UserInfo{
							UID:         "user123",
							Email:       "test@example.com",
							DisplayName: "Test User",
							PhotoURL:    "http://pic.url",
						},
					}, nil)
				} else {
					mockClient.On("GetUser", ctx, "user123").Return((*auth.UserRecord)(nil), tt.mockUserErr)
				}
			}

			service := internalservices.NewFirebaseAuthService(mockClient)

			user, err := service.VerifyToken(ctx, "valid-token")
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, user)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestFirebaseAuthService_SendPasswordReset(t *testing.T) {
	tests := []struct {
		name        string
		mockLink    string
		mockError   error
		expectedErr bool
	}{
		{
			name:     "success",
			mockLink: "https://reset.link",
		},
		{
			name:        "firebase error",
			mockError:   errors.New("firebase error"),
			expectedErr: true,
		},
		{
			name:        "empty link",
			mockLink:    "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockFirebaseClient)
			ctx := context.Background()

			mockClient.On("PasswordResetLinkWithSettings", ctx, "test@example.com", (*auth.ActionCodeSettings)(nil)).
				Return(tt.mockLink, tt.mockError)

			service := internalservices.NewFirebaseAuthService(mockClient)

			err := service.SendPasswordReset(ctx, "test@example.com")
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}
