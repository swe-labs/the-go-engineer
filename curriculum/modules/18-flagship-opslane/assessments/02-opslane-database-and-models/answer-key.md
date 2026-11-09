# Opslane Milestone 02 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) To enforce strict multi-tenant boundary isolation and prevent horizontal privilege escalation (cross-tenant data leakage).**  
*Explanation:* In multi-tenant systems, querying by entity ID alone allows an authenticated user from Tenant A to inspect records belonging to Tenant B simply by guessing or acquiring Tenant B's order ID. Including `WHERE tenant_id = ? AND id = ?` guarantees the database returns rows only if they belong to the requesting tenant.

**2. b) Prevents the hashed password from ever being serialized into JSON HTTP responses, preventing credential exposure.**  
*Explanation:* The `json:"-"` struct tag instructs Go's `encoding/json` encoder to skip the field entirely during marshaling, ensuring that password hashes never leak across public HTTP API responses.

**3. b) Under high concurrent load, Go will spawn thousands of concurrent database connections, exhausting database server memory and triggering connection rejection errors.**  
*Explanation:* By default, Go's `*sql.DB` has no limit on open connections (`MaxOpenConns = 0`). Under traffic spikes, thousands of simultaneous queries each open a new physical socket to the database, rapidly exceeding PostgreSQL or MySQL connection limits (e.g. `max_connections = 100`).

**4. b) To ensure atomicity: if any DDL statement fails, changes are completely rolled back and no partial schema or incorrect migration version is recorded.**  
*Explanation:* Wrapping migrations in `sql.Tx` ensures transactional atomicity. If a syntax error or constraint conflict occurs halfway through a migration, `tx.Rollback()` reverts all executed statements, preventing the database from ending up in a corrupted or half-migrated state.

**5. b) `row.Scan()` returns a runtime error (`sql: Scan error on column index ... converting NULL to string is unsupported`).**  
*Explanation:* Standard Go primitive types (`string`, `int`, `time.Time`) cannot represent SQL `NULL`. Scanning a `NULL` column into a `string` variable results in an immediate scan error. Using `sql.NullString` allows safe inspection of `notes.Valid` and `notes.String`.

**6. a) To prevent duplicate payment charges when network retries occur from client checkout attempts.**  
*Explanation:* If a client sends a payment request and the network drops before receiving the response, the client will retry. The database unique constraint on `idempotency_key` rejects the duplicate insert, guaranteeing the customer is charged exactly once.

---

## Short Answer (2 points each)

**7. Go's `*sql.DB` connection pool settings:**
- `MaxOpenConns`: The absolute ceiling on the number of open connections (active + idle) allowed to the database. Prevents connection exhaustion on the database server.
- `MaxIdleConns`: The maximum number of unused, open connections retained in the idle pool waiting for new queries. Keeps warm sockets ready to eliminate TCP/TLS handshake latency.
- `ConnMaxLifetime`: The maximum duration a connection may be reused before being closed. Recycles aging connections before firewalls or load balancers silently terminate stale idle connections.
*Grading: 1 pt for explaining open vs idle bounds; 1 pt for explaining max lifetime / stale socket recycling.*

**8. Composite index mechanics and ordering:**
- **Purpose:** Enables index lookups for queries that filter orders by tenant and order status (e.g. `WHERE tenant_id = 't1' AND status = 'PENDING'`).
- **Column Ordering:** B-tree composite indices can only be queried using left-to-right prefixes. Having `tenant_id` as the leading column allows queries filtering by `tenant_id` alone to use the index, or queries filtering by `(tenant_id, status)` together. An index starting with `status` would not be efficient for queries scoped strictly by `tenant_id`.
*Grading: 1 pt for explaining query filtering efficiency; 1 pt for explaining the left-prefix rule of composite B-tree indices.*

**9. Purpose of `context.Context` in repository methods:**
Passing `context.Context` to `ExecContext` and `QueryRowContext` enables query cancellation and timeout propagation. If the upstream HTTP client disconnects or an upstream timeout occurs, the context is cancelled. The Go database driver intercepts this cancellation signal and terminates the running SQL query on the database server immediately, reclaiming CPU and lock resources.
*Grading: 1 pt for explaining cancellation / timeout propagation; 1 pt for explaining reclaiming database server resources.*

**10. Migration races in multi-replica deployments:**
- **Problem:** If multiple service replicas boot concurrently in Kubernetes, they may all detect unapplied migrations simultaneously, attempt to execute duplicate DDL statements, and deadlock or crash with table-already-exists errors.
- **Enterprise Solutions:**
  1. **Distributed Migration Locks:** Using database-level advisory locks (e.g. `pg_advisory_lock` in PostgreSQL) so only one replica acquires the migration lock while others wait.
  2. **Dedicated CI/CD or Init Containers:** Decoupling migrations from application startup by executing migrations as a discrete Kubernetes `Job` or deployment step before new pods are launched.
*Grading: 1 pt for describing the concurrency collision / race condition; 1 pt for identifying advisory locks or CI/CD job separation as a solution.*
