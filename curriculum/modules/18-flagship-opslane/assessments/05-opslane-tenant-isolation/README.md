# Opslane tenant isolation

## Mission

Enforce strict multi-tenant data isolation across the persistence layer, combining context-bound tenant propagation, composite primary keys `(tenant_id, id)`, and automatic query scoping to prevent cross-tenant data leakage.

## Prerequisites

- `opslane-04`

## Mental Model

In multi-tenant SaaS architectures, data leaks between corporate accounts represent existential business and legal failures. A naive design where repository methods take an unverified `tenantID` string as an argument relies on developer discipline and is prone to human error. Opslane enforces an architectural defense-in-depth model: repository methods extract the verified tenant identity strictly from `context.Context` via `TenantFromContext()`. If the context lacks tenant identity, the operation fails immediately (`ErrMissingTenantContext`). At the database level, composite primary keys `(tenant_id, id)` and mandatory `WHERE tenant_id = ?` clauses ensure that even if an attacker acquires an entity ID from a competing tenant, the database returns no rows.

## Visual Model

```text
┌─────────────────────────────────────────────────────────────┐
│ Authenticated HTTP Request Context                          │
│   ctx = ContextWithTenant(ctx, "tenant-acme")               │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ IsolatedStore Data Access Methods                           │
│   Create, GetByID, List, Update, Delete                     │
│   ├─► tenantID, err := TenantFromContext(ctx)               │
│   └─► If err != nil -> Abort: ErrMissingTenantContext       │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Relational SQL Execution                                    │
│   - Primary Key: (tenant_id, id)                            │
│   - SELECT ... WHERE tenant_id = ? AND id = ?               │
│   - UPDATE ... WHERE tenant_id = ? AND id = ?               │
│   - DELETE ... WHERE tenant_id = ? AND id = ?               │
└─────────────────────────────────────────────────────────────┘
```

## Machine View

Under the hood:
1. `TenantFromContext` retrieves the unexported `tenantContextKey` from the active context. Missing or blank keys abort before executing database transactions.
2. The table `tenant_orders` defines `PRIMARY KEY (tenant_id, id)`. This B-tree index layout clusters all rows belonging to a tenant together physically on disk, optimizing cache locality and range scans.
3. Every mutating statement (`UPDATE`, `DELETE`) inspects `res.RowsAffected()`. If a caller attempts to update or delete an order with ID `"ord-1"` but their context specifies `tenant-beta` while the order belongs to `tenant-alpha`, `RowsAffected` evaluates to `0`, preventing unauthorized modification.

## Run Instructions

```bash
# Run tenant isolation test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/05-opslane-tenant-isolation

# Run standalone verification demo
go run ./curriculum/modules/18-flagship-opslane/assessments/05-opslane-tenant-isolation
```

## Code Walkthrough

1. **Context Management (lines 14–32)**: Implements `ContextWithTenant` and `TenantFromContext` with sentinel `ErrMissingTenantContext`.
2. **`TenantOrder` Struct (lines 35–41)**: Models orders with explicit `tenant_id` association.
3. **`IsolatedStore` & Schema (lines 44–64)**: Sets up `tenant_orders` table with composite primary key `(tenant_id, id)`.
4. **Isolated CRUD Operations (lines 67–178)**:
   - `Create`: Extracts tenant from context; binds to new order.
   - `GetByID`: Enforces `tenant_id = ? AND id = ?`; returns `nil, nil` on `sql.ErrNoRows`.
   - `List`: Scopes orders strictly to `WHERE tenant_id = ? ORDER BY created_at DESC`.
   - `Update` & `Delete`: Modifies rows strictly within the tenant boundary and verifies `RowsAffected() > 0`.

## Try It

1. Run the test suite: verify that tenant A cannot access, update, or delete tenant B's orders.
2. Attempt an operation with `context.Background()` lacking tenant context: verify `ErrMissingTenantContext` is returned.
3. Inspect database row count: verify that multi-tenant orders coexist in the same table without collisions.

Common mistakes:
- Passing tenant IDs as function arguments rather than deriving them from validated context tokens.
- Omission of `tenant_id` in `UPDATE` or `DELETE` statements, which can accidentally overwrite or purge all tenants' data.
- Relying on client-supplied headers (e.g. `X-Tenant-ID`) without verifying cryptographic JWT signatures.

## In Production

In production PostgreSQL deployments, organizations often supplement application-level context isolation with Row-Level Security (RLS) policies (`CREATE POLICY tenant_isolation ON tenant_orders USING (tenant_id = current_setting('app.current_tenant'))`). This ensures defense-in-depth: even if an engineer writes a raw SQL query forgetting `WHERE tenant_id = ?`, the database engine automatically restricts visibility.

## Thinking Questions

1. Why is deriving `tenant_id` from `context.Context` more secure than passing it as a URL query parameter or function argument?
2. What are the query performance benefits of composite primary key `(tenant_id, id)` compared to a standalone UUID `id`?
3. How does PostgreSQL Row-Level Security (RLS) reinforce application-level tenant isolation?

## Next Step

Proceed to [Milestone 06 — HTTP API Layer](../06-opslane-http-api-layer/README.md) to integrate routing, request decoding, standard JSON envelopes, and RESTful CRUD endpoints.
