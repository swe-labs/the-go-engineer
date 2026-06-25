package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			isbn TEXT NOT NULL UNIQUE,
			author_id INTEGER NOT NULL REFERENCES authors(id)
		)
	`)
	return err
}

type Author struct {
	ID   int64
	Name string
}

type Book struct {
	ID       int64
	Title    string
	ISBN     string
	AuthorID int64
}

func insertAuthor(db *sql.DB, a *Author) (int64, error) {
	res, err := db.Exec(`INSERT INTO authors (name) VALUES (?)`, a.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertBook(db *sql.DB, b *Book) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO books (title, isbn, author_id) VALUES (?, ?, ?)`,
		b.Title, b.ISBN, b.AuthorID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	authorID, err := insertAuthor(db, &Author{Name: "Jane Austen"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Author inserted with pk=%d\n", authorID)

	bookID, err := insertBook(db, &Book{Title: "Pride and Prejudice", ISBN: "9780141439518", AuthorID: authorID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Book inserted with pk=%d\n", bookID)

	var name string
	var title string
	err = db.QueryRow(`
		SELECT a.name, b.title FROM authors a JOIN books b ON a.id = b.author_id WHERE b.id = ?
	`, bookID).Scan(&name, &title)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s wrote %s\n", name, title)

	// Example with UUID primary key
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS courses (id TEXT PRIMARY KEY, name TEXT NOT NULL)`)
	if err != nil {
		log.Fatal(err)
	}
	courseID := newUUID()
	_, err = db.Exec(`INSERT INTO courses (id, name) VALUES (?, ?)`, courseID, "Go 101")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Course inserted with uuid pk=%s\n", courseID)
}
