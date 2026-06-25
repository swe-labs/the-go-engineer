package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestMigrationUp(t *testing.T) {
	db := newTestDB(t)
	m := NewMigrator(db)

	if err := m.Up(); err != nil {
		t.Fatal(err)
	}

	version, err := m.CurrentVersion()
	if err != nil {
		t.Fatal(err)
	}
	if version != 2 {
		t.Errorf("expected version 2, got %d", version)
	}
}

func TestMigrationIdempotent(t *testing.T) {
	db := newTestDB(t)
	m := NewMigrator(db)

	if err := m.Up(); err != nil {
		t.Fatal(err)
	}
	if err := m.Up(); err != nil {
		t.Fatal(err)
	}

	version, err := m.CurrentVersion()
	if err != nil {
		t.Fatal(err)
	}
	if version != 2 {
		t.Errorf("expected version 2 after second up, got %d", version)
	}
}

func TestMigrationDown(t *testing.T) {
	db := newTestDB(t)
	m := NewMigrator(db)

	if err := m.Up(); err != nil {
		t.Fatal(err)
	}

	if err := m.Down(); err != nil {
		t.Fatal(err)
	}

	version, err := m.CurrentVersion()
	if err != nil {
		t.Fatal(err)
	}
	if version != 1 {
		t.Errorf("expected version 1 after one down, got %d", version)
	}
}

func TestMigrationFullCycle(t *testing.T) {
	db := newTestDB(t)
	m := NewMigrator(db)

	if err := m.Up(); err != nil {
		t.Fatal(err)
	}
	if err := m.Down(); err != nil {
		t.Fatal(err)
	}
	if err := m.Down(); err != nil {
		t.Fatal(err)
	}

	version, err := m.CurrentVersion()
	if err != nil {
		t.Fatal(err)
	}
	if version != 0 {
		t.Errorf("expected version 0 after two downs, got %d", version)
	}
}

func TestUsersTableCreated(t *testing.T) {
	db := newTestDB(t)
	m := NewMigrator(db)

	m.Up()

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if count != 1 {
		t.Error("users table was not created")
	}
}

func TestAgeColumnAdded(t *testing.T) {
	db := newTestDB(t)
	m := NewMigrator(db)

	m.Up()

	var hasAge bool
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull bool
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			t.Fatal(err)
		}
		if name == "age" {
			hasAge = true
		}
	}
	if !hasAge {
		t.Error("age column was not added")
	}
}
