package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func SaveAndLoad(db *sql.DB, msg string) (int, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY, content TEXT NOT NULL)`)
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(`INSERT INTO notes (content) VALUES (?)`, msg)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (id INTEGER NOT NULL PRIMARY KEY, text TEXT NOT NULL);`)
	if err != nil {
		log.Fatalf("create table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO messages (text) VALUES (?)`, "Hello from SQLite!")
	if err != nil {
		log.Fatalf("insert: %v", err)
	}

	var id int
	var text string
	err = db.QueryRow(`SELECT id, text FROM messages WHERE id = ?`, 1).Scan(&id, &text)
	if err != nil {
		log.Fatalf("query: %v", err)
	}
	fmt.Printf("Row: id=%d, text=%q\n", id, text)

	id2, err := SaveAndLoad(db, "practice note")
	if err != nil {
		log.Fatalf("SaveAndLoad: %v", err)
	}
	fmt.Printf("Saved note with id=%d\n", id2)
}
