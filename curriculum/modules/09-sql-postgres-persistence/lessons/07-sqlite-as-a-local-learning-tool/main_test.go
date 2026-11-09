package main

import (
	"database/sql"
	"os"
	"path/filepath"
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

func TestInMemoryCreateAndQuery(t *testing.T) {
	db := openMemory(t)

	_, err := db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, val TEXT)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO test (val) VALUES ('hello')`)
	if err != nil {
		t.Fatal(err)
	}

	var val string
	err = db.QueryRow(`SELECT val FROM test WHERE id = 1`).Scan(&val)
	if err != nil {
		t.Fatal(err)
	}
	if val != "hello" {
		t.Fatalf("got %q, want %q", val, "hello")
	}
}

func TestFileBackedCreateAndQuery(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, val TEXT)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO test (val) VALUES ('persistent')`)
	if err != nil {
		t.Fatal(err)
	}

	db.Close()

	db2, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()

	var val string
	err = db2.QueryRow(`SELECT val FROM test WHERE id = 1`).Scan(&val)
	if err != nil {
		t.Fatal(err)
	}
	if val != "persistent" {
		t.Fatalf("got %q, want %q", val, "persistent")
	}
}

func TestWalMode(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "wal.db")
	db := openMemory(t)

	var mode string
	err := db.QueryRow(`PRAGMA journal_mode`).Scan(&mode)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("default journal mode: %s", mode)

	db2, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()

	_, err = db2.Exec(`PRAGMA journal_mode=WAL`)
	if err != nil {
		t.Fatal(err)
	}

	var mode2 string
	err = db2.QueryRow(`PRAGMA journal_mode`).Scan(&mode2)
	if err != nil {
		t.Fatal(err)
	}
	if mode2 != "wal" {
		t.Fatalf("expected journal_mode=wal, got %q", mode2)
	}
}

func TestSQLiteVersion(t *testing.T) {
	db := openMemory(t)

	var version string
	err := db.QueryRow(`SELECT sqlite_version()`).Scan(&version)
	if err != nil {
		t.Fatal(err)
	}
	if version == "" {
		t.Fatal("expected non-empty version")
	}
	t.Logf("SQLite version: %s", version)
}

func TestFileCleanup(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "cleanup.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	// Force connection creation so the file is written
	db.Ping()
	db.Close()

	exists, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		t.Fatal("expected db file to exist after ping + close")
	}
	_ = exists
}
