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

func enableFK(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestForeignKeyViolation(t *testing.T) {
	db := openMemory(t)
	enableFK(t, db)

	_, err := db.Exec(`
		CREATE TABLE departments (id INTEGER PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE employees (id INTEGER PRIMARY KEY, name TEXT NOT NULL, dept_id INTEGER NOT NULL REFERENCES departments(id))
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES ('Orphan', 999)`)
	if err == nil {
		t.Fatal("expected foreign key violation")
	}
}

func TestCascadeDelete(t *testing.T) {
	db := openMemory(t)
	enableFK(t, db)

	_, err := db.Exec(`
		CREATE TABLE departments (id INTEGER PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE employees (id INTEGER PRIMARY KEY, name TEXT NOT NULL, dept_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE CASCADE)
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO departments (name) VALUES ('Eng')`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES ('Alice', 1)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`DELETE FROM departments WHERE id = 1`)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM employees`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0 employees after cascade, got %d", count)
	}
}

func TestNoCascadeWithoutFK(t *testing.T) {
	db := openMemory(t)

	_, err := db.Exec(`
		CREATE TABLE departments (id INTEGER PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE employees (id INTEGER PRIMARY KEY, name TEXT NOT NULL, dept_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE CASCADE)
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO departments (name) VALUES ('Eng')`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES ('Alice', 1)`)
	if err != nil {
		t.Fatal(err)
	}

	// Without PRAGMA foreign_keys = ON, the insert with invalid FK succeeds
	_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES ('Bob', 999)`)
	if err != nil {
		t.Fatal("expected no error without FK pragma, got:", err)
	}
}

func TestCascadeMultipleChildren(t *testing.T) {
	db := openMemory(t)
	enableFK(t, db)

	_, err := db.Exec(`
		CREATE TABLE departments (id INTEGER PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE employees (id INTEGER PRIMARY KEY, name TEXT NOT NULL, dept_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE CASCADE)
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO departments (name) VALUES ('Eng')`)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"A", "B", "C"} {
		_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES (?, 1)`, name)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = db.Exec(`DELETE FROM departments WHERE id = 1`)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM employees`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0 employees, got %d", count)
	}
}
