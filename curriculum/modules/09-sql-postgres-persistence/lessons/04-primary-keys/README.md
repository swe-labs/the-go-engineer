# Primary keys

## Learning objective

Choose and implement primary keys in SQL tables — auto-increment integers, UUIDs, or composite keys — and write Go code that correctly inserts, retrieves, and references rows by primary key.

## Why this matters

Every row in a table needs a unique identifier. Without one, you cannot reliably update, delete, or reference a specific record. The choice of primary key strategy (auto-increment `int`, UUID, natural key) affects join performance, distributed system compatibility, and API design. Go engineers must understand primary keys to build correct APIs that expose resource identifiers and maintain referential integrity across tables.

## Mental model

A primary key is the row's passport. It uniquely identifies that row among all others in the table. Like a passport number, it should never change and should be assigned once at creation time. When another table references a row (via foreign key), it copies the primary key value — not the person's name or other mutable attribute.

## Core idea

A **primary key** is a column (or set of columns) that uniquely identifies each row in a table. It must satisfy:

- **Uniqueness**: no two rows have the same primary key value.
- **Non-null**: the primary key column(s) cannot contain `NULL`.
- **Immutable**: the primary key should never change once assigned.

SQLite enforces these automatically when you write `INTEGER PRIMARY KEY`. Any column (or combination) declared `PRIMARY KEY` gets a unique index.

Common primary key strategies:

| Strategy | Example | Pros | Cons |
|---|---|---|---|
| Auto-increment integer | `id INTEGER PRIMARY KEY` | Fast, small (8 bytes), simple | Hard to merge databases, exposes row count |
| UUID | `id TEXT PRIMARY KEY` | Globally unique, merge-friendly | Larger (36 bytes as string), slower index |
| Natural key | `isbn TEXT PRIMARY KEY` | Meaningful, no extra column | Can change, may be long or complex |
| Composite | `(user_id, role_id) PRIMARY KEY` | Natural for join tables | Larger index, more complex joins |

## Under the hood

In SQLite, declaring `INTEGER PRIMARY KEY` makes the column an alias for the **rowid**, a 64-bit integer that SQLite uses internally to identify each row in a B-tree. Inserting a row without specifying the `id` column auto-assigns the next unused rowid. Specifying a value explicitly overrides the auto-assignment.

UUIDs are not native to SQLite. They are stored as `TEXT` (or `BLOB` for compactness). The Go program generates the UUID and inserts it as a string. SQLite creates a B-tree index on the UUID column for fast lookups.

## How Go uses it

Go's `database/sql` retrieves primary key values after INSERT via `Result.LastInsertId()`. This works for `INTEGER PRIMARY KEY` in SQLite. For UUIDs, the Go code generates the UUID (using `github.com/google/uuid`) and inserts it as a string. The `LastInsertId` method is not meaningful for UUID primary keys.

## Go example

```go
package main

import (
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
}
```

## Step-by-step execution

1. Two tables are created. `authors.id INTEGER PRIMARY KEY` is an auto-increment integer primary key. `books.id` is the same. `books.isbn` has a `UNIQUE` constraint (a natural-key candidate).
2. An author row is inserted without specifying `id`. SQLite assigns the next rowid — in this case `1`.
3. A book row is inserted. The `author_id` column references the author's primary key. This links the book to the author.
4. The JOIN query uses primary keys (`a.id = b.author_id`) to retrieve related data efficiently.

## Common mistakes

- **Using a mutable natural key as primary key**. Example: using `email` as the primary key for a `users` table. If the user changes their email, all referencing rows must be updated. Fix: use an auto-increment integer or UUID as the surrogate key and put a `UNIQUE` constraint on `email`.
- **Forgetting that `INTEGER PRIMARY KEY` is special in SQLite**. Only `INTEGER` (exact word) aliases rowid. `INT PRIMARY KEY` creates a regular non-rowid integer column, which does not auto-increment.
- **Relying on `LastInsertId` for non-integer keys**. For UUID primary keys, read the generated value back with a `RETURNING` clause (SQLite 3.35+) or by selecting after insert.
- **Assuming primary key values are sequential without gaps**. Deleting rows creates gaps. SQLite may reuse rowids from deleted rows.

## Debugging walkthrough

```go
res, err := db.Exec(`INSERT INTO authors (name) VALUES (?)`, "Tolstoy")
if err != nil {
	log.Fatal(err)
}
id, _ := res.LastInsertId()
fmt.Println(id) // always 0
```

**Symptom**: LastInsertId returns 0.

**Root cause**: The INSERT might have failed silently (if `err` is not nil) or the table has an `INT` (not `INTEGER`) primary key. Only `INTEGER PRIMARY KEY` aliases rowid.

**Fix**: Check the error before using `LastInsertId`, and verify the column definition is exactly `INTEGER PRIMARY KEY`.

## Production notes

- Use `INTEGER PRIMARY KEY` for SQLite prototypes; switch to `UUID` or `BIGSERIAL` for PostgreSQL production schemas that may need to merge databases or shard.
- Expose public-facing IDs (e.g., in URLs) as UUIDs or hashids, never expose auto-increment integers to prevent enumeration attacks.
- In distributed systems, UUID v7 (time-ordered) is preferred over UUID v4 because it keeps B-tree indexes clustered and avoids page splits.
- Set `PRAGMA foreign_keys = ON` at connection time in SQLite to enforce foreign key constraints.

## Performance implications

- Integer primary key lookups are O(log n) in the B-tree and very fast (typically 1–3 B-tree page traversals).
- UUID primary key lookups are also O(log n), but the index is larger and comparison is slower than integer comparison. On a table with millions of rows, UUID primary keys can be 2-3x slower for range scans.
- Composite primary keys create a multi-column B-tree index. The order of columns matters: put the most selective column first.

## Practice task

Create a `courses` table with a UUID primary key (stored as TEXT). Also create an `enrollments` table with a composite primary key `(course_id, student_id)`. Write Go code to:
1. Insert two courses (generate UUIDs in Go using a simple random hex string).
2. Insert an enrollment linking a student to a course.
3. Query all enrollments for a given course using the composite key.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/04-primary-keys
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/04-primary-keys
```

## Review questions

1. Why is `INT PRIMARY KEY` not the same as `INTEGER PRIMARY KEY` in SQLite?
2. What is the difference between a surrogate key and a natural key?
3. Why should auto-increment integer IDs not be used in public-facing URLs?
4. What is a composite primary key and when would you use one?
5. How does `LastInsertId` behave for a UUID primary key?

## NEXT UP

Foreign keys — connecting tables and enforcing referential integrity.
