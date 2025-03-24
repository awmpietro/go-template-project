package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
)

type Pgsql struct {
	DB  *sql.DB
	Dsn string
}

var SQLOpen = sql.Open

func NewPGSql() *Pgsql {
	requiredEnvs := []string{"PG_USER", "PG_PASSWORD", "PG_HOST", "PG_DATABASE", "PG_PORT", "PG_SSL_MODE"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("[Postgres] Missing required environment variable: %s", env)
		}
	}
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	host := os.Getenv("PG_HOST")
	dbName := os.Getenv("PG_DATABASE")
	port := os.Getenv("PG_PORT")
	sslMode := os.Getenv("PG_SSL_MODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=America/Sao_Paulo",

		user, password, host, port, dbName, sslMode)

	return &Pgsql{Dsn: dsn}
}

// InitDB uses a default backoff config
func (d *Pgsql) InitDB() error {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 30 * time.Second
	return d.InitDBWithBackoff(expBackoff)
}

// InitDBWithBackoff allows injecting custom backoff (useful for tests)
func (d *Pgsql) InitDBWithBackoff(b backoff.BackOff) error {
	operation := func() error {
		db, err := SQLOpen("postgres", d.Dsn) // Use the overridable SQLOpen
		if err != nil {
			log.Println("[Postgres] Failed to open connection:", err)
			return err
		}

		// connection pool settings...
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err = db.Ping(); err != nil {
			log.Println("[Postgres] Connection ping failed:", err)
			return err
		}

		log.Println("[Postgres] Connected successfully")
		d.DB = db
		return nil
	}

	if err := backoff.Retry(operation, b); err != nil {
		log.Println("[Postgres] Failed to connect after retries:", err)
		return err
	}

	return nil
}

func (d *Pgsql) GetDB() *sql.DB {
	return d.DB
}

func (d *Pgsql) CloseDB() {
	if d.DB != nil {
		if err := d.DB.Close(); err != nil {
			log.Println("[Postgres] Error closing DB connection:", err)
		} else {
			log.Println("[Postgres] Disconnected successfully")
		}
	}
}

func (d *Pgsql) Stats() sql.DBStats {
	if d.DB != nil {
		return d.DB.Stats()
	}
	return sql.DBStats{}
}
