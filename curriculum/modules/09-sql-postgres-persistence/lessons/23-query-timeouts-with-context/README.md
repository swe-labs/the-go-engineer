# Query timeouts with context

## Learning objective

Use Go's context package with `context.WithTimeout`, `db.ExecContext`, and `db.QueryContext` to set deadlines on database operations and handle cancellation gracefully.

## Why this matters

A database query that hangs indefinitely holds a connection from the pool, blocks a goroutine, and can cascade into a full application outage. Every external call — especially database queries — must have a timeout. Go's `context.Context` is the standard mechanism for propagating deadlines and cancellation across API boundaries.

## Mental model

A context is a timer that starts when the query begins. If the timer expires before the query completes, the context is canceled, and the database operation is aborted.

```
Start query → context.WithTimeout(5s)
                  │
          ┌───────┴───────┐
          │               │
    Query completes    Timer fires
    before 5s          (5s elapsed)
          │               │
    Return result    Cancel context
          │          Abort query
          │          Return "deadline exceeded"
```

The `*sql.DB` methods with `Context` suffix (`ExecContext`, `QueryContext`, `QueryRowContext`, `PrepareContext`) all accept a context. If the context is canceled, the in-progress database operation is interrupted.

## Core idea

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

rows, err := db.QueryContext(ctx, "SELECT * FROM slow_query")
if errors.Is(err, context.DeadlineExceeded) {
    log.Println("Query timed out, consider optimizing or increasing timeout")
}
```

Pattern:
1. Create a context with timeout using `context.WithTimeout`.
2. Pass it to `QueryContext`, `ExecContext`, etc.
3. Check for `context.DeadlineExceeded` or `context.Canceled` in error handling.
4. Always call the `cancel` function to release resources.

## Under the hood

When you call `QueryContext(ctx, ...)`:

1. The `database/sql` package starts a goroutine that monitors `ctx.Done()`.
2. If the context is canceled, the goroutine calls `db.cancelRequest(connID)`.
3. The driver (e.g., modernc.org/sqlite) receives an interrupt signal.
4. The in-progress query is aborted, and the connection is returned to the pool.
5. `QueryContext` returns `context.DeadlineExceeded` or `context.Canceled`.

With SQLite specifically: `sqlite3_interrupt()` is called on the database connection, which causes the current `sqlite3_step()` to return `SQLITE_INTERRUPT`.

## How Go uses it

- **HTTP handlers**: The request context already has a timeout (from middleware). Pass it directly to DB calls.
- **Batch jobs**: Set a per-query timeout to avoid one slow query stalling the entire batch.
- **Background workers**: Use `context.WithTimeout` for each unit of work.
- **Health checks**: Short timeout (1-2s) to fail fast.

## Go example

```go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE numbers (id INTEGER PRIMARY KEY, value TEXT NOT NULL)`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert some data
	for i := 0; i < 100; i++ {
		db.Exec("INSERT INTO numbers (value) VALUES (?)", fmt.Sprintf("number-%d", i))
	}

	// Simulate a slow query using load_extension or a busy wait
	// In SQLite, we use a deliberately expensive query
	slowQuery := `SELECT a.value || b.value AS combined
		FROM numbers a
		CROSS JOIN numbers b
		CROSS JOIN numbers c
		LIMIT 1000000`

	// Query with a very short timeout (1ms) to demonstrate cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err = db.QueryContext(ctx, slowQuery)
	duration := time.Since(start)

	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Printf("Query timed out after %v: %v\n", duration, err)
	} else if err != nil {
		fmt.Printf("Query error: %v\n", err)
	} else {
		fmt.Println("Query succeeded (unexpected with 1ms timeout)")
	}

	// Successful query with adequate timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	start = time.Now()
	var count int
	err = db.QueryRowContext(ctx2, "SELECT COUNT(*) FROM numbers").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Count query completed in %v: %d rows\n", time.Since(start), count)
}
```

## Step-by-step execution

For `QueryContext(ctx, "SELECT COUNT(*) FROM numbers")`:

1. `db.QueryContext(ctx, ...)` creates a `txCtx` linked to the caller's context.
2. A connection is acquired from the pool.
3. The query starts executing.
4. If the query completes before the deadline: rows are returned, the goroutine monitoring the context exits.
5. If the deadline expires: the context is canceled, the monitoring goroutine calls `sqlite3_interrupt()`, the query returns `context.DeadlineExceeded`.

## Common mistakes

- Forgetting `defer cancel()`: `context.WithTimeout` allocates resources. Without `defer cancel()`, the context and its timer leak until the timeout expires.
- Using `context.Background()` for every query: No timeout = no protection. Always derive a context with timeout.
- Not checking `ctx.Err()`: After timeout, the context is done. Check `errors.Is(err, context.DeadlineExceeded)` for proper error handling.
- Using `Exec` instead of `ExecContext` in handlers: The HTTP request context contains the deadline from the client. Pass it through.
- Setting timeouts too short: A timeout that fires during normal load causes false positives. Monitor query latency and set timeouts at the 99th percentile + buffer.

## Debugging walkthrough

**Scenario**: An API endpoint occasionally returns "context deadline exceeded" even though individual queries are fast.

```go
func handler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.QueryContext(r.Context(), "SELECT * FROM users")
    // ...
}
```

**Root cause**: The HTTP request's context includes the client's connection timeout. If the client has a 5-second timeout and the handler does 3 sequential queries of 2 seconds each, by the third query the total exceeds 5 seconds.

**Fix**: Use a derived context with a per-operation timeout, independent of the request context:

```go
opCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
defer cancel()
rows, err := db.QueryContext(opCtx, "SELECT * FROM users")
```

Or increase the client-side timeout.

## Production notes

- **Timeout values**: Start with 5s for OLTP queries, 30s for reporting queries. Adjust based on p99 latency.
- **Context trees**: `r.Context()` → `opCtx` → `db.QueryContext(opCtx, ...)`. If the client disconnects, both the request and the query are canceled.
- **Pool exhaustion**: When queries time out, connections are returned to the pool (unlike stuck queries). Timeouts prevent pool exhaustion.
- **Retry logic**: Do NOT retry on `DeadlineExceeded` — the database may be overloaded. Retry on serialization errors or transient network failures.
- **Logging**: Log the query, duration, and timeout value when a timeout occurs. This helps identify slow queries.

## Performance implications

- Context monitoring adds minimal overhead (a select on a channel per query).
- A timed-out query may have consumed significant database resources before being interrupted. Timeouts limit damage but don't eliminate wasted work.
- SQLite processes queries on a single thread. An interrupted query frees that thread for the next query.
- PostgreSQL can cancel a query in progress, freeing server resources immediately.

## Practice task

Write a function `QueryWithRetry(ctx context.Context, db *sql.DB, query string, args ...any) (*sql.Rows, error)` that:
1. Executes the query with a 2-second timeout.
2. If `DeadlineExceeded`, retries once with a fresh timeout.
3. Returns the result or the final error.
4. Uses `QueryContext`.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/23-query-timeouts-with-context
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/23-query-timeouts-with-context
```

## Review questions

1. Why must you always call `defer cancel()` after `context.WithTimeout`?
2. What is the difference between `context.Background()` and `context.WithCancel()`?
3. How does `QueryContext` differ from `Query` in terms of cancellation?
4. What error should you check to detect a query timeout?
5. Why should you NOT retry on `context.DeadlineExceeded`?

## NEXT UP

Repository and service seam — abstracting data access behind interfaces for testability and clean architecture.
