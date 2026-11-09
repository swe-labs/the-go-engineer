package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

type User struct {
	ID    int64
	Name  string
	Email string
	Age   int
}

func createUserTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			age INTEGER NOT NULL CHECK(age >= 0 AND age <= 150)
		)
	`)
	return err
}

func insertUser(db *sql.DB, u *User) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO users (name, email, age) VALUES (?, ?, ?)`,
		u.Name, u.Email, u.Age,
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

	if err := createUserTable(db); err != nil {
		log.Fatal(err)
	}

	users := []User{
		{Name: "Alice", Email: "alice@example.com", Age: 30},
		{Name: "Bob", Email: "bob@example.com", Age: 25},
	}

	for _, u := range users {
		id, err := insertUser(db, &u)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				fmt.Printf("Skipped duplicate email: %s\n", u.Email)
				continue
			}
			log.Fatal(err)
		}
		fmt.Printf("Inserted %s with id=%d\n", u.Name, id)
	}

	_, err = db.Exec(`INSERT INTO users (name, email, age) VALUES (NULL, 'null@test.com', 20)`)
	if err != nil {
		fmt.Println("NOT NULL violation:", err)
	}

	_, err = db.Exec(`INSERT INTO users (name, email, age) VALUES ('Test', 'test@test.com', -5)`)
	if err != nil {
		fmt.Println("CHECK violation:", err)
	}

	_, err = insertUser(db, &User{Name: "Alice2", Email: "alice@example.com", Age: 35})
	if err != nil {
		fmt.Println("UNIQUE violation:", err)
	}
}
