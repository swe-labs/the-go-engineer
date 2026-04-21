// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./06-backend-db/01-web-and-database/database-migrations/1-embedded-migrations
package main

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	// Blank import the postgres driver to register it natively
	_ "github.com/lib/pq"
	// Blank import to ensure the `file://` scheme works (though we use iofs)
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// ============================================================================
// Stage 06: Database Migrations — Embedded Schema Evolution
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Bundling `.sql` migration files natively into the compiled Go binary
//   - Instantiating a golang-migrate driver and connecting to PostgreSQL
//   - Running `Up()` to apply database changes automatically
//
// ANALOGY:
//   Think of migrations as "Git tracking for your database".
//   It allows developers to keep their SQL schemas synchronized reliably,
//   rolling back mistakes via DOWN scripts or pushing forward via UP scripts.
//
// ENGINEERING DEPTH:
//   Using `//go:embed` means your Go backend binary includes its own setup instructions.
//   You don't need a separate CI step to run SQL migrations or upload heavy ZIP
//   files. When the binary starts, it inspects the DB, checks the current schema version,
//   and executes only the missing `.up.sql` scripts.
// ============================================================================

// --- 1. Embed the migrations folder ---

//go:embed migrations/*.sql
var migrationFiles embed.FS

// --- 2. Migration Execution Function ---

func RunEmbeddedMigrations(dbUrl string) error {
	// 1. Connect to PostgreSQL to give 'migrate' access to the metadata table
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	defer db.Close() // Keep this clean

	// 2. Initialize the Postgres Driver for golang-migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// 3. Load the folder of `.sql` files from the Embedded Filesystem (iofs)
	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	// 4. Create the Migration engine
	m, err := migrate.NewWithInstance(
		"iofs",       // Virtual folder scheme
		sourceDriver, // Read the embedded SQL files
		"postgres",   // Database engine
		driver,       // The active db connection
	)
	if err != nil {
		return err
	}

	// 5. Execute the UP migrations!
	fmt.Println("🚀 Executing Pre-flight Database Migrations...")
	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		// A fatal error occurred during schema execution
		return fmt.Errorf("migration failure: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("✅ Database schema is up-to-date (no new migrations found).")
	} else {
		fmt.Println("✅ Database schema evolution successful!")
	}

	return nil
}

func main() {
	// Normally this connection string would come from an os.Getenv("DATABASE_URL")
	fmt.Println("=== Automatic Embedded Migrations Boot Sequence ===")

	// We intentionally leave the DB URL broken here because this is just an example script.
	// In Section 22 (Enterprise Capstone), we will use this exact function
	// connected to a live Dockerized PostgreSQL instance!
	dbUrl := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"

	fmt.Println("   Note: Attempting to connect to dummy URL:", dbUrl)

	err := RunEmbeddedMigrations(dbUrl)
	if err != nil {
		fmt.Printf("   Caught expected error: %v\n", err)
		fmt.Println("\n   (Proceed to Section 22 to see this run perfectly under Docker Compose!)")
	}
}
