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

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserting 10,000 users...")
	tx, _ := db.Begin()
	for i := 0; i < 10000; i++ {
		tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)",
			fmt.Sprintf("User %d", i),
			fmt.Sprintf("user%d@example.com", i))
	}
	tx.Commit()
	fmt.Println("Done inserting.")

	fmt.Println("--- Query without index ---")
	explainQueryPlan(db, "SELECT * FROM users WHERE email = 'user5000@example.com'")

	start := time.Now()
	var name string
	db.QueryRow("SELECT name FROM users WHERE email = 'user5000@example.com'").Scan(&name)
	fmt.Printf("Lookup without index: %v\n\n", time.Since(start))

	db.Exec("CREATE INDEX idx_users_email ON users(email)")
	fmt.Println("Created index on email.")

	fmt.Println("--- Query with index ---")
	explainQueryPlan(db, "SELECT * FROM users WHERE email = 'user5000@example.com'")

	start = time.Now()
	db.QueryRow("SELECT name FROM users WHERE email = 'user5000@example.com'").Scan(&name)
	fmt.Printf("Lookup with index: %v\n", time.Since(start))
}

func explainQueryPlan(db *sql.DB, query string) {
	rows, err := db.Query("EXPLAIN QUERY PLAN " + query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, parent int
		var detail string
		var notused int
		rows.Scan(&id, &parent, &notused, &detail)
		fmt.Printf("  %s\n", detail)
	}
}
