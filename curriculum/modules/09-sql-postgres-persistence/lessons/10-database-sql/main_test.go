package main

import (
	"database/sql"
	"fmt"
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

func TestInitDB(t *testing.T) {
	db, err := InitDB("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestExecAndQuery(t *testing.T) {
	db := openMemory(t)

	_, err := db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, val TEXT)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO test (val) VALUES (?)`, "a")
	if err != nil {
		t.Fatal(err)
	}

	var val string
	err = db.QueryRow(`SELECT val FROM test WHERE id = 1`).Scan(&val)
	if err != nil {
		t.Fatal(err)
	}
	if val != "a" {
		t.Fatalf("got %q, want %q", val, "a")
	}
}

func TestMultipleRows(t *testing.T) {
	db := openMemory(t)

	db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, val TEXT)`)
	for i := 0; i < 3; i++ {
		db.Exec(`INSERT INTO test (val) VALUES (?)`, fmt.Sprintf("row%d", i))
	}

	rows, err := db.Query(`SELECT val FROM test ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("row%d", count)
		if val != expected {
			t.Fatalf("got %q, want %q", val, expected)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("got %d rows, want 3", count)
	}
}

func TestPoolConfig(t *testing.T) {
	db, err := InitDB("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	stats := db.Stats()
	if stats.MaxOpenConnections != 5 {
		t.Fatalf("expected MaxOpenConns=5, got %d", stats.MaxOpenConnections)
	}
}

func TestPing(t *testing.T) {
	db := openMemory(t)
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
}
