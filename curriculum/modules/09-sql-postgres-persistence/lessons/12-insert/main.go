package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY,
		customer TEXT NOT NULL,
		amount INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending'
	)`)
	return err
}

func InsertOrder(db *sql.DB, customer string, amount int) (int64, error) {
	res, err := db.Exec(`INSERT INTO orders (customer, amount) VALUES (?, ?)`, customer, amount)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func InsertOrdersBatch(db *sql.DB, orders []struct {
	Customer string
	Amount   int
}) (int64, error) {
	if len(orders) == 0 {
		return 0, nil
	}
	placeholders := make([]string, len(orders))
	args := make([]interface{}, 0, len(orders)*2)
	for i, o := range orders {
		placeholders[i] = "(?, ?)"
		args = append(args, o.Customer, o.Amount)
	}
	query := fmt.Sprintf("INSERT INTO orders (customer, amount) VALUES %s", strings.Join(placeholders, ", "))
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := CreateTable(db); err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(`INSERT INTO orders (customer, amount) VALUES (?, ?)`, "Alice", 100)
	if err != nil {
		log.Fatal(err)
	}
	aliceID, _ := res.LastInsertId()
	fmt.Printf("Alice order inserted with id=%d\n", aliceID)

	stmt, err := db.Prepare(`INSERT INTO orders (customer, amount) VALUES (?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	extraOrders := []struct {
		customer string
		amount   int
	}{
		{"Bob", 200},
		{"Carol", 300},
	}
	for _, o := range extraOrders {
		res, err := stmt.Exec(o.customer, o.amount)
		if err != nil {
			log.Fatal(err)
		}
		id, _ := res.LastInsertId()
		fmt.Printf("%s order inserted with id=%d\n", o.customer, id)
	}

	batchOrders := []struct {
		Customer string
		Amount   int
	}{
		{"Dave", 400},
		{"Eve", 500},
	}
	count, err := InsertOrdersBatch(db, batchOrders)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Batch inserted %d orders\n", count)

	rows, err := db.Query(`SELECT id, customer, amount, status FROM orders ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("All orders:")
	for rows.Next() {
		var id int
		var customer, status string
		var amount int
		if err := rows.Scan(&id, &customer, &amount, &status); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %d: %s ($%d) [%s]\n", id, customer, amount, status)
	}
}
