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
	if err := createProductTable(db); err != nil {
		t.Fatal(err)
	}
}

func TestInsertAndQuery(t *testing.T) {
	db := openMemory(t)
	if err := createProductTable(db); err != nil {
		t.Fatal(err)
	}

	id, err := insertProduct(db, &Product{Name: "Test", Price: 1.99, Stock: 10})
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("got id %d, want 1", id)
	}

	products, err := queryProducts(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 1 {
		t.Fatalf("got %d products, want 1", len(products))
	}
	if products[0].Name != "Test" {
		t.Fatalf("got name %q, want %q", products[0].Name, "Test")
	}
}

func TestDefaultStock(t *testing.T) {
	db := openMemory(t)
	if err := createProductTable(db); err != nil {
		t.Fatal(err)
	}

	_, err := db.Exec(`INSERT INTO products (name, price) VALUES (?, ?)`, "NoStock", 5.00)
	if err != nil {
		t.Fatal(err)
	}

	var stock int
	err = db.QueryRow(`SELECT stock FROM products WHERE name = ?`, "NoStock").Scan(&stock)
	if err != nil {
		t.Fatal(err)
	}
	if stock != 0 {
		t.Fatalf("expected default stock 0, got %d", stock)
	}
}

func TestCheckConstraint(t *testing.T) {
	db := openMemory(t)
	if err := createProductTable(db); err != nil {
		t.Fatal(err)
	}

	_, err := db.Exec(`INSERT INTO products (name, price) VALUES (?, ?)`, "Negative", -1.00)
	if err == nil {
		t.Fatal("expected CHECK constraint to reject negative price")
	}
}

func TestMultipleProducts(t *testing.T) {
	db := openMemory(t)
	if err := createProductTable(db); err != nil {
		t.Fatal(err)
	}

	names := []string{"A", "B", "C"}
	for _, n := range names {
		_, err := insertProduct(db, &Product{Name: n, Price: 1.00, Stock: 1})
		if err != nil {
			t.Fatal(err)
		}
	}

	products, err := queryProducts(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 3 {
		t.Fatalf("got %d products, want 3", len(products))
	}
}
