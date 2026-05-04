//go:build ignore

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Migration runner for Opslane
// Usage: go run ./11-flagship/01-opslane/scripts/migrate.go [up|down|redo|status]

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
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbDsn      = flag.String("dsn", "", "Database connection string (or use OPSLANE_DB_DSN)")
	migrations = flag.String("migrations", "./migrations", "Path to migrations directory")
	direction  = flag.String("direction", "up", "Migration direction: up, down, or redo")
)

type migration struct {
	version  int
	name     string
	upFile   string
	downFile string
}

func main() {
	flag.Parse()

	// Get DSN from flag or environment
	if *dbDsn == "" {
		*dbDsn = os.Getenv("OPSLANE_DB_DSN")
	}

	if *dbDsn == "" {
		fmt.Println("Error: Database DSN is required. Set OPSLANE_DB_DSN or use -dsn flag.")
		flag.Usage()
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", *dbDsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Ensure migrations table exists
	if err := ensureMigrationsTable(ctx, db); err != nil {
		log.Fatalf("Failed to ensure migrations table: %v", err)
	}

	// Get available migrations
	migs, err := loadMigrations(*migrations)
	if err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		log.Fatalf("Failed to get applied migrations: %v", err)
	}

	switch *direction {
	case "up":
		if err := runUpMigrations(ctx, db, migs, applied); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "down":
		if err := runDownMigration(ctx, db, migs, applied); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "redo":
		if err := redoLastMigration(ctx, db, migs, applied); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "status":
		printMigrationStatus(migs, applied)
	default:
		log.Fatalf("Unknown direction: %s", *direction)
	}
}

func ensureMigrationsTable(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := db.ExecContext(ctx, query)
	return err
}

func loadMigrations(path string) ([]migration, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var migs []migration
	seen := make(map[int]bool)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		// Parse version from filename (e.g., "001_create_tenants.up.sql")
		parts := strings.Split(name, "_")
		if len(parts) == 0 {
			continue
		}

		versionStr := strings.TrimPrefix(parts[0], "0")
		if versionStr == "" {
			versionStr = "0"
		}
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			continue
		}

		if seen[version] {
			continue
		}
		seen[version] = true

		// Find corresponding down migration
		downName := strings.Replace(name, ".up.sql", ".down.sql", 1)
		downPath := filepath.Join(path, downName)

		migs = append(migs, migration{
			version:  version,
			name:     strings.TrimSuffix(name, ".up.sql"),
			upFile:   filepath.Join(path, name),
			downFile: downPath,
		})
	}

	sort.Slice(migs, func(i, j int) bool {
		return migs[i].version < migs[j].version
	})

	return migs, nil
}

func getAppliedMigrations(ctx context.Context, db *sql.DB) ([]int, error) {
	rows, err := db.QueryContext(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	return versions, rows.Err()
}

func runUpMigrations(ctx context.Context, db *sql.DB, migs []migration, applied []int) error {
	appliedSet := make(map[int]bool)
	for _, v := range applied {
		appliedSet[v] = true
	}

	for _, m := range migs {
		if appliedSet[m.version] {
			continue
		}

		fmt.Printf("Applying migration %d: %s\n", m.version, m.name)

		content, err := os.ReadFile(m.upFile)
		if err != nil {
			return fmt.Errorf("failed to read migration %d: %w", m.version, err)
		}

		// Execute migration in transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		_, err = tx.ExecContext(ctx, string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %d: %w", m.version, err)
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", m.version, m.name)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", m.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.version, err)
		}

		fmt.Printf("  ✓ Migration %d applied\n", m.version)
	}

	fmt.Println("\n✓ All migrations complete!")
	return nil
}

func runDownMigration(ctx context.Context, db *sql.DB, migs []migration, applied []int) error {
	if len(applied) == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Find the last applied migration
	lastVersion := applied[len(applied)-1]
	var lastMig migration
	for _, m := range migs {
		if m.version == lastVersion {
			lastMig = m
			break
		}
	}

	if lastMig.downFile == "" {
		return fmt.Errorf("no down migration for version %d", lastVersion)
	}

	fmt.Printf("Rolling back migration %d: %s\n", lastVersion, lastMig.name)

	content, err := os.ReadFile(lastMig.downFile)
	if err != nil {
		return fmt.Errorf("failed to read down migration: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, string(content))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM schema_migrations WHERE version = $1", lastVersion)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	fmt.Printf("  ✓ Migration %d rolled back\n", lastVersion)
	return nil
}

func redoLastMigration(ctx context.Context, db *sql.DB, migs []migration, applied []int) error {
	// Get current applied to know what to redo
	oldApplied := make([]int, len(applied))
	copy(oldApplied, applied)

	if err := runDownMigration(ctx, db, migs, applied); err != nil {
		return err
	}

	// Re-apply migrations up to the last one that was applied
	// Get fresh applied list (should be one less now)
	applied, _ = getAppliedMigrations(ctx, db)
	return runUpMigrations(ctx, db, migs, applied)
}

func printMigrationStatus(migs []migration, applied []int) {
	appliedSet := make(map[int]bool)
	for _, v := range applied {
		appliedSet[v] = true
	}

	fmt.Println("Migration Status:")
	fmt.Println(strings.Repeat("-", 50))

	if len(migs) == 0 {
		fmt.Println("  No migration files found")
		return
	}

	for _, m := range migs {
		status := "pending"
		if appliedSet[m.version] {
			status = "applied"
		}
		fmt.Printf("  %d %-40s [%s]\n", m.version, m.name, status)
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Total: %d migrations, %d applied, %d pending\n",
		len(migs), len(applied), len(migs)-len(applied))
}
