package postgres_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/nuhorizon/go-project-template/services/template/internal/infra/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPgsql_GetDB(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	pg := &postgres.Pgsql{DB: db}
	assert.Equal(t, db, pg.GetDB())
}

func TestPgsql_CloseDB(t *testing.T) {
	t.Run("close success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectClose()

		pg := &postgres.Pgsql{DB: db}
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()

		pg.CloseDB()
	})

	t.Run("nil DB", func(t *testing.T) {
		pg := &postgres.Pgsql{DB: nil}
		pg.CloseDB() // Should not panic
	})
}

func TestPgsql_Stats(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	pg := &postgres.Pgsql{DB: db}
	stats := pg.Stats()
	assert.NotNil(t, stats)
}

func TestPgsql_InitDBWithBackoff(t *testing.T) {
	t.Run("successfully connects", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		originalSQLOpen := postgres.SQLOpen
		postgres.SQLOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return db, nil
		}
		defer func() { postgres.SQLOpen = originalSQLOpen }()

		pg := &postgres.Pgsql{Dsn: "ignored_dsn_for_mock"} // O dsn não importa mais
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 2 * time.Second

		err = pg.InitDBWithBackoff(b)
		assert.NoError(t, err)
		assert.NotNil(t, pg.DB)
	})

	t.Run("connection fails and retries until timeout", func(t *testing.T) {
		originalSQLOpen := postgres.SQLOpen
		postgres.SQLOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return nil, errors.New("connection failed")
		}
		defer func() { postgres.SQLOpen = originalSQLOpen }()

		pg := &postgres.Pgsql{Dsn: "ignored_dsn_for_mock"}
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 2 * time.Second

		err := pg.InitDBWithBackoff(b)
		assert.Error(t, err)
	})
}
