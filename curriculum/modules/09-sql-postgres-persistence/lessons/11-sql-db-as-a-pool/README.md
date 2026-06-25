# sql.DB as a pool

## Learning objective

Configure a `sql.DB` connection pool with `SetMaxOpenConns`, `SetMaxIdleConns`, and `SetConnMaxLifetime`, and explain how pool settings affect application performance and reliability.

## Why this matters

A misconfigured connection pool is one of the most common sources of production database issues. Too few connections causes request queuing and timeouts. Too many connections overwhelms the database. Stale connections cause spurious errors. Go engineers who understand pool internals can diagnose "connection pool exhausted" errors, tune for throughput, and prevent database outages before they happen.

## Mental model

`sql.DB` is a pool of reusable TCP connections to the database. The pool has three dials:

1. **MaxOpenConns**: The garage size — maximum number of connections that can exist at once.
2. **MaxIdleConns**: How many spare cars to keep ready in the garage, even when not in use.
3. **ConnMaxLifetime**: How long a car can stay in service before being sent to the scrapyard.

When a goroutine needs a connection, the pool hands out an idle one or opens a new one (up to MaxOpenConns). When the goroutine finishes, the connection returns to the idle pool (up to MaxIdleConns) or is closed.

## Core idea

The pool configuration methods:

```go
db.SetMaxOpenConns(n int)       // Maximum simultaneous connections (0 = unlimited)
db.SetMaxIdleConns(n int)       // Maximum idle connections kept open (default = 2)
db.SetConnMaxLifetime(d time.Duration) // Maximum age of a connection (0 = forever)
db.SetConnMaxIdleTime(d time.Duration) // Maximum idle time before close (0 = forever)
```

| Setting | Low value | High value |
|---|---|---|
| MaxOpenConns | Limits throughput, causes waiting | Overloads database, consumes memory |
| MaxIdleConns | Creates/destroys connections frequently | Wastes memory on idle connections |
| ConnMaxLifetime | More connection churn | Risk of stale connections from network drops |
| ConnMaxIdleTime | Closes idle connections aggressively | Keeps idle connections alive unnecessarily |

## Under the hood

The pool is implemented as a channel-based semaphore with a free list:

```
sql.DB {
    free:   [conn, conn, conn]  // idle connections
    wait:   [req, req]          // goroutines waiting for a connection
    maxOpen: 25                 // hard limit
}
```

When `db.Exec` is called:

1. Check if `free` has a connection. If yes, use it.
2. If `free` is empty and `total < maxOpen`, create a new connection.
3. If `total >= maxOpen`, block on the `wait` queue.
4. When a connection is released, return it to `free` (if `free < maxIdle`) or close it.
5. If `wait` has queued requests, wake the first one.

`SetConnMaxLifetime` works via a timer: when a connection is created, a goroutine waits for `maxLifetime`, then marks the connection as expired. The connection is closed when returned to the pool.

## How Go uses it

Pool configuration happens once during application initialization:

```go
db, _ := sql.Open("sqlite", ":memory:")
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(1 * time.Minute)
```

Settings should be based on:
- **Database limits**: PostgreSQL default `max_connections` is 100. Set Go pool below this.
- **Query latency**: For fast queries (1 ms), fewer connections suffice. For slow queries (100 ms), more connections are needed for throughput.
- **Concurrent goroutines**: If 100 goroutines make database calls simultaneously, MaxOpenConns should be close to 100 (if the database allows).

## Go example

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(30 * time.Second)
	db.SetConnMaxIdleTime(10 * time.Second)
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	configurePool(db)

	db.Exec(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, ts TEXT, msg TEXT)`)

	// Simulate concurrent work
	for i := 0; i < 5; i++ {
		go func(n int) {
			_, err := db.Exec(`INSERT INTO events (ts, msg) VALUES (datetime('now'), ?)`,
				fmt.Sprintf("event %d", n))
			if err != nil {
				log.Printf("insert error: %v", err)
			}
		}(i)
	}

	time.Sleep(100 * time.Millisecond)

	stats := db.Stats()
	fmt.Printf("Pool stats:\n")
	fmt.Printf("  MaxOpen:     %d\n", stats.MaxOpenConnections)
	fmt.Printf("  Open:        %d\n", stats.OpenConnections)
	fmt.Printf("  InUse:       %d\n", stats.InUse)
	fmt.Printf("  Idle:        %d\n", stats.Idle)
	fmt.Printf("  WaitCount:   %d\n", stats.WaitCount)
	fmt.Printf("  WaitDuration: %v\n", stats.WaitDuration)

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&count)
	fmt.Printf("Total events: %d\n", count)
}
```

## Step-by-step execution

1. `configurePool` sets the pool limits. For this in-memory SQLite, the limits are not binding (SQLite handles everything in-process), but they demonstrate the configuration pattern.
2. Five goroutines each insert an event concurrently. The pool manages the five concurrent requests, handing out connections from the pool.
3. `db.Stats()` returns current pool statistics: how many connections are open, in use, idle, and how many goroutines are waiting.
4. The program prints the pool stats and the total event count.

## Common mistakes

- **Setting MaxOpenConns too high**: Each connection consumes database resources (memory, file descriptors). Setting MaxOpenConns to 1000 on a PostgreSQL with `max_connections=100` causes connection refusal.
- **Not setting ConnMaxLifetime**: Network interruptions, firewall timeouts, and database restarts can silently kill connections. Without a max lifetime, stale connections stay in the pool and cause "connection reset by peer" errors.
- **Setting MaxIdleConns higher than MaxOpenConns**: The pool silently caps MaxIdleConns to MaxOpenConns, but relying on silent behavior is confusing. Always set MaxIdleConns <= MaxOpenConns.
- **Forgetting to configure the pool entirely**: Defaults are MaxOpenConns=0 (unlimited), MaxIdleConns=2, ConnMaxLifetime=0 (forever). This is not suitable for production.

## Debugging walkthrough

```go
db, _ := sql.Open("postgres", dsn)
// No pool configuration

for i := 0; i < 1000; i++ {
    go func() {
        for {
            db.Exec(`SELECT 1`)
        }
    }()
}
```

**Symptom**: PostgreSQL logs show "FATAL: sorry, too many clients already". Go logs show "connection refused" errors.

**Root cause**: MaxOpenConns defaults to 0 (unlimited). 1000 goroutines create 1000 connections, exceeding PostgreSQL's `max_connections` (default 100).

**Fix**: Set `db.SetMaxOpenConns(50)` — well under the PostgreSQL limit of 100.

```go
db.SetMaxOpenConns(50)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)
```

## Production notes

- Start with: MaxOpenConns = 25, MaxIdleConns = 5, ConnMaxLifetime = 5 minutes. Adjust based on load testing.
- Monitor `db.Stats()` in production via metrics (Prometheus, OpenTelemetry). Key metrics: `InUse` (should be below MaxOpenConns), `WaitCount` (should be low), `WaitDuration` (should be < query latency).
- For serverless (AWS Lambda), set MaxOpenConns low (1-2) because Lambda instances are short-lived. Do not reuse `*sql.DB` across invocations.
- Set `ConnMaxLifetime` below any network middlebox timeout (e.g., AWS ALB idle timeout is 60s; set ConnMaxLifetime to 30s).

## Performance implications

- More idle connections = faster first query (no TCP handshake). But idle connections consume memory (~3 MB each in PostgreSQL).
- Higher MaxOpenConns increases database CPU and memory usage. The database must manage more connections, more context switches, and more locks.
- Connection creation is expensive: TCP handshake + TLS negotiation + authentication = 10-100 ms. The pool amortizes this cost by reusing connections.
- `ConnMaxLifetime` slightly reduces pool efficiency by retiring connections early, but prevents using stale connections that cause errors.

## Practice task

Write a function `MonitorPool(db *sql.DB) *sql.DBStats` that:
1. Returns `db.Stats()`.
2. Logs a warning if `WaitCount > 100` (indicating pool contention).
3. Returns the stats.

Then write a `main()` that:
1. Opens an in-memory SQLite database.
2. Configures the pool: MaxOpenConns=5, MaxIdleConns=2, ConnMaxLifetime=10s.
3. Launches 10 goroutines that each do 10 INSERTs.
4. Calls `MonitorPool` periodically.
5. Reports the final stats.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/11-sql-db-as-a-pool
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/11-sql-db-as-a-pool
```

## Review questions

1. What is the default value of `MaxOpenConns` and why is it dangerous for production?
2. What happens when a goroutine calls `db.Exec` and all connections are in use?
3. Why is `ConnMaxLifetime` important even if your network is reliable?
4. What is the relationship between `MaxIdleConns` and `MaxOpenConns`?
5. How do you monitor the pool's health in production?

## NEXT UP

INSERT — how to insert rows with Go, use `LastInsertId`, and batch inserts efficiently.
