package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	memDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer memDB.Close()

	memDB.Exec(`CREATE TABLE memo (id INTEGER PRIMARY KEY, note TEXT)`)
	memDB.Exec(`INSERT INTO memo (note) VALUES ('in-memory note')`)

	var note string
	memDB.QueryRow(`SELECT note FROM memo WHERE id = 1`).Scan(&note)
	fmt.Println("Memory DB:", note)

	tmpDir := os.TempDir()
	dbPath := filepath.Join(tmpDir, "sqlite_example.db")
	defer os.Remove(dbPath)

	fileDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileDB.Close()

	fileDB.Exec(`CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, title TEXT)`)
	fileDB.Exec(`INSERT INTO tasks (title) VALUES ('persistent task')`)
	fileDB.QueryRow(`SELECT title FROM tasks WHERE id = 1`).Scan(&note)
	fmt.Println("File DB:", note)

	fileDB.Exec(`PRAGMA journal_mode=WAL`)

	var version string
	memDB.QueryRow(`SELECT sqlite_version()`).Scan(&version)
	fmt.Println("SQLite version:", version)
}
