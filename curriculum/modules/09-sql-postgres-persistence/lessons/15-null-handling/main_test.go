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
		description TEXT,
		price REAL NOT NULL,
		stock INTEGER
	)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO products VALUES
		(1, 'Widget', 'A useful widget', 9.99, 100),
		(2, 'Gadget', NULL, 24.99, 0),
		(3, 'Doohickey', 'Premium doohickey', 49.99, NULL)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestGetProducts(t *testing.T) {
	db := setupTestDB(t)
	products, err := GetProducts(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 3 {
		t.Fatalf("expected 3 products, got %d", len(products))
	}

	tests := []struct {
		id         int
		name       string
		descValid  bool
		descValue  string
		stockValid bool
		stockValue int64
	}{
		{1, "Widget", true, "A useful widget", true, 100},
		{2, "Gadget", false, "", true, 0},
		{3, "Doohickey", true, "Premium doohickey", false, 0},
	}

	for _, tc := range tests {
		p := products[tc.id-1]
		if p.ID != tc.id || p.Name != tc.name {
			t.Errorf("product %d: got (%d, %s), want (%d, %s)", tc.id, p.ID, p.Name, tc.id, tc.name)
		}
		if p.Description.Valid != tc.descValid {
			t.Errorf("product %d Description.Valid = %v, want %v", tc.id, p.Description.Valid, tc.descValid)
		}
		if p.Description.Valid && p.Description.String != tc.descValue {
			t.Errorf("product %d Description.String = %q, want %q", tc.id, p.Description.String, tc.descValue)
		}
		if p.Stock.Valid != tc.stockValid {
			t.Errorf("product %d Stock.Valid = %v, want %v", tc.id, p.Stock.Valid, tc.stockValid)
		}
		if p.Stock.Valid && p.Stock.Int64 != tc.stockValue {
			t.Errorf("product %d Stock.Int64 = %d, want %d", tc.id, p.Stock.Int64, tc.stockValue)
		}
	}
}

func TestNonNullProducts(t *testing.T) {
	db := setupTestDB(t)
	products, err := NonNullProducts(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 2 {
		t.Fatalf("expected 2 non-null products, got %d", len(products))
	}
	for _, p := range products {
		if !p.Description.Valid {
			t.Errorf("product %d has invalid description, expected non-null", p.ID)
		}
	}
}

func TestNullScanIntoPlainType(t *testing.T) {
	db := setupTestDB(t)

	var stock int
	row := db.QueryRow("SELECT stock FROM products WHERE id = ?", 3)
	err := row.Scan(&stock)
	if err == nil {
		t.Error("expected error scanning NULL into int, got nil")
	}
}
