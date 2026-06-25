# SQL injection prevention

## Learning objective

Write and defend Go database queries against SQL injection by always using parameterized queries, recognizing unsafe string concatenation patterns, and understanding why input escaping alone is insufficient.

## Why this matters

SQL injection is the most critical web application vulnerability. It has caused some of the largest data breaches in history -- including the 2017 Equifax breach that exposed 147 million records. In Go, the standard `database/sql` package makes parameterized queries easy and effective, but developers can still concatenate user input into SQL strings. A single `fmt.Sprintf` on a query string can sink a company. Professional Go engineers must never write an unsafe query and must catch them in code review.

## Mental model

Imagine you are a bank teller. A customer hands you a withdrawal slip that says "Pay $500 to bearer." You process it. Now imagine the slip says "Pay $500 to bearer. Also transfer $1,000,000 from the vault to my car." If you treat the entire slip as instructions, you are compromised. Parameterized queries are like having a pre-printed form with a blank for the amount: the customer fills in only the blank, and the instructions are fixed.

SQL injection works because string concatenation mixes code (SQL instructions) with data (user input). Parameterized queries send the code template and the data values on separate channels, so the database never confuses the two.

## Core idea

SQL injection occurs when untrusted user input is placed directly into a SQL query string, allowing an attacker to modify the query structure. The core defense is to use parameterized queries (also called prepared statements) where the SQL structure is fixed and user values are passed as separate parameters.

In Go, parameterized queries look like:

```go
db.Query("SELECT * FROM users WHERE id = ?", userID)
```

The `?` is a parameter marker. The database driver safely substitutes `userID` as a data value, never as executable SQL.

| Approach | Security | Example |
|---|---|---|
| String concatenation | Unsafe | `fmt.Sprintf("WHERE name = '%s'", input)` |
| Parameterized query | Safe | `db.Query("WHERE name = ?", input)` |
| ORM method | Safe | `db.Where("name = ?", input).Find(&user)` |
| Stored procedure | Safe | `CALL get_user(?)` |

## Under the hood

When a database receives a parameterized query, it first parses and compiles the SQL template into an execution plan. Only then does it bind the parameter values. Because the SQL structure is already compiled before the parameters arrive, user input cannot alter the query logic.

In `database/sql`, the `Prepare` method sends the SQL template to the database server, which returns a statement handle. Subsequent `Exec` or `Query` calls with different parameter values reuse the same compiled plan. This also provides a performance benefit for repeated queries.

```text
Unsafe:  "SELECT * FROM users WHERE name = '" + input + "'"
         -> Server parses: SELECT * FROM users WHERE name = 'alice' OR '1'='1'
         -> The OR '1'='1' becomes part of the SQL structure

Safe:    "SELECT * FROM users WHERE name = ?"  +  param("alice' OR '1'='1")
         -> Server compiles: SELECT * FROM users WHERE name = $1
         -> Server binds: $1 = "alice' OR '1'='1"
         -> The input is treated as a literal string value, not SQL
```

## How Go uses it

The `database/sql` package supports parameterized queries with the `?` marker (for most drivers) or driver-specific markers like `$1` (PostgreSQL) and `@p1` (SQL Server).

Key functions that accept parameters:

- `db.Query(query, args...)` -- returns rows
- `db.QueryRow(query, args...)` -- returns a single row
- `db.Exec(query, args...)` -- executes without returning rows
- `db.Prepare(query)` -- returns a `*Stmt` that can be reused with `stmt.Exec(args...)`

Most Go ORMs (GORM, sqlx, Ent) also use parameterized queries internally. However, they also offer raw query methods that bypass protection if used with string concatenation.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

func unsafeQuery(template, input string) string {
	return strings.ReplaceAll(template, "'{input}'", "'"+input+"'")
}

func escapeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func safeQuery(template, input string) string {
	escaped := escapeSQLString(input)
	return strings.ReplaceAll(template, "?", "'"+escaped+"'")
}

func isInjection(query string) bool {
	dangerous := []string{" OR ", " UNION ", " DROP ", " -- ", "/*", ";", "1'='1"}
	upper := strings.ToUpper(query)
	for _, d := range dangerous {
		if strings.Contains(upper, d) {
			return true
		}
	}
	return false
}

func simulateDBLookup(query string) string {
	if isInjection(query) {
		return "BREACHED: All users returned!"
	}
	if strings.Contains(query, "alice") {
		return "alice@example.com"
	}
	return "no results"
}

func main() {
	base := "SELECT email FROM users WHERE username = '{input}'"

	fmt.Println("=== Unsafe query ===")
	q1 := unsafeQuery(base, "alice' OR '1'='1")
	fmt.Println("Query:", q1)
	fmt.Println("Result:", simulateDBLookup(q1))

	fmt.Println("\n=== Safe query (parameterized) ===")
	q2 := safeQuery("SELECT email FROM users WHERE username = ?", "alice' OR '1'='1")
	fmt.Println("Query:", q2)
	fmt.Println("Result:", simulateDBLookup(q2))
}
```

## Step-by-step execution

For the unsafe query `SELECT email FROM users WHERE username = 'alice' OR '1'='1'`:

1. User submits username `alice' OR '1'='1`.
2. Code formats the query: `fmt.Sprintf("SELECT email FROM users WHERE username = '%s'", input)`.
3. Resulting string: `SELECT email FROM users WHERE username = 'alice' OR '1'='1'`.
4. Database parses the entire string as SQL. The `WHERE` clause becomes `username = 'alice' OR '1'='1'`.
5. Since `'1'='1'` is always true, the query returns every row in the users table.
6. The `QueryRow` returns the first row (alice), leaking data the user should not see.

For the safe parameterized query `SELECT email FROM users WHERE username = ?`:

1. Go sends the SQL template to the database: `SELECT email FROM users WHERE username = ?`.
2. Database compiles this into an execution plan with `?` as a parameter slot.
3. Go sends the parameter value: `alice' OR '1'='1`.
4. Database binds the value as a string literal in the parameter slot.
5. The database searches for a username whose literal value is `alice' OR '1'='1`.
6. No such user exists, so `QueryRow` returns `sql.ErrNoRows`.

## Common mistakes

- Mistake: Escaping quotes instead of using parameterized queries.
  - Why it happens: Developers think escaping user input is sufficient.
  - Fix: Never write custom escaping functions. Use parameterized queries. Escaping fails against second-order injection, alternate encodings, and edge cases in the escaping logic.

- Mistake: Using `fmt.Sprintf` to build queries with "trusted" parts.
  - Why it happens: The developer believes some inputs are safe (e.g., column names or sort directions).
  - Fix: Column names and SQL keywords cannot be parameterized. Use an allowlist: `if dir != "ASC" && dir != "DESC" { return error }`.

- Mistake: Using an ORM but calling raw query methods with string concatenation.
  - Why it happens: ORMs like GORM provide `.Exec()` and `.Raw()` that accept raw SQL strings.
  - Fix: Always use the ORM's query builder methods. When using `.Raw()` or `.Exec()`, pass values as arguments, never interpolated into the string.

- Mistake: Assuming stored procedures prevent injection.
  - Why it happens: Stored procedures contain SQL inside them, and if they concatenate inputs internally, they are vulnerable.
  - Fix: Write stored procedures with parameterized statements internally too.

- Mistake: Logging or returning full SQL queries in error messages.
  - Why it happens: Developers log `err.Error()` which may contain the full query with user input.
  - Fix: Log only the error message without the query, or use a structured logger that redacts sensitive fields.

## Debugging walkthrough

Consider this broken code:

```go
func searchProducts(db *sql.DB, name string, category string) string {
	query := fmt.Sprintf("SELECT * FROM products WHERE name LIKE '%%%s%%' AND category = '%s'", name, category)
	rows, err := db.Query(query)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		count++
	}
	return fmt.Sprintf("Found %d products", count)
}
```

Symptom: The function works normally until a user searches for `'; DROP TABLE products; --`.

Investigation: Add a print statement to see the generated query:

```go
fmt.Println("Query:", query)
```

Output: `SELECT * FROM products WHERE name LIKE '%'; DROP TABLE products; --%' AND category = 'shoes'`

Root cause: The input `'; DROP TABLE products; --` breaks out of the single-quoted string, terminates the current statement, and executes a DROP command. The `--` comments out the remainder of the query. This is a classic SQL injection attack.

Fix: Use parameterized queries for both `name` and `category`:

```go
func searchProducts(db *sql.DB, name string, category string) (string, error) {
	query := "SELECT * FROM products WHERE name LIKE ? AND category = ?"
	pattern := "%" + name + "%"
	rows, err := db.Query(query, pattern, category)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		count++
	}
	return fmt.Sprintf("Found %d products", count), nil
}
```

With parameterized queries, the `name` value `'; DROP TABLE products; --` is treated as a literal search string. The database searches for products whose name contains that exact text. No DROP command executes.

## Production notes

- Enforce parameterized queries in code review. A linter like `sqlvet` or a custom `go vet` analysis can catch `fmt.Sprintf` patterns used with `db.Query`.
- Use a query builder library (such as `sq` or `squirrel`) for dynamic queries to avoid string concatenation entirely.
- For column names or SQL keywords that cannot be parameterized, use a strict allowlist: map user-facing identifiers to known-safe column names.
- Rate-limit login endpoints regardless of query safety. SQL injection is not the only attack; credential stuffing and brute force are also threats.
- Database users in production should have the minimum necessary privileges. The application's database user should not have `DROP` or `CREATE` permissions.

## Performance implications

- Parameterized queries incur a one-time cost for preparing the statement. Reusing a prepared statement across multiple calls avoids re-parsing and re-optimizing the SQL, improving throughput.
- String concatenation queries cannot benefit from cached execution plans because each query string is different, forcing the database to re-parse every time.
- For one-off queries, `db.Query` with parameters internally prepares and executes in a single round trip (most drivers). There is no meaningful performance penalty compared to concatenation.
- The security benefit of parameterized queries far outweighs any marginal performance difference.

## Practice task

Write a function `queryUser(db *sql.DB, id int, fields []string) (map[string]interface{}, error)` that safely fetches a user by ID but allows the caller to specify which columns to return. Use an allowlist for the column names and a parameterized query for the ID.

Column allowlist: `id`, `username`, `email`, `created_at`. Reject any field not in this list with an error.

Then write a `main()` that:
1. Opens an in-memory SQLite database.
2. Creates a `users` table with columns: `id INTEGER, username TEXT, email TEXT, created_at TEXT`.
3. Inserts at least two test rows.
4. Calls `queryUser` with valid and invalid column lists.
5. Prints results or errors.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/14-sql-injection-prevention
go test ./curriculum/modules/10-auth-security/lessons/14-sql-injection-prevention
```

The existing tests verify that unsafe queries are vulnerable to injection and safe queries properly block injection attempts.

## Review questions

1. What is the fundamental difference between a parameterized query and a string-concatenated query in terms of how the database processes them?
2. Can input validation alone prevent SQL injection? Why or why not?
3. In Go's `database/sql`, what does the `?` parameter marker represent and how does the driver handle it?
4. Why can't column names or SQL keywords be parameterized? What is the correct alternative?
5. How does a prepared statement improve performance for repeated queries beyond just security?

## NEXT UP

XSS -- Cross-Site Scripting attacks and how Go's `html/template` package provides automatic context-sensitive escaping to prevent script injection.
