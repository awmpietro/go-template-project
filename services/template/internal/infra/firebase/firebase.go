package firebase

import (
	"context"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/services"
	"google.golang.org/api/option"
)

// FirebaseNewAppFunc é a função injetável (para facilitar testes)
var FirebaseNewAppFunc = func(ctx context.Context, config *firebase.Config, opts ...option.ClientOption) (services.AppInterface, error) {
	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// NewFirebaseClient inicializa e retorna o *auth.Client do Firebase
func NewFirebaseClient() (*auth.Client, error) {
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_JSON")
	if credentialsPath == "" {
		return nil, os.ErrNotExist
	}

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := FirebaseNewAppFunc(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}
