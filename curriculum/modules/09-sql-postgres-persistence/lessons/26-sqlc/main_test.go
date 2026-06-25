package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestQueries(t *testing.T) *Queries {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL
	)`)
	return NewQueries(db)
}

func TestCreateProduct(t *testing.T) {
	q := newTestQueries(t)

	result, err := q.CreateProduct(context.Background(), "Widget", 9.99)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := result.LastInsertId()
	if id == 0 {
		t.Error("expected non-zero id")
	}
}

func TestGetProduct(t *testing.T) {
	q := newTestQueries(t)

	q.CreateProduct(context.Background(), "Widget", 9.99)
	p, err := q.GetProduct(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Widget" {
		t.Errorf("name = %q, want %q", p.Name, "Widget")
	}
	if p.Price != 9.99 {
		t.Errorf("price = %.2f, want 9.99", p.Price)
	}
}

func TestGetProductNotFound(t *testing.T) {
	q := newTestQueries(t)

	_, err := q.GetProduct(context.Background(), 999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestListProducts(t *testing.T) {
	q := newTestQueries(t)

	q.CreateProduct(context.Background(), "A", 1.0)
	q.CreateProduct(context.Background(), "B", 2.0)
	q.CreateProduct(context.Background(), "C", 3.0)

	products, err := q.ListProducts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 3 {
		t.Errorf("expected 3, got %d", len(products))
	}
}

func TestListProductsEmpty(t *testing.T) {
	q := newTestQueries(t)

	products, err := q.ListProducts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 0 {
		t.Errorf("expected 0, got %d", len(products))
	}
}

func TestUpdateProductPrice(t *testing.T) {
	q := newTestQueries(t)

	q.CreateProduct(context.Background(), "Widget", 9.99)
	q.UpdateProductPrice(context.Background(), 1, 14.99)

	p, _ := q.GetProduct(context.Background(), 1)
	if p.Price != 14.99 {
		t.Errorf("price = %.2f, want 14.99", p.Price)
	}
}

func TestDeleteProduct(t *testing.T) {
	q := newTestQueries(t)

	q.CreateProduct(context.Background(), "Widget", 9.99)
	q.DeleteProduct(context.Background(), 1)

	_, err := q.GetProduct(context.Background(), 1)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestDBTXWithTx(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL
	)`)

	tx, _ := db.Begin()
	q := NewQueries(tx)

	q.CreateProduct(context.Background(), "Widget", 9.99)
	tx.Commit()

	p, _ := NewQueries(db).GetProduct(context.Background(), 1)
	if p.Name != "Widget" {
		t.Errorf("name = %q, want %q", p.Name, "Widget")
	}
}

func TestDBTXWithTxRollback(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL
	)`)

	tx, _ := db.Begin()
	q := NewQueries(tx)

	q.CreateProduct(context.Background(), "Widget", 9.99)
	tx.Rollback()

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 after rollback, got %d", count)
	}
}
