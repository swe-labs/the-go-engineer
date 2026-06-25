# INSERT

## Learning objective

Insert rows into SQL tables using Go's `db.Exec`, retrieve the last inserted ID, use prepared statements for efficient repeated inserts, and batch insert multiple rows.

## Why this matters

INSERT is the most common write operation in any application. Doing it efficiently separates a production-quality system from a toy. Go engineers must understand when to use `Exec` vs. prepared statements, how to handle auto-generated IDs, and how to batch inserts to minimize round trips.

## Mental model

Think of INSERT as adding a new card to a filing cabinet. You can:
- Add one card at a time (`INSERT INTO t VALUES (1)`).
- Add multiple cards in one trip (`INSERT INTO t VALUES (1), (2), (3)`).
- Prepare a template card and fill it in many times (prepared statement).

Each trip to the cabinet has overhead (opening the drawer, finding the right section). Fewer trips = faster overall insertion.

## Core idea

Three approaches to INSERT in Go:

| Method | Use case | Example |
|---|---|---|
| `db.Exec` | Single insert, no preparation | `db.Exec("INSERT INTO t (name) VALUES (?)", name)` |
| `db.Prepare` + `stmt.Exec` | Repeated inserts with same SQL | `stmt, _ := db.Prepare("INSERT INTO t (name) VALUES (?)")` |
| Multi-value INSERT | Batch insert, known values | `db.Exec("INSERT INTO t (name) VALUES (?), (?)", a, b)` |

`Result` from `db.Exec` provides:
- `LastInsertId()`: The auto-increment ID of the inserted row (driver-dependent).
- `RowsAffected()`: Number of rows inserted (always 1 for a single-value INSERT).

## Under the hood

When `db.Exec` is called:

1. Go acquires a connection from the pool.
2. The SQL text and parameters are sent to the database.
3. The database parses the SQL, binds parameters, and executes the insert.
4. The database returns the row count and last insert ID.
5. The connection returns to the pool.

With a prepared statement:

1. `db.Prepare` sends the SQL without parameters to the database, which parses and plans it. The database returns a statement handle.
2. Each `stmt.Exec` sends only the parameter values — no parsing or planning needed.
3. This reduces per-insert overhead, especially for complex INSERT statements.

## How Go uses it

```go
// Single insert
res, err := db.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, "Alice", "alice@example.com")
id, _ := res.LastInsertId()

// Prepared statement for repeated inserts
stmt, err := db.Prepare(`INSERT INTO users (name, email) VALUES (?, ?)`)
for _, u := range users {
    _, err := stmt.Exec(u.Name, u.Email)
}

// Batch insert with multiple values
res, err := db.Exec(`INSERT INTO users (name, email) VALUES (?, ?), (?, ?)`, "A", "a@x.com", "B", "b@x.com")
```

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

	db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE)`)

	// 1. Single insert with LastInsertId
	res, err := db.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, "Alice", "alice@example.com")
	if err != nil {
		log.Fatal(err)
	}
	aliceID, _ := res.LastInsertId()
	fmt.Printf("Alice inserted with id=%d\n", aliceID)

	// 2. Prepared statement for repeated inserts
	stmt, err := db.Prepare(`INSERT INTO users (name, email) VALUES (?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	extraUsers := []struct{ name, email string }{
		{"Bob", "bob@example.com"},
		{"Carol", "carol@example.com"},
	}
	for _, u := range extraUsers {
		res, err := stmt.Exec(u.name, u.email)
		if err != nil {
			log.Fatal(err)
		}
		id, _ := res.LastInsertId()
		fmt.Printf("%s inserted with id=%d\n", u.name, id)
	}

	// 3. Multi-value INSERT
	res, err = db.Exec(
		`INSERT INTO users (name, email) VALUES (?, ?), (?, ?)`,
		"Dave", "dave@example.com", "Eve", "eve@example.com",
	)
	if err != nil {
		log.Fatal(err)
	}
	count, _ := res.RowsAffected()
	fmt.Printf("Batch inserted %d users\n", count)

	// Verify all users
	rows, err := db.Query(`SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %d: %s <%s>\n", id, name, email)
	}
}
```

## Step-by-step execution

1. Single INSERT adds Alice and returns her auto-generated ID (1).
2. A prepared statement is created for `INSERT INTO users ...`. This sends the SQL to the database for parsing and planning once.
3. Bob and Carol are inserted using `stmt.Exec`, which sends only the parameters. The prepared statement is cached.
4. A multi-value INSERT adds Dave and Eve in one round trip.
5. `RowsAffected()` returns 2 for the batch insert.
6. All rows are queried and printed to verify.

## Common mistakes

- **Ignoring the error from `db.Exec`**: An INSERT can fail due to constraints, unique violations, or type mismatches. Always check the error.
- **Using `LastInsertId` with non-SEQUENCE tables**: In PostgreSQL with `lib/pq`, `LastInsertId` is not supported (use `RETURNING id` with `QueryRow` instead). SQLite supports it for `INTEGER PRIMARY KEY`.
- **Not closing prepared statements**: `db.Prepare` returns a `*sql.Stmt` that must be closed with `stmt.Close()` to free database resources.
- **Assuming `RowsAffected` always returns the number of rows**: Some drivers return 0 for statements that don't produce a row count.

## Debugging walkthrough

```go
for _, u := range users {
    res, err := db.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, u.Name, u.Email)
    if err != nil {
        log.Fatal(err)
    }
    id, _ := res.LastInsertId()
    fmt.Println(id) // 0 for all but the first insert
}
```

**Symptom**: LastInsertId returns 0 after first insert.

**Root cause**: Some drivers reset the last insert ID after each `Exec` if the table changes. But more commonly: the driver doesn't support `LastInsertId`, or the table uses a non-integer primary key.

**Fix**: For PostgreSQL with `pgx`, use `RETURNING id`:

```go
var id int64
db.QueryRow(`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`, u.Name, u.Email).Scan(&id)
```

## Production notes

- Use prepared statements for any INSERT that executes more than once. The planning overhead is saved on every execution.
- For bulk inserts (thousands of rows), use a transaction around the batch: `BEGIN; INSERT ...; INSERT ...; COMMIT;`. This reduces WAL flushes from one per insert to one per transaction.
- Monitor insert latency and volume. Unexpected spikes may indicate a DDoS or a bug in a batch job.
- Set `db.SetMaxOpenConns` appropriately: bulk inserts benefit from fewer, dedicated connections to avoid overwhelming the database.

## Performance implications

- Prepared statements: 20-50% faster for repeated inserts because SQL parsing and planning are done once.
- Transaction-wrapped batch inserts: 10-100x faster than individual inserts because each transaction flush is expensive.
- Multi-value INSERT (single statement, multiple values): fastest for known-in-advance data, but limited by SQL statement length.
- In SQLite, `INSERT` within a transaction is dramatically faster because each INSERT normally triggers a WAL flush. Always batch in a transaction for bulk loads.

## Practice task

Write functions to:
1. `CreateTable(db *sql.DB)` — creates an `orders` table (id INTEGER PRIMARY KEY, customer TEXT NOT NULL, amount INTEGER NOT NULL, status TEXT NOT NULL DEFAULT 'pending').
2. `InsertOrder(db *sql.DB, customer string, amount int) (int64, error)` — inserts one order and returns the id.
3. `InsertOrdersBatch(db *sql.DB, orders []struct{Customer string; Amount int}) (int64, error)` — inserts multiple orders in a single multi-value INSERT, returns total rows affected.

Write a `main()` that inserts 5 orders using both methods and prints the results.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/12-insert
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/12-insert
```

## Review questions

1. What does `res.LastInsertId()` return and when is it zero?
2. Why are prepared statements faster for repeated inserts?
3. How does wrapping inserts in a transaction improve performance?
4. What is the difference between `db.Exec` and `db.Prepare` + `stmt.Exec`?
5. How do you get the inserted ID in PostgreSQL with pgx?

## NEXT UP

SELECT — querying data with `db.Query` and iterating result sets.
