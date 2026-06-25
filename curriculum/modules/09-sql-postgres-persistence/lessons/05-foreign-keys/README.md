# Foreign keys

## Learning objective

Define foreign key constraints to enforce referential integrity, choose appropriate `ON DELETE` and `ON UPDATE` actions, and write Go code that respects and tests foreign key rules.

## Why this matters

Foreign keys are the mechanism that keeps related data consistent. Without them, an application can create orphaned rows (e.g., orders for a deleted customer) or miss cascading updates. Databases enforce foreign keys at the engine level, which is far more reliable than application-level checks. Go engineers who use foreign keys correctly sleep better knowing the database will reject invalid state no matter how many bugs exist in the application code.

## Mental model

A foreign key is a promise: "every value in this column must exist in that other table's primary key column." It is like a library catalog that only assigns a book to a shelf that actually exists. If someone removes a shelf, the catalog must either prevent the removal (RESTRICT), remove the book too (CASCADE), or set the shelf reference to NULL (SET NULL).

## Core idea

A **foreign key** is a column (or set of columns) that references the primary key (or unique column) of another table. It ensures:

- **Referential integrity**: no row in the child table can reference a non-existent parent row.
- **Cascading actions**: when a parent row is deleted or updated, the database automatically handles child rows according to the specified action.

Actions:

| Action | Behavior on parent DELETE/UPDATE |
|---|---|
| `NO ACTION` (default) | Prevents the operation if child rows exist (deferred check in some engines) |
| `RESTRICT` | Immediately prevents the operation if child rows exist |
| `CASCADE` | Deletes/updates child rows automatically |
| `SET NULL` | Sets the foreign key column(s) to NULL in child rows |
| `SET DEFAULT` | Sets the foreign key to the column's default value |

In SQLite, foreign key enforcement must be enabled at runtime with `PRAGMA foreign_keys = ON`.

## Under the hood

When a foreign key constraint is defined, the database creates an index on the referencing column(s) (if one does not already exist). On every INSERT or UPDATE of the child table, the database checks that the new foreign key value exists in the parent table. On every DELETE or UPDATE of the parent table, the database scans for referencing child rows and applies the configured action.

This check happens inside the same transaction as the DML statement. If the check fails, the statement is rolled back and an error is returned.

## How Go uses it

Go enables `PRAGMA foreign_keys = ON` via `db.Exec` immediately after connecting to SQLite. For PostgreSQL, foreign keys are enforced by default. The Go code does not need special handling for foreign keys — the `database/sql` error from a constraint violation is returned as a regular `error` from `Exec` or `Query`.

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

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE departments (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
		CREATE TABLE employees (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			dept_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO departments (name) VALUES ('Engineering')`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES ('Alice', 1)`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted department and employee successfully")

	// Try inserting with invalid dept_id
	_, err = db.Exec(`INSERT INTO employees (name, dept_id) VALUES ('Bob', 999)`)
	if err != nil {
		fmt.Println("Foreign key violation caught:", err)
	}

	// Cascade delete
	_, err = db.Exec(`DELETE FROM departments WHERE id = 1`)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM employees`).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Employees remaining after cascade delete: %d\n", count)
}
```

## Step-by-step execution

1. `PRAGMA foreign_keys = ON` enables foreign key enforcement for this session. SQLite defaults to OFF.
2. Two tables are created with `dept_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE CASCADE`.
3. A department (id=1) and an employee (dept_id=1) are inserted successfully.
4. Inserting an employee with `dept_id=999` fails because no department with id=999 exists.
5. Deleting the department with `ON DELETE CASCADE` automatically deletes the employee. The count query returns 0.

## Common mistakes

- **Forgetting `PRAGMA foreign_keys = ON`**: Foreign key violations are silently ignored in SQLite unless this pragma is set. Always execute it right after `sql.Open`.
- **Using `ON DELETE CASCADE` without understanding the implications**: Cascading deletes can remove far more data than expected. Use `RESTRICT` or `NO ACTION` as the default and only use `CASCADE` when the child data is meaningless without the parent.
- **Defining a foreign key that references a non-unique column**: The target column must be `PRIMARY KEY` or have a `UNIQUE` constraint.
- **Creating circular foreign key references**: Table A references B and B references A. This prevents inserting into either table.

## Debugging walkthrough

```go
db, _ := sql.Open("sqlite", ":memory:")
db.Exec(`CREATE TABLE parent (id INTEGER PRIMARY KEY)`)
db.Exec(`CREATE TABLE child (id INTEGER PRIMARY KEY, pid INTEGER REFERENCES parent(id))`)
_, err := db.Exec(`INSERT INTO child (pid) VALUES (1)`)
fmt.Println(err) // nil — but we expected an error!
```

**Symptom**: Inserting a child row with a non-existent parent ID does not error.

**Root cause**: `PRAGMA foreign_keys = ON` was never executed. SQLite silently accepts the violation.

**Fix**: Always set `PRAGMA foreign_keys = ON` in the connection setup.

```go
db.Exec(`PRAGMA foreign_keys = ON`)
```

## Production notes

- In PostgreSQL, foreign keys are enforced by default and do not require a pragma.
- Use `ON DELETE CASCADE` sparingly. Prefer `ON DELETE RESTRICT` (or `NO ACTION`) and let the application handle deletion order explicitly.
- Always index foreign key columns. Without an index, every DELETE or UPDATE on the parent table triggers a full table scan of the child table.
- Document cascading behavior in your schema migration files — future developers (including yourself) will thank you.

## Performance implications

- Enforcing foreign keys adds a read on the parent table for every INSERT/UPDATE on the child table, and a scan on the child table for every DELETE/UPDATE on the parent table.
- An index on the foreign key column turns the scan into an index lookup, which is O(log n) instead of O(n).
- The overhead is usually negligible (microseconds) but can become significant at tens of thousands of writes per second.

## Practice task

Create three tables: `customers` (id, name), `orders` (id, customer_id REFERENCES customers ON DELETE CASCADE, total), `order_items` (id, order_id REFERENCES orders ON DELETE CASCADE, product_name, price). Write Go code to:
1. Insert a customer, an order, and two order items.
2. Verify that a query joining all three tables returns the correct data.
3. Delete the customer and verify all related orders and items are also deleted (cascade).

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/05-foreign-keys
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/05-foreign-keys
```

## Review questions

1. What does `ON DELETE CASCADE` do?
2. Why does SQLite require `PRAGMA foreign_keys = ON`?
3. What happens if you try to INSERT a row with a foreign key value that does not exist in the parent?
4. When would you use `ON DELETE SET NULL` instead of `CASCADE`?
5. Why should foreign key columns be indexed?

## NEXT UP

Constraints — NOT NULL, UNIQUE, CHECK, DEFAULT, and how they enforce data quality at the database level.
