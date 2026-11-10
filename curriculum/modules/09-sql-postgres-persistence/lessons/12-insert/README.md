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

## Visual model

```mermaid
flowchart LR
    QUERY["db.QueryRowContext(ctx,
  "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
  "Alice", "alice@example.com"
)"]
    QUERY -->|Parameterized execution ($1, $2)| PG["PostgreSQL Execution Engine"]
    PG -->|Returns auto-generated UUID/serial| SCAN["row.Scan(&newUserID)"]
```

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
	"strings"

	_ "modernc.org/sqlite"
)

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY,
		customer TEXT NOT NULL,
		amount INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending'
	)`)
	return err
}

func InsertOrder(db *sql.DB, customer string, amount int) (int64, error) {
	res, err := db.Exec(`INSERT INTO orders (customer, amount) VALUES (?, ?)`, customer, amount)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func InsertOrdersBatch(db *sql.DB, orders []struct {
	Customer string
	Amount   int
}) (int64, error) {
	if len(orders) == 0 {
		return 0, nil
	}
	placeholders := make([]string, len(orders))
	args := make([]interface{}, 0, len(orders)*2)
	for i, o := range orders {
		placeholders[i] = "(?, ?)"
		args = append(args, o.Customer, o.Amount)
	}
	query := fmt.Sprintf("INSERT INTO orders (customer, amount) VALUES %s", strings.Join(placeholders, ", "))
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := CreateTable(db); err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(`INSERT INTO orders (customer, amount) VALUES (?, ?)`, "Alice", 100)
	if err != nil {
		log.Fatal(err)
	}
	aliceID, _ := res.LastInsertId()
	fmt.Printf("Alice order inserted with id=%d\n", aliceID)

	stmt, err := db.Prepare(`INSERT INTO orders (customer, amount) VALUES (?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	extraOrders := []struct {
		customer string
		amount   int
	}{
		{"Bob", 200},
		{"Carol", 300},
	}
	for _, o := range extraOrders {
		res, err := stmt.Exec(o.customer, o.amount)
		if err != nil {
			log.Fatal(err)
		}
		id, _ := res.LastInsertId()
		fmt.Printf("%s order inserted with id=%d\n", o.customer, id)
	}

	batchOrders := []struct {
		Customer string
		Amount   int
	}{
		{"Dave", 400},
		{"Eve", 500},
	}
	count, err := InsertOrdersBatch(db, batchOrders)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Batch inserted %d orders\n", count)

	rows, err := db.Query(`SELECT id, customer, amount, status FROM orders ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("All orders:")
	for rows.Next() {
		var id int
		var customer, status string
		var amount int
		if err := rows.Scan(&id, &customer, &amount, &status); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %d: %s ($%d) [%s]\n", id, customer, amount, status)
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
- **Using `LastInsertId` where it is unsupported**: PostgreSQL has no "last insert id" primitive — with `lib/pq` the call errors, and `WITHOUT ROWID` SQLite tables (or inserts into UUID-keyed tables) report 0 (see lesson 04). Use `RETURNING id` with `QueryRow` for PostgreSQL, and `INTEGER PRIMARY KEY` tables in SQLite.
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

**Symptom**: LastInsertId returns 0 (or errors) although the INSERT succeeded.

**Root cause**: `LastInsertId` is driver- and table-dependent. PostgreSQL drivers have no last-insert-id primitive to query — `lib/pq` returns an error. In SQLite, the value exists only for rowid tables: `WITHOUT ROWID` tables report 0, and a `TEXT`/UUID primary key has no auto-increment id to report at all.

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

- Prepared statements: faster for repeated inserts because SQL parsing and planning happen once — the gain grows with statement complexity.
- Transaction-wrapped batch inserts: one to two orders of magnitude faster than individual auto-commit inserts, because each transaction flush is expensive (SQLite's own documentation measures this gap on bulk loads).
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

[Select](../13-select/README.md)
