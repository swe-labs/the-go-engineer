// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Database Migrations
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to manage database schema changes over time.
//   - How to embed SQL migration files into your Go binary using 'embed'.
//   - How to use 'golang-migrate' to automate schema updates.
//
// WHY THIS MATTERS:
//   - In a real project, the database schema changes constantly.
//     Manually running SQL scripts is dangerous and error-prone.
//     Migrations provide a version-controlled, automated way to keep
//     your database in sync with your code.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/database-migrations/1-embedded-migrations
//
// KEY TAKEAWAY:
//   - Treat your database schema like code: version it, automate it, and embed it.
// ============================================================================

package main

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	// Drivers
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Stage 06: Database Migrations - Embedded Schema Evolution
//
//   - //go:embed: Bundling SQL files into the binary
//   - UP/DOWN Migrations: Version-controlled evolution
//   - Automatic Boot Sequence: Running migrations on startup
//
// ENGINEERING DEPTH:
//   By embedding migrations, your Go binary becomes "Self-Contained".
//   When the binary starts in production, it can check the current
//   database version and apply only the necessary updates before
//   starting the web server. This ensures that your code and
//   your database schema are always in perfect sync.

// --- 1. Embed the migrations folder ---
// This tells Go to include all .sql files in the 'migrations' directory
// inside the compiled binary.
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunEmbeddedMigrations applies pending SQL changes to the database.
func RunEmbeddedMigrations(dbUrl string) error {
	// 1. Open a temporary connection for the migration engine
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	// 2. Initialize the Postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// 3. Create a source driver from our embedded filesystem
	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	// 4. Initialize the migration engine
	m, err := migrate.NewWithInstance(
		"iofs",       // Source name
		sourceDriver, // Source driver
		"postgres",   // Database name
		driver,       // Database driver
	)
	if err != nil {
		return err
	}

	// 5. Apply all 'UP' migrations
	fmt.Println("  🚀 Checking database schema version...")
	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failure: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("  ✔ Database is already up-to-date.")
	} else {
		fmt.Println("  ✔ Schema evolution successful!")
	}

	return nil
}

func main() {
	fmt.Println("=== Database Migrations Boot Sequence ===")
	fmt.Println()

	// Demo URL (this would normally come from an Environment Variable)
	dbUrl := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"

	fmt.Printf("  Connecting to: %s\n", dbUrl)
	fmt.Println("  (Note: This will fail as no Postgres is running locally.)")

	err := RunEmbeddedMigrations(dbUrl)
	if err != nil {
		fmt.Printf("\n  Caught expected error: %v\n", err)
		fmt.Println("  In a real environment (like the Flagship Project), this would apply your SQL changes automatically.")
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.1 web-masterclass")
	fmt.Println("Current: DM.1 (embedded-migrations)")
	fmt.Println("Previous: DB.8 (query-timeouts-via-context)")
	fmt.Println("---------------------------------------------------")
}
