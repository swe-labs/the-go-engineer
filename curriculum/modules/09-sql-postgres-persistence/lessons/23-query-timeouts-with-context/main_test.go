package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE numbers (id INTEGER PRIMARY KEY, value TEXT NOT NULL)`)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2000; i++ {
		db.Exec("INSERT INTO numbers (value) VALUES (?)", fmt.Sprintf("number-%d", i))
	}
	return db
}

func TestQueryTimeoutExceeded(t *testing.T) {
	db := setupTestDB(t)

	slowQuery := `SELECT COUNT(*) FROM numbers a
		CROSS JOIN numbers b
		CROSS JOIN numbers c
		CROSS JOIN numbers d`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()

	_, err := db.QueryContext(ctx, slowQuery)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got: %v", err)
	}
}

func TestQueryWithAdequateTimeout(t *testing.T) {
	db := setupTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM numbers").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2000 {
		t.Errorf("expected 2000, got %d", count)
	}
}

func TestQueryContextCancellation(t *testing.T) {
	db := setupTestDB(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := db.QueryContext(ctx, "SELECT * FROM numbers")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled, got: %v", err)
	}
}

func TestExecContextWithTimeout(t *testing.T) {
	db := setupTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO numbers (value) VALUES (?)", "test")
	if err != nil {
		t.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM numbers WHERE value = 'test'").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 row inserted, got %d", count)
	}
}

func TestQueryRowContextWithTimeout(t *testing.T) {
	db := setupTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var value string
	err := db.QueryRowContext(ctx, "SELECT value FROM numbers WHERE id = ?", 1).Scan(&value)
	if err != nil {
		t.Fatal(err)
	}
	if value != "number-0" {
		t.Errorf("expected 'number-0', got %s", value)
	}
}

func TestContextPropagation(t *testing.T) {
	db := setupTestDB(t)

	parentCtx, parentCancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer parentCancel()

	// Child context with longer timeout should still be capped by parent
	childCtx, childCancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer childCancel()

	slowQuery := `SELECT COUNT(*) FROM numbers a
		CROSS JOIN numbers b
		CROSS JOIN numbers c
		CROSS JOIN numbers d`

	_, err := db.QueryContext(childCtx, slowQuery)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded from parent context, got: %v", err)
	}
}
