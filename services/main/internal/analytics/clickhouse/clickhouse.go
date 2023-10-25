package clickhouse

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/raisultan/url-shortener/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/main/internal/analytics"
	"github.com/raisultan/url-shortener/services/main/internal/config"
	"golang.org/x/exp/slog"
)

type AnalyticsTracker struct {
	db *sql.DB
}

func NewClickHouseAnalyticsTracker(cfg config.ClickHouse) (*AnalyticsTracker, error) {
	conn, err := sql.Open("clickhouse", cfg.Dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	createTableIfNotExistsQuery := `
		CREATE TABLE IF NOT EXISTS testing.clicks (
			url_alias String,
			timestamp DateTime,
			user_agent String,
			ip String,
			referrer String,
			latency UInt64
		) ENGINE = MergeTree()
		ORDER BY timestamp
	`
	if _, err := conn.Exec(createTableIfNotExistsQuery); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &AnalyticsTracker{db: conn}, nil
}

func (tracker *AnalyticsTracker) Close(log *slog.Logger) {
	err := tracker.db.Close()
	if err != nil {
		log.Error("could not close storage", sl.Err(err))
	}
}

func (tracker *AnalyticsTracker) TrackClickEvent(event analytics.ClickEvent) error {
	query := `
		INSERT INTO testing.clicks (
			url_alias,
			timestamp,
			user_agent,
			ip,
			referrer,
			latency
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := tracker.db.Exec(
		query,
		event.URLAlias,
		event.Timestamp,
		event.UserAgent,
		event.IP,
		event.Referrer,
		event.Latency.Milliseconds(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert click event: %w", err)
	}

	return nil
}
