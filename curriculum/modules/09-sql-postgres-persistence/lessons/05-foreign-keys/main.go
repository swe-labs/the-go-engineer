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

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE departments (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
		CREATE TABLE employees (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			dept_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO departments (name) VALUES ('Engineering')`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES ('Alice', 1)`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted department and employee successfully")

	_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES ('Bob', 999)`)
	if err != nil {
		fmt.Println("Foreign key violation caught:", err)
	}

	_, err = db.Exec(`DELETE FROM departments WHERE id = 1`)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM employees`).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Employees remaining after cascade delete: %d\n", count)
}
