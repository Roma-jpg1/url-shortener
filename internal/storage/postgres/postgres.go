package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "postgres"
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS url (
    			id INTEGER PRIMARY KEY,
    alias TEXT UNIQUE NOT NULL UNIQUE,
    url TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
`)
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %s", op, err)
	}

	return &Storage{db: db}, nil

}
