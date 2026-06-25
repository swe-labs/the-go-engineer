# Null handling

## Learning objective

Scan nullable database columns into Go safely using `sql.NullString`, `sql.NullInt64`, and related types, distinguish NULL from zero values, and handle JSON serialization of nullable fields.

## Why this matters

SQL NULL is not zero, empty string, or false — it means "unknown" or "missing." If you scan a NULL column into a plain Go `string`, you get `""`. If you scan it into an `int`, you get `0`. Both silently corrupt your data. Production systems crash, send wrong invoices, or skip records because NULL was not handled explicitly. The `database/sql` package provides nullable types that force you to decide: is this value valid, or is it NULL?

## Mental model

Think of a nullable Go value as a box with an "is-valid" light. When you scan a column that is NULL, the light stays off. When you scan a non-NULL value, the light turns on and the value is stored inside.

```
┌─────────────┐      ┌─────────────┐
│ NULL column │ ──→  │ sql.NullString   │
│             │      │ Valid: false │
│             │      │ String: ""   │
└─────────────┘      └─────────────┘

┌─────────────┐      ┌─────────────┐
│ "hello"     │ ──→  │ sql.NullString   │
│             │      │ Valid: true  │
│             │      │ String: "hello"  │
└─────────────┘      └─────────────┘
```

Always check `Valid` before reading the value field. If you skip the check, you are treating missing data as real data.

## Core idea

Go's `database/sql` defines four nullable types in the standard library:

| Go type | SQL NULL equivalent | Zero value when Valid=false |
|---|---|---|
| `sql.NullString` | `TEXT NULL`, `VARCHAR NULL` | `String: ""` |
| `sql.NullInt64` | `INTEGER NULL`, `BIGINT NULL` | `Int64: 0` |
| `sql.NullFloat64` | `REAL NULL`, `DOUBLE NULL` | `Float64: 0.0` |
| `sql.NullBool` | `BOOLEAN NULL`, `INT NULL` | `Bool: false` |
| `sql.NullTime` | `TIMESTAMP NULL`, `DATETIME NULL` | `Time: time.Time{}` |
| `sql.Null[V]` (Go 1.22+) | Any nullable column | `V: zero value` |

Each type has a `Valid bool` field set to `true` when the column was non-NULL and `false` otherwise.

## Under the hood

When `db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&s)` executes:
1. The SQLite driver returns a byte slice or nil for the column value.
2. `Scan` calls the `sql.Scanner` interface on the target.
3. `sql.NullString.Scan` checks: if the raw value is nil, set `Valid = false`; otherwise convert bytes to string, set `String = result`, `Valid = true`.
4. For JSON, `sql.NullString.MarshalJSON` returns `null` when `Valid` is false, and the quoted string otherwise.

Internally, each Null type is a struct with an exported value field and a `Valid` field — no magic, no allocation on scan.

## How Go uses it

- **Table scans**: any column with a `NULL` constraint must be scanned into a nullable type.
- **JSON APIs**: `json.Marshal(sql.NullString{String: "hello", Valid: true})` → `"hello"`; `json.Marshal(sql.NullString{Valid: false})` → `null`.
- **Optional filter parameters**: nullable fields in WHERE clauses.
- **Partial updates**: only update fields where `Valid` is true.

## Go example

```go
package main

import (
	"database/sql"
	"encoding/json"
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

	_, err = db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		stock INTEGER
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO products VALUES
		(1, 'Widget', 'A useful widget', 9.99, 100),
		(2, 'Gadget', NULL, 24.99, 0),
		(3, 'Doohickey', 'Premium doohickey', 49.99, NULL)`)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, name, description, price, stock FROM products ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type Product struct {
		ID          int
		Name        string
		Description sql.NullString
		Price       float64
		Stock       sql.NullInt64
	}

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock)
		if err != nil {
			log.Fatal(err)
		}
		jsonBytes, _ := json.MarshalIndent(p, "", "  ")
		fmt.Printf("Product %d: %s\n", p.ID, jsonBytes)
	}
}
```

Output:

```
Product 1: {
  "ID": 1,
  "Name": "Widget",
  "Description": "A useful widget",
  "Price": 9.99,
  "Stock": 100
}
Product 2: {
  "ID": 2,
  "Name": "Gadget",
  "Description": null,
  "Price": 24.99,
  "Stock": 0
}
Product 3: {
  "ID": 3,
  "Name": "Doohickey",
  "Description": "Premium doohickey",
  "Price": 49.99,
  "Stock": null
}
```

Notice: Product 2 has `Description: null` (NULL in DB) and `Stock: 0` (real zero, not NULL). Product 3 has `Stock: null` (NULL in DB). JSON output distinguishes them correctly because we used nullable types.

## Step-by-step execution

For the scan `rows.Scan(&p.Description)` on the Gadget row:

1. SQLite returns a NULL for the `description` column.
2. The `database/sql` driver converts NULL to a nil `any`.
3. `sql.NullString.Scan` receives `nil`.
4. Since the value is nil, `Valid` is set to `false`.
5. The `String` field retains its zero value (`""`).
6. `json.Marshal` checks `Valid` → `false` → outputs `null`.

For the Widget row:

1. SQLite returns `"A useful widget"`.
2. The driver passes it as `[]byte("A useful widget")`.
3. `Scan` converts bytes to `string`, sets `String = "A useful widget"`, `Valid = true`.
4. `json.Marshal` checks `Valid` → `true` → outputs `"A useful widget"`.

## Common mistakes

- Scanning NULL into `string` or `int`: Go silently produces `""` or `0`, indistinguishable from real data. Always use `sql.Null*` for nullable columns.
- Forgetting to check `Valid`: reading `NullString.String` when `Valid` is false gives `""` — which might be a valid empty string. Check `Valid` first.
- Using `sql.NullString` for `NOT NULL` columns: unnecessary overhead. Use plain `string` when the column is guaranteed non-NULL.
- Confusing NULL with zero in business logic: `stock = 0` means "out of stock." `stock IS NULL` means "stock unknown." They are not the same.

## Debugging walkthrough

**Scenario**: A pricing report shows `$0.00` for some products. Expected: those products should be excluded from the report.

```go
type Product struct {
	ID    int
	Name  string
	Price float64 // BUG: should be sql.NullFloat64
}

for rows.Next() {
	var p Product
	rows.Scan(&p.ID, &p.Name, &p.Price)
	fmt.Printf("%s: $%.2f\n", p.Name, p.Price)
}
```

**Root cause**: The `price` column is `REAL NOT NULL` in schema, so `float64` works — but the `stock` column is `INTEGER NULL`. If a product has `price NOT NULL` but the report query joined a table with NULL prices, or if a migration removed the NOT NULL constraint, NULL creeps in.

**Fix**: Use `sql.NullFloat64` for any column that might ever be NULL. Check `Valid`. Log or skip NULL prices.

```go
type Product struct {
	ID    int
	Name  string
	Price sql.NullFloat64
}
for rows.Next() {
	var p Product
	rows.Scan(&p.ID, &p.Name, &p.Price)
	if !p.Price.Valid {
		fmt.Printf("%s: price unknown, skipping\n", p.Name)
		continue
	}
	fmt.Printf("%s: $%.2f\n", p.Name, p.Price.Float64)
}
```

## Production notes

- **API contracts**: Document which fields can be null in your API schema (OpenAPI). Use `omitempty` with care — it omits zero values, not null.
- **JSON encoding**: `sql.NullString` marshals to `null` when `Valid` is false. Custom JSON encoding can strip null fields, include them as null, or skip them — decide per endpoint.
- **CSV exports**: NULL cells should be empty or the string `NULL`. Decide and document.
- **gRPC / protobuf**: Use `google.protobuf.StringValue` (wrapper types) for nullable fields; they map to pointers or `sql.Null*` on the Go side.

## Performance implications

- `sql.Null*` types add a single `bool` field — negligible memory overhead.
- JSON marshaling of nullable types checks `Valid` in a branch — sub-nanosecond cost.
- Scanning into nullable types is identical to scanning into plain types plus a nil check.
- Generics `sql.Null[V]` (Go 1.22+) may reduce boilerplate but currently has limited driver support. Stick with the named types for broad compatibility.

## Practice task

Write a function `GetProducts` that:
1. Opens an in-memory SQLite DB.
2. Creates a `products` table with nullable `description` and `stock`.
3. Inserts 3-5 rows, mixing NULL and non-NULL values.
4. Returns `[]Product` where `Product` uses `sql.NullString`, `sql.NullInt64`.
5. Prints each product as JSON, showing `null` for NULL fields.

Then write a second function `NonNullProducts` that filters out any product where `description` is NULL.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/15-null-handling
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/15-null-handling
```

## Review questions

1. What is the difference between `""` (empty string) and SQL NULL when scanned into `sql.NullString`?
2. Why does scanning a NULL column into a plain `int` produce 0, and why is that dangerous?
3. What does `json.Marshal(sql.NullString{Valid: false})` return?
4. When would you use `sql.NullTime` instead of `time.Time` for a `TIMESTAMP` column?
5. How does `Scan` decide whether to set `Valid` to true or false?

## NEXT UP

Prepared statements — parameterized queries that prevent SQL injection and improve performance.
