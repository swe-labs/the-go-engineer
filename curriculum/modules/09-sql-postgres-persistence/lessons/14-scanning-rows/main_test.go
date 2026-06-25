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

func setupEmployees(t *testing.T, db *sql.DB) {
	t.Helper()
	db.Exec(`CREATE TABLE employees (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		salary REAL NOT NULL,
		active INTEGER NOT NULL DEFAULT 1
	)`)
	db.Exec(`INSERT INTO employees (name, salary, active) VALUES ('Alice', 75000, 1)`)
	db.Exec(`INSERT INTO employees (name, salary, active) VALUES ('Bob', 82000, 1)`)
}

func TestScanEmployee(t *testing.T) {
	db := openMemory(t)
	setupEmployees(t, db)

	e, err := scanEmployee(db.QueryRow(`SELECT id, name, salary, active FROM employees WHERE id = 1`))
	if err != nil {
		t.Fatal(err)
	}
	if e.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", e.Name)
	}
	if e.Salary != 75000 {
		t.Fatalf("expected 75000, got %f", e.Salary)
	}
}

func TestScanEmployees(t *testing.T) {
	db := openMemory(t)
	setupEmployees(t, db)

	rows, err := db.Query(`SELECT id, name, salary, active FROM employees ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}

	employees, err := scanEmployees(rows)
	if err != nil {
		t.Fatal(err)
	}
	if len(employees) != 2 {
		t.Fatalf("expected 2 employees, got %d", len(employees))
	}
}

func TestScanNullString(t *testing.T) {
	db := openMemory(t)
	db.Exec(`CREATE TABLE t (id INTEGER PRIMARY KEY, val TEXT)`)
	db.Exec(`INSERT INTO t (id, val) VALUES (1, 'hello')`)
	db.Exec(`INSERT INTO t (id, val) VALUES (2, NULL)`)

	var ns sql.NullString
	err := db.QueryRow(`SELECT val FROM t WHERE id = 1`).Scan(&ns)
	if err != nil {
		t.Fatal(err)
	}
	if !ns.Valid || ns.String != "hello" {
		t.Fatalf("expected valid 'hello', got valid=%v val=%q", ns.Valid, ns.String)
	}

	err = db.QueryRow(`SELECT val FROM t WHERE id = 2`).Scan(&ns)
	if err != nil {
		t.Fatal(err)
	}
	if ns.Valid {
		t.Fatalf("expected invalid, got valid=true val=%q", ns.String)
	}
}

func TestScanPerson(t *testing.T) {
	db := openMemory(t)
	db.Exec(`CREATE TABLE IF NOT EXISTS people (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT,
		age INTEGER
	)`)
	db.Exec(`INSERT INTO people (name, email, age) VALUES ('Alice', 'alice@example.com', 30)`)
	db.Exec(`INSERT INTO people (name, email, age) VALUES ('Bob', NULL, NULL)`)

	p, err := scanPerson(db.QueryRow(`SELECT id, name, email, age FROM people WHERE id = 1`))
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", p.Name)
	}
	if !p.Email.Valid || p.Email.String != "alice@example.com" {
		t.Fatalf("expected valid email")
	}
	if !p.Age.Valid || p.Age.Int64 != 30 {
		t.Fatalf("expected age 30")
	}

	p, err = scanPerson(db.QueryRow(`SELECT id, name, email, age FROM people WHERE id = 2`))
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Bob" {
		t.Fatalf("expected Bob, got %s", p.Name)
	}
	if p.Email.Valid {
		t.Fatal("expected invalid email")
	}
	if p.Age.Valid {
		t.Fatal("expected invalid age")
	}
}

func TestScanError(t *testing.T) {
	db := openMemory(t)
	setupEmployees(t, db)

	var name string
	err := db.QueryRow(`SELECT id, name FROM employees WHERE id = 1`).Scan(&name)
	if err == nil {
		t.Fatal("expected error for mismatched column count")
	}
}
