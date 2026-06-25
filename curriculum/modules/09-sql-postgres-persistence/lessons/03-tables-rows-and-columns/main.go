package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Product struct {
	ID    int64
	Name  string
	Price float64
	Stock int
}

func createProductTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price REAL NOT NULL CHECK(price >= 0),
			stock INTEGER NOT NULL DEFAULT 0
		)
	`)
	return err
}

func insertProduct(db *sql.DB, p *Product) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO products (name, price, stock) VALUES (?, ?, ?)`,
		p.Name, p.Price, p.Stock,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func queryProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, price, stock FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createProductTable(db); err != nil {
		log.Fatal(err)
	}

	for _, p := range []Product{
		{Name: "Widget", Price: 9.99, Stock: 100},
		{Name: "Gadget", Price: 24.99, Stock: 50},
	} {
		id, err := insertProduct(db, &p)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted %s with id=%d\n", p.Name, id)
	}

	products, err := queryProducts(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Queried %d product(s):\n", len(products))
	for _, p := range products {
		fmt.Printf("  %d: %s ($%.2f, stock=%d)\n", p.ID, p.Name, p.Price, p.Stock)
	}
}
