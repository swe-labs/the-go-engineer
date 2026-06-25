# Joins

## Learning objective

Combine data from multiple tables using INNER JOIN, LEFT JOIN, and other join types, scan joined result sets in Go, and recognize the N+1 query problem.

## Why this matters

Relational databases are useful precisely because they let you split data across normalized tables and recombine it with joins. A single JOIN query replaces N+1 separate queries, reduces network round-trips, and lets the database optimizer choose the fastest access path. Every production application — from e-commerce orders to social graphs to billing systems — relies on joins.

## Mental model

A join is a virtual table created by pairing rows from two tables according to a condition.

```
Table A (users)          Table B (orders)
┌─────┬───────┐          ┌─────┬──────┬────────┐
│ id  │ name  │          │ id  │ uid  │ total  │
├─────┼───────┤          ├─────┼──────┼────────┤
│ 1   │ Alice │          │ 101 │ 1    │ 50.00  │
│ 2   │ Bob   │          │ 102 │ 1    │ 30.00  │
│ 3   │ Carol │          │ 103 │ 3    │ 20.00  │
└─────┴───────┘          └─────┴──────┴────────┘

INNER JOIN on users.id = orders.uid:
┌───────┬────────┐
│ name  │ total  │
├───────┼────────┤
│ Alice │ 50.00  │  ← Alice has 2 orders
│ Alice │ 30.00  │
│ Carol │ 20.00  │  ← Carol has 1 order
└───────┴────────┘
Bob has no orders → excluded from INNER JOIN, included with LEFT JOIN.
```

## Core idea

The four standard join types:

| Join type | Returns |
|---|---|
| `INNER JOIN` | Only rows where the condition matches in both tables |
| `LEFT JOIN` | All rows from the left table, plus matching right rows (NULL if no match) |
| `RIGHT JOIN` | All rows from the right table, plus matching left rows (SQLite does not support natively — swap tables and use LEFT JOIN) |
| `FULL JOIN` | All rows from both tables (not supported in SQLite — emulate with UNION) |

In Go, you scan joined results like any other query. Column order in `Scan` must match the SELECT list order.

## Under the hood

The database's query planner chooses the join strategy:
- **Nested loop join**: For each row in table A, scan table B for matches. O(n*m). Used when one table is small.
- **Hash join**: Build a hash table of one table, probe with the other. O(n+m). Used for large unsorted tables.
- **Merge join**: Sort both tables by the join key, then merge. O(n log n + m log m). Used when both tables are sorted.

SQLite uses nested loop joins primarily, with automatic index usage for the inner table.

## How Go uses it

- **GET /users/:id/orders**: Join users and orders to return user info with order history.
- **Reporting**: Join sales, products, and categories for aggregate reports.
- **Dashboard queries**: Multi-table joins for metrics (e.g., "orders by region with customer tier").
- **N+1 prevention**: One JOIN query replaces N queries in a loop.

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

	_, err = db.Exec(`
		CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE orders (id INTEGER PRIMARY KEY, user_id INTEGER NOT NULL, total REAL NOT NULL, FOREIGN KEY (user_id) REFERENCES users(id));
		INSERT INTO users VALUES (1, 'Alice'), (2, 'Bob'), (3, 'Carol');
		INSERT INTO orders VALUES (101, 1, 50.00), (102, 1, 30.00), (103, 3, 20.00);
	`)
	if err != nil {
		log.Fatal(err)
	}

	type OrderWithUser struct {
		OrderID int
		UserID  int
		Total   float64
		// LEFT JOIN produces NULL name for users without orders
		UserName sql.NullString
	}

	// INNER JOIN — only users with orders
	fmt.Println("=== INNER JOIN ===")
	rows, _ := db.Query(`SELECT orders.id, orders.user_id, orders.total, users.name
		FROM orders INNER JOIN users ON orders.user_id = users.id
		ORDER BY orders.id`)
	scanAndPrint(rows)

	// LEFT JOIN — all orders, even if user is missing (unlikely here)
	fmt.Println("\n=== LEFT JOIN (orders -> users) ===")
	rows, _ = db.Query(`SELECT orders.id, orders.user_id, orders.total, users.name
		FROM orders LEFT JOIN users ON orders.user_id = users.id
		ORDER BY orders.id`)
	scanAndPrint(rows)

	// LEFT JOIN — all users, with their orders (NULL if no order)
	fmt.Println("\n=== LEFT JOIN (users -> orders) ===")
	type UserWithOrder struct {
		UserID   int
		UserName string
		OrderID  sql.NullInt64
		Total    sql.NullFloat64
	}
	userRows, _ := db.Query(`SELECT users.id, users.name, orders.id, orders.total
		FROM users LEFT JOIN orders ON users.id = orders.user_id
		ORDER BY users.id, orders.id`)
	defer userRows.Close()
	for userRows.Next() {
		var uwo UserWithOrder
		userRows.Scan(&uwo.UserID, &uwo.UserName, &uwo.OrderID, &uwo.Total)
		if uwo.OrderID.Valid {
			fmt.Printf("User %s (id=%d): Order %d = $%.2f\n", uwo.UserName, uwo.UserID, uwo.OrderID.Int64, uwo.Total.Float64)
		} else {
			fmt.Printf("User %s (id=%d): No orders\n", uwo.UserName, uwo.UserID)
		}
	}
}

func scanAndPrint(rows *sql.Rows) {
	defer rows.Close()
	for rows.Next() {
		var o OrderWithUser
		rows.Scan(&o.OrderID, &o.UserID, &o.Total, &o.UserName)
		fmt.Printf("Order %d: user %d, $%.2f, name=%s\n", o.OrderID, o.UserID, o.Total, o.UserName.String)
	}
}
```

## Step-by-step execution

For `INNER JOIN orders ON users.id = orders.user_id`:

1. The query planner reads `orders` (driving table) and `users`.
2. For each order row, it looks up `users.id = order.user_id`.
3. If a match is found, the combined row is emitted.
4. `Bob` (id=2) has no orders — no combined row for Bob.
5. `Alice` appears twice because she has two orders.

For `LEFT JOIN` with `users LEFT JOIN orders`:
1. All users are in the result.
2. `Bob` appears once, with NULL for order fields.
3. Use `sql.NullInt64` and `sql.NullFloat64` to scan the potentially-null order columns.

## Common mistakes

- **Forgetting to handle NULLs from LEFT JOIN**: Right-side columns in a LEFT JOIN can be NULL. Always use nullable Go types for those columns.
- **N+1 queries**: Fetching orders for each user in a loop instead of one JOIN.
- **Ambiguous column names**: `SELECT id` fails when both tables have `id`. Use aliases: `SELECT users.id AS user_id, orders.id AS order_id`.
- **Forgetting JOIN condition**: `FROM users, orders` produces a Cartesian product. Always use `ON` or `WHERE` to specify the join condition.
- **RIGHT JOIN not supported in SQLite**: SQLite only supports LEFT JOIN. Use LEFT JOIN with swapped table order.

## Debugging walkthrough

**Scenario**: A dashboard query is slow and returns unexpected results.

```go
rows, _ := db.Query(`SELECT u.name, o.total
	FROM users u, orders o
	WHERE u.id = o.user_id`)
```

**Problem**: This is an implicit INNER JOIN. It works, but for complex queries explicit `JOIN ... ON` is clearer. Worse, if the WHERE clause is accidentally omitted (during refactoring), it becomes a Cartesian product returning millions of rows.

**Fix**: Use explicit JOIN syntax:

```go
rows, _ := db.Query(`SELECT u.name, o.total
	FROM users u
	INNER JOIN orders o ON u.id = o.user_id`)
```

## Production notes

- **Index join columns**: `orders.user_id` should be indexed. Without an index, the database scans the entire `orders` table for each user row.
- **Join order**: The query planner usually picks the optimal order, but with complex joins (5+ tables) you may need to guide it with `INNER JOIN` vs `LEFT JOIN` order hints.
- **Limit joined results**: When joining one-to-many, a single user can explode into many rows. Use `DISTINCT` or aggregate functions if you need one row per entity.
- **ORMs vs raw joins**: ORMs hide join complexity but often generate inefficient N+1 queries. For critical paths, write raw JOIN queries.

## Performance implications

- An index on the join column of the inner table reduces JOIN from O(n*m) to approximately O(n log m).
- LEFT JOIN is slightly more expensive than INNER JOIN because the database must also emit non-matching rows.
- Joining 3-4 tables is typical. Beyond 5-6, query plans become complex and may need optimization.
- Avoid joining on `LIKE` or function-wrapped columns — they defeat index usage.

## Practice task

Create three tables: `customers`, `orders`, `order_items`. Write queries using:
1. INNER JOIN to list orders with customer names.
2. LEFT JOIN to list all customers and their orders (including those with none).
3. A three-table JOIN to list order items with product names and customer names.

Scan all results in Go using appropriate nullable types.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/20-joins
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/20-joins
```

## Review questions

1. What is the difference between INNER JOIN and LEFT JOIN?
2. When scanning a LEFT JOIN result, why must you use `sql.NullString` for right-table columns?
3. What is the N+1 query problem and how do joins solve it?
4. Why doesn't SQLite support RIGHT JOIN, and what is the workaround?
5. How does an index on the join column affect query performance?

## NEXT UP

Indexes — speeding up queries with B-tree indexes, composite indexes, and understanding tradeoffs.
