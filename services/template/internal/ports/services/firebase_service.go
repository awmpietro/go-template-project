package services

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type AppInterface interface {
	Auth(ctx context.Context) (*auth.Client, error)
}

type FirebaseAuthService interface {
	VerifyToken(ctx context.Context, firebaseToken string) (*FirebaseUser, error)
	SendPasswordReset(ctx context.Context, email string) error
}

type FirebaseClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
	PasswordResetLinkWithSettings(ctx context.Context, email string, settings *auth.ActionCodeSettings) (string, error)
}

type FirebaseUser struct {
	UID       string
	Email     string
	Name      string
	Picture   string
	AppUserID string // Optional: for DB PK if you need to inject/generate
}
