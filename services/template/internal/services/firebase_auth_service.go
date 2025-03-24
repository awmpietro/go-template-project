package services

import (
	"context"
	"errors"

	"github.com/nuhorizon/go-project-template/services/template/internal/ports/services"
)

type firebaseAuthService struct {
	client services.FirebaseClient
}

func NewFirebaseAuthService(client services.FirebaseClient) services.FirebaseAuthService {
	return &firebaseAuthService{client: client}
}

func (f *firebaseAuthService) VerifyToken(ctx context.Context, firebaseToken string) (*services.FirebaseUser, error) {
	token, err := f.client.VerifyIDToken(ctx, firebaseToken)
	if err != nil {
		return nil, err
	}

	// Optionally fetch user details
	userRecord, err := f.client.GetUser(ctx, token.UID)
	if err != nil {
		return nil, err
	}

	return &services.FirebaseUser{
		UID:     userRecord.UID,
		Email:   userRecord.Email,
		Name:    userRecord.DisplayName,
		Picture: userRecord.PhotoURL,
	}, nil
}

func (f *firebaseAuthService) SendPasswordReset(ctx context.Context, email string) error {
	link, err := f.client.PasswordResetLinkWithSettings(ctx, email, nil)
	if err != nil {
		return err
	}
	// Aqui você pode enviar o link por e-mail ou logar no sistema
	if link == "" {
		return errors.New("failed to generate password reset link")
	}
	// Exemplo de log (ou enviar e-mail real)
	// log.Println("Password reset link:", link)
	return nil
}
