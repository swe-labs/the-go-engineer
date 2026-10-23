# Opslane Milestone 05 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) The method immediately returns `nil, ErrMissingTenantContext` without issuing a database query.**  
*Explanation:* `TenantFromContext` checks if the context contains a non-empty `tenantContextKey`. If missing, it returns `ErrMissingTenantContext`, failing fast before any database connection or query is initiated.

**2. b) It physically clusters records by `tenant_id` on disk and guarantees uniqueness within a tenant, while allowing efficient index range scans for tenant queries.**  
*Explanation:* A composite primary key `(tenant_id, id)` orders the index pages by `tenant_id` first. This groups all of a tenant's records sequentially in the B-tree index, maximizing CPU cache efficiency and minimizing disk seek times during tenant queries.

**3. b) By checking `res.RowsAffected() > 0`; if the order ID exists but belongs to another tenant, 0 rows are affected.**  
*Explanation:* The SQL update statement is `UPDATE ... WHERE tenant_id = ? AND id = ?`. If the order exists under another tenant, the `WHERE` condition matches 0 rows. Checking `RowsAffected() > 0` accurately reflects whether the caller's tenant possessed and modified the record.

**4. b) Any client or attacker can manipulate HTTP request headers, allowing them to impersonate any tenant in the system.**  
*Explanation:* HTTP headers are completely under client control. If an application blindly trusts `X-Tenant-ID`, any authenticated user could pass `X-Tenant-ID: target_corp` to view or mutate a competitor's data. Tenant identity must be extracted from cryptographically verified tokens.

**5. b) `nil, nil`, representing that the entity does not exist in this tenant's silo.**  
*Explanation:* Returning `nil, nil` allows the caller to treat missing records gracefully (e.g. mapping to HTTP 404 Not Found), cleanly distinguishing between "record does not exist in this tenant" and unexpected database errors.

**6. a) To return the newest records first for pagination and dashboard feeds.**  
*Explanation:* Chronological reverse ordering (`ORDER BY created_at DESC`) is standard for operational activity feeds, ensuring users see the most recent transactions first.

---

## Short Answer (2 points each)

**7. Security advantage of context-derived tenant ID over function arguments:**
- **Function Arguments:** If a repository signature is `GetByID(ctx, tenantID, orderID string)`, a developer can accidentally pass a tenant ID received from an unvalidated URL parameter or request body, enabling cross-tenant data leaks.
- **Context Derivation:** Deriving the tenant ID internally from `ctx` forces reliance on the upstream authentication middleware that verified the cryptographic signature of the caller's token. Developers cannot bypass this check, providing defense-in-depth across the entire codebase.
*Grading: 1 pt for explaining developer error with arguments; 1 pt for explaining cryptographic verification guarantee through context.*

**8. B-tree index clustering and composite primary keys:**
In relational databases (such as PostgreSQL or SQLite), table rows or index leaves are organized in B-tree structures. When `tenant_id` is the leading column of `(tenant_id, id)`, all rows for a given tenant share adjacent index leaf pages. When executing queries like `SELECT ... WHERE tenant_id = ?`, the database performs a single seek to the tenant's subtree and reads contiguous pages, minimizing random I/O and maximizing buffer cache hits.
*Grading: 1 pt for explaining B-tree physical ordering; 1 pt for explaining contiguous reads and cache locality.*

**9. PostgreSQL Row-Level Security (RLS) vs Application Filtering:**
- **How RLS Works:** PostgreSQL RLS policies evaluate an expression (e.g. `tenant_id = current_setting('app.tenant_id')`) directly in the database engine for every SQL statement, transparently filtering rows regardless of what SQL the application sends.
- **Pros:** True defense-in-depth: even a rogue query or SQL injection cannot view another tenant's rows.
- **Cons:** Added database complexity, requires setting session configuration parameters (`SET LOCAL app.tenant_id`) on every pooled connection checkout, and can complicate cross-tenant administrative reporting.
*Grading: 1 pt for describing database-level enforcement; 1 pt for analyzing pros (foolproof defense) and cons (connection pool complexity).*

**10. Insecure Direct Object Reference (IDOR) and Opslane Mitigation:**
- **IDOR Attack:** An attacker authenticated as Tenant A inspects their own order `ord-1234`. They notice the sequential or predictable ID pattern and issue an API call for `ord-1235` (belonging to Tenant B). If the system queries `SELECT * FROM orders WHERE id = 'ord-1235'`, Tenant A receives Tenant B's confidential order.
- **Opslane Mitigation:** Opslane binds the authenticated `tenant_id` to the query (`WHERE tenant_id = 'tenant_a' AND id = 'ord-1235'`). Because `ord-1235` does not belong to `tenant_a`, the query returns 0 rows, completely neutralizing the IDOR vulnerability.
*Grading: 1 pt for explaining the IDOR attack mechanism; 1 pt for demonstrating how the compound query prevents it.*
