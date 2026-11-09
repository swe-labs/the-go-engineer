package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func WaitForDB(dsn string, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("sqlite", dsn)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		err = db.Ping()
		if err == nil {
			return db, nil
		}
		db.Close()
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("could not connect after %d retries: %w", maxRetries, err)
}

func main() {
	fmt.Println("--- Docker Compose for PostgreSQL ---")
	fmt.Println()
	fmt.Println("To start PostgreSQL:")
	fmt.Println("  docker compose up -d")
	fmt.Println()
	fmt.Println("Connection string:")
	fmt.Println("  postgres://app:secret@localhost:5432/myapp?sslmode=disable")
	fmt.Println()

	db, err := WaitForDB(":memory:", 3)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)`)

	users := []string{"Alice", "Bob", "Carol"}
	for _, name := range users {
		db.Exec(`INSERT INTO users (name) VALUES (?)`, name)
	}

	rows, err := db.Query(`SELECT id, name FROM users ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Users in database:")
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %d: %s\n", id, name)
	}
}
