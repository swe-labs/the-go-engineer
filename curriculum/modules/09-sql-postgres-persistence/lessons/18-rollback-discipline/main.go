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

	_, err = db.Exec(`CREATE TABLE items (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		log.Fatal(err)
	}

	err = WithTx(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "Widget", 10); err != nil {
			return err
		}
		_, err := tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "Gadget", 5)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	fmt.Printf("Items inserted: %d\n", count)
}

func WithTx(db *sql.DB, fn func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
