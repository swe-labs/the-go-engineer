# Opslane Milestone 05 Assessment — Tenant Isolation: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** What happens if `IsolatedStore.GetByID` is called with a context that has no tenant attached?
- a) The method defaults to querying all tenants.
- b) The method immediately returns `nil, ErrMissingTenantContext` without issuing a database query.
- c) The database panics.
- d) It returns the first record in the database.

**2.** Why does `tenant_orders` use `PRIMARY KEY (tenant_id, id)` instead of `PRIMARY KEY (id)`?
- a) SQLite requires at least two columns in primary keys.
- b) It physically clusters records by `tenant_id` on disk and guarantees uniqueness within a tenant, while allowing efficient index range scans for tenant queries.
- c) To reduce database file size.
- d) To disable foreign key constraints.

**3.** In `IsolatedStore.Update`, how does the function determine if the update succeeded within the tenant's boundary?
- a) By checking if `err == nil`.
- b) By checking `res.RowsAffected() > 0`; if the order ID exists but belongs to another tenant, 0 rows are affected.
- c) By executing a second `SELECT` query immediately after.
- d) By inspecting the HTTP response code.

**4.** Why is trusting client-supplied headers like `X-Tenant-ID` without verifying cryptographic token claims an insecure multi-tenant design?
- a) `X-Tenant-ID` is not supported in HTTP/2.
- b) Any client or attacker can manipulate HTTP request headers, allowing them to impersonate any tenant in the system.
- c) Header values are limited to 8 characters.
- d) Headers are automatically stripped by modern web browsers.

**5.** When an entity is not found within the caller's tenant, what does `IsolatedStore.GetByID` return?
- a) An HTTP 500 error.
- b) `nil, nil`, representing that the entity does not exist in this tenant's silo.
- c) A blank `TenantOrder` struct with empty strings.
- d) `nil, sql.ErrNoRows`.

**6.** Why does `IsolatedStore.List` order by `created_at DESC`?
- a) To return the newest records first for pagination and dashboard feeds.
- b) Because SQL databases cannot sort in ascending order.
- c) To avoid using disk space.
- d) To bypass database locking.

---

## Short Answer (2 points each)

**7.** Explain why deriving `tenant_id` from `context.Context` (injected via verified authentication tokens) is more secure than accepting `tenant_id` as a function parameter in repository methods.

**8.** What is B-tree index clustering, and why does having `tenant_id` as the leading column in `PRIMARY KEY (tenant_id, id)` improve database performance?

**9.** How does PostgreSQL Row-Level Security (RLS) complement application-level tenant isolation? What are the pros and cons of relying on RLS vs application-level filtering?

**10.** Explain how an attacker might attempt IDOR (Insecure Direct Object Reference) against a multi-tenant API, and how Opslane's repository design mitigates it.
