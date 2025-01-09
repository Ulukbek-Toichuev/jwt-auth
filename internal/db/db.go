package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const create_schemas = `
CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
	role TEXT NOT NULL,
    created_date DATETIME NOT NULL,
    deleted_date DATETIME
);`

func NewDB(driverName, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
	_, err = db.Exec(create_schemas)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
	return db
}
