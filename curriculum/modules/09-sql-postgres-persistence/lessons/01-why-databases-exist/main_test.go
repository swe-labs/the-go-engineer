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

func TestOpenInMemory(t *testing.T) {
	db := openMemory(t)
	if err := db.Ping(); err != nil {
		t.Fatal("ping failed:", err)
	}
}

func TestCreateAndQuery(t *testing.T) {
	db := openMemory(t)

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS messages (id INTEGER NOT NULL PRIMARY KEY, text TEXT NOT NULL)`)
	if err != nil {
		t.Fatal("create:", err)
	}

	_, err = db.Exec(`INSERT INTO messages (text) VALUES (?)`, "hello")
	if err != nil {
		t.Fatal("insert:", err)
	}

	var text string
	err = db.QueryRow(`SELECT text FROM messages WHERE id = ?`, 1).Scan(&text)
	if err != nil {
		t.Fatal("query:", err)
	}
	if text != "hello" {
		t.Fatalf("got %q, want %q", text, "hello")
	}
}

func TestSaveAndLoad(t *testing.T) {
	db := openMemory(t)

	id, err := SaveAndLoad(db, "test note")
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("got id %d, want 1", id)
	}

	var content string
	err = db.QueryRow(`SELECT content FROM notes WHERE id = ?`, id).Scan(&content)
	if err != nil {
		t.Fatal("query saved note:", err)
	}
	if content != "test note" {
		t.Fatalf("got content %q, want %q", content, "test note")
	}
}

func TestSaveAndLoadMultiple(t *testing.T) {
	db := openMemory(t)

	id1, err := SaveAndLoad(db, "first")
	if err != nil {
		t.Fatal(err)
	}
	id2, err := SaveAndLoad(db, "second")
	if err != nil {
		t.Fatal(err)
	}
	if id1 != 1 || id2 != 2 {
		t.Fatalf("expected ids 1 and 2, got %d and %d", id1, id2)
	}
}
