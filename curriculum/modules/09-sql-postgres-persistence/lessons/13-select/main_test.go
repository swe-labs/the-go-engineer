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

func setupProducts(t *testing.T, db *sql.DB) {
	t.Helper()
	db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0
	)`)
	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Widget', 9.99, 100)`)
	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Gadget', 24.99, 50)`)
	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Gizmo', 14.99, 0)`)
}

func TestQueryAllProducts(t *testing.T) {
	db := openMemory(t)
	setupProducts(t, db)

	products, err := queryAllProducts(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 3 {
		t.Fatalf("expected 3 products, got %d", len(products))
	}
}

func TestQueryProductByID(t *testing.T) {
	db := openMemory(t)
	setupProducts(t, db)

	p, err := queryProductByID(db, 1)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Widget" {
		t.Fatalf("expected Widget, got %s", p.Name)
	}
}

func TestQueryProductByIDNotFound(t *testing.T) {
	db := openMemory(t)
	setupProducts(t, db)

	_, err := queryProductByID(db, 999)
	if err != sql.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestSearchProducts(t *testing.T) {
	db := openMemory(t)
	setupProducts(t, db)

	results, err := SearchProducts(db, 10.00, 20.00)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 product in range 10-20, got %d", len(results))
	}
	if results[0].Name != "Gizmo" {
		t.Fatalf("expected Gizmo, got %s", results[0].Name)
	}
}

func TestSearchProductsEmpty(t *testing.T) {
	db := openMemory(t)
	setupProducts(t, db)

	results, err := SearchProducts(db, 100.00, 200.00)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 products, got %d", len(results))
	}
}

func TestQueryRowScanError(t *testing.T) {
	db := openMemory(t)
	setupProducts(t, db)

	var name string
	var price float64
	err := db.QueryRow(`SELECT name, price FROM products WHERE id = ?`, 1).Scan(&name)
	if err == nil {
		t.Fatal("expected Scan error for mismatched columns")
	}
	_ = price
}
