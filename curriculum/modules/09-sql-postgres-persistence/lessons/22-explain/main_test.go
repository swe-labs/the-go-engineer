package main

import (
	"database/sql"
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
			email TEXT NOT NULL
		);
		INSERT INTO users VALUES (1, 'Alice', 'alice@example.com');
		INSERT INTO users VALUES (2, 'Bob', 'bob@example.com');
		INSERT INTO users VALUES (3, 'Carol', 'carol@example.com');
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestExplainPrimaryKeyLookup(t *testing.T) {
	db := setupTestDB(t)

	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM users WHERE id = 1")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var foundPK bool
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		if strings.Contains(detail, "PRIMARY KEY") {
			foundPK = true
		}
	}
	if !foundPK {
		t.Error("expected PRIMARY KEY lookup in query plan")
	}
}

func TestExplainFullScanWithoutIndex(t *testing.T) {
	db := setupTestDB(t)

	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM users WHERE email = 'alice@example.com'")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var foundScan bool
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		if strings.Contains(detail, "SCAN") {
			foundScan = true
		}
	}
	if !foundScan {
		t.Error("expected SCAN in query plan without index")
	}
}

func TestExplainIndexSearchWithIndex(t *testing.T) {
	db := setupTestDB(t)
	db.Exec("CREATE INDEX idx_users_email ON users(email)")

	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM users WHERE email = 'alice@example.com'")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var foundIndexSearch bool
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		if strings.Contains(detail, "SEARCH") && strings.Contains(detail, "idx_users_email") {
			foundIndexSearch = true
		}
	}
	if !foundIndexSearch {
		t.Error("expected SEARCH using index in query plan")
	}
}

func TestExplainOrderByWithoutIndexShowsTempBTree(t *testing.T) {
	db := setupTestDB(t)

	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM users ORDER BY name")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var foundTemp bool
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		if strings.Contains(detail, "TEMP") {
			foundTemp = true
		}
	}
	if !foundTemp {
		t.Error("expected TEMP B-TREE for ORDER BY without index")
	}
}

func TestExplainOrderByWithIndexNoTemp(t *testing.T) {
	db := setupTestDB(t)
	db.Exec("CREATE INDEX idx_users_name ON users(name)")

	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM users ORDER BY name")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var foundTemp bool
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		if strings.Contains(detail, "TEMP") {
			foundTemp = true
		}
	}
	if foundTemp {
		t.Error("expected no TEMP B-TREE for ORDER BY with index")
	}
}
