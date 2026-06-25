# database/sql

## Learning objective

Use Go's `database/sql` package to open a database connection, register a driver, execute queries, and handle connection lifecycle with `sql.Open`, `sql.DB`, and `db.Ping`.

## Why this matters

`database/sql` is the standard interface between Go and SQL databases. Every SQL driver (SQLite, PostgreSQL, MySQL, SQL Server) conforms to it. Understanding `database/sql` means you can switch databases by changing one import and one DSN — the rest of your code stays the same. This abstraction is what makes Go database code portable and maintainable across environments.

## Mental model

`database/sql` is a **connection pool manager**, not a single connection. Think of it as a taxi dispatcher: you tell the dispatcher where to find taxis (the DSN), and it maintains a pool of taxis (connections) ready to go. When you need to make a trip (execute a query), the dispatcher assigns you a taxi, you use it, and it returns to the pool. The dispatcher handles creation, reuse, and retirement of taxis automatically.

## Core idea

The `database/sql` package provides:

- **`sql.DB`**: A handle to a database connection pool. It is safe for concurrent use by multiple goroutines.
- **`sql.Open`**: Creates a `*sql.DB`. It does not connect — it only validates the driver name and DSN format.
- **`sql.Driver`**: An interface that database drivers implement. Drivers register themselves via `init()` functions with blank imports.
- **`db.Ping`**: Verifies a connection is alive by actually connecting to the database.

Key types:

| Type | Purpose |
|---|---|
| `*sql.DB` | Connection pool handle |
| `*sql.Rows` | Result set iterator |
| `*sql.Row` | Single-row result |
| `sql.Result` | Result of Exec (LastInsertId, RowsAffected) |
| `*sql.Stmt` | Prepared statement |
| `*sql.Tx` | Database transaction |

## Under the hood

When you call `sql.Open`, the package stores the driver name and DSN but does not establish a connection. The first actual connection is created lazily, when you call `Ping`, `Exec`, `Query`, or `Begin`.

The connection pool works as follows:

1. A goroutine needs a connection. It acquires a mutex.
2. If an idle connection exists in the pool, it is returned.
3. If no idle connection exists and `maxOpen` has not been reached, a new connection is opened.
4. If `maxOpen` has been reached, the goroutine blocks until a connection is returned or the context is cancelled.
5. When the goroutine finishes its query, the connection returns to the idle pool.
6. If the connection is stale or has exceeded `maxLifetime`, it is closed instead of reused.

The driver registration happens via blank imports:

```go
import _ "modernc.org/sqlite"
```

This calls the driver's `init()` function, which calls `sql.Register("sqlite", &driver{})`.

## How Go uses it

Every Go program that uses a SQL database follows this pattern:

```go
import (
	"database/sql"
	_ "driver/path" // driver registers itself
)

func main() {
	db, err := sql.Open("drivername", "datasource")
	if err != nil { /* handle */ }
	defer db.Close()

	if err := db.Ping(); err != nil { /* handle */ }

	// Use db.Exec, db.Query, db.QueryRow, db.Begin
}
```

The blank import is essential — it tells Go to include the driver package and run its `init()` function, which registers the driver name with `database/sql`.

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

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected (pool ready)")

	// Exec — for statements that don't return rows
	res, err := db.Exec(`CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, val TEXT)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Table created: %d rows affected\n", mustRowsAffected(res))

	// Exec with INSERT
	res, err = db.Exec(`INSERT INTO test (val) VALUES ('hello')`)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	fmt.Printf("Inserted row id=%d\n", id)

	// QueryRow — single row
	var val string
	err = db.QueryRow(`SELECT val FROM test WHERE id = ?`, id).Scan(&val)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Queried: %s\n", val)

	// Query — multiple rows
	rows, err := db.Query(`SELECT id, val FROM test ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var val string
		if err := rows.Scan(&id, &val); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Row: id=%d, val=%s\n", id, val)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func mustRowsAffected(res sql.Result) int64 {
	n, _ := res.RowsAffected()
	return n
}
```

## Step-by-step execution

1. `sql.Open("sqlite", ":memory:")` registers the driver name and DSN. No connection is made yet.
2. `db.Ping()` triggers the first actual connection to the in-memory SQLite database.
3. `db.Exec` runs a `CREATE TABLE` statement. It acquires a connection from the pool, executes the statement, and returns the connection.
4. A second `db.Exec` inserts a row. `LastInsertId` returns the auto-increment ID.
5. `db.QueryRow` fetches a single row and scans the value into a Go variable.
6. `db.Query` fetches multiple rows, iterates with `rows.Next()`, and scans each row.

## Common mistakes

- **Assuming `sql.Open` establishes a connection**: It does not. Always call `db.Ping()` or a query to verify connectivity.
- **Forgetting to close `*sql.Rows`**: Leaking `rows` holds a connection from the pool indefinitely. Always `defer rows.Close()`.
- **Using `db.Query` for statements that don't return rows**: Use `db.Exec` for INSERT, UPDATE, DELETE, CREATE. `db.Query` returns `*sql.Rows` that must be closed.
- **Ignoring `rows.Err()`**: After iterating with `rows.Next()`, check `rows.Err()` for iteration errors (network issues, context cancellation).

## Debugging walkthrough

```go
rows, err := db.Query(`SELECT id, val FROM test`)
if err != nil {
	log.Fatal(err)
}
// Forgot: defer rows.Close()

for rows.Next() {
	// process rows
}
if err := rows.Err(); err != nil {
	log.Fatal(err)
}
```

**Symptom**: After all rows are processed, subsequent queries block indefinitely.

**Root cause**: `rows.Close()` was never called. The connection used by `rows` was never returned to the pool. Eventually, all connections in the pool are held by unclosed `Rows` objects, and new queries block waiting for a connection.

**Fix**: Add `defer rows.Close()` immediately after the error check.

## Production notes

- Always set pool limits: `db.SetMaxOpenConns(25)`, `db.SetMaxIdleConns(5)`, `db.SetConnMaxLifetime(5 * time.Minute)`.
- Use `db.Stats()` to monitor pool health in production: `OpenConnections`, `InUse`, `Idle`, `WaitCount`.
- For read-heavy workloads, use `db.QueryContext` with a context timeout to prevent runaway queries.
- For write-heavy workloads, batch operations in transactions to reduce round trips.

## Performance implications

- `database/sql` adds a small overhead per call (acquire from pool, driver conversion). This is typically < 10µs.
- The pool eliminates the cost of creating new connections on every query. Connecting to PostgreSQL over TCP takes 1–10 ms; reusing connections from the pool avoids this cost.
- Prepared statements (`db.Prepare`) can be cached and reused for repeated queries, reducing query planning overhead.

## Practice task

Write a function `InitDB(driverName, dsn string) (*sql.DB, error)` that:
1. Opens a connection with the given driver and DSN.
2. Pings the database.
3. Sets max open connections to 5.
4. Returns the `*sql.DB`.

Then write a `main()` that:
1. Calls `InitDB` with `"sqlite"` and `":memory:"`.
2. Creates a table `logs (id INTEGER PRIMARY KEY, message TEXT, level TEXT)`.
3. Inserts 3 log entries.
4. Queries all entries and prints them.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/10-database-sql
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/10-database-sql
```

## Review questions

1. What does `sql.Open` actually do? What does it NOT do?
2. How does a blank import register a database driver?
3. What is the difference between `db.Exec` and `db.Query`?
4. What happens if you forget to close `*sql.Rows`?
5. What does `db.Ping` do that `sql.Open` does not?

## NEXT UP

sql.DB as a pool — how connection pooling works and how to configure it.
