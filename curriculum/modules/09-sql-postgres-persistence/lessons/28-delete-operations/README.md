# DELETE operations

## Learning objective

Remove rows from SQLite databases from Go using hard delete, soft delete with `deleted_at`, CASCADE DELETE via foreign keys, and batch delete strategies with `IN` clauses.

## Why this matters

DELETE is irreversible. An incorrect DELETE without a WHERE clause removes all data. In production, most systems avoid hard deletes entirely — they use soft deletes (marking rows as deleted) to preserve data for auditing, recovery, and analytics. Choosing the right delete strategy (hard vs. soft, CASCADE vs. manual) is a critical data integrity decision.

## Mental model

```
Hard delete:     Row is gone forever. SELECT finds nothing.
                 DELETE FROM users WHERE id = 1;

Soft delete:     Row is marked as deleted. Queries filter it out.
                 UPDATE users SET deleted_at = NOW() WHERE id = 1;
                 SELECT * FROM users WHERE deleted_at IS NULL;

CASCADE delete:  Deleting a parent row automatically deletes children.
                 DELETE FROM users WHERE id = 1;
                 → All orders for user 1 are deleted automatically.
```

## Core idea

Three delete strategies:

| Strategy | SQL | Recoverable | Performance | Use case |
|---|---|---|---|---|
| Hard delete | `DELETE FROM t WHERE id = ?` | No | Fast (O(log n)) | Temporary data, cache |
| Soft delete | `UPDATE t SET deleted_at = ? WHERE id = ?` | Yes | Same as UPDATE | User data, orders |
| CASCADE delete | `DELETE FROM t WHERE id = ?` + FK CASCADE | No | Depends on children | Strongly dependent data |

In Go:

```go
// Hard delete
result, err := db.Exec("DELETE FROM products WHERE id = ?", id)

// Soft delete
result, err := db.Exec("UPDATE products SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", id)

// Batch delete
result, err := db.Exec("DELETE FROM products WHERE id IN (?, ?, ?)", 1, 2, 3)
```

## Under the hood

When a hard DELETE executes:
1. SQLite locates the row(s) using the WHERE clause.
2. Acquires a write lock.
3. Removes the row from the table B-tree.
4. Removes corresponding entries from all indexes.
5. If foreign keys are enabled and CASCADE is set, cascaded deletes are performed.
6. The change is journaled for rollback.
7. The free page is marked as available for reuse.

In SQLite, deleted rows don't immediately release space back to the OS (unless `VACUUM` is run). The space is reused by subsequent INSERTs.

## How Go uses it

- **Account deletion**: Soft delete user data, keep for 30 days, then hard delete (GDPR compliance).
- **Order cancellation**: Soft delete so cancelled orders still appear in reports.
- **Cache cleanup**: Hard delete expired cache entries.
- **Batch cleanup**: Delete all expired sessions: `DELETE FROM sessions WHERE expires_at < ?`.

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

	// Enable foreign keys
	db.Exec("PRAGMA foreign_keys = ON")

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			deleted_at TIMESTAMP
		);
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			total REAL NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		INSERT INTO users VALUES (1, 'Alice', NULL);
		INSERT INTO users VALUES (2, 'Bob', NULL);
		INSERT INTO users VALUES (3, 'Carol', NULL);
		INSERT INTO orders VALUES (101, 1, 50.00);
		INSERT INTO orders VALUES (102, 1, 30.00);
		INSERT INTO orders VALUES (103, 3, 20.00);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Soft delete Bob (id=2)
	result, _ := db.Exec("UPDATE users SET deleted_at = ? WHERE id = ?", time.Now().UTC(), 2)
	rows, _ := result.RowsAffected()
	fmt.Printf("Soft deleted %d user(s)\n", rows)

	// Query only active users
	countActive(db)
	countDeleted(db)

	// Hard delete by IN clause (batch)
	result, _ = db.Exec("DELETE FROM orders WHERE user_id = ? AND total < ?", 1, 40.00)
	rows, _ = result.RowsAffected()
	fmt.Printf("Hard deleted %d order(s) for user 1 with total < 40\n", rows)

	// Demonstrate cascade is NOT automatic without FK declaration
	queryOrders(db, 1)

	// Restore soft-deleted user (undelete)
	result, _ = db.Exec("UPDATE users SET deleted_at = NULL WHERE id = ?", 2)
	rows, _ = result.RowsAffected()
	fmt.Printf("Restored %d user(s)\n", rows)

	countActive(db)
}

func countActive(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&count)
	fmt.Printf("Active users: %d\n", count)
}

func countDeleted(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NOT NULL").Scan(&count)
	fmt.Printf("Deleted users: %d\n", count)
}

func queryOrders(db *sql.DB, userID int) {
	rows, _ := db.Query("SELECT id, total FROM orders WHERE user_id = ?", userID)
	defer rows.Close()
	fmt.Printf("Orders for user %d:\n", userID)
	for rows.Next() {
		var id int
		var total float64
		rows.Scan(&id, &total)
		fmt.Printf("  Order %d: $%.2f\n", id, total)
	}
}
```

## Step-by-step execution

For `DELETE FROM orders WHERE user_id = 1 AND total < 40`:

1. SQLite scans the `orders` table (or uses an index on `user_id`).
2. For each row with `user_id = 1`, checks `total < 40`.
3. Matches order 102 (total=30.00).
4. Acquires write lock.
5. Removes order 102 from the table and any indexes.
6. `RowsAffected()` returns 1.
7. Lock is released.

## Common mistakes

- **DELETE without WHERE**: `DELETE FROM users` removes all rows. Always double-check.
- **Not using transactions for batch delete**: If a batch delete fails partway through, some rows are deleted and some aren't. Wrap in a transaction.
- **Foreign key violations**: Deleting a parent row that has children fails if foreign keys are enabled. Use CASCADE, SET NULL, or delete children first.
- **Soft delete without filtering**: Soft-deleted rows still appear in SELECT * queries. Always add `WHERE deleted_at IS NULL` to user-facing queries.
- **Not cleaning up soft-deleted data**: Soft-deleted rows accumulate. Schedule a job to hard-delete old soft-deleted rows.

## Debugging walkthrough

**Scenario**: Deleting a user returns an error "FOREIGN KEY constraint failed".

```go
db.Exec("DELETE FROM users WHERE id = 1")
```

**Root cause**: The user has orders referencing them. Foreign keys are enabled (`PRAGMA foreign_keys = ON`), and no CASCADE is defined.

**Solutions**:

1. Delete children first:
```go
tx, _ := db.Begin()
tx.Exec("DELETE FROM orders WHERE user_id = 1")
tx.Exec("DELETE FROM users WHERE id = 1")
tx.Commit()
```

2. Use CASCADE in schema:
```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
```

3. Soft delete (recommended for user data):
```go
db.Exec("UPDATE users SET deleted_at = ? WHERE id = 1", time.Now())
```

## Production notes

- **Soft delete pattern**: Add `deleted_at TIMESTAMP` to tables. Create views or query helpers that filter `WHERE deleted_at IS NULL`.
- **Hard delete cleanup**: Schedule periodic cleanup of soft-deleted rows older than a retention period (e.g., 90 days).
- **Audit log**: Log all DELETE operations with user ID, timestamp, and affected rows for compliance.
- **Foreign key handling**: SQLite requires `PRAGMA foreign_keys = ON` per connection. It's off by default. Always enable it in your application.
- **Batch delete limits**: Very large IN clauses (10,000+ IDs) may hit SQLite's limit. Use a temporary table with JOIN instead.

## Performance implications

- DELETE removes rows from the table B-tree and all indexes — O(log n) per index.
- Soft delete (UPDATE on `deleted_at`) is identical cost to UPDATE — no index rebuild on `deleted_at` if indexed.
- CASCADE deletes compound the cost: deleting one parent may delete 1000 children.
- Batch DELETE with IN clause is faster than individual DELETE statements (fewer round-trips and transactions).
- Deleted pages are reused. SQLite does not automatically shrink the database file. Use `VACUUM` to reclaim space after bulk deletes.

## Practice task

Write a function `SafeDeleteUser(db *sql.DB, userID int) error` that:
1. Begins a transaction.
2. Soft-deletes the user (sets `deleted_at`).
3. Hard-deletes the user's orders older than 30 days.
4. Transfers remaining orders to a "ghost" user.
5. Commits on success, rolls back on any error.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/28-delete-operations
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/28-delete-operations
```

## Review questions

1. What is the difference between hard delete and soft delete?
2. Why does `DELETE FROM users` without a WHERE clause delete all rows?
3. How does CASCADE DELETE work? What are its risks?
4. Why must soft-deleted rows be filtered with `WHERE deleted_at IS NULL`?
5. How can you reclaim disk space after a bulk DELETE in SQLite?

## NEXT UP

Congratulations on completing Module 09! You now understand SQL, PostgreSQL, and building persistent Go applications. Next up: Module 10 — Auth and Security.
