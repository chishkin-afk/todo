package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func MigrateUP(ctx context.Context, db *sql.DB) error {
	migrations, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read dir with migrations: %w", err)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name() < migrations[j].Name()
	})

	for _, migration := range migrations {
		if filepath.Ext(migration.Name()) != ".sql" {
			continue
		}

		content, err := migrationFS.ReadFile("migrations/" + migration.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration's file: %w", err)
		}

		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	}

	return nil
}
