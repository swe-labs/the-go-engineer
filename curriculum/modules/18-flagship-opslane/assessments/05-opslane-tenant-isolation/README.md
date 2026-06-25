# Opslane tenant isolation

## Mission

Understand and apply Opslane tenant isolation in the context of professional Go software engineering.

## Prerequisites

- opslane-04

## Mental Model

Tenant isolation is a data firewall between tenants. Every database query must include a tenant_id filter. Every cache key must include the tenant_id. Every queue message must include the tenant_id. Every API response must include only the requesting tenant's data. The tenant_id is not a parameter — it is a security boundary enforced at the infrastructure layer (database, cache, queue) and validated at the application layer (middleware, service, repository).

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Tenant isolation in Go is enforced at the repository layer using context-scoped values. The middleware stores the tenant ID in context using a typed key (type tenantKeyType struct{} to prevent key collisions). The repository base type (or a query builder wrapper) extracts the tenant ID and appends it to every SQL query's WHERE clause. For PostgreSQL, a composite index on (tenant_id, ...) ensures efficient filtering even with millions of rows per tenant. For Redis, cache keys include the tenant_id prefix (e.g., acme-corp:order:123). For message queues, the message header includes the tenant_id, and consumers filter by tenant_id before processing. Testing tenant isolation requires two tenants with identical data shapes and randomized data — the test asserts that queries scoped to tenant A return only A's data even when B's IDs are known.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/05-opslane-tenant-isolation
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Adding WHERE tenant_id = ? to queries but forgetting it in one query — a single un-scoped query leaks all tenants' data to every tenant, and the bug is invisible until a customer reports seeing another tenant's data.
- Using tenant_id from the URL path instead of the authenticated session — a tenant A user can change the URL path to tenant B's ID and access B's data if the middleware does not validate ownership.
- Sharing a single connection pool across all tenants — one tenant's slow query blocks the pool, starving all other tenants and causing cross-tenant latency impact.
- Including tenant_id in error messages or logs — a stack trace with tenant-specific identifiers in a shared log aggregator leaks tenant information to the on-call engineer who may not have authorization for that tenant.
- Using a single database schema for all tenants without a tenant_id column — adding tenant isolation after launch requires a full data migration with downtime.

## In Production

Every SaaS platform implementing multi-tenancy depends on tenant isolation. GitHub uses organization-level isolation — every query includes the org ID. Slack uses workspace-level isolation. Opslane uses tenant-level isolation with a tenant_id in every table and every API response. A tenant leak at GitHub, Slack, or Opslane would be a P0 security incident with legal and regulatory consequences.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-06`.
