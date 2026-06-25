# UPDATE operations

## Learning objective

Update existing rows in SQLite from Go using `UPDATE SET`, `db.Exec`, and prepared statements, implement optimistic locking with version numbers, manage `updated_at` timestamps, and handle partial updates.

## Why this matters

Writing correct UPDATE queries is harder than INSERT or SELECT because you must identify exactly which row(s) to modify, set only the intended columns, and handle concurrent modifications. A missing WHERE clause updates every row. A stale UPDATE overwrites another transaction's change. Production systems use optimistic locking and partial update patterns to prevent data corruption.

## Mental model

An UPDATE is a targeted mutation: find a row (WHERE), change specific columns (SET), and verify the expected state before updating.

```
BEFORE:                  AFTER:
┌──────┬───────┬──────┐  ┌──────┬───────┬──────┐
│ id   │ name  │ qty  │  │ id   │ name  │ qty  │
├──────┼───────┼──────┤  ├──────┼───────┼──────┤
│ 1    │Widget │  10  │→ │ 1    │Widget │  15  │  ← qty updated
│ 2    │Gadget │   5  │  │ 2    │Gadget │   5  │
└──────┴───────┴──────┘  └──────┴───────┴──────┘

SQL: UPDATE products SET qty = 15 WHERE id = 1;
                          \______/         \_____/
                          columns to       which row
                          change
```

## Core idea

Basic UPDATE in Go:

```go
result, err := db.Exec("UPDATE products SET price = ? WHERE id = ?", newPrice, id)
rows, _ := result.RowsAffected()
// rows == 1 if the row existed, 0 if not found
```

Optimistic locking with version:

```go
result, err := db.Exec(
    "UPDATE products SET price = ?, version = version + 1 WHERE id = ? AND version = ?",
    newPrice, id, expectedVersion)
rows, _ := result.RowsAffected()
if rows == 0 {
    // Someone else modified the row — retry or return conflict error
}
```

Partial updates (only update fields that changed):

```go
updates := []string{}
args := []any{}
if p.Name != "" {
    updates = append(updates, "name = ?")
    args = append(args, p.Name)
}
if p.Price != 0 {
    updates = append(updates, "price = ?")
    args = append(args, p.Price)
}
query := "UPDATE products SET " + strings.Join(updates, ", ") + " WHERE id = ?"
args = append(args, id)
```

## Under the hood

An UPDATE in SQLite:
1. SQLite locates the row(s) using the WHERE clause (with or without index).
2. It acquires a write lock on the database (or the relevant pages in WAL mode).
3. It modifies the row data in-place (or marks the old row as deleted and writes a new one).
4. Any indexes on the updated columns are updated.
5. The change is written to the journal (for atomicity).

`db.Exec` returns `sql.Result` with:
- `LastInsertId()` — for INSERT, not meaningful for UPDATE.
- `RowsAffected()` — number of rows modified (0 if WHERE matched nothing).

## How Go uses it

- **PATCH endpoints**: Partial update of resource fields.
- **Inventory management**: Decrement stock on order.
- **Profile updates**: Change email, name, or password.
- **Optimistic locking**: Prevent lost updates in concurrent systems.
- **Batch updates**: Update many rows with the same SET clause (e.g., mark all notifications as read).

## Go example

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Product struct {
	ID      int
	Name    string
	Price   float64
	Stock   int
	Version int
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE products (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			stock INTEGER NOT NULL DEFAULT 0,
			version INTEGER NOT NULL DEFAULT 1
		);
		INSERT INTO products VALUES (1, 'Widget', 9.99, 10, 1);
		INSERT INTO products VALUES (2, 'Gadget', 24.99, 5, 1);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Simple update
	result, err := db.Exec("UPDATE products SET price = ? WHERE id = ?", 12.99, 1)
	if err != nil {
		log.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("Updated %d row(s)\n", rows)

	// Optimistic locking update
	expectedVersion := 1
	newStock := 15
	result, err = db.Exec(
		"UPDATE products SET stock = ?, version = version + 1 WHERE id = ? AND version = ?",
		newStock, 1, expectedVersion)
	if err != nil {
		log.Fatal(err)
	}
	rows, _ = result.RowsAffected()
	if rows == 0 {
		fmt.Println("Optimistic lock failed — row was modified by another transaction")
	} else {
		fmt.Println("Optimistic lock succeeded")
	}

	// Partial update
	type PartialUpdate struct {
		Name  *string
		Price *float64
		Stock *int
	}
	update := PartialUpdate{
		Name:  nil,
		Price: ptr(19.99),
		Stock: ptr(20),
	}
	rowsAffected, err := PartialUpdateProduct(db, 2, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Partial update affected %d row(s)\n", rowsAffected)

	printProduct(db, 1)
	printProduct(db, 2)
}

func PartialUpdateProduct(db *sql.DB, id int, p PartialUpdate) (int64, error) {
	query := "UPDATE products SET "
	var updates []string
	var args []any

	if p.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *p.Name)
	}
	if p.Price != nil {
		updates = append(updates, "price = ?")
		args = append(args, *p.Price)
	}
	if p.Stock != nil {
		updates = append(updates, "stock = ?")
		args = append(args, *p.Stock)
	}

	if len(updates) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	query += join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func join(items []string, sep string) string {
	result := ""
	for i, s := range items {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func ptr[T any](v T) *T { return &v }

func printProduct(db *sql.DB, id int) {
	var p Product
	err := db.QueryRow("SELECT id, name, price, stock, version FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Product %d: %s $%.2f stock=%d version=%d\n", p.ID, p.Name, p.Price, p.Stock, p.Version)
}

var _ = context.Background
```

## Step-by-step execution

For `UPDATE products SET price = 12.99 WHERE id = 1`:

1. SQLite's query planner finds `id = 1` (using primary key index).
2. The row with `id=1` is located (1 page read).
3. A write lock is acquired on the database page.
4. The `price` field is updated to 12.99.
5. Any indexes on `price` are updated.
6. The change is written to the journal.
7. `RowsAffected()` returns 1.
8. The lock is released.

If `WHERE id = 999` matched nothing: `RowsAffected()` returns 0. No lock is acquired (row-level lock not needed for no rows).

## Common mistakes

- **Missing WHERE clause**: `UPDATE products SET price = 0` updates every row. Always double-check WHERE.
- **Not checking RowsAffected**: An UPDATE that matches zero rows is not an error. Check `RowsAffected()` to verify the intended row was updated.
- **Stale updates (no optimistic locking)**: Two concurrent requests can overwrite each other's changes. Use `version` column and `WHERE version = ?`.
- **Partial update building SQL with string concatenation**: Risk of SQL injection if field values contain SQL. Use parameterized queries.
- **Forgetting `updated_at`**: Many tables benefit from an `updated_at` timestamp. Set it in the UPDATE: `SET price = ?, updated_at = CURRENT_TIMESTAMP`.

## Debugging walkthrough

**Scenario**: An UPDATE silently affects 0 rows, but the row exists.

```go
result, err := db.Exec("UPDATE products SET stock = ? WHERE id = ?", newStock, id)
```

**Root cause**: Time passes between selecting the row and updating it. Another transaction deleted the row (soft delete with `deleted_at IS NOT NULL`). The SELECT included `WHERE deleted_at IS NULL`, but the UPDATE did not.

**Fix**: Use the same WHERE clause in the UPDATE, or use optimistic locking with the retrieved version.

```go
result, err := db.Exec(
    "UPDATE products SET stock = ? WHERE id = ? AND deleted_at IS NULL",
    newStock, id)
```

## Production notes

- **Batch updates**: Use a transaction for multiple UPDATEs. Each UPDATE in auto-commit mode is a separate transaction.
- **Index maintenance**: Updating an indexed column forces an index update. Measure the write workload.
- **Lock contention**: Long-running transactions holding UPDATE locks block other writers. Keep transactions short.
- **Returning data**: PostgreSQL supports `UPDATE ... RETURNING *` to return the updated row. SQLite does not — you must re-query.
- **Partial update patterns**: In REST APIs, use `PATCH` with JSON Merge Patch (RFC 7396) or JSON Patch (RFC 6902) for partial updates.

## Performance implications

- Updating a row with an index on the changed column requires O(log n) index maintenance.
- Updating many rows with a single `UPDATE ... WHERE id IN (...)` is faster than individual UPDATEs.
- SQLite's WAL mode allows concurrent reads during an UPDATE.
- `RowsAffected()` requires a round-trip to the driver. For single-row UPDATEs, it's cheap.

## Practice task

Write a function `UpdateProductStock(db *sql.DB, id, newStock, expectedVersion int) (bool, error)` that:
1. Uses optimistic locking with `version`.
2. Returns `true` if the update succeeded, `false` if version conflict.
3. Returns the error if the query failed.
4. Test with concurrent goroutines simulating conflicting updates.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/27-update-operations
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/27-update-operations
```

## Review questions

1. What does `RowsAffected()` return when an UPDATE's WHERE clause matches no rows?
2. How does optimistic locking prevent lost updates?
3. Why is dynamic SQL generation for partial updates potentially dangerous?
4. What happens to indexes when an indexed column is updated?
5. How is a PATCH endpoint (partial update) different from a PUT endpoint (full replace)?

## NEXT UP

DELETE operations — hard delete, soft delete, CASCADE, and batch delete strategies.
