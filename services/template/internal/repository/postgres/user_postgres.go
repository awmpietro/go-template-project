package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/repositories"
)

type userPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) repositories.UserRepository {
	return &userPostgres{db: db}
}

func (r *userPostgres) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.FirebaseUID, user.Email, user.Name, user.PictureURL, user.PlanType, user.PremiumSince, user.PlanExpiry,
	)
	return err
}

func (r *userPostgres) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET email=$1, name=$2, picture_url=$3, plan_type=$4, premium_since=$5, plan_expiry=$6, updated_at=NOW() WHERE id=$7`
	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.Name, user.PictureURL, user.PlanType, user.PremiumSince, user.PlanExpiry, user.ID,
	)
	return err
}

func (r *userPostgres) UpsertByFirebaseUID(ctx context.Context, user *domain.User) (*domain.User, error) {
	existing, err := r.FindByFirebaseUID(ctx, user.FirebaseUID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing == nil {
		if err := r.Create(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}
	// Update existing user
	existing.Email = user.Email
	existing.Name = user.Name
	existing.PictureURL = user.PictureURL
	existing.PlanType = user.PlanType
	existing.PremiumSince = user.PremiumSince
	existing.PlanExpiry = user.PlanExpiry
	if err := r.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (r *userPostgres) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE id=$1`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanUser(row)
}

func (r *userPostgres) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	query := `SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE firebase_uid=$1`
	row := r.db.QueryRowContext(ctx, query, firebaseUID)
	return scanUser(row)
}

func (r *userPostgres) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE email=$1`
	row := r.db.QueryRowContext(ctx, query, email)
	return scanUser(row)
}

// Helper to scan SQL row into domain.User
func scanUser(row *sql.Row) (*domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.FirebaseUID,
		&user.Email,
		&user.Name,
		&user.PictureURL,
		&user.PlanType,
		&user.PremiumSince,
		&user.PlanExpiry,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
