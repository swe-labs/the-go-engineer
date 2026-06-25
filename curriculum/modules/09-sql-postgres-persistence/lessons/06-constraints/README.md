# Constraints

## Learning objective

Apply column and table constraints (`NOT NULL`, `UNIQUE`, `CHECK`, `DEFAULT`) to enforce data quality at the database level, and contrast database constraints with Go-side validation.

## Why this matters

Application code has bugs. APIs receive malformed input. Migration scripts run out of order. Data constraints at the database level are the last line of defense against corrupt data. A `NOT NULL` constraint is worth a thousand nil checks. A `CHECK` constraint prevents a negative price even if the frontend forgets to validate. Go engineers who layer database constraints under application validation build systems that stay correct despite chaos.

## Mental model

Think of constraints as bouncers at a club. Each bouncer enforces one rule:
- `NOT NULL`: "You must have this item to enter."
- `UNIQUE`: "Only one person can have this ID."
- `CHECK`: "You must be at least 21 years old."
- `DEFAULT`: "If you don't have an ID, we'll give you a temporary one."

If any bouncer rejects a row, the entire row is turned away (the INSERT or UPDATE is rolled back).

## Core idea

SQL constraints are rules attached to columns or tables that the database engine enforces on every INSERT, UPDATE, or DELETE.

| Constraint | Enforces | Example |
|---|---|---|
| `NOT NULL` | Column cannot contain NULL | `name TEXT NOT NULL` |
| `UNIQUE` | All values in column(s) must be distinct | `email TEXT UNIQUE` |
| `PRIMARY KEY` | NOT NULL + UNIQUE | `id INTEGER PRIMARY KEY` |
| `FOREIGN KEY` | Value must exist in another table | `dept_id INT REFERENCES departments(id)` |
| `CHECK` | Arbitrary boolean expression | `CHECK (price >= 0)` |
| `DEFAULT` | Default value when none provided | `stock INTEGER DEFAULT 0` |

Constraints can be named for clearer error messages and easier management:

```sql
CREATE TABLE products (
    id INTEGER PRIMARY KEY,
    price REAL NOT NULL CONSTRAINT positive_price CHECK (price >= 0)
);
```

## Under the hood

Constraints are stored in SQLite's `sqlite_master` table (or PostgreSQL's `pg_constraint`). On every DML statement, the database evaluates each applicable constraint in order:

1. `NOT NULL` — fast pointer comparison (is the value nil?).
2. `CHECK` — evaluates the expression (may involve function calls).
3. `UNIQUE` — searches the unique index (O(log n)).
4. `FOREIGN KEY` — probes the parent table index (O(log n)).
5. `PRIMARY KEY` — combination of NOT NULL and unique index.

If any constraint fails, the statement returns an error and all changes in the statement are rolled back. No partial modifications reach the table.

## How Go uses it

Go handles constraint violations as errors from `db.Exec` or `db.Query`. The error type is driver-specific: SQLite returns a `sqlite.Error` with a message like `UNIQUE constraint failed: users.email`. In production code, check for constraint violation errors and map them to appropriate HTTP status codes or business logic responses.

Go-side validation (e.g., checking `len(name) > 0` before INSERT) and database constraints are complementary, not alternatives. Go validation catches errors early (before a round trip to the database), while database constraints are the safety net.

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

type User struct {
	ID    int64
	Name  string
	Email string
	Age   int
}

func createUserTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			age INTEGER NOT NULL CHECK(age >= 0 AND age <= 150)
		)
	`)
	return err
}

func insertUser(db *sql.DB, u *User) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO users (name, email, age) VALUES (?, ?, ?)`,
		u.Name, u.Email, u.Age,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createUserTable(db); err != nil {
		log.Fatal(err)
	}

	users := []User{
		{Name: "Alice", Email: "alice@example.com", Age: 30},
		{Name: "Bob", Email: "bob@example.com", Age: 25},
	}

	for _, u := range users {
		id, err := insertUser(db, &u)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				fmt.Printf("Skipped duplicate email: %s\n", u.Email)
				continue
			}
			log.Fatal(err)
		}
		fmt.Printf("Inserted %s with id=%d\n", u.Name, id)
	}

	// Try violating NOT NULL
	_, err = db.Exec(`INSERT INTO users (name, email, age) VALUES (NULL, 'null@test.com', 20)`)
	if err != nil {
		fmt.Println("NOT NULL violation:", err)
	}

	// Try violating CHECK
	_, err = db.Exec(`INSERT INTO users (name, email, age) VALUES ('Test', 'test@test.com', -5)`)
	if err != nil {
		fmt.Println("CHECK violation:", err)
	}

	// Try violating UNIQUE
	_, err = insertUser(db, &User{Name: "Alice2", Email: "alice@example.com", Age: 35})
	if err != nil {
		fmt.Println("UNIQUE violation:", err)
	}
}
```

## Step-by-step execution

1. The `users` table is created with `name TEXT NOT NULL`, `email TEXT NOT NULL UNIQUE`, and `age INTEGER NOT NULL CHECK(age >= 0 AND age <= 150)`.
2. Two valid users are inserted successfully.
3. An attempt to insert `NULL` name triggers the `NOT NULL` constraint. The database returns an error.
4. An attempt to insert `age = -5` triggers the `CHECK` constraint. The database returns an error.
5. An attempt to insert a duplicate email triggers the `UNIQUE` constraint. The database returns an error.

## Common mistakes

- **Relying only on Go validation without DB constraints**: A bug in the validation logic lets invalid data through. Always add database constraints as a safety net.
- **Using `CHECK` with functions that are not deterministic**: SQLite allows `CHECK( random() > 0.5 )`, which is evaluated per row and may produce inconsistent results. Avoid non-deterministic expressions.
- **Creating a `UNIQUE` constraint on a nullable column**: SQLite allows multiple NULLs in a UNIQUE column (NULL != NULL). This surprises developers coming from PostgreSQL, where this is configurable.
- **Putting constraints on optional columns inappropriately**: A column with `DEFAULT 0` and `NOT NULL` prevents users from omitting the column. Only use `NOT NULL` when every row truly must have a value.

## Debugging walkthrough

```go
_, err := db.Exec(`INSERT INTO users (name, email) VALUES ('Test', 'test@test.com')`)
```

**Error**: `NOT NULL constraint failed: users.age`

**Root cause**: The `age` column has `NOT NULL` but no `DEFAULT` value. The INSERT omits `age`, so SQLite tries to insert NULL, which violates the constraint.

**Fix**: Either provide a `DEFAULT` value for `age` or include it in the INSERT:

```go
_, err := db.Exec(`INSERT INTO users (name, email, age) VALUES ('Test', 'test@test.com', 0)`)
```

Or:
```sql
CREATE TABLE users (..., age INTEGER NOT NULL DEFAULT 0 CHECK(age >= 0))
```

## Production notes

- Name constraints explicitly: `CONSTRAINT positive_price CHECK (price >= 0)` makes error messages clearer.
- Use `UNIQUE` constraints on columns that should be logically unique (email, username) even if application code also checks uniqueness. The database constraint prevents race conditions where two concurrent requests insert the same value.
- `CHECK` constraints are a good place for business rules that should never change (e.g., `CHECK (status IN ('active', 'inactive'))`). For rules that change frequently, use application-level validation.
- Monitor constraint violation errors in production. A high rate of unique violations may indicate a bot attack or a bug in the client.

## Performance implications

- `NOT NULL` and `DEFAULT` have zero runtime cost — they are metadata checks.
- `CHECK` evaluates the expression once per row; simple comparisons are free, function calls add overhead.
- `UNIQUE` and `PRIMARY KEY` maintain an index. Index updates add O(log n) cost per write. The read benefit usually outweighs the write cost.
- Multiple constraints on the same column are checked independently; each adds its own overhead.

## Practice task

Create a table `accounts` with constraints:
- `id INTEGER PRIMARY KEY`
- `owner TEXT NOT NULL`
- `balance INTEGER NOT NULL CHECK(balance >= 0)` (store cents)
- `currency TEXT NOT NULL CHECK(currency IN ('USD', 'EUR', 'GBP'))`
- `created_at TEXT NOT NULL DEFAULT (datetime('now'))`

Write Go code that:
1. Inserts valid accounts.
2. Attempts inserts that violate each constraint and verifies the errors.
3. Queries all accounts and prints them.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/06-constraints
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/06-constraints
```

## Review questions

1. What is the difference between `NOT NULL` and `DEFAULT`?
2. Why would you add a `UNIQUE` constraint in the database when your Go code already validates uniqueness?
3. Can a `CHECK` constraint reference other rows or tables?
4. What happens to an INSERT that violates multiple constraints?
5. What is the benefit of naming constraints?

## NEXT UP

SQLite as a local learning tool — why SQLite is ideal for development and testing.
