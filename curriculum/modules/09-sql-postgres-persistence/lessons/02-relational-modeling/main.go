package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

type Order struct {
	ID     int64
	UserID int64
	Total  float64
	Status string
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			total REAL NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
	`)
	return err
}

func insertUser(db *sql.DB, u *User) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)`,
		u.Name, u.Email, u.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertOrder(db *sql.DB, o *Order) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO orders (user_id, total, status) VALUES (?, ?, ?)`,
		o.UserID, o.Total, o.Status,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	userID, err := insertUser(db, &User{Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted user with id=%d\n", userID)

	orderID, err := insertOrder(db, &Order{UserID: userID, Total: 49.99, Status: "completed"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted order with id=%d\n", orderID)

	var userName string
	var orderTotal float64
	err = db.QueryRow(`
		SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id WHERE o.id = ?
	`, orderID).Scan(&userName, &orderTotal)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User %s placed order for $%.2f\n", userName, orderTotal)
}
