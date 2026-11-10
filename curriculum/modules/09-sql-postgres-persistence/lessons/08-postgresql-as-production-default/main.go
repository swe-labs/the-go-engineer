package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// ConnectDB selects the driver from the DSN shape: PostgreSQL URLs use the
// pgx stdlib driver ("pgx"); anything else falls back to SQLite.
func ConnectDB(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
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
	driver := "sqlite"
	if dsn == "" {
		dsn = ":memory:"
	}
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		driver = "pgx"
	}

	db, err := ConnectDB(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Connected to database")

	// DDL differs by dialect: SQLite auto-generates INTEGER PRIMARY KEY from
	// the rowid; PostgreSQL needs an explicit identity (SERIAL).
	createSQL := `CREATE TABLE IF NOT EXISTS items (id SERIAL PRIMARY KEY, name TEXT)`
	if driver == "sqlite" {
		createSQL = `CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, name TEXT)`
	}
	db.Exec(createSQL)

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
