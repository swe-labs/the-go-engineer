# Relational modeling

## Learning objective

Design a relational schema by identifying entities, attributes, and relationships, apply normalization rules up to third normal form, and map the result to Go structs that mirror the relational model.

## Why this matters

A database schema that mirrors the problem domain is easier to query, maintain, and extend. Relational modeling gives you a disciplined process for translating real-world concepts (users, orders, products) into tables and columns before writing a single `CREATE TABLE`. Go engineers who skip this step end up with denormalized tables, duplicated data, and JOIN queries that are slow and confusing. Learning relational modeling is the difference between a database that is a liability and one that is an asset.

## Mental model

Imagine you are designing a filing cabinet for a company. Each drawer holds one kind of thing (employees, customers, invoices). The drawer's label is the entity name. Within the drawer, each folder is a row, and each tab on the folder is an attribute (name, date, amount). A relationship between drawers is a cross-reference: a folder in the invoices drawer might have a tab that says "customer id: 42", pointing to a folder in the customers drawer.

Normalization is the process of ensuring that every piece of information lives in exactly one drawer. If you find yourself writing the customer name on every invoice folder, you have denormalized data — fix it by keeping only the customer ID on the invoice and looking up the name in the customer drawer.

## Core idea

**Relational modeling** is a top-down design process:

1. **Entities**: identify the nouns in your problem domain (User, Product, Order).
2. **Attributes**: for each entity, list the properties that describe it (User: name, email, created_at).
3. **Relationships**: determine how entities connect (User has many Orders; Order belongs to User).
4. **Keys**: choose a primary key for each entity and foreign keys to express relationships.
5. **Normalization**: apply rules to eliminate redundancy and update anomalies.

The result is a set of **relations** (tables) that can be queried with SQL.

Normal forms (simplified):

| Normal form | Rule |
|---|---|
| 1NF | Every column holds atomic values, not lists or sets. |
| 2NF | Every non-key column depends on the whole primary key, not just part of it. |
| 3NF | Every non-key column depends only on the primary key, not on another non-key column. |

## Under the hood

Relational model theory, developed by E. F. Codd in 1970, is grounded in set theory and first-order predicate logic. A table is a mathematical relation: a set of tuples (rows) over a set of attributes (columns). SQL is the concrete language that implements relational algebra operations: SELECT (σ, projection), WHERE (σ, selection), JOIN (⋈), GROUP BY (aggregation).

When a Go program issues `SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id`, the database engine:

1. Scans the users table for matching rows.
2. For each matching user row, probes the orders index for matching order rows.
3. Produces a result set that is the relational join of the two relations.
4. Returns rows to the Go driver, which populates `sql.Rows`.

The Go struct is the application-level analog of a relation. A well-designed Go struct mirrors its corresponding table column by column.

## How Go uses it

Go structs map directly to relational tables in a pattern called **Object-Relational Mapping (ORM)** — though Go's philosophy prefers explicit mapping over magic. A typical mapping:

```go
type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

type Order struct {
	ID       int64   `db:"id"`
	UserID   int64   `db:"user_id"`
	Total    float64 `db:"total"`
	Status   string  `db:"status"`
}
```

The struct tag (e.g., `db:"name"`) documents which column maps to which field. Code generators like `sqlc` automate this mapping, but the conceptual bridge between relational model and Go struct is the same.

## Go example

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

type Order struct {
	ID     int64
	UserID int64
	Total  float64
	Status string
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			total REAL NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
	`)
	return err
}

func insertUser(db *sql.DB, u *User) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)`,
		u.Name, u.Email, u.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertOrder(db *sql.DB, o *Order) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO orders (user_id, total, status) VALUES (?, ?, ?)`,
		o.UserID, o.Total, o.Status,
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

	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	userID, err := insertUser(db, &User{Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted user with id=%d\n", userID)

	orderID, err := insertOrder(db, &Order{UserID: userID, Total: 49.99, Status: "completed"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted order with id=%d\n", orderID)

	var userName string
	var orderTotal float64
	err = db.QueryRow(`
		SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id WHERE o.id = ?
	`, orderID).Scan(&userName, &orderTotal)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User %s placed order for $%.2f\n", userName, orderTotal)
}
```

## Step-by-step execution

1. `createTables` executes two `CREATE TABLE` statements that define the relational model: `users` and `orders`. The `orders.user_id` column references `users.id` — this is the foreign key relationship.
2. An `INSERT` into `users` stores a user row. The `created_at` column stores an ISO 8601 timestamp string (SQLite has no native timestamp type).
3. An `INSERT` into `orders` references the user's primary key (`user_id`). This expresses the relationship: "Alice has an order for $49.99."
4. The `SELECT ... JOIN` query reconstructs the relationship by combining data from both tables. The relational model ensures that Alice's name is stored once in `users` and referenced many times in `orders` — no duplication.

## Common mistakes

- **Storing lists in a single column** (violating 1NF). Example: `tags TEXT` containing `"go,sql,api"`. Fix: create a separate `tags` table with one row per tag, linked by foreign key.
- **Copying the same data across rows** (violating 3NF). Example: storing `user_name` and `user_email` in the `orders` table. Fix: store only `user_id` and JOIN to get the name.
- **Modeling many-to-many without a junction table**. Example: students and courses need a `enrollments` table with `student_id` and `course_id` as a composite primary key.
- **Using float for monetary values**. Example: `total REAL`. Floating-point rounding causes accounting errors. Use `INTEGER` (cents) or a `NUMERIC` type.

## Debugging walkthrough

Consider this schema design:

```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    customer_name TEXT,
    customer_email TEXT,
    product_name TEXT,
    product_price REAL,
    quantity INTEGER
);
```

**Symptom**: Updating a customer's email requires updating every row in `orders` that references that customer. A bug leaves some rows with the old email.

**Root cause**: Denormalization. `customer_name` and `customer_email` belong in a `customers` table. `product_name` and `product_price` belong in a `products` table.

**Fix**: Split into normalized tables:

```sql
CREATE TABLE customers (id INTEGER PRIMARY KEY, name TEXT, email TEXT UNIQUE);
CREATE TABLE products (id INTEGER PRIMARY KEY, name TEXT, price REAL);
CREATE TABLE orders (id INTEGER PRIMARY KEY, customer_id INTEGER REFERENCES customers(id));
CREATE TABLE order_items (id INTEGER PRIMARY KEY, order_id INTEGER REFERENCES orders(id), product_id INTEGER REFERENCES products(id), quantity INTEGER);
```

Now an email update is a single UPDATE on `customers`.

## Production notes

- Start with a pen-and-paper or whiteboard ER diagram before writing any code. It is cheap to fix design errors at this stage.
- Use a dedicated migration tool (golang-migrate, goose) to apply schema changes. Never create tables manually in production.
- Enforce foreign keys in PostgreSQL; in SQLite they must be enabled at runtime with `PRAGMA foreign_keys = ON`.
- Go struct tags should match the column naming convention (snake_case in SQL, camelCase in Go is conventional).

## Performance implications

- Proper normalization reduces data duplication, which reduces storage and improves cache efficiency.
- JOINs have a cost: every JOIN adds a table scan or index lookup. Denormalization (intentional redundancy for read performance) is sometimes necessary at scale, but it should be a deliberate trade-off, not the default.
- Indexes on foreign key columns are critical for JOIN performance. Without them, every JOIN becomes a full table scan.

## Practice task

Design a relational schema for a blog. Entities: Author (name, email), Post (title, body, published_at), Tag (name). An author has many posts. A post has many tags (many-to-many). Write Go structs for each entity, then write functions to:
1. Create all tables.
2. Insert an author, a post, and tag the post.
3. Query all posts by an author with their tags.

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/02-relational-modeling
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/02-relational-modeling
```

The tests verify the relational integrity of INSERT, JOIN queries, and the Go struct mapping.

## Review questions

1. What is the difference between 1NF, 2NF, and 3NF?
2. Why is storing `customer_name` in the `orders` table a bad idea?
3. How would you model a many-to-many relationship (e.g., students and courses) in SQL?
4. What is the relationship between a Go struct and a relational table?
5. What happens to a JOIN query when there is no index on the foreign key column?

## NEXT UP

Tables, rows, and columns — the concrete SQL primitives that implement the relational model.
