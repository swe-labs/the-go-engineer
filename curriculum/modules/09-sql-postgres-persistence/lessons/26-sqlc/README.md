# sqlc

## Learning objective

Generate type-safe Go code from SQL queries using sqlc, understand the SQL-first approach, and write queries that produce idiomatic Go code with proper error handling.

## Why this matters

Writing raw `database/sql` code is repetitive and error-prone: you write SQL strings, manually `Scan` columns, and handle NULLs. ORMs hide SQL but add complexity and performance overhead. sqlc is the middle ground: you write SQL, and it generates type-safe Go functions that handle scanning, NULL types, and parameter binding automatically. The generated code is idiomatic, compilable, and fast.

## Mental model

```
SQL query → sqlc → Go code

query.sql:
  -- name: GetUser :one
  SELECT * FROM users WHERE id = ?;

users.sql.go (generated):
  func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
      row := q.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)
      // ... auto-generated Scan
  }
```

You write and maintain SQL. sqlc generates the boilerplate. The generated code is checked into your repository and reviewed like any other code.

## Core idea

sqlc reads annotated SQL queries and generates Go functions. Key annotations:

| Annotation | Meaning |
|---|---|
| `:one` | Query returns a single row (generates `QueryRowContext`) |
| `:many` | Query returns multiple rows (generates `QueryContext`) |
| `:exec` | Query returns no rows (generates `ExecContext`) |
| `:execrows` | Query returns number of affected rows |
| `:copyfrom` | Bulk insert using `COPY FROM` (PostgreSQL) |

For each query, sqlc generates:
- A function with typed parameters and return values.
- A `Scan` implementation for the result.
- Proper `sql.Null*` usage based on column nullability.

## Under the hood

sqlc generates code that uses `database/sql` internally. It produces:

1. A `Queries` struct holding a `db` interface (`DBTX`).
2. One method per query with typed parameters.
3. A `models.go` with Go structs matching your schema.
4. Scan helpers that iterate over `*sql.Rows`.

The generated code follows Go idioms: returns `(T, error)` for `:one`, `([]T, error)` for `:many`, and `(sql.Result, error)` for `:exec`.

## How Go uses it

- **SQL-first development**: Write and optimize SQL in your DB GUI, then add it to a `.sql` file and run sqlc.
- **Code generation in CI**: Run sqlc as a step to verify all queries compile.
- **Type-safe parameters**: No more `any` slices — parameters are typed structs or positional args.
- **Migration + sqlc**: Schema migrations define the tables, sqlc generates the Go types.

## Go example

Since sqlc is a code generation tool, this example shows the equivalent manually-written Go that sqlc would generate. In practice, you would run `sqlc generate` to produce this code.

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// models.go (generated from schema)
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Total     string     `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
}

// db.go (generated)
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queries struct {
	db DBTX
}

func NewQueries(db DBTX) *Queries {
	return &Queries{db: db}
}

// users.sql.go (generated from --name: CreateUser :exec)
func (q *Queries) CreateUser(ctx context.Context, name, email string) (sql.Result, error) {
	return q.db.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", name, email)
}

// users.sql.go (generated from --name: GetUser :one)
func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, "SELECT id, name, email, created_at FROM users WHERE id = ?", id)
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	return u, err
}

// users.sql.go (generated from --name: ListUsers :many)
func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT id, name, email, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, u)
	}
	return items, rows.Err()
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	q := NewQueries(db)

	q.CreateUser(context.Background(), "Alice", "alice@example.com")
	q.CreateUser(context.Background(), "Bob", "bob@example.com")

	u, _ := q.GetUser(context.Background(), 1)
	fmt.Printf("User 1: %s <%s> created at %s\n", u.Name, u.Email, u.CreatedAt.Format(time.RFC3339))

	users, _ := q.ListUsers(context.Background())
	fmt.Printf("Total users: %d\n", len(users))
	for _, u := range users {
		fmt.Printf("  %d: %s\n", u.ID, u.Name)
	}
}
```

## Step-by-step execution

For the generated `GetUser`:

1. `q.GetUser(ctx, 1)` calls `QueryRowContext` with the SQL and parameter.
2. The SQL `SELECT ... WHERE id = ?` is parameterized — `1` is bound via the driver.
3. `row.Scan` reads the result columns into the `User` struct fields.
4. If no row is found, `Scan` returns `sql.ErrNoRows`.
5. If the query fails, the error is returned.

## Common mistakes

- **Not checking `sql.ErrNoRows`**: Generated `:one` functions return `sql.ErrNoRows` when no row matches. Handle this explicitly.
- **Forgetting to run sqlc after schema changes**: The generated code is stale until re-run. Integrate sqlc into your build pipeline.
- **Using `:one` for queries that may return zero rows**: Use `:one` only when exactly one row is expected. For optional results, use a wrapper or handle `ErrNoRows`.
- **Not reviewing generated code**: sqlc generates good code, but you should review it. Check for correct NULL handling and proper types.
- **Complex queries**: sqlc handles most SQL, but very dynamic queries (computed column lists, dynamic JOINs) may need manual `database/sql`.

## Debugging walkthrough

**Scenario**: The generated `GetUser` returns `sql.ErrNoRows` for a user that exists.

```sql
-- name: GetUser :one
SELECT * FROM users WHERE id = ?;
```

**Problem**: The table has a `deleted_at` column, and the user was soft-deleted. The SQL doesn't filter `WHERE deleted_at IS NULL`.

**Fix**: Update the SQL and re-run sqlc:

```sql
-- name: GetUser :one
SELECT * FROM users WHERE id = ? AND deleted_at IS NULL;
```

## Production notes

- **sqlc configuration**: A `sqlc.yaml` file defines the schema, queries, and output directory.
- **Version control**: Check in both the SQL source files and the generated Go files. This makes code review easier and ensures builds are reproducible without running sqlc.
- **CI integration**: Run `sqlc diff` in CI to detect when generated code is out of sync with source SQL.
- **Database-specific SQL**: sqlc supports PostgreSQL, MySQL, and SQLite. Some SQL syntax differs — tag queries with the database driver.
- **NULL handling**: sqlc generates `sql.Null*` types for nullable columns. Verify your schema marks columns as `NOT NULL` when appropriate.

## Performance implications

- Generated code is as fast as hand-written `database/sql` — no ORM overhead.
- sqlc avoids the `any` reflection of some ORMs by generating concrete types.
- The `DBTX` interface allows swapping the underlying connection (useful for transactions).
- No runtime query parsing — SQL is compiled into Go at generation time.

## Practice task

Manually write the equivalent of sqlc-generated code for a `Product` table with:
1. `CreateProduct` (`:exec`)
2. `GetProduct` (`:one`)
3. `ListProducts` (`:many`)
4. `UpdateProductPrice` (`:exec`)
5. `DeleteProduct` (`:exec`)

Write a `Queries` struct and `DBTX` interface matching the sqlc pattern, then use them in `main()`.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/26-sqlc
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/26-sqlc
```

## Review questions

1. How does sqlc differ from an ORM like GORM?
2. What does the `:one` annotation tell sqlc about the expected result?
3. Why is it a best practice to check in generated sqlc code?
4. How does sqlc handle nullable columns in the generated Go code?
5. What interface does sqlc use to allow both `*sql.DB` and `*sql.Tx` as the database handle?

## NEXT UP

UPDATE operations — modifying existing rows with SET, optimistic locking, and partial updates.
