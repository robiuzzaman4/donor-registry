package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	if err := ensureMigrationsTable(ctx, pool); err != nil {
		return err
	}

	files, err := readMigrationFiles(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		applied, err := isMigrationApplied(ctx, pool, file)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		if err := applyMigrationFile(ctx, pool, filepath.Join(migrationsDir, file), file); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			filename TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}
	return nil
}

func readMigrationFiles(migrationsDir string) ([]string, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations dir: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".sql") {
			files = append(files, name)
		}
	}

	sort.Strings(files)
	return files, nil
}

func isMigrationApplied(ctx context.Context, pool *pgxpool.Pool, filename string) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE filename = $1)`, filename).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check migration %s: %w", filename, err)
	}
	return exists, nil
}

func applyMigrationFile(ctx context.Context, pool *pgxpool.Pool, path string, filename string) error {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read migration %s: %w", filename, err)
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin migration tx for %s: %w", filename, err)
	}

	if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to apply migration %s: %w", filename, err)
	}

	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1)`, filename); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to record migration %s: %w", filename, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migration %s: %w", filename, err)
	}

	return nil
}
