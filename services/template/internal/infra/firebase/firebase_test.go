package firebase

import (
	"context"
	"errors"
	"os"
	"testing"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

// MockFirebaseApp mocks ports.AppInterface
type MockFirebaseApp struct {
	authClient *auth.Client
	authError  error
}

func (m *MockFirebaseApp) Auth(ctx context.Context) (*auth.Client, error) {
	return m.authClient, m.authError
}

func TestFirebaseNewAppFunc(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully creates Firebase App with valid credentials", func(t *testing.T) {
		credPath := "./test_firebase_credentials.json"
		_ = os.WriteFile(credPath, []byte(`{"type": "service_account"}`), 0644)
		defer os.Remove(credPath)

		opt := option.WithCredentialsFile(credPath)
		app, err := FirebaseNewAppFunc(ctx, nil, opt)

		assert.NoError(t, err)
		assert.NotNil(t, app)
	})

	t.Run("creates Firebase App even with invalid credential path (lazy error)", func(t *testing.T) {
		invalidOpt := option.WithCredentialsFile("/invalid/path/to/creds.json")
		app, err := FirebaseNewAppFunc(ctx, nil, invalidOpt)

		// Firebase does not check credential file at this stage
		assert.NoError(t, err)
		assert.NotNil(t, app)

		// Real error comes when using the app
		_, authErr := app.Auth(ctx)
		assert.Error(t, authErr, "Expected error when initializing Auth with invalid credentials")
	})
}

func TestNewFirebaseClient_Success(t *testing.T) {
	original := FirebaseNewAppFunc
	defer func() { FirebaseNewAppFunc = original }()

	mockAuthClient := &auth.Client{} // Use a real pointer to avoid nil issues
	mockApp := &MockFirebaseApp{authClient: mockAuthClient}

	FirebaseNewAppFunc = func(ctx context.Context, config *firebase.Config, opts ...option.ClientOption) (services.AppInterface, error) {
		return mockApp, nil
	}

	// Mock environment variable
	os.Setenv("FIREBASE_CREDENTIALS_JSON", "/fake/path/creds.json")
	defer os.Unsetenv("FIREBASE_CREDENTIALS_JSON")

	client, err := NewFirebaseClient()
	require.NoError(t, err)
	assert.Equal(t, mockAuthClient, client)
}

func TestNewFirebaseClient_MissingEnv(t *testing.T) {
	os.Unsetenv("FIREBASE_CREDENTIALS_JSON")
	client, err := NewFirebaseClient()
	assert.Nil(t, client)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestNewFirebaseClient_AppCreationFails(t *testing.T) {
	original := FirebaseNewAppFunc
	defer func() { FirebaseNewAppFunc = original }()

	FirebaseNewAppFunc = func(ctx context.Context, config *firebase.Config, opts ...option.ClientOption) (services.AppInterface, error) {
		return nil, errors.New("failed to create app")
	}

	os.Setenv("FIREBASE_CREDENTIALS_JSON", "/fake/path/creds.json")
	defer os.Unsetenv("FIREBASE_CREDENTIALS_JSON")

	client, err := NewFirebaseClient()
	assert.Nil(t, client)
	assert.EqualError(t, err, "failed to create app")
}

func TestNewFirebaseClient_AuthFails(t *testing.T) {
	original := FirebaseNewAppFunc
	defer func() { FirebaseNewAppFunc = original }()

	mockApp := &MockFirebaseApp{authError: errors.New("auth failed")}
	FirebaseNewAppFunc = func(ctx context.Context, config *firebase.Config, opts ...option.ClientOption) (services.AppInterface, error) {
		return mockApp, nil
	}

	os.Setenv("FIREBASE_CREDENTIALS_JSON", "/fake/path/creds.json")
	defer os.Unsetenv("FIREBASE_CREDENTIALS_JSON")

	client, err := NewFirebaseClient()
	assert.Nil(t, client)
	assert.EqualError(t, err, "auth failed")
}
