# Tables, rows, and columns

## Learning objective

Write `CREATE TABLE` statements with appropriate data types and column constraints, explain how rows and columns map to Go types, and query table metadata from SQLite's internal schema.

## Why this matters

The `CREATE TABLE` statement is the most important SQL statement you will write. Every query, join, index, and constraint depends on a correct table definition. Choosing the wrong data type (e.g., `TEXT` for a date) or omitting a constraint (e.g., `NOT NULL`) leads to data corruption, slow queries, and application bugs. Go engineers who master table design write efficient, self-documenting schemas that make application code simpler and more reliable.

## Mental model

A table is a spreadsheet with a rigid header. Each column has a fixed name and a single allowed type. Each row is one record. Unlike a spreadsheet, SQL enforces rules per column: "this column must never be empty" (`NOT NULL`), "every value here must be unique" (`UNIQUE`), "values must be between 0 and 100" (`CHECK`). The table definition is the contract between the application and the database.

## Core idea

A **table** is a collection of **rows** (records). Each row has the same set of **columns** (fields). Each column has:

- A **name** (e.g., `email`).
- A **data type** (e.g., `TEXT`).
- Optional **constraints** (e.g., `NOT NULL`, `UNIQUE`).

SQLite's five core storage classes:

| Storage class | Go analog | Example |
|---|---|---|
| `NULL` | `nil` / pointer | `NULL` |
| `INTEGER` | `int`, `int64` | `42` |
| `REAL` | `float64` | `3.14` |
| `TEXT` | `string` | `"hello"` |
| `BLOB` | `[]byte` | binary data |

SQLite is flexible: you can store any type in any column (flexible typing), but it is best practice to use **strict typing** via `CREATE TABLE ... STRICT` or column constraints.

## Under the hood

SQLite stores each table as a B-tree in a single file. The root page of the B-tree contains the table schema. Column metadata (name, type, NOT NULL, default value) is stored in `sqlite_master`, an internal table that SQLite maintains automatically.

When you create a table:

```sql
CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL);
```

SQLite inserts a row into `sqlite_master` with `type='table'`, `name='users'`, and `sql` containing the full CREATE statement. You can query it:

```sql
SELECT sql FROM sqlite_master WHERE type='table' AND name='users';
```

## How Go uses it

Go maps SQLite types to Go types in a straightforward way:

| SQL type | Go scan target | Notes |
|---|---|---|
| `INTEGER` | `int`, `int64`, `uint64` | `int` is 64-bit on modern platforms |
| `REAL` | `float64` | |
| `TEXT` | `string` | |
| `BLOB` | `[]byte` | |
| `NULL` | `*T` or `sql.Null[T]` | See lesson 14 |

When scanning a row, `database/sql` converts the column value to the target Go type. If the conversion fails (e.g., scanning a non-numeric string into an `int`), `Scan` returns an error.

## Go example

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Product struct {
	ID    int64
	Name  string
	Price float64
	Stock int
}

func createProductTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price REAL NOT NULL CHECK(price >= 0),
			stock INTEGER NOT NULL DEFAULT 0
		)
	`)
	return err
}

func insertProduct(db *sql.DB, p *Product) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO products (name, price, stock) VALUES (?, ?, ?)`,
		p.Name, p.Price, p.Stock,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func queryProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, price, stock FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createProductTable(db); err != nil {
		log.Fatal(err)
	}

	for _, p := range []Product{
		{Name: "Widget", Price: 9.99, Stock: 100},
		{Name: "Gadget", Price: 24.99, Stock: 50},
	} {
		id, err := insertProduct(db, &p)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted %s with id=%d\n", p.Name, id)
	}

	products, err := queryProducts(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Queried %d product(s):\n", len(products))
	for _, p := range products {
		fmt.Printf("  %d: %s ($%.2f, stock=%d)\n", p.ID, p.Name, p.Price, p.Stock)
	}
}
```

## Step-by-step execution

1. `createProductTable` defines the `products` table with four columns. `price REAL NOT NULL CHECK(price >= 0)` ensures no negative prices at the database level. `stock INTEGER NOT NULL DEFAULT 0` defaults missing values to zero.
2. Two rows are inserted. Each row gets an auto-incrementing primary key.
3. `queryProducts` issues `SELECT` to fetch all rows. `rows.Next()` iterates through the result set. `rows.Scan` maps each column to the corresponding Go struct field.
4. The program prints the products. Column order in `Scan` must match the SELECT list order.

## Common mistakes

- **Wrong column order in Scan**: `rows.Scan(&p.Name, &p.ID, ...)` when SELECT is `id, name, ...` causes type conversion errors or silently wrong values.
- **Forgetting `defer rows.Close()`**: Leaking `sql.Rows` holds a database connection from the pool, eventually exhausting it.
- **Using `INTEGER PRIMARY KEY` without understanding auto-increment**: In SQLite, `INTEGER PRIMARY KEY` auto-increments by default. The value is the rowid, which is unique even if rows are deleted.
- **Using REAL for money**: `REAL` is a floating-point type. `0.10 + 0.20` evaluates to `0.30000000000000004`. Store monetary values as INTEGER (cents) or NUMERIC.

## Debugging walkthrough

```go
rows, err := db.Query(`SELECT id, name, price FROM products`)
if err != nil {
	log.Fatal(err)
}
// forgot: defer rows.Close()

for rows.Next() {
	var p Product
	err := rows.Scan(&p.ID, &p.Price, &p.Name) // wrong order
	if err != nil {
		log.Fatal(err)
	}
}
```

**Symptom**: The first row scans correctly, but subsequent rows fail silently or the program panics.

**Root cause 1**: Missing `rows.Close()`. Each `rows.Next()` call advances the cursor; without closing, the connection is never released.

**Root cause 2**: Scan order mismatch. The SELECT is `id, name, price` but Scan expects `id, price, name`. If `price` is REAL and `name` is TEXT, the scan will fail with a type conversion error.

**Fix**: Always `defer rows.Close()` immediately after the error check, and match Scan order to SELECT order exactly.

## Production notes

- Prefer `TEXT` for timestamps in SQLite (ISO 8601 format). PostgreSQL has native `TIMESTAMPTZ`.
- Use `INTEGER PRIMARY KEY` for SQLite auto-increment, and `BIGSERIAL` / `UUID` for PostgreSQL.
- Document the schema with comments in migration files. SQLite stores the `sql` text in `sqlite_master`, so the original CREATE is recoverable.
- In production PostgreSQL, use `VARCHAR(n)` instead of `TEXT` for fields with a natural maximum length.

## Performance implications

- Column order in the table definition affects storage alignment. Put `INTEGER` columns before `TEXT` columns for slightly better B-tree page density.
- `NOT NULL` and `CHECK` constraints are evaluated on every INSERT/UPDATE but have negligible cost.
- Selecting only the columns you need (`SELECT name` instead of `SELECT *`) reduces I/O and memory.
- `INTEGER PRIMARY KEY` in SQLite is the rowid alias and is stored in the B-tree interior pages, making lookups by primary key extremely fast.

## Practice task

Create a table `employees` with columns: id (INTEGER PRIMARY KEY), name (TEXT NOT NULL), salary (REAL NOT NULL CHECK > 0), department (TEXT NOT NULL), hired_at (TEXT NOT NULL). Write Go code to:
1. Insert three employees.
2. Query all employees in a given department.
3. Query the total salary cost per department using `SUM` and `GROUP BY`.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/03-tables-rows-and-columns
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/03-tables-rows-and-columns
```

## Review questions

1. What is the difference between `TEXT` and `VARCHAR(n)` in SQLite?
2. What does `CHECK(price >= 0)` guarantee?
3. What happens if you SELECT `id, name` but Scan into `&name, &id`?
4. Why should you `defer rows.Close()`?
5. How does SQLite store table metadata?

## NEXT UP

Primary keys — choosing the right identifier for every row.
