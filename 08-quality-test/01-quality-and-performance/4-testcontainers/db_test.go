// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

//go:build integration

// RUN: go test -v -tags=integration ./08-quality-test/01-quality-and-performance/4-testcontainers
package testcontainers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ============================================================================
// Stage 08: Quality & Performance — Integration Testing with Testcontainers
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to spin up a real PostgreSQL database inside a Docker container from Go tests.
//   - How to run raw SQL queries against it to verify your repository layers.
//   - Automatic container teardown so tests don't leave messy state behind.
//
// ENGINEERING DEPTH:
//   Mocking is great for unit testing business logic,
//   but it completely fails to test whether your SQL queries are structurally valid.
//   An application can have 100% test coverage with mocks, yet crash in production
//   because of a typo in an INSERT statement!
//
//   Testcontainers solves this by programmatically launching real infrastructure
//   (DBs, Redis, Kafka) using Docker during `go test`.
// ============================================================================

// --- 1. The Real Repository ---
// This is typical of what we'd test.
type UserRepository struct {
	db *sql.DB
}

func (repo *UserRepository) CreateTable(ctx context.Context) error {
	_, err := repo.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL);")
	return err
}

func (repo *UserRepository) InsertUser(ctx context.Context, name string) (int, error) {
	var id int
	err := repo.db.QueryRowContext(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING id;", name).Scan(&id)
	return id, err
}

// --- 2. The Integration Test ---

func TestUserRepositoryIntegration(t *testing.T) {
	// Skip if Docker is not available in the environment (e.g. some CI pipelines)
	// For this lesson, we assume Docker is running.

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Step 1: Spin up a PostgreSQL container using the official module
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	// Clean up the container fully after the test finishes
	defer func() {
		if err := postgresContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate postgres container: %s", err)
		}
	}()

	// Step 2: Grab the dynamic connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// Step 3: Connect to the DB
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to db: %s", err)
	}
	defer db.Close()

	// Wait for it to be actually ready
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("db ping failed: %s", err)
	}

	// Step 4: Execute the test against a REAL database!
	repo := &UserRepository{db: db}

	err = repo.CreateTable(ctx)
	assert.NoError(t, err)

	insertedId, err := repo.InsertUser(ctx, "Alice Testcontainers")
	assert.NoError(t, err)
	assert.Equal(t, 1, insertedId)

	fmt.Println("✅ Integration test passed against a live, ephemeral PostgreSQL database!")
}
