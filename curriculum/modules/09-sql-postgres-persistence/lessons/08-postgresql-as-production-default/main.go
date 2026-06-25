package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func ConnectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	return db, nil
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = ":memory:"
	}

	db, err := ConnectDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Connected to database")

	db.Exec(`CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, name TEXT)`)

	_, err = db.Exec(`INSERT INTO items (name) VALUES ('production item')`)
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = db.QueryRow(`SELECT name FROM items WHERE id = 1`).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Queried: %s\n", name)

	fmt.Println("\n--- PostgreSQL connection string example ---")
	fmt.Println("Driver: pgx")
	fmt.Println("DSN:    postgres://user:password@localhost:5432/mydb?sslmode=disable")
}
