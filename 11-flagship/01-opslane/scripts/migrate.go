//go:build ignore

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Migration runner for Opslane
// Usage: go run ./11-flagship/01-opslane/scripts/migrate.go [up|down|redo|status]

package main

import (
	"crypto/sha256"
	"context"
	"database/sql"
	"encoding/hex"
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
	checksum string
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
	if err := validateAppliedChecksums(ctx, db, migs); err != nil {
		log.Fatalf("Migration checksum validation failed: %v", err)
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
			checksum TEXT NOT NULL DEFAULT '',
			dirty BOOLEAN NOT NULL DEFAULT FALSE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	_, _ = db.ExecContext(ctx, `ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS checksum TEXT NOT NULL DEFAULT ''`)
	_, _ = db.ExecContext(ctx, `ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS dirty BOOLEAN NOT NULL DEFAULT FALSE`)
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
			checksum: "",
		})
	}

	sort.Slice(migs, func(i, j int) bool {
		return migs[i].version < migs[j].version
	})

	for i := range migs {
		sum, err := computeFileChecksum(migs[i].upFile)
		if err != nil {
			return nil, fmt.Errorf("checksum failed for %s: %w", migs[i].upFile, err)
		}
		migs[i].checksum = sum
	}
	if err := validateStrictOrder(migs); err != nil {
		return nil, err
	}

	return migs, nil
}

func getAppliedMigrations(ctx context.Context, db *sql.DB) ([]int, error) {
	var dirtyCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE dirty = TRUE").Scan(&dirtyCount); err == nil && dirtyCount > 0 {
		return nil, fmt.Errorf("database has dirty migration state; resolve before continuing")
	}

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

		_, err = tx.ExecContext(ctx, "INSERT INTO schema_migrations (version, name, checksum, dirty) VALUES ($1, $2, $3, FALSE)", m.version, m.name, m.checksum)
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

func computeFileChecksum(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:]), nil
}

func validateStrictOrder(migs []migration) error {
	if len(migs) == 0 {
		return nil
	}
	for i := 1; i < len(migs); i++ {
		if migs[i].version != migs[i-1].version+1 {
			return fmt.Errorf("migration order gap: expected version %d but found %d", migs[i-1].version+1, migs[i].version)
		}
	}
	return nil
}

func validateAppliedChecksums(ctx context.Context, db *sql.DB, migs []migration) error {
	rows, err := db.QueryContext(ctx, "SELECT version, checksum FROM schema_migrations ORDER BY version")
	if err != nil {
		return err
	}
	defer rows.Close()

	byVersion := make(map[int]string, len(migs))
	for _, m := range migs {
		byVersion[m.version] = m.checksum
	}

	for rows.Next() {
		var version int
		var checksum string
		if err := rows.Scan(&version, &checksum); err != nil {
			return err
		}
		expected, ok := byVersion[version]
		if !ok {
			return fmt.Errorf("applied migration version %d has no corresponding migration file", version)
		}
		if checksum != "" && checksum != expected {
			return fmt.Errorf("checksum mismatch for version %d", version)
		}
	}
	return rows.Err()
}
