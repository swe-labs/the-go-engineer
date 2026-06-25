package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Product struct {
	ID          int
	Name        string
	Description sql.NullString
	Price       float64
	Stock       sql.NullInt64
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		stock INTEGER
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO products VALUES
		(1, 'Widget', 'A useful widget', 9.99, 100),
		(2, 'Gadget', NULL, 24.99, 0),
		(3, 'Doohickey', 'Premium doohickey', 49.99, NULL)`)
	if err != nil {
		log.Fatal(err)
	}

	products, err := GetProducts(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range products {
		b, _ := json.MarshalIndent(p, "", "  ")
		fmt.Printf("%s\n", b)
	}
}

func GetProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, description, price, stock FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func NonNullProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, description, price, stock FROM products WHERE description IS NOT NULL ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}
