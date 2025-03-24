package infrastructure

import "database/sql"

type SQLConnector interface {
	InitDB() error
	GetDB() *sql.DB
	CloseDB()
	Stats() sql.DBStats
}
