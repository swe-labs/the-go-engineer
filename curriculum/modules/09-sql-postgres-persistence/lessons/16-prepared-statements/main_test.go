package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL
	)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestBatchInsertProducts(t *testing.T) {
	db := setupTestDB(t)

	products := []Product{
		{Name: "Widget", Price: 9.99},
		{Name: "Gadget", Price: 24.99},
		{Name: "Doohickey", Price: 49.99},
	}

	err := BatchInsertProducts(db, products)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count != 3 {
		t.Fatalf("expected 3 products, got %d", count)
	}
}

func TestBatchInsertProductsEmpty(t *testing.T) {
	db := setupTestDB(t)

	err := BatchInsertProducts(db, []Product{})
	if err != nil {
		t.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count != 0 {
		t.Fatalf("expected 0 products, got %d", count)
	}
}

func TestBatchInsertProductsDuplicateFails(t *testing.T) {
	db := setupTestDB(t)

	_, err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_name ON products(name)")
	if err != nil {
		t.Fatal(err)
	}

	products := []Product{
		{Name: "Widget", Price: 9.99},
		{Name: "Widget", Price: 14.99},
	}

	err = BatchInsertProducts(db, products)
	if err == nil {
		t.Fatal("expected error for duplicate name, got nil")
	}
}

func TestPrepareOutsideLoop(t *testing.T) {
	db := setupTestDB(t)

	stmt, err := db.Prepare("INSERT INTO products (name, price) VALUES (?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < 100; i++ {
		_, err := stmt.Exec("test", 1.0)
		if err != nil {
			t.Fatal(err)
		}
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count != 100 {
		t.Fatalf("expected 100 products, got %d", count)
	}
}
