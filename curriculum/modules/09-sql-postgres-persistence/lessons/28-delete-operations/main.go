package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("PRAGMA foreign_keys = ON")

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			deleted_at TIMESTAMP
		);
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			total REAL NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		INSERT INTO users VALUES (1, 'Alice', NULL);
		INSERT INTO users VALUES (2, 'Bob', NULL);
		INSERT INTO users VALUES (3, 'Carol', NULL);
		INSERT INTO orders VALUES (101, 1, 50.00, '2024-01-15');
		INSERT INTO orders VALUES (102, 1, 30.00, '2024-06-20');
		INSERT INTO orders VALUES (103, 3, 20.00, '2024-03-10');
	`)
	if err != nil {
		log.Fatal(err)
	}

	result, _ := db.Exec("UPDATE users SET deleted_at = ? WHERE id = ?", time.Now().UTC(), 2)
	rows, _ := result.RowsAffected()
	fmt.Printf("Soft deleted %d user(s)\n", rows)

	countActive(db)
	countDeleted(db)

	result, _ = db.Exec("DELETE FROM orders WHERE user_id = ? AND total < ?", 1, 40.00)
	rows, _ = result.RowsAffected()
	fmt.Printf("Hard deleted %d order(s) for user 1 with total < 40\n", rows)

	queryOrders(db, 1)

	result, _ = db.Exec("UPDATE users SET deleted_at = NULL WHERE id = ?", 2)
	rows, _ = result.RowsAffected()
	fmt.Printf("Restored %d user(s)\n", rows)

	countActive(db)
}

func countActive(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&count)
	fmt.Printf("Active users: %d\n", count)
}

func countDeleted(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NOT NULL").Scan(&count)
	fmt.Printf("Deleted users: %d\n", count)
}

func queryOrders(db *sql.DB, userID int) {
	rows, err := db.Query("SELECT id, total FROM orders WHERE user_id = ?", userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Printf("Orders for user %d:\n", userID)
	for rows.Next() {
		var id int
		var total float64
		if err := rows.Scan(&id, &total); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  Order %d: $%.2f\n", id, total)
	}
}
