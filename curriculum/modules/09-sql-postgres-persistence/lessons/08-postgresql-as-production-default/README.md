# PostgreSQL as production default

## Learning objective

Describe PostgreSQL's key features that make it the recommended production database, construct a connection string, and connect to PostgreSQL from Go using the `pgx` driver.

## Why this matters

SQLite is ideal for learning and local development, but production Go services need a database that handles concurrent writes, enforces robust authentication, supports network access, and scales horizontally. PostgreSQL is the industry standard for Go backends — it is feature-rich, battle-tested, and has first-class Go driver support. Understanding how to connect and interact with PostgreSQL prepares you for real-world deployment.

## Mental model

If SQLite is a library (embedded in-process), PostgreSQL is a service (separate process, network-accessible). You connect to PostgreSQL over TCP, authenticate with a username and password, and send SQL queries just like you do to SQLite. The key difference is that PostgreSQL runs as a daemon, manages its own memory and disk, and handles thousands of concurrent connections with full MVCC (Multi-Version Concurrency Control).

## Core idea

PostgreSQL is an **open-source, object-relational database system** with over 30 years of development. Key features that make it the production default:

| Feature | Benefit |
|---|---|
| Full ACID compliance | Guaranteed data integrity |
| MVCC | Concurrent reads and writes without blocking |
| Rich data types | JSONB, arrays, ranges, UUID, hstore |
| Extensibility | Custom functions, data types, extensions (PostGIS, pgvector) |
| Replication | Streaming replication, logical replication, hot standby |
| Role-based auth | Fine-grained access control |
| `pgx` driver | High-performance, pure-Go PostgreSQL driver |

## Under the hood

PostgreSQL uses a **process-per-connection** model. Each client connection spawns a backend process. These processes communicate via shared memory for caching and locks. The main process, `postmaster`, listens on a TCP port (default 5432) and forks a new process for each incoming connection.

When you execute a query:

1. The client sends the SQL text over TCP to the backend process.
2. The backend parses, analyzes, and plans the query.
3. The executor reads/writes data pages in shared buffers.
4. WAL records are written to the write-ahead log for durability.
5. Results are serialized and sent back over TCP.

PostgreSQL's query planner uses cost-based optimization, considering table statistics, index presence, and join strategies.

## How Go uses it

Go connects to PostgreSQL primarily via the **`pgx`** driver (`github.com/jackc/pgx/v5`), which is fully compatible with `database/sql`:

```go
import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

db, err := sql.Open("pgx", "postgres://user:pass@localhost:5432/mydb?sslmode=disable")
```

The connection string format: `postgres://username:password@host:port/database?params`.

Connection string parameters:
- `sslmode`: `disable`, `require`, `verify-full` (production).
- `connect_timeout`: connection timeout in seconds.
- `pool_max_conns`: max connections in pgx pool (when using pgx directly).

## Go example

Since PostgreSQL is not available in this environment, the example connects to SQLite but shows the PostgreSQL connection string and code pattern. The logic is identical — `database/sql` abstracts the difference.

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// In production, use:
// import _ "github.com/jackc/pgx/v5/stdlib"
// db, err := sql.Open("pgx", "postgres://user:pass@localhost:5432/mydb?sslmode=disable")

func main() {
	// Simulating the pattern with SQLite
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database (simulated with SQLite)")

	db.Exec(`CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, name TEXT)`)

	// Insert and query — same code works for PostgreSQL
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

	// Check PostgreSQL-specific feature support
	var level string
	db.QueryRow(`SELECT sqlite_version()`).Scan(&level)
	fmt.Println("DB version:", level)

	fmt.Println("\n--- PostgreSQL connection string example ---")
	fmt.Println("Driver: pgx")
	fmt.Println("DSN:    postgres://user:password@localhost:5432/mydb?sslmode=disable")
}
```

## Step-by-step execution

1. `sql.Open("sqlite", ":memory:")` opens a SQLite database. To use PostgreSQL, change the driver to `"pgx"` and the DSN to a PostgreSQL connection string.
2. `db.Ping()` verifies the connection is alive. For SQLite in-memory, this always succeeds.
3. The rest of the code (CREATE TABLE, INSERT, SELECT) uses the same `database/sql` API that works with PostgreSQL.
4. The program prints the PostgreSQL connection string pattern that would be used in production.

## Common mistakes

- **Using `lib/pq` instead of `pgx`**: `lib/pq` is maintained but `pgx` is faster, more feature-rich, and the actively recommended driver. Use `github.com/jackc/pgx/v5/stdlib`.
- **Hardcoding connection strings**: Store them in environment variables or a config file, never in source code.
- **Using `sslmode=disable` in production**: This sends credentials and data in plaintext. Use `sslmode=verify-full` with proper certificates.
- **Assuming PostgreSQL defaults are secure**: By default, PostgreSQL trusts local connections. Configure `pg_hba.conf` for your environment.

## Debugging walkthrough

```go
db, err := sql.Open("pgx", "postgres://localhost:5432/mydb")
```

**Symptom**: `sql.Open` returns nil error, but `db.Ping()` fails with "connection refused" or "role does not exist".

**Root cause**: `sql.Open` does not connect — it only validates the DSN format. `db.Ping()` is the first actual network call.

**Fix**: Check the error from `db.Ping()` (or the first query) and inspect the PostgreSQL server logs.

```go
if err := db.Ping(); err != nil {
	log.Fatalf("cannot connect: %v", err)
}
```

## Production notes

- Always use connection pooling: `db.SetMaxOpenConns(25)` and `db.SetMaxIdleConns(5)` (covered in lesson 11).
- Use environment-specific connection strings via environment variables: `DATABASE_URL=postgres://...`.
- Run migrations with a tool like `golang-migrate/migrate` or `pressly/goose`.
- Monitor `pg_stat_activity` to detect long-running queries and connection leaks.
- Set `statement_timeout` in PostgreSQL to kill queries that run too long.
- Use `LISTEN`/`NOTIFY` for real-time event notifications from the database.

## Performance implications

- PostgreSQL over localhost adds 0.1–0.5 ms of network latency per query compared to SQLite's in-process access.
- Over a network (even same datacenter), add 1–5 ms per query round trip. Batch queries when possible.
- PostgreSQL's query planner generates efficient execution plans for complex JOINs and aggregations that SQLite handles poorly.
- `pgx` is a high-performance driver that uses binary protocol for data transfer, reducing serialization overhead.

## Practice task

Write a function `ConnectPostgres(dsn string) (*sql.DB, error)` that:
1. Opens a connection using `pgx` driver (but test with SQLite).
2. Calls `Ping` to verify connectivity.
3. Sets `SetMaxOpenConns(10)` and `SetMaxIdleConns(3)`.
4. Returns the `*sql.DB`.

Then write a `main()` that:
1. Reads a DSN from an environment variable `DATABASE_URL` (fallback to `:memory:` for SQLite).
2. Calls `ConnectPostgres`.
3. Creates a table, inserts a row, and queries it.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/08-postgresql-as-production-default
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/08-postgresql-as-production-default
```

## Review questions

1. What is the main architectural difference between SQLite and PostgreSQL?
2. What does `sslmode=verify-full` protect against?
3. Why is `pgx` preferred over `lib/pq`?
4. What does `db.Ping()` do that `sql.Open` does not?
5. Name three PostgreSQL features not available in SQLite.

## NEXT UP

Docker Compose for Postgres — how to run PostgreSQL locally with Docker Compose for development.
