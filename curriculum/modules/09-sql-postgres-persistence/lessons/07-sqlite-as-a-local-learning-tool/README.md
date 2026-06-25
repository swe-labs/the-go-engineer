# SQLite as a local learning tool

## Learning objective

Use SQLite as a local, zero-configuration database for development and testing, understand its in-memory and file-backed modes, and recognize its limitations compared to PostgreSQL.

## Why this matters

SQLite is the most widely deployed database engine in the world — every smartphone, browser, and many desktop applications use it. For Go developers, SQLite is the ideal learning database because it requires no installation, no daemon, no authentication, and no network configuration. You can start writing SQL in Go within seconds. Understanding SQLite's strengths and limitations helps you choose the right database for each phase of development: SQLite for prototyping and testing, PostgreSQL for production.

## Mental model

SQLite is an embedded database: it runs in the same process as your Go application. There is no separate server to start, no port to configure, and no connection string with host and port. The database is either a single file on disk (`./data.db`) or pure memory (`:memory:`). Think of it as a library (like `encoding/json`) rather than a service (like a web server).

## Core idea

SQLite is a **serverless, zero-configuration, transactional SQL database engine** delivered as a C library. In Go, it is accessed through `modernc.org/sqlite`, a pure-Go port that requires no CGO.

Key features:

| Feature | SQLite | PostgreSQL |
|---|---|---|
| Setup | Zero — just import the driver | Install, start daemon, create role/DB |
| Process model | Embedded (in-process) | Client-server (separate process) |
| Concurrency | WAL mode allows concurrent reads + one writer | Full MVCC, many concurrent writers |
| Storage | Single file or `:memory:` | Data directory managed by server |
| Network | No network, file-based | TCP connection (localhost or remote) |
| Users/auth | None (file permissions) | Role-based authentication |

## Under the hood

SQLite organizes a database file into **pages** (usually 4096 bytes). The first page is the **header**, which stores the database version, page size, and schema encoding. The rest of the file contains B-trees for each table and index.

When you open a SQLite database in WAL (Write-Ahead Log) mode, all writes go to a separate `.db-wal` file. Readers continue reading from the main file while the WAL accumulates writes. At a checkpoint, the WAL is merged back into the main database file. This allows concurrent reads during writes.

The `:memory:` database lives entirely in RAM. No file is created. It is ideal for tests because it is fast and leaves no cleanup.

## How Go uses it

Go uses SQLite through the `database/sql` interface. The import path is `modernc.org/sqlite`:

```go
import _ "modernc.org/sqlite"

db, err := sql.Open("sqlite", ":memory:")
db, err := sql.Open("sqlite", "./data.db")
```

The driver name is `"sqlite"`. The second argument is the data source name (DSN), which can be:
- `:memory:` — pure in-memory database.
- `./path/to/file.db` — file-backed database.
- `file:path.db?cache=shared&mode=rwc` — URI-style with options.

Additional PRAGMAs can be set via the DSN or executed after opening:

```go
db.Exec("PRAGMA journal_mode=WAL")
db.Exec("PRAGMA foreign_keys=ON")
```

## Go example

```go
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
	// In-memory database
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

	// File-backed database
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

	// Enable WAL mode
	fileDB.Exec(`PRAGMA journal_mode=WAL`)

	// Query SQLite version
	var version string
	memDB.QueryRow(`SELECT sqlite_version()`).Scan(&version)
	fmt.Println("SQLite version:", version)
}
```

## Step-by-step execution

1. An in-memory database is opened with `sql.Open("sqlite", ":memory:")`. No file is created.
2. A table is created and a row inserted. This data exists only while the `*sql.DB` is open.
3. A file-backed database is opened at a temporary path. A file `sqlite_example.db` is created.
4. A table is created and a row inserted. This data persists on disk after the connection closes.
5. `PRAGMA journal_mode=WAL` switches the file database to WAL mode for better concurrent read performance.
6. `SELECT sqlite_version()` queries the SQLite library version.

## Common mistakes

- **Using `:memory:` for production data**: Everything is lost on process restart. Use a file path for any data that matters.
- **Assuming file path is relative to the project root**: SQLite resolves relative paths from the working directory. Use absolute paths or `os.Getwd` plus `filepath.Join`.
- **Forgetting that each `:memory:` connection is a separate database**: Two `sql.Open("sqlite", ":memory:")` calls create two independent databases. To share, use `file::memory:?cache=shared`.
- **Not setting WAL mode for concurrent access**: The default journal mode (DELETE) locks the database file during writes, causing `database is locked` errors in concurrent Go programs.

## Debugging walkthrough

```go
db, _ := sql.Open("sqlite", "./data.db")
db.Exec(`CREATE TABLE t (id INT PRIMARY KEY)`)
// ... later in another function:
db2, _ := sql.Open("sqlite", "./data.db")
db2.Exec(`INSERT INTO t VALUES (1)`)
// Success — both connections see the same file
```

But if you use `:memory:`:

```go
db, _ := sql.Open("sqlite", ":memory:")
db.Exec(`CREATE TABLE t (id INT PRIMARY KEY)`)
db2, _ := sql.Open("sqlite", ":memory:")
db2.Exec(`INSERT INTO t VALUES (1)`)
// Error: "no such table: t"
```

**Root cause**: Each `:memory:` DSN creates a separate, private in-memory database. `db2` does not see tables created by `db`.

**Fix**: Use `file::memory:?cache=shared` for shared in-memory databases, or pass the same `*sql.DB` instance.

## Production notes

- SQLite is suitable for production in specific scenarios: embedded applications, client-side storage, single-server tools, read-heavy workloads with infrequent writes.
- For multi-server deployments, high write volume, or complex access control, use PostgreSQL.
- Always enable WAL mode for production SQLite: `PRAGMA journal_mode=WAL;`
- Set a busy timeout to avoid `database is locked` errors: `PRAGMA busy_timeout=5000;`
- SQLite does not support `GRANT` or user authentication. Use file system permissions.
- Backup SQLite databases with `.backup` command or `VACUUM INTO 'backup.db'`.

## Performance implications

- In-memory SQLite: 50,000+ simple queries per second on modern hardware.
- File-backed SQLite (WAL mode): 10,000–50,000 reads/sec, 500–5,000 writes/sec depending on disk speed.
- SQLite write throughput is fundamentally limited by the single-writer constraint. Concurrent writers queue.
- `BEGIN IMMEDIATE` transaction mode starts a write transaction immediately instead of upgrading from a read, avoiding some `database is locked` errors.
- Benchmark your specific workload — SQLite performance varies dramatically with schema, data size, and access pattern.

## Practice task

Write a Go program that:
1. Creates a file-backed SQLite database (e.g., `practice.db`) in the OS temp directory.
2. Enables WAL mode and foreign keys.
3. Creates a `sessions` table (id INTEGER PRIMARY KEY, token TEXT UNIQUE, created_at TEXT DEFAULT CURRENT_TIMESTAMP).
4. Inserts 1000 sessions in a transaction.
5. Queries and counts them.
6. Clean up the database file at the end.
7. Run the same operations on an in-memory database and compare the approach.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/07-sqlite-as-a-local-learning-tool
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/07-sqlite-as-a-local-learning-tool
```

## Review questions

1. What is the difference between `:memory:` and a file-backed SQLite database?
2. What does WAL mode do and why is it important for concurrent access?
3. Why does each `:memory:` connection create an independent database?
4. What is the main limitation of SQLite compared to PostgreSQL?
5. How do you enable foreign key enforcement in SQLite?

## NEXT UP

PostgreSQL as production default — why PostgreSQL is the recommended production database for Go services.
