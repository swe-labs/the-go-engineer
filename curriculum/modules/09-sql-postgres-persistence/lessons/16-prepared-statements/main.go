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
		price REAL NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	products := []Product{
		{Name: "Widget", Price: 9.99},
		{Name: "Gadget", Price: 24.99},
		{Name: "Doohickey", Price: 49.99},
		{Name: "Thingamajig", Price: 14.99},
		{Name: "Whatsit", Price: 3.99},
	}

	if err := BatchInsertProducts(db, products); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, name, price FROM products ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Inserted products:")
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %d: %s ($%.2f)\n", p.ID, p.Name, p.Price)
	}
}

func BatchInsertProducts(db *sql.DB, products []Product) error {
	stmt, err := db.Prepare("INSERT INTO products (name, price) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range products {
		if _, err := stmt.Exec(p.Name, p.Price); err != nil {
			return err
		}
	}
	return nil
}
