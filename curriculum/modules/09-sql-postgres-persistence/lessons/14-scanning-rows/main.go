package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Employee struct {
	ID     int
	Name   string
	Salary float64
	Active bool
}

type Person struct {
	ID    int
	Name  string
	Email sql.NullString
	Age   sql.NullInt64
}

func scanEmployee(row *sql.Row) (Employee, error) {
	var e Employee
	err := row.Scan(&e.ID, &e.Name, &e.Salary, &e.Active)
	return e, err
}

func scanEmployees(rows *sql.Rows) ([]Employee, error) {
	defer rows.Close()
	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Salary, &e.Active); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, rows.Err()
}

func scanPerson(row *sql.Row) (Person, error) {
	var p Person
	err := row.Scan(&p.ID, &p.Name, &p.Email, &p.Age)
	return p, err
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE employees (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		salary REAL NOT NULL,
		active INTEGER NOT NULL DEFAULT 1
	)`)

	db.Exec(`INSERT INTO employees (name, salary, active) VALUES ('Alice', 75000, 1)`)
	db.Exec(`INSERT INTO employees (name, salary, active) VALUES ('Bob', 82000, 1)`)
	db.Exec(`INSERT INTO employees (name, salary, active) VALUES ('Carol', 0, 0)`)

	fmt.Println("=== Scan into struct fields ===")
	var e Employee
	err = db.QueryRow(`SELECT id, name, salary, active FROM employees WHERE id = 1`).Scan(
		&e.ID, &e.Name, &e.Salary, &e.Active,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Employee: %d: %s ($%.0f, active=%v)\n", e.ID, e.Name, e.Salary, e.Active)

	fmt.Println("\n=== Scan multiple rows with helper ===")
	rows, err := db.Query(`SELECT id, name, salary, active FROM employees ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}

	employees, err := scanEmployees(rows)
	if err != nil {
		log.Fatal(err)
	}
	for _, emp := range employees {
		fmt.Printf("  %d: %s ($%.0f, active=%v)\n", emp.ID, emp.Name, emp.Salary, emp.Active)
	}

	fmt.Println("\n=== NULL handling demo ===")
	db.Exec(`CREATE TABLE IF NOT EXISTS contacts (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT
	)`)
	db.Exec(`INSERT INTO contacts (name, phone) VALUES ('Dave', '555-0100')`)
	db.Exec(`INSERT INTO contacts (name, phone) VALUES ('Eve', NULL)`)

	rows, err = db.Query(`SELECT id, name, phone FROM contacts ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var phone sql.NullString
		if err := rows.Scan(&id, &name, &phone); err != nil {
			log.Fatal(err)
		}
		if phone.Valid {
			fmt.Printf("  %d: %s (phone=%s)\n", id, name, phone.String)
		} else {
			fmt.Printf("  %d: %s (no phone)\n", id, name)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Person scan with NULL handling ===")
	db.Exec(`CREATE TABLE IF NOT EXISTS people (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT,
		age INTEGER
	)`)
	db.Exec(`INSERT INTO people (name, email, age) VALUES ('Alice', 'alice@example.com', 30)`)
	db.Exec(`INSERT INTO people (name, email, age) VALUES ('Bob', NULL, NULL)`)

	for i := 1; i <= 2; i++ {
		p, err := scanPerson(db.QueryRow(`SELECT id, name, email, age FROM people WHERE id = ?`, i))
		if err != nil {
			log.Fatal(err)
		}
		if p.Email.Valid {
			fmt.Printf("  %d: %s (email=%s, age=", p.ID, p.Name, p.Email.String)
		} else {
			fmt.Printf("  %d: %s (email=none, age=", p.ID, p.Name)
		}
		if p.Age.Valid {
			fmt.Printf("%d)\n", p.Age.Int64)
		} else {
			fmt.Printf("unknown)\n")
		}
	}
}
