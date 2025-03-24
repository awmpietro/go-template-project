package usecases

import (
	"context"

	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/repositories"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/services"

	"github.com/google/uuid"
)

// AuthUseCase defines the behavior for authentication flows
type AuthUseCase interface {
	LoginOrRegister(ctx context.Context, firebaseToken string) (*domain.User, string, error)
	ResetPassword(ctx context.Context, email string) error
}

type authUseCase struct {
	userRepo     repositories.UserRepository
	firebaseAuth services.FirebaseAuthService
	jwtService   services.JWTService
}

func NewAuthUseCase(
	userRepo repositories.UserRepository,
	firebaseAuth services.FirebaseAuthService,
	jwtService services.JWTService,
) AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		firebaseAuth: firebaseAuth,
		jwtService:   jwtService,
	}
}

// LoginOrRegister validates Firebase token, syncs user, and generates app JWT
func (a *authUseCase) LoginOrRegister(ctx context.Context, firebaseToken string) (*domain.User, string, error) {
	// Step 1: Validate Firebase Token and extract claims
	fbUser, err := a.firebaseAuth.VerifyToken(ctx, firebaseToken)
	if err != nil {
		return nil, "", err
	}

	// Step 2: Upsert user based on Firebase UID
	user := &domain.User{
		ID:          fbUser.AppUserID, // Can be generated if empty
		FirebaseUID: fbUser.UID,
		Email:       fbUser.Email,
		Name:        fbUser.Name,
		PictureURL:  fbUser.Picture,
	}

	if user.ID == "" {
		user.ID = uuid.NewString()
	}

	dbUser, err := a.userRepo.UpsertByFirebaseUID(ctx, user)
	if err != nil {
		return nil, "", err
	}

	// Step 3: Issue App JWT
	token, err := a.jwtService.GenerateToken(dbUser)
	if err != nil {
		return nil, "", err
	}

	return dbUser, token, nil
}

// ResetPassword triggers Firebase password reset
func (a *authUseCase) ResetPassword(ctx context.Context, email string) error {
	return a.firebaseAuth.SendPasswordReset(ctx, email)
}
