package postgres

import (
	"database/sql"
	"fmt"

	"github.com/chishkin-afk/todo/internal/common/config"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d sslmode=%s user=%s password=%s dbname=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.SSLMode,
		cfg.Postgres.Auth.User,
		cfg.Postgres.Auth.Password,
		cfg.Postgres.Auth.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection with postgres: %w", err)
	}

	db.SetMaxIdleConns(cfg.Postgres.Conn.MaxIdles)
	db.SetMaxOpenConns(cfg.Postgres.Conn.MaxOpens)
	db.SetConnMaxIdleTime(cfg.Postgres.Conn.MaxIdleTime)
	db.SetConnMaxLifetime(cfg.Postgres.Conn.MaxLifetime)

	return db, nil
}

func Close(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close connection with postgres: %w", err)
	}

	return nil
}
