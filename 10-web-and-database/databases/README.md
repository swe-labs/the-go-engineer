# Section 12: Databases

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| SQL Driver Binding | Beginner | High | Blank imports `_ "github.com/..."` |
| Queries | Intermediate | High | Connection pooling, executing contexts |
| Transactions | Advanced | Medium | ACID state rollovers |
| Repository Pattern | Advanced | **Critical** | Interface driven data-access layer |

## Engineering Depth
The `database/sql` standard package does not map objects directly (No built-in ORM). Instead, it heavily relies on **Connection Pooling**. When you execute `db.Query()`, Go checks out a connection from its underlying thread-safe connection pool. 
- Critical bug zone: Every `*sql.Rows` instance returned by `db.Query()` holds a live database connection. If `rows.Close()` is not called (usually via `defer`), you will exhaust the database connection pool, leading to catastrophic thread-locking under load.

## References
1. **[Go Database]** [Accessing relational databases](https://go.dev/doc/database/)

---

## 🏗 Exercise: CRUD SQLite Repository (`6-repository`)

### Step-by-Step Instructions & Hints
1. **Architecture Definition:** Define an entity `User{}` and a `UserRepository{}` interface outlining `Create`, `GetByID`, etc.
2. **Concrete Implementation:** Build an `SQLiteUserRepository` struct wrapping an initialized `*sql.DB` connection pool.
3. **Execute Safely:** Inside `GetByID`, use `db.QueryRow()` to pull a single result mapped instantly to the struct pointers.
   - *Hint:* Distinguish between `sql.ErrNoRows` (Clean empty response) vs global system errors.
4. **Leverage the Context:** Always use `QueryContext` or `ExecContext` passing a generated `context.Context` to allow application timeout shutdowns to natively kill dangling database requests!


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| DB.1 | [connecting](./1-connecting-to-db) | Blank import · sql.Open lazy · db.Ping · connection pool | 🟢 entry |
| DB.2 | [query — INSERT](./2-query) | db.Exec · ? parameterisation · LastInsertId · bcrypt hash | DB.1 |
| DB.3 | [query — SELECT](./3-select) | db.QueryRow · rows.Scan · rows.Close · rows.Err | DB.1, DB.2 |
| DB.4 | [prepared statements](./4-prepare) | db.Prepare · stmt.ExecContext · when to prepare explicitly | DB.2, DB.3 |
| DB.5 | [transactions](./5-transactions) | BeginTx · defer Rollback · Commit · ACID all-or-nothing | DB.1, DB.2, DB.3 |
