package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

func SearchProducts(db *sql.DB, minPrice, maxPrice float64) ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, price, stock FROM products WHERE price >= ? AND price <= ? ORDER BY price`, minPrice, maxPrice)
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

func queryAllProducts(db *sql.DB) ([]Product, error) {
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

func queryProductByID(db *sql.DB, id int) (*Product, error) {
	var p Product
	err := db.QueryRow(`SELECT id, name, price, stock FROM products WHERE id = ?`, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0
	)`)

	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Widget', 9.99, 100)`)
	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Gadget', 24.99, 50)`)
	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Gizmo', 14.99, 0)`)

	fmt.Println("=== Query: all products ===")
	products, err := queryAllProducts(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range products {
		fmt.Printf("  %d: %s ($%.2f, stock=%d)\n", p.ID, p.Name, p.Price, p.Stock)
	}

	fmt.Println("\n=== Query: single product by ID ===")
	p, err := queryProductByID(db, 2)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("  Product not found")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("  %d: %s ($%.2f, stock=%d)\n", p.ID, p.Name, p.Price, p.Stock)
	}

	fmt.Println("\n=== Query: in-stock products ===")
	rows, err := db.Query(`SELECT id, name, price, stock FROM products WHERE stock > 0 ORDER BY price`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: $%.2f (%d in stock)\n", p.Name, p.Price, p.Stock)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Search: products $10-$30 ===")
	results, err := SearchProducts(db, 10.00, 30.00)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d product(s):\n", len(results))
	for _, p := range results {
		fmt.Printf("  %s ($%.2f)\n", p.Name, p.Price)
	}
}
