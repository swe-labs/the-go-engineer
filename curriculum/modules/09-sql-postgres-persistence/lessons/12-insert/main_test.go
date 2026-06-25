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
	if err := CreateTable(db); err != nil {
		t.Fatal(err)
	}
}

func TestInsertOrder(t *testing.T) {
	db := openMemory(t)
	if err := CreateTable(db); err != nil {
		t.Fatal(err)
	}

	id, err := InsertOrder(db, "Test Customer", 5000)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
}

func TestDefaultStatus(t *testing.T) {
	db := openMemory(t)
	CreateTable(db)

	_, err := db.Exec(`INSERT INTO orders (customer, amount) VALUES ('Test', 100)`)
	if err != nil {
		t.Fatal(err)
	}

	var status string
	err = db.QueryRow(`SELECT status FROM orders WHERE id = 1`).Scan(&status)
	if err != nil {
		t.Fatal(err)
	}
	if status != "pending" {
		t.Fatalf("expected status 'pending', got %q", status)
	}
}

func TestInsertOrdersBatch(t *testing.T) {
	db := openMemory(t)
	CreateTable(db)

	orders := []struct {
		Customer string
		Amount   int
	}{
		{"A", 10},
		{"B", 20},
		{"C", 30},
	}
	count, err := InsertOrdersBatch(db, orders)
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("expected 3 rows affected, got %d", count)
	}

	var total int
	err = db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&total)
	if err != nil {
		t.Fatal(err)
	}
	if total != 3 {
		t.Fatalf("expected 3 total orders, got %d", total)
	}
}

func TestEmptyBatch(t *testing.T) {
	db := openMemory(t)
	CreateTable(db)

	count, err := InsertOrdersBatch(db, nil)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0, got %d", count)
	}
}

func TestPreparedStatement(t *testing.T) {
	db := openMemory(t)
	CreateTable(db)

	stmt, err := db.Prepare(`INSERT INTO orders (customer, amount) VALUES (?, ?)`)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < 5; i++ {
		_, err := stmt.Exec("Customer", 100)
		if err != nil {
			t.Fatal(err)
		}
	}

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&count)
	if count != 5 {
		t.Fatalf("expected 5 orders, got %d", count)
	}
}
