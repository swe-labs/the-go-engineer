package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type OrderWithUser struct {
	OrderID  int
	UserID   int
	Total    float64
	UserName sql.NullString
}

type UserWithOrder struct {
	UserID   int
	UserName string
	OrderID  sql.NullInt64
	Total    sql.NullFloat64
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE orders (id INTEGER PRIMARY KEY, user_id INTEGER NOT NULL, total REAL NOT NULL, FOREIGN KEY (user_id) REFERENCES users(id));
		INSERT INTO users VALUES (1, 'Alice'), (2, 'Bob'), (3, 'Carol');
		INSERT INTO orders VALUES (101, 1, 50.00), (102, 1, 30.00), (103, 3, 20.00);
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== INNER JOIN ===")
	rows, err := db.Query(`SELECT orders.id, orders.user_id, orders.total, users.name
		FROM orders INNER JOIN users ON orders.user_id = users.id ORDER BY orders.id`)
	if err != nil {
		log.Fatal(err)
	}
	printOrderWithUser(rows)

	fmt.Println("\n=== LEFT JOIN (users -> orders) ===")
	userRows, err := db.Query(`SELECT users.id, users.name, orders.id, orders.total
		FROM users LEFT JOIN orders ON users.id = orders.user_id ORDER BY users.id, orders.id`)
	if err != nil {
		log.Fatal(err)
	}
	defer userRows.Close()
	for userRows.Next() {
		var uwo UserWithOrder
		if err := userRows.Scan(&uwo.UserID, &uwo.UserName, &uwo.OrderID, &uwo.Total); err != nil {
			log.Fatal(err)
		}
		if uwo.OrderID.Valid {
			fmt.Printf("User %s (id=%d): Order %d = $%.2f\n", uwo.UserName, uwo.UserID, uwo.OrderID.Int64, uwo.Total.Float64)
		} else {
			fmt.Printf("User %s (id=%d): No orders\n", uwo.UserName, uwo.UserID)
		}
	}
}

func printOrderWithUser(rows *sql.Rows) {
	defer rows.Close()
	for rows.Next() {
		var o OrderWithUser
		if err := rows.Scan(&o.OrderID, &o.UserID, &o.Total, &o.UserName); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Order %d: user %d, $%.2f, name=%s\n", o.OrderID, o.UserID, o.Total, o.UserName.String)
	}
}
