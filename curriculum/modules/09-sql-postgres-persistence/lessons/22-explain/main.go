package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		);
		INSERT INTO users VALUES (1, 'Alice', 'alice@example.com');
		INSERT INTO users VALUES (2, 'Bob', 'bob@example.com');
		INSERT INTO users VALUES (3, 'Carol', 'carol@example.com');
	`)
	if err != nil {
		log.Fatal(err)
	}

	queries := []string{
		"SELECT * FROM users WHERE id = 1",
		"SELECT * FROM users WHERE email = 'alice@example.com'",
		"SELECT * FROM users ORDER BY name",
	}

	fmt.Println("=== Without Indexes ===")
	for _, q := range queries {
		fmt.Printf("\nQuery: %s\n", q)
		printPlan(db, q)
	}

	db.Exec("CREATE INDEX idx_users_email ON users(email)")
	db.Exec("CREATE INDEX idx_users_name ON users(name)")

	fmt.Println("\n=== With Indexes ===")
	for _, q := range queries {
		fmt.Printf("\nQuery: %s\n", q)
		printPlan(db, q)
	}
}

func printPlan(db *sql.DB, query string) {
	rows, err := db.Query("EXPLAIN QUERY PLAN " + query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, parent, notused int
		var detail string
		if err := rows.Scan(&id, &parent, &notused, &detail); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  id=%d parent=%d %s\n", id, parent, detail)
	}
}
