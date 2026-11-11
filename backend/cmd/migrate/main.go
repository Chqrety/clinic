package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"clinic-api/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Main CLI tool for managing PostgreSQL database migrations and Chart of Accounts seeders.
func main() {
	action := flag.String("action", "up", "Migration action: up, down, status, seed-coa")
	flag.Parse()

	cfg := config.LoadConfig()
	db, err := sql.Open("pgx", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("[FATAL] Failed to open PostgreSQL connection: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("[FATAL] Failed to ping PostgreSQL on %s:%d: %v", cfg.Database.Host, cfg.Database.Port, err)
	}

	log.Printf("[MIGRATE] Connected to database %s on %s:%d", cfg.Database.DBName, cfg.Database.Host, cfg.Database.Port)
	log.Printf("[MIGRATE] Running migration command with action: %s", *action)

	// Ensure migrations table exists
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`); err != nil {
		log.Fatalf("[FATAL] Failed to ensure schema_migrations table: %v", err)
	}

	switch *action {
	case "up":
		if err := applyMigrations(ctx, db); err != nil {
			log.Fatalf("[FATAL] Migration failed: %v", err)
		}
	case "down":
		if err := rollbackMigration(ctx, db); err != nil {
			log.Fatalf("[FATAL] Rollback failed: %v", err)
		}
	case "seed-coa":
		if err := seedCOA(ctx, db); err != nil {
			log.Fatalf("[FATAL] Seeding failed: %v", err)
		}
	default:
		log.Fatalf("[MIGRATE] Unknown action: %s", *action)
	}

	log.Println("[MIGRATE] Migration task finished successfully.")
}

func applyMigrations(ctx context.Context, db *sql.DB) error {
	entries, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("reading migrations dir: %w", err)
	}

	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}
	sort.Strings(upFiles)

	for _, file := range upFiles {
		version := strings.TrimSuffix(file, ".up.sql")

		var exists bool
		_ = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
		if exists {
			log.Printf("[MIGRATE] %s already applied, skipping", file)
			continue
		}

		log.Printf("[MIGRATE] Applying %s ...", file)
		content, err := os.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", file, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("starting tx for %s: %w", file, err)
		}

		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("executing %s: %w", file, err)
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("recording %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing %s: %w", file, err)
		}
		log.Printf("[MIGRATE] Successfully applied %s", file)
	}

	return nil
}

func rollbackMigration(ctx context.Context, db *sql.DB) error {
	var lastVersion string
	err := db.QueryRowContext(ctx, "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&lastVersion)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[MIGRATE] No migrations to roll back")
			return nil
		}
		return fmt.Errorf("getting last migration: %w", err)
	}

	downFile := lastVersion + ".down.sql"
	log.Printf("[MIGRATE] Rolling back %s ...", downFile)
	content, err := os.ReadFile(filepath.Join("migrations", downFile))
	if err != nil {
		return fmt.Errorf("reading down migration %s: %w", downFile, err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting tx: %w", err)
	}

	if _, err := tx.ExecContext(ctx, string(content)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("executing %s: %w", downFile, err)
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM schema_migrations WHERE version = $1", lastVersion); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("removing version %s: %w", lastVersion, err)
	}

	return tx.Commit()
}

func seedCOA(ctx context.Context, db *sql.DB) error {
	seedFiles := []string{
		"seeds/branches.sql",
		"seeds/chart_of_accounts.sql",
	}

	for _, file := range seedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		log.Printf("[MIGRATE] Seeding from %s ...", file)
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading seed %s: %w", file, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("starting tx: %w", err)
		}

		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("executing seed %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing seed %s: %w", file, err)
		}
		log.Printf("[MIGRATE] Successfully seeded %s", file)
	}

	return nil
}
