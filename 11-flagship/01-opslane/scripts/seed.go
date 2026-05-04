//go:build ignore

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Seed script for Opslane development
// Usage: go run ./11-flagship/01-opslane/scripts/seed.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	dbDsn = flag.String("dsn", "", "Database connection string (or use OPSLANE_DB_DSN)")
	reset = flag.Bool("reset", false, "Reset (delete) existing seed data before seeding")
)

func main() {
	flag.Parse()

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

	if *reset {
		if err := resetSeedData(ctx, db); err != nil {
			log.Fatalf("Failed to reset seed data: %v", err)
		}
		fmt.Println("✓ Seed data reset")
	}

	if err := seedData(ctx, db); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	fmt.Println("✓ Seed data applied successfully")
}

func resetSeedData(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM payments WHERE tenant_id = 1")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM orders WHERE tenant_id = 1")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE tenant_id = 1")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM tenants WHERE id = 1")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func seedData(ctx context.Context, db *sql.DB) error {
	// Generate password hash
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert tenant
	_, err = tx.ExecContext(ctx, `
		INSERT INTO tenants (id, name, slug, created_at) 
		VALUES (1, 'Demo Organization', 'demo', $1)
		ON CONFLICT (slug) DO NOTHING
	`, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert tenant: %w", err)
	}

	// Insert users
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, tenant_id, email, display_name, password_hash, role, created_at) 
		VALUES 
			(1, 1, 'admin@demo.com', 'Demo Admin', $1, 'admin', $2),
			(2, 1, 'user@demo.com', 'Demo User', $1, 'user', $2)
		ON CONFLICT (tenant_id, email) DO NOTHING
	`, string(hash), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}

	// Insert orders
	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders (id, tenant_id, user_id, status, total_cents, currency, idempotency_key, created_at, updated_at)
		VALUES 
			(1, 1, 1, 'completed', 10000, 'USD', 'seed-order-001', $1, $1),
			(2, 1, 1, 'pending', 25000, 'USD', 'seed-order-002', $1, $1),
			(3, 1, 2, 'completed', 5000, 'USD', 'seed-order-003', $1, $1)
		ON CONFLICT (tenant_id, idempotency_key) DO NOTHING
	`, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert orders: %w", err)
	}

	// Insert payments
	_, err = tx.ExecContext(ctx, `
		INSERT INTO payments (id, tenant_id, order_id, status, provider_reference, amount_cents, failure_reason, created_at, updated_at)
		VALUES 
			(1, 1, 1, 'succeeded', 'seed-pi-001', 10000, '', $1, $1),
			(2, 1, 2, 'pending', 'seed-pi-002', 25000, '', $1, $1),
			(3, 1, 3, 'succeeded', 'seed-pi-003', 5000, '', $1, $1)
		ON CONFLICT (tenant_id, provider_reference) DO NOTHING
	`, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert payments: %w", err)
	}

	return tx.Commit()
}
