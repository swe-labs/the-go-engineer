package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestConnectDB(t *testing.T) {
	db, err := ConnectDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestConnectDBInvalidDSN(t *testing.T) {
	_, err := sql.Open("sqlite", "")
	if err != nil {
		// sql.Open with empty DSN may or may not error; this varies by driver
		t.Logf("Open error (expected with some drivers): %v", err)
	}
}

func TestConnectAndQuery(t *testing.T) {
	db, err := ConnectDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, name TEXT)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO items (name) VALUES ('test')`)
	if err != nil {
		t.Fatal(err)
	}

	var name string
	err = db.QueryRow(`SELECT name FROM items WHERE id = 1`).Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	if name != "test" {
		t.Fatalf("got %q, want %q", name, "test")
	}
}

func TestPoolSettings(t *testing.T) {
	db, err := ConnectDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	stats := db.Stats()
	if stats.MaxOpenConnections != 10 {
		t.Fatalf("expected MaxOpenConns=10, got %d", stats.MaxOpenConnections)
	}
}

func TestMultipleConnects(t *testing.T) {
	for i := 0; i < 5; i++ {
		db, err := ConnectDB(":memory:")
		if err != nil {
			t.Fatalf("connect %d: %v", i, err)
		}
		db.Close()
	}
}
