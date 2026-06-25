# SELECT

## Learning objective

Query data from a SQL table using Go's `db.Query`, iterate over result rows with `sql.Rows`, and handle query errors gracefully.

## Why this matters

SELECT is the most frequently executed SQL statement in any read-heavy application. Every API endpoint, report, dashboard, and data export is built on SELECT queries. Mastering `db.Query`, `sql.Rows` iteration, and proper resource cleanup is essential for building reliable, leak-free database code.

## Mental model

Think of `db.Query` as ordering from a menu. You tell the kitchen (database) what you want (the SELECT statement), and they prepare a tray (sql.Rows) with all the items matching your order. You take items off the tray one at a time (rows.Next), examine each (rows.Scan), and when you are done, you return the tray (rows.Close) so it can be washed and reused.

## Core idea

The SELECT workflow in Go:

1. **`db.Query(sql, args...)`**: Executes the query and returns `*sql.Rows`.
2. **`defer rows.Close()`**: Ensures the result set and connection are released.
3. **`rows.Next()`**: Advances to the next row. Returns `false` when iteration is complete.
4. **`rows.Scan(&dest...)`**: Copies the current row's columns into Go variables.
5. **`rows.Err()`**: Checks for iteration errors (connection drop, context cancellation).

Variations:

| Method | Returns | When to use |
|---|---|---|
| `db.Query` | `*sql.Rows` | Multiple rows expected |
| `db.QueryRow` | `*sql.Row` | Exactly one row expected |
| `db.QueryContext` | `*sql.Rows` | With context timeout/cancellation |

## Under the hood

When `db.Query` is called:

1. A connection is acquired from the pool.
2. The SQL is sent to the database. The database opens a **portal** (a server-side cursor) for the result set.
3. As `rows.Next()` is called, the driver fetches rows from the portal, typically in batches (default 1 row per fetch for SQLite, configurable for PostgreSQL).
4. When `rows.Close()` is called (explicitly or via `defer`), the portal is closed and the connection returns to the pool.
5. If `rows.Next()` returns `false` and all rows have been consumed, the portal is automatically closed and the connection is released — but it is safer to always call `rows.Close()`.

## How Go uses it

```go
// Multiple rows
rows, err := db.Query("SELECT id, name FROM users WHERE active = ?", true)
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil {
        return err
    }
    fmt.Println(id, name)
}
return rows.Err()

// Single row
var name string
err := db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)
if err == sql.ErrNoRows {
    fmt.Println("not found")
} else if err != nil {
    log.Fatal(err)
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

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0
	)`)

	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Widget', 9.99, 100)`)
	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Gadget', 24.99, 50)`)
	db.Exec(`INSERT INTO products (name, price, stock) VALUES ('Gizmo', 14.99, 0)`)

	fmt.Println("=== Query: all products ===")
	products, err := queryAllProducts(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range products {
		fmt.Printf("  %d: %s ($%.2f, stock=%d)\n", p.ID, p.Name, p.Price, p.Stock)
	}

	fmt.Println("\n=== Query: single product by ID ===")
	p, err := queryProductByID(db, 2)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("  Product not found")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("  %d: %s ($%.2f, stock=%d)\n", p.ID, p.Name, p.Price, p.Stock)
	}

	fmt.Println("\n=== Query: in-stock products ===")
	rows, err := db.Query(`SELECT id, name, price, stock FROM products WHERE stock > 0 ORDER BY price`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: $%.2f (%d in stock)\n", p.Name, p.Price, p.Stock)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func queryAllProducts(db *sql.DB) ([]Product, error) {
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

func queryProductByID(db *sql.DB, id int) (*Product, error) {
	var p Product
	err := db.QueryRow(`SELECT id, name, price, stock FROM products WHERE id = ?`, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
```

## Step-by-step execution

1. Three products are inserted into the `products` table.
2. `queryAllProducts` uses `db.Query` to fetch all rows, iterates with `rows.Next()`, scans each row into a `Product` struct, and returns the slice.
3. `queryProductByID` uses `db.QueryRow` to fetch a single product. If no row matches, `Scan` returns `sql.ErrNoRows`.
4. An ad-hoc query filters for in-stock products (`stock > 0`) and orders by price.
5. Each result set is properly closed with `defer rows.Close()`, and `rows.Err()` is checked after iteration.

## Common mistakes

- **Not checking `rows.Err()`**: After the `for rows.Next()` loop, there may be an iteration error (network issue, context cancellation). Always check `rows.Err()`.
- **Using `db.Query` when `db.QueryRow` suffices**: `db.Query` allocates a `*sql.Rows` that must be closed. For a single row, use `db.QueryRow`.
- **Scanning before `rows.Next()`**: Calling `rows.Scan` without first calling `rows.Next()` scans garbage data (or panics with nil pointer). Always check `rows.Next()` first, or use `QueryRow`.
- **Reusing Scan destinations**: If you scan into the same variable in a loop, all appended rows share the last value. Create new variables or dereference pointers in each iteration.

## Debugging walkthrough

```go
var products []Product
rows, err := db.Query(`SELECT id, name FROM products`)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var p Product
    err := rows.Scan(&p.ID, &p.Name, &p.Price) // 3 columns expected, only 2 in SELECT
    if err != nil {
        log.Fatal(err)
    }
    products = append(products, p)
}
```

**Symptom**: Scan error: `sql: expected 3 destination arguments in Scan, not 2`.

**Root cause**: The SELECT list has 2 columns (`id, name`) but `Scan` expects 3 (`&p.ID, &p.Name, &p.Price`). The number of Scan arguments must match the number of SELECT columns.

**Fix**: Match Scan to SELECT exactly: `rows.Scan(&p.ID, &p.Name)`, or update SELECT to include `price`.

## Production notes

- Use `SELECT` with specific column names, not `SELECT *`. This makes the contract between Go and SQL explicit and avoids breakage when columns are added.
- For paginated queries, use `LIMIT ? OFFSET ?` and pass parameters from the request.
- Use `QueryContext` with a context timeout to prevent queries from running indefinitely in production.
- Log slow queries. Set `log.Println` around queries that exceed a threshold (e.g., 100 ms).

## Performance implications

- `rows.Next()` fetches one row at a time from the database by default (for SQLite). For large result sets, this minimizes memory but adds per-row network overhead.
- For PostgreSQL, the driver can be configured with `PreferSimpleProtocol` to use the simple query protocol (text-based) or the extended protocol (binary, faster).
- `SELECT *` fetches all columns over the network, even unused ones. Be selective.
- Filter as much as possible in SQL (`WHERE`, `LIMIT`) rather than in Go. Databases are optimized for filtering.

## Practice task

Write a function `SearchProducts(db *sql.DB, minPrice, maxPrice float64) ([]Product, error)` that:
1. Queries products with price between minPrice and maxPrice (inclusive).
2. Returns the matching products ordered by price ascending.
3. Uses `db.Query` with parameters.

Then write a `main()` that: inserts 5 products and searches for products in the $10–$30 range.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/13-select
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/13-select
```

## Review questions

1. What is the difference between `db.Query` and `db.QueryRow`?
2. Why must `rows.Close()` be called (and deferred immediately)?
3. What does `rows.Err()` check for?
4. What is `sql.ErrNoRows` and when does it occur?
5. Why is `SELECT *` discouraged in production Go code?

## NEXT UP

Scanning rows — how to map SQL result columns to Go variables with `rows.Scan`.
