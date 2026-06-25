# Transactions

## Learning objective

Group multiple SQL operations into atomic transactions using `db.Begin`, `tx.Commit`, and `tx.Rollback`, and understand transaction isolation in Go's `database/sql`.

## Why this matters

A bank transfer debits one account and credits another. If the debit succeeds but the credit fails, money disappears. Transactions guarantee that either both operations succeed or neither does. In production systems, transactions are the foundation of data integrity — they prevent partial updates, maintain referential consistency, and allow safe concurrent access.

## Mental model

A transaction is a staging area. You perform work on a copy of the data. No one else sees the changes until you commit. If something goes wrong, you roll back and the data reverts to its original state.

```
Normal flow:  Begin → Operation 1 → Operation 2 → Commit    → changes visible
Error flow:   Begin → Operation 1 → Operation 2 → Rollback  → all changes discarded
```

Inside a transaction, you use `tx.Exec` and `tx.Query` instead of `db.Exec` and `db.Query`. All operations use the same database connection.

## Core idea

`db.Begin()` returns a `*sql.Tx` representing a transaction. The `*sql.Tx` has its own `Exec`, `Query`, `QueryRow`, `Prepare`, and `Stmt` methods. After all operations complete, call `tx.Commit()`. On any error, call `tx.Rollback()`.

```go
tx, err := db.Begin()
if err != nil {
    return err
}
// Use tx.Exec, tx.Query, etc.
err = tx.Commit()   // or tx.Rollback()
```

## Under the hood

1. `db.Begin()` acquires a connection from the pool and sends `BEGIN TRANSACTION` (or `BEGIN`).
2. All subsequent `tx.Exec`/`tx.Query` calls use that same connection — it is pinned for the transaction's duration.
3. `tx.Commit()` sends `COMMIT` and releases the connection back to the pool.
4. `tx.Rollback()` sends `ROLLBACK` and releases the connection.
5. If the transaction goes out of scope without Commit or Rollback, the connection is not released until garbage collection, causing a pool leak.

With SQLite, modernc.org/sqlite implements transactions using `sqlite3_exec(db, "BEGIN TRANSACTION")` and `sqlite3_exec(db, "COMMIT")`. SQLite's default isolation level is serializable (deferred until write).

## How Go uses it

- **Banking/financial**: Debit one account, credit another in one transaction.
- **Order processing**: Create order, decrement inventory, charge payment — all or nothing.
- **Batch imports**: Insert thousands of rows in a single transaction for performance.
- **Schema migrations**: Apply multiple DDL statements atomically.

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

	_, err = db.Exec(`CREATE TABLE accounts (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		balance REAL NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO accounts VALUES
		(1, 'Alice', 1000.00),
		(2, 'Bob', 500.00)`)
	if err != nil {
		log.Fatal(err)
	}

	err = Transfer(db, 1, 2, 200.00)
	if err != nil {
		log.Fatal(err)
	}

	var aliceBalance, bobBalance float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&aliceBalance)
	db.QueryRow("SELECT balance FROM accounts WHERE id = 2").Scan(&bobBalance)
	fmt.Printf("Alice: $%.2f, Bob: $%.2f\n", aliceBalance, bobBalance)
}

func Transfer(db *sql.DB, fromID, toID int, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var fromBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", fromID).Scan(&fromBalance)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("read from account: %w", err)
	}
	if fromBalance < amount {
		tx.Rollback()
		return fmt.Errorf("insufficient funds: %.2f < %.2f", fromBalance, amount)
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("debit from account %d: %w", fromID, err)
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("credit to account %d: %w", toID, err)
	}

	return tx.Commit()
}
```

## Step-by-step execution

For `Transfer(db, 1, 2, 200.00)`:

1. `db.Begin()` sends `BEGIN TRANSACTION`, pins connection C1.
2. `tx.QueryRow("SELECT balance FROM accounts WHERE id = 1")` runs on C1. Returns 1000.00.
3. Check passes (1000 >= 200).
4. `tx.Exec("UPDATE accounts SET balance = balance - 200 WHERE id = 1")` runs on C1. Alice's balance is now 800.00 (visible only within the transaction).
5. `tx.Exec("UPDATE accounts SET balance = balance + 200 WHERE id = 2")` runs on C1. Bob's balance is now 700.00 (visible only within the transaction).
6. `tx.Commit()` sends `COMMIT`. Both changes become visible. Connection C1 returns to pool.

If step 4 succeeds but step 5 fails (e.g., invalid account), the rollback in step 5's error handler reverts the debit. The database is consistent.

## Common mistakes

- Not rolling back on error: The transaction stays open, holding the connection. Always `Rollback()` in error paths.
- Using `db.Exec` instead of `tx.Exec` inside a transaction: `db.Exec` uses a different connection, bypassing the transaction.
- Long-running transactions: Holding a transaction open blocks concurrent access on that row (or table, depending on isolation level). Keep transactions short.
- Nested transactions: Go's `database/sql` does not support savepoints natively. You must manage nested logic within a single `*sql.Tx`.
- Ignoring the error from `tx.Commit`: Commit can fail (e.g., serialization error). Check it.

## Debugging walkthrough

**Scenario**: A transfer silently succeeds but the balances don't change.

```go
func Transfer(db *sql.DB, fromID, toID int, amount float64) error {
	tx, _ := db.Begin() // error ignored
	tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
	tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
	return tx.Commit()
}
```

**Symptom**: No error, but balances are unchanged.

**Root cause**: The accounts table may not exist — `db.Exec("CREATE TABLE ...")` failed earlier (error ignored). The UPDATEs silently affect 0 rows. `tx.Commit()` succeeds because no changes were made.

**Fix**: Check every error. Use a transaction helper (next lesson) to centralize error handling.

```go
func Transfer(db *sql.DB, fromID, toID int, amount float64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// ... operations, setting err on failure
	return tx.Commit()
}
```

## Production notes

- **Isolation levels**: PostgreSQL supports `READ COMMITTED` (default), `REPEATABLE READ`, `SERIALIZABLE`. SQLite defaults to `DEFERRED`. Use `db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})` when needed.
- **Retry logic**: Serializable transactions may fail with `40001` (serialization failure). Retry the entire transaction up to 3 times.
- **Monitoring**: Track transaction duration, commit rate, and rollback rate. A high rollback rate indicates application errors or contention.
- **Connection pool impact**: Transactions hold one connection for their entire duration. Long transactions starve the pool.

## Performance implications

- Wrapping multiple INSERTs in a single transaction is much faster than auto-commit per row (SQLite can be 100x faster).
- Each commit in SQLite triggers a `fsync`. Batch commits for throughput.
- Transactions in WAL mode (SQLite) allow concurrent readers during a write transaction.
- Pessimistic vs. optimistic locking: transactions use pessimistic locking by default. For high-contention workloads, consider optimistic concurrency with retry.

## Practice task

Write a function `CreateOrder(db *sql.DB, userID int, productIDs []int) error` that:
1. Begins a transaction.
2. Inserts a row into `orders` (user_id, total).
3. For each product, inserts into `order_items` and decrements `stock` in `products`.
4. Commits on success, rolls back on any error.
5. Returns an error if any product has insufficient stock.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/17-transactions
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/17-transactions
```

## Review questions

1. What happens to a connection held by a transaction if you forget to Commit or Rollback?
2. Can you use `db.Exec` inside a transaction? What would happen?
3. Why must you check the error from `tx.Commit()`?
4. What is the default transaction isolation level in PostgreSQL? In SQLite?
5. How does wrapping 1000 INSERTs in a single transaction improve performance?

## NEXT UP

Rollback discipline — patterns for safe, predictable transaction rollback with deferred functions.
