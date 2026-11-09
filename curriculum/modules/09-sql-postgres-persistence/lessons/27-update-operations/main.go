package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Product struct {
	ID      int
	Name    string
	Price   float64
	Stock   int
	Version int
}

type PartialUpdate struct {
	Name  *string
	Price *float64
	Stock *int
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		log.Fatal(err)
	}

	result, err := db.Exec("UPDATE products SET price = ? WHERE id = ?", 12.99, 1)
	if err != nil {
		log.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("Simple update affected %d row(s)\n", rows)

	expectedVersion := 1
	result, err = db.Exec(
		"UPDATE products SET stock = ?, version = version + 1 WHERE id = ? AND version = ?",
		15, 1, expectedVersion)
	if err != nil {
		log.Fatal(err)
	}
	rows, _ = result.RowsAffected()
	if rows == 0 {
		fmt.Println("Optimistic lock failed")
	} else {
		fmt.Println("Optimistic lock succeeded")
	}

	update := PartialUpdate{
		Name:  nil,
		Price: ptr(19.99),
		Stock: ptr(20),
	}
	rowsAffected, err := PartialUpdateProduct(db, 2, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Partial update affected %d row(s)\n", rowsAffected)

	printProduct(db, 1)
	printProduct(db, 2)
}

func PartialUpdateProduct(db *sql.DB, id int, p PartialUpdate) (int64, error) {
	query := "UPDATE products SET "
	var updates []string
	var args []any

	if p.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *p.Name)
	}
	if p.Price != nil {
		updates = append(updates, "price = ?")
		args = append(args, *p.Price)
	}
	if p.Stock != nil {
		updates = append(updates, "stock = ?")
		args = append(args, *p.Stock)
	}

	if len(updates) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	query += join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func join(items []string, sep string) string {
	result := ""
	for i, s := range items {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func ptr[T any](v T) *T { return &v }

func printProduct(db *sql.DB, id int) {
	var p Product
	err := db.QueryRow("SELECT id, name, price, stock, version FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Product %d: %s $%.2f stock=%d version=%d\n", p.ID, p.Name, p.Price, p.Stock, p.Version)
}
