# Why databases exist

## Learning objective

Explain why applications need persistent storage, contrast flat-file storage with a relational database across correctness, concurrency, and consistency dimensions, and write a Go program that uses SQLite to persist and query data.

## Why this matters

Every production Go application outlives a single process. If you restart a web server and all user data disappears, the application is useless. Databases exist to solve three fundamental problems that flat files cannot: concurrent access without corruption, crash recovery, and efficient querying at scale. Go engineers reach for databases whether they build a CLI tool, a REST API, or a streaming pipeline. Understanding the "why" behind databases informs every later decision about schema design, driver selection, transaction boundaries, and connection pooling.

## Mental model

Think of a flat file as a notebook that only one person can write in at a time. If two people write simultaneously, pages get lost or overwritten. A database is like a librarian: it ensures only one person modifies a page at a time, keeps a log of every change, and can find any page in milliseconds even when the archive holds millions of records. The librarian also guarantees that a partial update (e.g., transferring money from account A to B) either completes fully or leaves both accounts unchanged.

## Core idea

A database is a managed, durable, concurrent storage system that guarantees:

- **Persistence**: data survives process restarts and power loss.
- **Atomicity**: a set of writes either all happen or none happen.
- **Consistency**: data obeys application-defined rules (constraints).
- **Isolation**: concurrent transactions do not interfere with each other.
- **Durability**: completed writes survive crashes.

These properties are collectively known as **ACID**. Flat files provide none of them by default.

## Under the hood

When a Go program writes to a file using `os.WriteFile`, the operating system may buffer the write in memory for seconds before flushing to disk. If the process crashes during that window, the write is lost. Concurrent writes from goroutines or separate processes interleave bytes, corrupting the file. Even a single `Write` call that spans multiple disk blocks can be partially applied after a power failure.

A database engine (SQLite, PostgreSQL) wraps the raw file descriptor with:

1. **Write-ahead log (WAL)**: every mutation is first appended to a sequential log file. If the process crashes, the log is replayed on restart to restore a consistent state.
2. **Page cache**: data pages are cached in memory and flushed to disk at checkpoint boundaries controlled by the engine.
3. **Lock manager**: readers and writers coordinate via fine-grained locks (page-level, row-level) to prevent corruption.
4. **Transaction log**: a commit record is written to the WAL only after all changes for a transaction are safely logged. This guarantees atomicity and durability.

SQLite, which we use in this module's examples, implements all of these in a single C library (compiled to pure Go via `modernc.org/sqlite`). PostgreSQL does the same with a client-server architecture.

## How Go uses it

Go's standard library provides `database/sql`, a generic interface over SQL database drivers. A Go program:

1. Registers a driver (e.g., `modernc.org/sqlite`, `github.com/jackc/pgx/v5/stdlib`).
2. Opens a connection pool with `sql.Open`.
3. Executes SQL statements with `db.Exec`, `db.Query`, or `db.QueryRow`.
4. Scans result rows into Go structs with `rows.Scan`.

Go does not include a built-in database; instead, it exposes a seam that any compliant driver can implement. This means the same `database/sql` API works for SQLite, PostgreSQL, MySQL, and many others.

## Go example

```go
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

	sqlStmt := `CREATE TABLE IF NOT EXISTS messages (id INTEGER NOT NULL PRIMARY KEY, text TEXT NOT NULL);`
	_, err = db.Exec(sqlStmt)
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
}
```

## Step-by-step execution

1. `sql.Open("sqlite", ":memory:")` creates an in-memory SQLite database — no file is written to disk. The driver is registered by the blank import `_ "modernc.org/sqlite"`.
2. `db.Exec(createSQL)` executes the CREATE TABLE statement. The database creates the schema in memory.
3. `db.Exec(insertSQL, "Hello from SQLite!")` inserts one row. The `?` is a parameter marker — SQLite binds the value safely.
4. `db.QueryRow(selectSQL, 1).Scan(&id, &text)` retrieves the row where `id = 1` and maps the columns to Go variables `id` and `text`.
5. The program prints the result and exits. Because the database is in-memory, all data is lost when the connection closes — this demonstrates why persistence to disk files (or a remote database) matters for production data.

## Common mistakes

- **Treating `sql.Open` as a connection attempt**: `sql.Open` only validates the driver name and data source name. It does not connect. Use `db.Ping()` to verify the database is reachable.
- **Forgetting to close the database**: `defer db.Close()` should always follow a successful `sql.Open`. Leaking connections exhausts file descriptors.
- **Assuming files are ACID by default**: Writing to a file with `os.WriteFile` does not guarantee atomicity or durability. A crash during a write can produce a corrupted file.
- **Using string formatting instead of parameterized queries**: `fmt.Sprintf("INSERT INTO t VALUES (%d)", id)` opens the door to SQL injection. Always use `?` parameter markers.

## Debugging walkthrough

Consider this broken program:

```go
db, err := sql.Open("sqlite", ":memory:")
if err != nil {
	log.Fatal(err)
}
db.Close()
_, err = db.Exec(`CREATE TABLE t (id INT)`)
if err != nil {
	log.Fatal(err)
}
```

**Symptom**: The `db.Exec` call fails with "database closed" or "sql: database is closed".

**Root cause**: `db.Close()` was called before using the database. `sql.DB` is a pool; closing it makes all future operations fail.

**Fix**: Use `defer db.Close()` after the initial `sql.Open`, not before executing statements.

```go
db, err := sql.Open("sqlite", ":memory:")
if err != nil {
	log.Fatal(err)
}
defer db.Close()

_, err = db.Exec(`CREATE TABLE t (id INT)`)
if err != nil {
	log.Fatal(err)
}
```

## Production notes

- In production, never use `:memory:` for data you cannot afford to lose. The in-memory mode is for testing and prototyping.
- For file-backed SQLite, use a path like `./data.db`. The database engine handles crash recovery via the WAL.
- For PostgreSQL, the data source name includes host, port, user, password, dbname, and sslmode: `postgres://user:pass@localhost:5432/mydb?sslmode=disable`.
- Always set connection pool limits via `db.SetMaxOpenConns` and `db.SetMaxIdleConns` before issuing queries (covered in detail in lesson 11).
- Monitor `sql.DB.Stats()` in production to detect pool exhaustion.

## Performance implications

- SQLite in-memory mode is extremely fast (microsecond latencies for simple queries) but uses process memory. File-backed SQLite adds disk I/O latency (typically 1–10 ms per write).
- PostgreSQL adds network round-trip latency (sub-millisecond on localhost, 1–30 ms over a network).
- Parameterized queries (`?` markers) allow the database to cache query plans, improving execution time on repeated queries.
- The `database/sql` interface introduces a small overhead per call compared to raw driver usage, but the abstraction is worth the cost for portability and maintainability.

## Practice task

Write a function `SaveAndLoad(db *sql.DB, msg string) (int, error)` that:
1. Inserts `msg` into a `notes` table (id INTEGER PRIMARY KEY, content TEXT).
2. Returns the `id` of the inserted row.
3. Creates the table if it does not exist (use `CREATE TABLE IF NOT EXISTS`).

Then write a `main()` that opens an in-memory SQLite database, calls `SaveAndLoad` twice with different messages, queries all rows, and prints them.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/01-why-databases-exist
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/01-why-databases-exist
```

The tests verify that a database can be opened, queried, and that concurrent operations do not corrupt data.

## Review questions

1. What are the five ACID properties and which one does a flat-file `Write` call violate most obviously?
2. Why does `sql.Open("sqlite", ":memory:")` not connect to the database immediately?
3. What is the purpose of a write-ahead log (WAL)?
4. Name two things a database provides that a flat file does not.
5. What happens if you call `db.Close()` and then try to `db.Exec(...)`?

## NEXT UP

Relational modeling — how to organize data into entities, attributes, and relationships before writing SQL.
