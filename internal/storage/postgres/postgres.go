package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"URL-shortener/internal/storage"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "postgres"
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: ping db: %w", op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url (
			id BIGSERIAL PRIMARY KEY,
			alias TEXT UNIQUE NOT NULL,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: create schema: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urltoSave string, alias string) (int64, error) {
	const op = "postgres.SaveURL"

	const q = `INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id`

	var id int64
	err := s.db.QueryRow(q, urltoSave, alias).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, storage.ErrURLExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "postgres.GetURL"
	const q = `SELECT url FROM url WHERE alias = $1`

	var res string
	err := s.db.QueryRow(q, alias).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "postgres.DeleteURL"
	const q = `DELETE FROM url WHERE alias = $1`

	res, err := s.db.Exec(q, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affected: %w", op, err)
	}
	if affected == 0 {
		return storage.ErrURLNotFound
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
