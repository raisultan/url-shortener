package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/raisultan/url-shortener/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/alias-gen/internal/config"
	"golang.org/x/exp/slog"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Postgres) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createTableIfDoesNotExistStmt := `
        CREATE TABLE IF NOT EXISTS counter (
            id SERIAL PRIMARY KEY,
            value BIGINT NOT NULL DEFAULT 1
        );
    `
	_, err = db.Exec(createTableIfDoesNotExistStmt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Ensure counter row exists
	insertRowStmt := `
        INSERT INTO counter (id, value) 
        VALUES (1, 1) 
        ON CONFLICT (id) 
        DO NOTHING;
    `
	_, err = db.Exec(insertRowStmt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) Close(log *slog.Logger) {
	err := s.db.Close()
	if err != nil {
		log.Error("could not close storage", sl.Err(err))
	}
}

func (s *Storage) IncrementCounter() (int64, error) {
	const op = "storage.postgres.IncrementCounter"

	var count int64
	incrementStmt := `
        UPDATE counter 
        SET value = value + 1 
        WHERE id = 1 
        RETURNING value;
    `
	err := s.db.QueryRow(incrementStmt).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
