# Opslane Milestone 02 Assessment — Database and Models: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** Why is `tenant_id` included in `OrderRepository.GetByID(ctx, tenantID, orderID)` when `orderID` is already a primary key?
- a) To satisfy SQLite index formatting syntax.
- b) To enforce strict multi-tenant boundary isolation and prevent horizontal privilege escalation (cross-tenant data leakage).
- c) Because primary keys in SQL databases cannot be queried individually.
- d) To reduce query cache utilization.

**2.** In `User` struct definition, what does the tag `json:"-"` on `PasswordHash` achieve?
- a) Tells SQLite to omit the column.
- b) Prevents the hashed password from ever being serialized into JSON HTTP responses, preventing credential exposure.
- c) Disables password hashing in tests.
- d) Requires users to supply a blank password.

**3.** What is the danger of omitting `db.SetMaxOpenConns(...)` on a production `*sql.DB` connection pool?
- a) Queries will run in read-only mode.
- b) Under high concurrent load, Go will spawn thousands of concurrent database connections, exhausting database server memory and triggering connection rejection errors.
- c) SQLite will delete the database file.
- d) Migrations will run in reverse order.

**4.** Why does `RunMigrations` execute each migration's DDL and record insertion inside a database transaction (`db.Begin()`)?
- a) To speed up disk write operations by 50%.
- b) To ensure atomicity: if any DDL statement fails, changes are completely rolled back and no partial schema or incorrect migration version is recorded.
- c) Because DDL commands cannot be run outside of transactions in SQL.
- d) To prevent concurrent HTTP reads.

**5.** What happens if `notes` in `orders` is NULL in the database and is scanned into a regular Go `string` variable without `sql.NullString`?
- a) Go assigns `""` to the string silently.
- b) `row.Scan()` returns a runtime error (`sql: Scan error on column index ... converting NULL to string is unsupported`).
- c) The entire row is deleted.
- d) The database connection closes.

**6.** Why does `payments` define a unique constraint on `idempotency_key`?
- a) To prevent duplicate payment charges when network retries occur from client checkout attempts.
- b) To link the payment to the user's credit card company.
- c) To enforce that payments are greater than zero.
- d) To sort payments chronologically.

---

## Short Answer (2 points each)

**7.** Explain how Go's `*sql.DB` manages its internal connection pool. Contrast `MaxOpenConns`, `MaxIdleConns`, and `ConnMaxLifetime`.

**8.** What is the purpose of composite index `CREATE INDEX idx_orders_tenant_status ON orders(tenant_id, status)`? Why is index column ordering important?

**9.** Why is `context.Context` accepted as the first parameter to repository methods (`Create`, `GetByID`) and passed to `ExecContext` / `QueryRowContext`?

**10.** In a multi-replica containerized deployment, what problem arises if all replicas run `RunMigrations()` concurrently upon startup? How is this resolved in enterprise systems?
