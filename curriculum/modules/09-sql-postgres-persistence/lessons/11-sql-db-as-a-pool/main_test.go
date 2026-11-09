package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

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

func TestConfigurePool(t *testing.T) {
	db := openMemory(t)
	configurePool(db)

	stats := db.Stats()
	if stats.MaxOpenConnections != 10 {
		t.Fatalf("expected MaxOpenConns=10, got %d", stats.MaxOpenConnections)
	}
}

func TestSerialInsert(t *testing.T) {
	db := openMemory(t)
	configurePool(db)

	db.Exec(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, ts TEXT, msg TEXT)`)

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`INSERT INTO events (msg) VALUES (?)`, fmt.Sprintf("evt%d", i))
		if err != nil {
			t.Fatal(err)
		}
	}

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&count)
	if count != 10 {
		t.Fatalf("expected 10 events, got %d", count)
	}
}

func TestPoolStats(t *testing.T) {
	db := openMemory(t)
	configurePool(db)

	stats := MonitorPool(db)
	if stats == nil {
		t.Fatal("expected non-nil stats")
	}
}

func TestIdleConnections(t *testing.T) {
	db := openMemory(t)
	configurePool(db)

	db.SetMaxIdleConns(5)

	db.Exec(`SELECT 1`)

	stats := db.Stats()
	if stats.Idle > 5 {
		t.Fatalf("idle connections %d exceeded max idle 5", stats.Idle)
	}
}

func TestConnMaxLifetime(t *testing.T) {
	dbPath := t.TempDir() + "/lifetime.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	db.SetConnMaxLifetime(5 * time.Second)

	db.Exec(`CREATE TABLE IF NOT EXISTS t (id INTEGER PRIMARY KEY)`)
	db.Exec(`INSERT INTO t VALUES (1)`)

	var id int
	err = db.QueryRow(`SELECT id FROM t WHERE id = 1`).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("expected 1, got %d", id)
	}
}
