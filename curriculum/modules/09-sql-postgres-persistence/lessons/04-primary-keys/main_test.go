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

func TestCreateTables(t *testing.T) {
	db := openMemory(t)
	if err := createTables(db); err != nil {
		t.Fatal(err)
	}
}

func TestAutoIncrementPK(t *testing.T) {
	db := openMemory(t)
	if err := createTables(db); err != nil {
		t.Fatal(err)
	}

	id, err := insertAuthor(db, &Author{Name: "Test Author"})
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
}

func TestLastInsertId(t *testing.T) {
	db := openMemory(t)
	if err := createTables(db); err != nil {
		t.Fatal(err)
	}

	aid, err := insertAuthor(db, &Author{Name: "Author"})
	if err != nil {
		t.Fatal(err)
	}
	bid, err := insertBook(db, &Book{Title: "Book", ISBN: "1234567890", AuthorID: aid})
	if err != nil {
		t.Fatal(err)
	}
	if bid == 0 {
		t.Fatal("expected non-zero book id")
	}
}

func TestUniqueConstraint(t *testing.T) {
	db := openMemory(t)
	if err := createTables(db); err != nil {
		t.Fatal(err)
	}

	aid, err := insertAuthor(db, &Author{Name: "Author"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = insertBook(db, &Book{Title: "Book1", ISBN: "1111111111", AuthorID: aid})
	if err != nil {
		t.Fatal(err)
	}
	_, err = insertBook(db, &Book{Title: "Book2", ISBN: "1111111111", AuthorID: aid})
	if err == nil {
		t.Fatal("expected UNIQUE constraint violation")
	}
}

func TestUUIDPK(t *testing.T) {
	db := openMemory(t)
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS courses (id TEXT PRIMARY KEY, name TEXT NOT NULL)`)
	if err != nil {
		t.Fatal(err)
	}

	uid := newUUID()
	_, err = db.Exec(`INSERT INTO courses (id, name) VALUES (?, ?)`, uid, "Test Course")
	if err != nil {
		t.Fatal(err)
	}

	var name string
	err = db.QueryRow(`SELECT name FROM courses WHERE id = ?`, uid).Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	if name != "Test Course" {
		t.Fatalf("got %q, want %q", name, "Test Course")
	}
}

func TestNewUUIDUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		u := newUUID()
		if seen[u] {
			t.Fatal("duplicate UUID generated")
		}
		seen[u] = true
	}
}
