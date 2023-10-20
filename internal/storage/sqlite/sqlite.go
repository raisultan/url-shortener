package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/raisultan/url-shortener/internal/config"
	"github.com/raisultan/url-shortener/internal/lib/logger/sl"
	"golang.org/x/exp/slog"

	"github.com/mattn/go-sqlite3"
	"github.com/raisultan/url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(config config.Storages, _ context.Context) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", config.SQLite.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Close(_ context.Context, log *slog.Logger) {
	err := s.db.Close()
	if err != nil {
		log.Error("could not close storage", sl.Err(err))
	}
}

func (s *Storage) SaveUrl(_ context.Context, urlToSave string, alias string) error {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		ok := errors.As(err, &sqliteErr)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
		}
	}

	return nil
}

func (s *Storage) GetUrl(_ context.Context, alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement %w", op, err)
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrUrlNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resUrl, nil
}
