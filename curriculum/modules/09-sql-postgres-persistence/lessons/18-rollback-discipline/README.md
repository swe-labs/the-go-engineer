# Rollback discipline

## Learning objective

Apply deferred rollback patterns, handle no-op rollback after commit, and write reusable transaction helper functions that eliminate rollback boilerplate.

## Why this matters

The most common transaction bug in Go is forgetting to roll back on error. The second most common is calling rollback after commit (which is harmless but noisy). Without a disciplined pattern, every transaction function repeats the same error-handling logic — and a single missed rollback leaks a connection forever. A deferred rollback wrapper makes transactions correct by construction.

## Mental model

A transaction is a resource that must be closed, like a file or a network connection. The `defer` keyword exists precisely for this pattern: schedule cleanup the moment you acquire the resource.

```
Acquire:  tx := db.Begin()
          ↓
Schedule: defer rollback if not committed
          ↓
Use:      tx.Exec / tx.Query
          ↓
Release:  tx.Commit()    → marks rollback as no-op
          OR error        → defer triggers rollback
```

The key insight: schedule the rollback *before* doing any work, but make it conditional on whether commit has already happened.

## Core idea

The standard pattern uses a named return value and `defer`:

```go
func DoSomething(db *sql.DB) (err error) {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()  // only rolls back if err is non-nil
        }
    }()
    // ... do work, setting err on failure ...
    return tx.Commit()
}
```

When `tx.Commit()` succeeds, `err` is nil, so the deferred rollback does nothing. When any operation fails, `err` is set, and the deferred rollback runs. This guarantees cleanup without duplicating rollback calls.

## Under the hood

`tx.Rollback()` after `tx.Commit()` is safe. The Go `database/sql` driver tracks whether the transaction has been committed or rolled back. A second rollback returns `sql.ErrTxDone`. This means the deferred rollback can safely run unconditionally — it becomes a no-op after commit.

```go
tx.Commit()          // marks tx as done
tx.Rollback()        // returns sql.ErrTxDone — harmless
```

The `defer` pattern above checks `err` first, avoiding even the no-op call, but calling it unconditionally is also safe:

```go
defer tx.Rollback()  // safe: no-op if already committed
```

## How Go uses it

- **Every transactional function**: Use the deferred rollback pattern as a standard idiom.
- **Transaction helper**: Extract the pattern into a reusable `WithTx` function.
- **Nested transactions**: Combine with savepoints or restructure to avoid nesting.
- **Read-only transactions**: Use `db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})` with deferred rollback.

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

	_, err = db.Exec(`CREATE TABLE items (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		log.Fatal(err)
	}

	WithTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "Widget", 10)
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO items (name, quantity) VALUES (?, ?)", "Gadget", 5)
		return err
	})

	var count int
	db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	fmt.Printf("Items inserted: %d\n", count)
}

func WithTx(db *sql.DB, fn func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
```

The `WithTx` helper centralizes:
- Beginning the transaction
- Deferred rollback on error or panic
- Committing on success
- Returning the error

Any function that needs a transaction just passes its logic as a closure.

## Step-by-step execution

For `WithTx(db, fn)` where `fn` succeeds:

1. `db.Begin()` acquires connection, returns `tx`.
2. `defer` closure captures `err` (currently nil).
3. `fn(tx)` runs: inserts succeed, returns nil.
4. `err = fn(tx)` → `err` is nil.
5. `return tx.Commit()` → commit succeeds, returns nil.
6. Deferred function runs: `err` is nil, so rollback is skipped.
7. `WithTx` returns nil.

For `WithTx(db, fn)` where `fn` fails:

1. `db.Begin()` acquires connection.
2. `defer` closure captures `err`.
3. `fn(tx)` fails on second insert: returns error.
4. `err = fn(tx)` → `err` is non-nil.
5. `return tx.Commit()` is not reached; returns `err`.
6. Deferred function runs: `err` is non-nil → `tx.Rollback()` runs.
7. `WithTx` returns the error.

## Common mistakes

- Not using `defer` for rollback: Every error path must manually call `tx.Rollback()`. One missed path = leaked connection.
- Calling `tx.Rollback()` unconditionally after commit: Returns `sql.ErrTxDone`. Harmless but noisy. Check `err` first.
- Forgetting to include panic recovery in `defer`: A panic before commit leaves the transaction open. The `defer` pattern should recover panics and roll back.
- Using `db.Begin()` instead of `db.BeginTx(ctx)` in production: `BeginTx` allows context cancellation to abort the transaction.
- Passing `*sql.DB` to the closure instead of `*sql.Tx`: Inside the closure, use `tx.Exec`, not `db.Exec`.

## Debugging walkthrough

**Scenario**: The application suddenly has no available connections. `db.Stats().InUse` shows all connections busy.

```go
func UpdateBalance(db *sql.DB, id int, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// OOPS: no deferred rollback
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, id)
	if err != nil {
		return err // transaction leaked!
	}
	return tx.Commit()
}
```

**Root cause**: When `tx.Exec` fails, the function returns immediately without rolling back. The transaction holds a connection forever. After enough failures, the pool exhausts.

**Fix**: Add the deferred rollback pattern:

```go
func UpdateBalance(db *sql.DB, id int, amount float64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}
```

Now every error path automatically rolls back.

## Production notes

- **`WithTx` in libraries**: Consider providing a `WithTx` helper in your internal packages to enforce consistent transaction handling across the codebase.
- **Context support**: Use `db.BeginTx(ctx, nil)` and pass `ctx` to the closure for cancellation support.
- **Nested transactions**: Go does not support savepoints. If you need nesting, restructure the code or use a driver-specific escape. The `WithTx` pattern does not nest well — avoid calling `WithTx` inside `WithTx`.
- **Read-only transactions**: Use `db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})` for read-only work. PostgreSQL can use this to avoid acquiring locks.

## Performance implications

- The `defer` overhead is negligible (nanoseconds).
- The rollback itself is a round-trip to the database, but it only happens on error paths — not the happy path.
- Using a `WithTx` helper reduces code duplication and makes transactions easier to review.
- A rolled-back transaction wastes the work already done. Keep transactions short to minimize wasted work on failure.

## Practice task

Write a `WithTxV2` function that additionally:
1. Supports `context.Context` via `BeginTx`.
2. Logs a warning (using `log.Printf`) when a rollback occurs.
3. Returns the rollback error wrapped if commit succeeds but rollback fails.

Then use it in a function that inserts an order with multiple items, deliberately failing on the third item to verify rollback.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/18-rollback-discipline
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/18-rollback-discipline
```

## Review questions

1. Why does the deferred rollback pattern use a named return value `err`?
2. What happens if you call `tx.Rollback()` after `tx.Commit()`?
3. Why should the `defer` function check for panics?
4. How does `WithTx` prevent connection leaks?
5. Can you call `WithTx` inside another `WithTx`? What are the risks?

## NEXT UP

Migrations — schema changes applied in versioned, repeatable steps.
