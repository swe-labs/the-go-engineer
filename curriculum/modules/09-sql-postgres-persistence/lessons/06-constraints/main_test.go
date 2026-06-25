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

func TestCreateTable(t *testing.T) {
	db := openMemory(t)
	if err := createUserTable(db); err != nil {
		t.Fatal(err)
	}
}

func TestValidInsert(t *testing.T) {
	db := openMemory(t)
	if err := createUserTable(db); err != nil {
		t.Fatal(err)
	}

	id, err := insertUser(db, &User{Name: "Alice", Email: "a@b.com", Age: 25})
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
}

func TestNotNullViolation(t *testing.T) {
	db := openMemory(t)
	if err := createUserTable(db); err != nil {
		t.Fatal(err)
	}

	_, err := db.Exec(`INSERT INTO users (name, email, age) VALUES (NULL, 'x@y.com', 20)`)
	if err == nil {
		t.Fatal("expected NOT NULL violation")
	}
}

func TestCheckViolation(t *testing.T) {
	db := openMemory(t)
	if err := createUserTable(db); err != nil {
		t.Fatal(err)
	}

	_, err := db.Exec(`INSERT INTO users (name, email, age) VALUES ('X', 'x@y.com', -5)`)
	if err == nil {
		t.Fatal("expected CHECK violation")
	}
}

func TestUniqueViolation(t *testing.T) {
	db := openMemory(t)
	if err := createUserTable(db); err != nil {
		t.Fatal(err)
	}

	_, err := insertUser(db, &User{Name: "A", Email: "dup@test.com", Age: 20})
	if err != nil {
		t.Fatal(err)
	}
	_, err = insertUser(db, &User{Name: "B", Email: "dup@test.com", Age: 30})
	if err == nil {
		t.Fatal("expected UNIQUE violation")
	}
}

func TestAgeBoundary(t *testing.T) {
	db := openMemory(t)
	if err := createUserTable(db); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		age   int
		valid bool
	}{
		{"zero", 0, true},
		{"one fifty", 150, true},
		{"negative", -1, false},
		{"one fifty one", 151, false},
	}
	for _, tc := range tests {
		_, err := db.Exec(`INSERT INTO users (name, email, age) VALUES (?, ?, ?)`,
			tc.name, tc.name+"@test.com", tc.age)
		got := err == nil
		if got != tc.valid {
			t.Fatalf("age=%d: valid=%v, want=%v (err=%v)", tc.age, got, tc.valid, err)
		}
	}
}
