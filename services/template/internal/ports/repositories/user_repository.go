package repositories

import (
	"context"

	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
)

// UserRepository defines the methods the use case expects from any User repository implementation (DB, mock, etc.)
type UserRepository interface {
	// Create a new user in the database
	Create(ctx context.Context, user *domain.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// UpsertByFirebaseUID creates or updates a user based on Firebase UID (used in login/sync)
	UpsertByFirebaseUID(ctx context.Context, user *domain.User) (*domain.User, error)

	// FindByID retrieves a user by internal system UUID
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// FindByFirebaseUID retrieves a user by Firebase UID
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)

	// FindByEmail retrieves a user by email (useful if supporting direct login)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
