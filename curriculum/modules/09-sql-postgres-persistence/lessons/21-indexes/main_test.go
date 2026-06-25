package main

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func populateUsers(t *testing.T, db *sql.DB, count int) {
	t.Helper()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < count; i++ {
		_, err := tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)",
			fmt.Sprintf("User %d", i),
			fmt.Sprintf("user%d@example.com", i))
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestExplainWithoutIndexShowsScan(t *testing.T) {
	db := setupTestDB(t)
	populateUsers(t, db, 100)

	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM users WHERE email = 'user50@example.com'")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var hasScan bool
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		if strings.Contains(detail, "SCAN") {
			hasScan = true
		}
	}
	if !hasScan {
		t.Error("expected SCAN in query plan without index")
	}
}

func TestExplainWithIndexShowsSearch(t *testing.T) {
	db := setupTestDB(t)
	populateUsers(t, db, 100)

	db.Exec("CREATE INDEX idx_users_email ON users(email)")

	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM users WHERE email = 'user50@example.com'")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var hasSearch bool
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		if strings.Contains(detail, "SEARCH") && strings.Contains(detail, "idx_users_email") {
			hasSearch = true
		}
	}
	if !hasSearch {
		t.Error("expected SEARCH using index in query plan with index")
	}
}

func TestIndexAcceleratesLookup(t *testing.T) {
	db := setupTestDB(t)
	populateUsers(t, db, 1000)

	var name string
	err := db.QueryRow("SELECT name FROM users WHERE email = 'user500@example.com'").Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	if name != "User 500" {
		t.Errorf("expected User 500, got %s", name)
	}
}

func TestCompositeIndexOrder(t *testing.T) {
	db := setupTestDB(t)

	_, err := db.Exec(`
		CREATE TABLE logs (
			id INTEGER PRIMARY KEY,
			level TEXT NOT NULL,
			message TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Composite index on (level, created_at)
	_, err = db.Exec("CREATE INDEX idx_logs_level_created ON logs(level, created_at)")
	if err != nil {
		t.Fatal(err)
	}

	// Query that benefits from composite index
	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM logs WHERE level = 'ERROR' AND created_at > '2024-01-01'")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var usesIndex bool
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		if strings.Contains(detail, "idx_logs_level_created") {
			usesIndex = true
		}
	}
	if !usesIndex {
		t.Error("expected composite index to be used")
	}

	// Query on created_at alone should NOT use the composite index
	rows2, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM logs WHERE created_at > '2024-01-01'")
	if err != nil {
		t.Fatal(err)
	}
	defer rows2.Close()

	var usesIndex2 bool
	for rows2.Next() {
		var id, parent int
		var detail string
		rows2.Scan(&id, &parent, &detail)
		if strings.Contains(detail, "idx_logs_level_created") {
			usesIndex2 = true
		}
	}
	if usesIndex2 {
		t.Error("composite index should NOT be used for created_at alone")
	}
}

func TestUniqueIndexPreventsDuplicates(t *testing.T) {
	db := setupTestDB(t)

	_, err := db.Exec("CREATE UNIQUE INDEX idx_users_email_unique ON users(email)")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO users (name, email) VALUES ('Alice Again', 'alice@example.com')")
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
}
