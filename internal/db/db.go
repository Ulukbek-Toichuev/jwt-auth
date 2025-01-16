package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB(driverName, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}

	m, err := migrate.New(
		"file:///Users/ulukbek_toichuev/Documents/jwt-auth/internal/db/migrations",
		"sqlite3://"+dataSourceName,
	)

	if err != nil {
		log.Printf("Failed to initialized migration: %v", err)
		os.Exit(1)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Failed to initialized migration: %v", err)
		os.Exit(1)
	}

	return db
}
