# Prepared statements

## Learning objective

Use `db.Prepare` to create prepared statements, execute them with `Exec` and `Query`, understand how they prevent SQL injection, and manage statement lifecycle with `Close`.

## Why this matters

Every time you write a SQL query string with `fmt.Sprintf("SELECT * FROM users WHERE id = %s", userInput)`, you create a SQL injection vulnerability. Prepared statements separate SQL structure from data: the database compiles the query plan once, and each execution sends only the parameters. This prevents injection, improves performance for repeated queries, and is the standard way to interact with SQL databases in production Go applications.

## Mental model

Think of a prepared statement as a template with slots.

```
SQL template:  INSERT INTO users (name, email) VALUES (?, ?)
                                      slot 1     slot 2

Execution 1:   slot 1 = "Alice"    slot 2 = "alice@example.com"
Execution 2:   slot 1 = "Bob"      slot 2 = "bob@example.com"
```

The database parses and optimizes the template once. Each execution fills the slots with literal values that are never re-interpreted as SQL. The `?` markers are not string interpolation — they are protocol-level parameter binding.

## Core idea

`db.Prepare(sql)` sends the SQL string to the database, which parses, validates, and compiles it into an execution plan. It returns a `*sql.Stmt`. You can then call `stmt.Exec(args...)` or `stmt.Query(args...)` multiple times with different arguments. When done, call `stmt.Close()` to free the server-side resource.

```go
stmt, err := db.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
// ... handle err
defer stmt.Close()
stmt.Exec("Alice", "alice@example.com")
stmt.Exec("Bob", "bob@example.com")
```

## Under the hood

When `Prepare` is called:
1. The Go `database/sql` package sends a `COM_STMT_PREPARE` command (MySQL) or `Parse` message (PostgreSQL) to the server.
2. The database server parses the SQL, validates permissions, and generates an execution plan (a prepared statement object).
3. The server returns a statement ID to the client.
4. On `Exec`/`Query`, the client sends `COM_STMT_EXECUTE` with the statement ID and the parameters as binary values.
5. The server substitutes parameters directly into the compiled plan — they never pass through the SQL parser again.

With SQLite (via modernc.org/sqlite), `Prepare` calls `sqlite3_prepare_v2` internally, which compiles the SQL into bytecode. Each `Exec` call binds parameters with `sqlite3_bind_*` and steps through with `sqlite3_step`.

## How Go uses it

- **Bulk inserts**: Prepare once, execute many times with different values.
- **User-facing queries**: Any query that includes user input (search, filters, IDs) must use parameterized queries.
- **Batch operations**: Update many rows with different values in a loop.
- **Dynamic filtering**: Build queries with a fixed structure and variable WHERE clauses using parameterized markers.

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

	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		log.Fatal(err)
	}

	users := []struct {
		name  string
		email string
	}{
		{"Alice", "alice@example.com"},
		{"Bob", "bob@example.com"},
		{"Charlie", "charlie@example.com"},
	}

	stmt, err := db.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, u := range users {
		result, err := stmt.Exec(u.name, u.email)
		if err != nil {
			log.Fatal(err)
		}
		id, _ := result.LastInsertId()
		fmt.Printf("Inserted user %d: %s <%s>\n", id, u.name, u.email)
	}

	row := db.QueryRow("SELECT COUNT(*) FROM users")
	var count int
	row.Scan(&count)
	fmt.Printf("Total users: %d\n", count)
}
```

**SQL injection demonstration** — the wrong way and the right way:

```go
// DANGEROUS: string interpolation with user input
userInput := "1 OR 1=1"
query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userInput)
// Executes: SELECT * FROM users WHERE id = 1 OR 1=1  -- returns ALL rows!

// SAFE: parameterized query
stmt, _ := db.Prepare("SELECT * FROM users WHERE id = ?")
stmt.Query(userInput)  // treats "1 OR 1=1" as a literal string value
```

## Step-by-step execution

For `stmt.Exec("Bob", "bob@example.com")`:

1. `database/sql` calls the driver's `Stmt.Exec` with `["Bob", "bob@example.com"]`.
2. The driver (modernc.org/sqlite) calls `sqlite3_bind_text(stmt, 1, "Bob")` and `sqlite3_bind_text(stmt, 2, "bob@example.com")`.
3. `sqlite3_step` executes the compiled bytecode with bound values.
4. The row is inserted. SQLite returns `SQLITE_DONE`.
5. `database/sql` wraps the result, extracts `LastInsertId` and `RowsAffected`.
6. Your Go code receives the `sql.Result`.

## Common mistakes

- Not closing prepared statements: Every `Prepare` uses server-side resources. Always `defer stmt.Close()`.
- Preparing the same statement repeatedly inside a loop: Prepare once outside the loop, execute inside.
- Using `?` for identifiers (table names, column names): Parameter markers only work for values. Use string validation or allowlists for identifiers.
- Not checking errors from `Prepare`: A typo in SQL becomes a prepare-time error. Always check.
- Assuming `Prepare` caches across `db` instances: Each `*sql.DB` connection pool manages its own prepared statements.

## Debugging walkthrough

**Scenario**: A bulk insert function is slow and seems to leak connections.

```go
func InsertUsers(db *sql.DB, users []User) error {
	for _, u := range users {
		stmt, err := db.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
		if err != nil {
			return err
		}
		stmt.Exec(u.name, u.email)
		stmt.Close()
	}
	return nil
}
```

**Problem**: `Prepare` is called inside the loop for every user. Each prepare sends a round-trip to the database. For 10,000 users, that's 10,000 prepare + 10,000 execute = 20,000 round-trips.

**Fix**: Prepare once before the loop:

```go
func InsertUsers(db *sql.DB, users []User) error {
	stmt, err := db.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, u := range users {
		if _, err := stmt.Exec(u.name, u.email); err != nil {
			return err
		}
	}
	return nil
}
```

This reduces round-trips to 1 prepare + N executes = N+1. For bulk operations, consider wrapping in a transaction to reduce to N+2.

## Production notes

- **Statement pooling**: `database/sql` internally caches prepared statements per connection. Calling `Prepare` multiple times with the same SQL reuses the cached statement on that connection.
- **Connection pinning**: After `Prepare`, subsequent `Exec`/`Query` calls must use the same connection. This can reduce connection pool utilization under high concurrency.
- **PostgreSQL specifics**: Use `$1, $2` instead of `?` for parameter markers. The `database/sql` interface abstracts this via `db.Rebind()` in some drivers.
- **Monitoring**: Track `stmt.Prepare` count and `stmt.Close` in production to detect leaks.

## Performance implications

- Prepared statements eliminate SQL parsing overhead on repeated executions.
- Parameter binding sends values in binary format (more efficient than text encoding in SQL strings).
- For one-off queries, `db.Exec` with inline SQL is fine. For hot-path queries, prepare and cache.
- The `database/sql` connection pool's prepared statement cache may grow unbounded if you prepare many distinct SQL strings. Use `db.Stats()` to monitor.

## Practice task

Write a function `BatchInsertProducts(db *sql.DB, products []Product) error` that:
1. Creates a `products` table with `id INTEGER PRIMARY KEY`, `name TEXT`, `price REAL`.
2. Prepares an INSERT statement once.
3. Executes it for each product in the slice.
4. Returns the first error, if any.
5. In `main()`, insert 5 products and query them back to verify.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/16-prepared-statements
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/16-prepared-statements
```

## Review questions

1. Why does `fmt.Sprintf` for SQL queries create a security vulnerability?
2. What happens if you call `stmt.Close()` before all executions complete?
3. Can you use `?` parameter markers for table or column names? Why or why not?
4. What is the performance difference between preparing inside a loop vs. preparing once outside?
5. How does `database/sql` handle prepared statements across multiple connections in the pool?

## NEXT UP

Transactions — grouping multiple operations into atomic units with commit and rollback.
