package postgres_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	"github.com/nuhorizon/go-project-template/services/template/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
)

func TestUserPostgres_Create(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
					INSERT INTO users (id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
				`)).
					WithArgs("user-id", "firebase-uid", "test@example.com", "Test User", "", "", nil, nil).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
					INSERT INTO users (id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
				`)).
					WithArgs("user-id", "firebase-uid", "test@example.com", "Test User", "", "", nil, nil).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.setupMock(mock)

			repo := postgres.NewUserPostgres(db)
			err := repo.Create(context.Background(), &domain.User{
				ID:          "user-id",
				FirebaseUID: "firebase-uid",
				Email:       "test@example.com",
				Name:        "Test User",
			})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPostgres_Update(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
					UPDATE users SET email=$1, name=$2, picture_url=$3, plan_type=$4, premium_since=$5, plan_expiry=$6, updated_at=NOW() WHERE id=$7
				`)).
					WithArgs("test@example.com", "Test User", "", "", nil, nil, "user-id").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
					UPDATE users SET email=$1, name=$2, picture_url=$3, plan_type=$4, premium_since=$5, plan_expiry=$6, updated_at=NOW() WHERE id=$7
				`)).
					WithArgs("test@example.com", "Test User", "", "", nil, nil, "user-id").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.setupMock(mock)

			repo := postgres.NewUserPostgres(db)
			err := repo.Update(context.Background(), &domain.User{
				ID:    "user-id",
				Email: "test@example.com",
				Name:  "Test User",
			})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPostgres_FindByID(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE id=$1
				`)).
					WithArgs("user-id").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "firebase_uid", "email", "name", "picture_url", "plan_type", "premium_since", "plan_expiry", "created_at", "updated_at",
					}).AddRow("user-id", "firebase-uid", "test@example.com", "Test User", "", "free", now, now, now, now))
			},
		},
		{
			name: "db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE id=$1
				`)).
					WithArgs("user-id").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.setupMock(mock)

			repo := postgres.NewUserPostgres(db)
			_, err := repo.FindByID(context.Background(), "user-id")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPostgres_FindByFirebaseUID(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE firebase_uid=$1
				`)).
					WithArgs("firebase-uid").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "firebase_uid", "email", "name", "picture_url", "plan_type", "premium_since", "plan_expiry", "created_at", "updated_at",
					}).AddRow("user-id", "firebase-uid", "test@example.com", "Test User", "", "free", now, now, now, now))
			},
		},
		{
			name: "db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE firebase_uid=$1
				`)).
					WithArgs("firebase-uid").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.setupMock(mock)

			repo := postgres.NewUserPostgres(db)
			_, err := repo.FindByFirebaseUID(context.Background(), "firebase-uid")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPostgres_FindByEmail(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE email=$1
				`)).
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "firebase_uid", "email", "name", "picture_url", "plan_type", "premium_since", "plan_expiry", "created_at", "updated_at",
					}).AddRow("user-id", "firebase-uid", "test@example.com", "Test User", "", "free", now, now, now, now))
			},
		},
		{
			name: "db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, firebase_uid, email, name, picture_url, plan_type, premium_since, plan_expiry, created_at, updated_at FROM users WHERE email=$1
				`)).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.setupMock(mock)

			repo := postgres.NewUserPostgres(db)
			_, err := repo.FindByEmail(context.Background(), "test@example.com")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
