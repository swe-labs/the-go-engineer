package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func openMemory(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestWaitForDBSuccess(t *testing.T) {
	db, err := WaitForDB(":memory:", 5)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
}

func TestWaitForDBRetry(t *testing.T) {
	_, err := WaitForDB(":memory:", 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateAndQueryUsers(t *testing.T) {
	db := openMemory(t)

	db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)`)
	db.Exec(`INSERT INTO users (name) VALUES ('Alice')`)

	var name string
	err := db.QueryRow(`SELECT name FROM users WHERE id = 1`).Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	if name != "Alice" {
		t.Fatalf("got %q, want %q", name, "Alice")
	}
}

func TestMultipleUsers(t *testing.T) {
	db := openMemory(t)

	db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)`)
	names := []string{"A", "B", "C"}
	for _, n := range names {
		_, err := db.Exec(`INSERT INTO users (name) VALUES (?)`, n)
		if err != nil {
			t.Fatal(err)
		}
	}

	rows, err := db.Query(`SELECT name FROM users ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		if name != names[count] {
			t.Fatalf("row %d: got %q, want %q", count, name, names[count])
		}
		count++
	}
	if count != 3 {
		t.Fatalf("got %d rows, want 3", count)
	}
}

func TestSerialUpdates(t *testing.T) {
	db := openMemory(t)

	db.Exec(`CREATE TABLE IF NOT EXISTS counters (id INTEGER PRIMARY KEY, val INTEGER NOT NULL DEFAULT 0)`)
	db.Exec(`INSERT INTO counters (val) VALUES (0)`)

	for j := 0; j < 10; j++ {
		_, err := db.Exec(`UPDATE counters SET val = val + 1 WHERE id = 1`)
		if err != nil {
			t.Fatal(err)
		}
	}

	var val int
	err := db.QueryRow(`SELECT val FROM counters WHERE id = 1`).Scan(&val)
	if err != nil {
		t.Fatal(err)
	}
	if val != 10 {
		t.Fatalf("expected 10, got %d", val)
	}
}
