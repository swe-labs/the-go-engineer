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

	_, err = db.Exec(`
		CREATE TABLE products (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			stock INTEGER NOT NULL DEFAULT 0,
			version INTEGER NOT NULL DEFAULT 1
		);
		INSERT INTO products VALUES (1, 'Widget', 9.99, 10, 1);
		INSERT INTO products VALUES (2, 'Gadget', 24.99, 5, 1);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSimpleUpdate(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec("UPDATE products SET price = ? WHERE id = ?", 12.99, 1)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 1 {
		t.Errorf("expected 1 row affected, got %d", rows)
	}

	var price float64
	db.QueryRow("SELECT price FROM products WHERE id = 1").Scan(&price)
	if price != 12.99 {
		t.Errorf("price = %.2f, want 12.99", price)
	}
}

func TestUpdateNoMatchingRow(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec("UPDATE products SET price = ? WHERE id = ?", 12.99, 999)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 0 {
		t.Errorf("expected 0 rows affected, got %d", rows)
	}
}

func TestOptimisticLockSuccess(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec(
		"UPDATE products SET stock = ?, version = version + 1 WHERE id = ? AND version = ?",
		20, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 1 {
		t.Errorf("expected 1 row affected, got %d", rows)
	}

	var version int
	db.QueryRow("SELECT version FROM products WHERE id = 1").Scan(&version)
	if version != 2 {
		t.Errorf("version = %d, want 2", version)
	}
}

func TestOptimisticLockConflict(t *testing.T) {
	db := setupTestDB(t)

	// First update succeeds
	db.Exec("UPDATE products SET stock = ?, version = version + 1 WHERE id = ? AND version = ?", 20, 1, 1)

	// Second update with stale version should fail
	result, err := db.Exec(
		"UPDATE products SET stock = ?, version = version + 1 WHERE id = ? AND version = ?",
		30, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 0 {
		t.Errorf("expected 0 rows affected for stale update, got %d", rows)
	}
}

func TestPartialUpdate(t *testing.T) {
	db := setupTestDB(t)

	update := PartialUpdate{
		Price: ptr(19.99),
		Stock: nil,
		Name:  nil,
	}
	rows, err := PartialUpdateProduct(db, 1, update)
	if err != nil {
		t.Fatal(err)
	}
	if rows != 1 {
		t.Errorf("expected 1 row, got %d", rows)
	}

	var p Product
	db.QueryRow("SELECT id, name, price, stock, version FROM products WHERE id = 1").
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Version)
	if p.Price != 19.99 {
		t.Errorf("price = %.2f, want 19.99", p.Price)
	}
	if p.Stock != 10 {
		t.Errorf("stock = %d, want 10 (unchanged)", p.Stock)
	}
}

func TestPartialUpdateNoFields(t *testing.T) {
	db := setupTestDB(t)

	update := PartialUpdate{}
	_, err := PartialUpdateProduct(db, 1, update)
	if err == nil {
		t.Fatal("expected error for no fields to update")
	}
}

func TestBatchUpdate(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.Exec("UPDATE products SET stock = stock + ? WHERE id IN (?, ?)", 5, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 2 {
		t.Errorf("expected 2 rows affected, got %d", rows)
	}

	var stock1, stock2 int
	db.QueryRow("SELECT stock FROM products WHERE id = 1").Scan(&stock1)
	db.QueryRow("SELECT stock FROM products WHERE id = 2").Scan(&stock2)
	if stock1 != 15 || stock2 != 10 {
		t.Errorf("stock1 = %d (want 15), stock2 = %d (want 10)", stock1, stock2)
	}
}
