# Indexes

## Learning objective

Create B-tree indexes in SQLite, understand how indexes accelerate queries, design composite indexes, use `EXPLAIN QUERY PLAN` to verify index usage, and evaluate the tradeoffs of index maintenance.

## Why this matters

Without an index, finding a row in a million-row table requires scanning all million rows (sequential scan). With a B-tree index, the same lookup takes ~20 reads. Indexes are the difference between a 10ms response and a 10-second timeout. But indexes aren't free — they slow down writes and consume disk space. Understanding when and how to create indexes is essential for production database performance.

## Mental model

Think of an index like a phone book. If you want to find "Smith, John," you don't read every page. You open to the S section (the index on last name) and find the exact entry.

```
Table (unsorted rows):        Index on last_name (B-tree):
┌──────┬───────────┬───────┐   ┌───────┬──────────┐
│ id   │ last_name │ city  │   │ Adams │ rowid 1  │
│ 1    │ Adams     │ NY    │   │ Baker │ rowid 3  │
│ 2    │ Smith     │ LA    │   │ Smith │ rowid 2  │  ← 3 reads
│ 3    │ Baker     │ CHI   │   └───────┴──────────┘
└──────┴───────────┴───────┘
3 reads (full scan)            vs       3 reads (index lookup)
```

SQLite uses B-tree indexes by default. B-trees are balanced, sorted, and support fast search, insert, and delete (O(log n)).

## Core idea

```sql
CREATE INDEX idx_users_email ON users(email);
```

Creates an index on the `email` column of the `users` table. The database maintains this index automatically — when you INSERT, UPDATE, or DELETE a row, the index is updated.

```go
db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)")
```

Composite indexes cover multiple columns:

```sql
CREATE INDEX idx_orders_user_created ON orders(user_id, created_at);
```

The order of columns in a composite index matters. The index is sorted by the first column, then the second, then the third. It can speed up queries on `user_id` alone, or `user_id + created_at`, but NOT `created_at` alone.

## Under the hood

A B-tree index in SQLite:
1. Internal nodes store key ranges and pointers to child pages.
2. Leaf nodes store key values and rowids.
3. Each page is 4KB (default). A page holds ~50-100 entries.
4. Height of the tree: for 1M rows, height ≈ 4 (page size / entries per page → log_50(1M) ≈ 4).
5. Lookup cost: 4 pages read (one per tree level) + 1 page read for the actual row data = 5 total reads.

Without an index: 1M pages read (full table scan, but compressed into fewer pages if using B-tree for the table itself — SQLite's table is a B-tree on the rowid).

## How Go uses it

- **Migration files**: Index creation is part of schema migrations.
- **Query optimization**: Add indexes based on slow query analysis.
- **Unique indexes**: `CREATE UNIQUE INDEX` enforces uniqueness at the database level.
- **Partial indexes**: `CREATE INDEX idx_active_users ON users(email) WHERE active = 1` — only indexes active users.

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

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert 10,000 users
	fmt.Println("Inserting 10,000 users...")
	start := time.Now()
	tx, _ := db.Begin()
	for i := 0; i < 10000; i++ {
		tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)",
			fmt.Sprintf("User %d", i),
			fmt.Sprintf("user%d@example.com", i))
	}
	tx.Commit()
	fmt.Printf("Inserted in %v\n", time.Since(start))

	// Query without index
	fmt.Println("\n--- Query without index ---")
	explainQueryPlan(db, "SELECT * FROM users WHERE email = 'user5000@example.com'")

	start = time.Now()
	var name string
	db.QueryRow("SELECT name FROM users WHERE email = 'user5000@example.com'").Scan(&name)
	fmt.Printf("Lookup without index: %v\n", time.Since(start))

	// Create index
	db.Exec("CREATE INDEX idx_users_email ON users(email)")
	fmt.Println("\nCreated index on email.")

	// Query with index
	fmt.Println("\n--- Query with index ---")
	explainQueryPlan(db, "SELECT * FROM users WHERE email = 'user5000@example.com'")

	start = time.Now()
	db.QueryRow("SELECT name FROM users WHERE email = 'user5000@example.com'").Scan(&name)
	fmt.Printf("Lookup with index: %v\n", time.Since(start))
}

func explainQueryPlan(db *sql.DB, query string) {
	rows, err := db.Query("EXPLAIN QUERY PLAN " + query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, parent int
		var detail string
		rows.Scan(&id, &parent, &detail)
		fmt.Printf("  %s\n", detail)
	}
}
```

On a small in-memory database, the difference may be negligible, but on a disk-based database with millions of rows, the index makes a dramatic difference.

## Step-by-step execution

For `SELECT * FROM users WHERE email = 'user5000@example.com'`:

**Without index:**
1. SQLite reads the root page of the `users` B-tree.
2. It scans each leaf page sequentially until finding the matching email.
3. For 10,000 rows, it may read 20-50 pages (depending on row size).
4. Returns the matching row.

**With index `idx_users_email`:**
1. SQLite reads the root page of the index B-tree.
2. Navigates down the tree to find `user5000@example.com` (~3-4 pages).
3. Reads the rowid from the index leaf.
4. Does a single lookup in the `users` table by rowid.
5. Total: ~5 pages read (vs. 20-50 without index, much more for millions of rows).

## Common mistakes

- Over-indexing: Every index slows down INSERT, UPDATE, DELETE. Don't index every column — index only columns used in WHERE, JOIN, and ORDER BY.
- Indexing low-cardinality columns: An index on a boolean column (true/false) is nearly useless. It doesn't narrow down the search enough.
- Wrong column order in composite indexes: `INDEX(a, b)` helps queries on `a` and `(a, b)`, but NOT `b` alone.
- Not indexing foreign keys: `orders.user_id` should be indexed if you join orders with users.
- Indexing without verifying: Always use `EXPLAIN QUERY PLAN` to confirm the index is being used.

## Debugging walkthrough

**Scenario**: A query that joins orders with users is slow.

```sql
SELECT * FROM orders JOIN users ON orders.user_id = users.id WHERE users.email = 'alice@example.com';
```

**Check**: `EXPLAIN QUERY PLAN SELECT ...`

**Problem**: The query plan shows `SCAN orders` and `SEARCH users USING INTEGER PRIMARY KEY (rowid=?)`. This means it's scanning ALL orders and looking up each user — fine if few orders, bad if millions.

**Fix**: Ensure there's an index on `orders.user_id`:

```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

Re-check the query plan. It should show `SEARCH orders USING INDEX idx_orders_user_id (user_id=?)`.

## Production notes

- **Monitoring**: Identify slow queries via database logs or application tracing. Add indexes based on actual query patterns, not guesses.
- **Index maintenance**: In PostgreSQL, `CREATE INDEX` blocks writes (table-level lock). Use `CREATE INDEX CONCURRENTLY` in production.
- **Index size**: In SQLite, an index is about 1.5x the size of the indexed data. Factor this into storage planning.
- **Covering indexes**: An index that includes all columns needed by a query (covering index) eliminates the need to read the table entirely — the query is served from the index alone.
- **Drop unused indexes**: Periodically review index usage statistics and drop indexes that are never used.

## Performance implications

| Operation | Without index | With B-tree index |
|---|---|---|
| Point lookup (WHERE id = ?) | O(n) | O(log n) |
| Range scan (WHERE id BETWEEN ? AND ?) | O(n) | O(log n + k) |
| INSERT | O(1) | O(log n) — must update index |
| DELETE | O(n) | O(log n) — must remove from index |
| UPDATE indexed column | O(n) | O(log n) — must update index |

An index on a column that's updated frequently adds overhead to every update. Measure the write workload before adding indexes.

## Practice task

1. Create a table `logs` with columns `id`, `level`, `message`, `created_at`.
2. Insert 100,000 rows with random levels (`INFO`, `WARN`, `ERROR`) and dates.
3. Query without index: count of `ERROR` logs in a date range. Measure time.
4. Create a composite index on `(level, created_at)`.
5. Re-run the query. Measure time. Compare.
6. Use `EXPLAIN QUERY PLAN` before and after.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/21-indexes
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/21-indexes
```

## Review questions

1. What data structure does SQLite use for indexes?
2. Why doesn't a boolean column benefit much from an index?
3. For a composite index on `(a, b)`, which queries benefit?
4. What is the tradeoff of adding an index to a column that is frequently updated?
5. How can you verify that a query is using an index?

## NEXT UP

EXPLAIN — reading query plans to understand and optimize query execution.
