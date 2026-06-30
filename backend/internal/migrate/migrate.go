package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Run(ctx context.Context, db *sql.DB, migrationsDir string) error {
	if err := ensureMigrationsTable(ctx, db); err != nil {
		return err
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		applied, err := isApplied(ctx, db, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		path := filepath.Join(migrationsDir, name)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", name, err)
		}

		if _, err := tx.ExecContext(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("exec migration %s: %w", name, err)
		}

		if _, err := tx.ExecContext(ctx,
			"INSERT INTO schema_migrations (filename) VALUES ($1)", name); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", name, err)
		}

		fmt.Printf("applied migration: %s\n", name)
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func isApplied(ctx context.Context, db *sql.DB, filename string) (bool, error) {
	var count int
	err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM schema_migrations WHERE filename = $1", filename).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
