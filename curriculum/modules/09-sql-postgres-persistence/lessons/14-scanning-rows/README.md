# Scanning rows

## Learning objective

Use `rows.Scan` to map SQL result columns to Go variables, handle column ordering and type conversions, and debug common scanning errors.

## Why this matters

`Scan` is the bridge between the database's world of rows and columns and Go's world of typed variables. A mismatch in column count, order, or type produces runtime errors. Mastering Scan means you can reliably convert any query result into Go structs, handle NULL values, and write reusable scanning functions.

## Mental model

`Scan` is like a factory conveyor belt. Each row arrives as a sequence of untyped values. `Scan` picks each value off the belt and places it into a waiting Go variable, converting the type as needed. If the belt has 4 items but you only provided 3 boxes, the conversion fails. If an item is a "NULL" sticker but you provided a non-pointer box, the conversion fails.

## Core idea

`rows.Scan(dest ...interface{})` copies the current row's column values into the provided Go variables. Key rules:

1. **Column order must match**: The first column in SELECT goes to the first argument, second to the second, etc.
2. **Type conversion is automatic**: SQLite `INTEGER` → Go `int`, `TEXT` → `string`, `REAL` → `float64`.
3. **NULL handling**: Scanning a NULL into a non-pointer type (like `string`) returns an error. Use pointer types or `sql.NullXxx` types.
4. **Extra columns cause errors**: If SELECT returns 3 columns but Scan has 2 arguments, Scan returns an error.

Common Scan destinations:

| SQL type | Go scan target |
|---|---|
| `INTEGER` | `int`, `int64`, `uint64`, `sql.NullInt64` |
| `REAL` | `float64`, `sql.NullFloat64` |
| `TEXT` | `string`, `sql.NullString` |
| `BLOB` | `[]byte` |
| `NULL` any | `*T` (pointer), `sql.Null[T]` |

## Under the hood

When `Scan` is called, the driver converts the database-native value to a Go `interface{}` and passes it to `Scan`. The `Scan` function uses reflection to determine the target type and performs the conversion:

1. If the target is `*string` and the value is a `[]byte` (SQLite TEXT), it converts `[]byte` to `string`.
2. If the target is `*int64` and the value is `int64` (SQLite INTEGER), it assigns directly.
3. If the target is `*float64` and the value is `int64`, it converts.
4. If the target is `*sql.NullString` and the value is nil, it sets `Valid = false`.

The reflection-based conversion is safe but has a small overhead (hundreds of nanoseconds per column).

## How Go uses it

Go programs typically scan into struct fields or local variables:

```go
// Struct scan
type User struct {
    ID    int
    Name  string
    Email string
}
var u User
err := row.Scan(&u.ID, &u.Name, &u.Email)

// Anonymous scan
var id int
var name string
err := row.Scan(&id, &name)
```

For reusable scanning, write helper functions:

```go
func scanUser(row *sql.Row) (User, error) {
    var u User
    err := row.Scan(&u.ID, &u.Name, &u.Email)
    return u, err
}
```

## Go example

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Employee struct {
	ID     int
	Name   string
	Salary float64
	Active bool
}

func scanEmployee(row *sql.Row) (Employee, error) {
	var e Employee
	err := row.Scan(&e.ID, &e.Name, &e.Salary, &e.Active)
	return e, err
}

func scanEmployees(rows *sql.Rows) ([]Employee, error) {
	defer rows.Close()
	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Salary, &e.Active); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, rows.Err()
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE employees (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		salary REAL NOT NULL,
		active INTEGER NOT NULL DEFAULT 1
	)`)

	db.Exec(`INSERT INTO employees (name, salary, active) VALUES ('Alice', 75000, 1)`)
	db.Exec(`INSERT INTO employees (name, salary, active) VALUES ('Bob', 82000, 1)`)
	db.Exec(`INSERT INTO employees (name, salary, active) VALUES ('Carol', 0, 0)`)

	fmt.Println("=== Scan into struct fields ===")
	var e Employee
	err = db.QueryRow(`SELECT id, name, salary, active FROM employees WHERE id = 1`).Scan(
		&e.ID, &e.Name, &e.Salary, &e.Active,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Employee: %d: %s ($%.0f, active=%v)\n", e.ID, e.Name, e.Salary, e.Active)

	fmt.Println("\n=== Scan multiple rows with helper ===")
	rows, err := db.Query(`SELECT id, name, salary, active FROM employees ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}

	employees, err := scanEmployees(rows)
	if err != nil {
		log.Fatal(err)
	}
	for _, emp := range employees {
		fmt.Printf("  %d: %s ($%.0f, active=%v)\n", emp.ID, emp.Name, emp.Salary, emp.Active)
	}

	fmt.Println("\n=== NULL handling demo ===")
	db.Exec(`CREATE TABLE IF NOT EXISTS contacts (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT
	)`)
	db.Exec(`INSERT INTO contacts (name, phone) VALUES ('Dave', '555-0100')`)
	db.Exec(`INSERT INTO contacts (name, phone) VALUES ('Eve', NULL)`)

	rows, err = db.Query(`SELECT id, name, phone FROM contacts ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var phone sql.NullString
		if err := rows.Scan(&id, &name, &phone); err != nil {
			log.Fatal(err)
		}
		if phone.Valid {
			fmt.Printf("  %d: %s (phone=%s)\n", id, name, phone.String)
		} else {
			fmt.Printf("  %d: %s (no phone)\n", id, name)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
```

## Step-by-step execution

1. Three employees are inserted, including one with `salary=0` and `active=0`.
2. `db.QueryRow(...).Scan(&e.ID, &e.Name, &e.Salary, &e.Active)` maps the four SELECT columns to the four struct fields. Column order must match: `id → ID`, `name → Name`, `salary → Salary`, `active → Active`.
3. `scanEmployees` is a reusable helper that iterates multiple rows and scans each into an `Employee` struct.
4. The NULL handling example uses `sql.NullString` for the nullable `phone` column. When the database value is NULL, `phone.Valid` is `false`. When non-NULL, `phone.String` contains the value and `phone.Valid` is `true`.

## Common mistakes

- **Column count mismatch**: SELECT returns 4 columns but Scan provides 3 arguments → `sql: expected 4 destination arguments, not 3`.
- **Column order mismatch**: SELECT is `id, name, salary` but Scan is `&name, &id, &salary` → type conversion error (name is string, id is int).
- **Scanning NULL into non-pointer type**: If `phone TEXT` is NULL and you scan into `var phone string`, Scan returns `sql: Scan error on column index 2: converting NULL to string is unsupported`. Fix: use `sql.NullString`.
- **Reusing the same variable in a loop without creating new structs**: When appending to a slice, each appended element shares the same memory if you scan into the same struct. Fix: create a new struct in each iteration.

## Debugging walkthrough

```go
rows, err := db.Query(`SELECT name, id FROM employees`)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var e Employee
    err := rows.Scan(&e.ID, &e.Name) // name → ID (string → int), id → Name (int → string)
    if err != nil {
        log.Fatal(err) // fails here
    }
}
```

**Symptom**: `Scan error on column index 0: converting TEXT "Alice" to int is unsupported`.

**Root cause**: SELECT is `name, id` but Scan expects `id, name`. Column 0 (name = "Alice") is being scanned into `&e.ID` (int), which fails.

**Fix**: Match Scan order to SELECT order: `SELECT id, name FROM employees` and `Scan(&e.ID, &e.Name)`.

## Production notes

- Always scan into named struct fields, not anonymous variables. This makes the mapping explicit and readable.
- Use `sql.NullString`, `sql.NullInt64`, `sql.NullFloat64`, `sql.NullBool`, `sql.NullTime` for nullable columns. Go 1.21+ also has `sql.Null[T]`.
- Write helper functions (`scanEmployee`) to avoid repeating Scan calls.
- For large result sets, scan in batches (e.g., scan 100 rows, process, scan next 100) to bound memory usage.
- In PostgreSQL, use `$1`, `$2` parameter markers. In SQLite, use `?`. The Scan code is identical.

## Performance implications

- Scanning is fast: ~100-500 ns per column, dominated by reflection for type conversion.
- Scanning into `sql.NullXxx` types adds a small overhead (checking `Valid` flag) compared to non-nullable types.
- Scanning into `interface{}` (e.g., `var v interface{}; row.Scan(&v)`) is slower because the value is not type-asserted.
- For maximum performance in hot paths, scan into concrete types and avoid `sql.NullXxx` by ensuring columns are `NOT NULL`.

## Practice task

Write a function `scanPerson(row *sql.Row) (Person, error)` that scans a row into:
```go
type Person struct {
    ID    int
    Name  string
    Email sql.NullString
    Age   sql.NullInt64
}
```

Then write a `main()` that:
1. Creates a `people` table (id, name, email, age).
2. Inserts two rows: one with all values, one with NULL email and NULL age.
3. Queries both rows and scans them using `scanPerson`.
4. Prints the results, handling NULLs gracefully.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/14-scanning-rows
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/14-scanning-rows
```

## Review questions

1. What happens if SELECT returns 4 columns but Scan provides 3 arguments?
2. Why does scanning a NULL value into a `string` variable fail?
3. How do you handle a nullable column in Go?
4. What is the difference between `sql.NullString` and a `*string` pointer?
5. What causes the error "converting TEXT to int is unsupported"?

## NEXT UP

Null handling — comprehensive coverage of NULL values in SQL and their Go counterparts.
