# Multi-tenancy architecture

## Learning objective

Implement tenant isolation models (silent, shared, hybrid) in Go, propagate tenant context through the application using context.Context, and scope database queries by tenant identifier at the row level.

## Why this matters

Every B2B SaaS platform serves multiple customers (tenants) from a single deployment. Each tenant's data must be isolated: tenant A should never see tenant B's data. The architecture decision -- how to isolate tenants -- affects database schema, connection pooling, query complexity, deployment cost, and operational burden. Choosing wrong means either overpaying for separate databases or creating cross-tenant data leaks. Go engineers building SaaS products must understand tenant isolation to make the right tradeoff.

## Mental model

A multi-tenant building has separate offices on each floor. In the **silent** model, each tenant gets their own building (separate database). No noise between tenants, but expensive. In the **shared** model, all tenants share the same open floor plan (same database, same tables). Cheap, but every query must check the tenant badge. In the **hybrid** model, small tenants share a floor and large tenants get their own floor. This maps directly to database isolation strategies.

The tenant identifier flows through every request like a badge: from HTTP request to middleware, into `context.Context`, through business logic, and into every SQL query as a `WHERE tenant_id = ?` clause.

## Core idea

Three isolation models:

| Model | Database per tenant | Shared tables | Shared infrastructure |
|---|---|---|---|
| Silent | Yes | No | No |
| Shared | No | Yes | Yes |
| Hybrid | Large tenants: yes; small tenants: no | Yes | Yes |

**Silent (database-per-tenant)**: each tenant gets their own database instance (or schema). Maximum isolation. No risk of cross-tenant data leaks. Best for regulated industries (healthcare, finance). Drawbacks: expensive to operate, migrations must run N times.

**Shared (row-level isolation)**: all tenants share the same tables. Every row has a `tenant_id` column. Every query includes `WHERE tenant_id = ?`. Cheap to operate, simple migrations. Risk: one buggy query without the tenant filter leaks data.

**Hybrid**: combines both. Large tenants with compliance requirements get dedicated databases. Small tenants share. The application routes to the right database based on the tenant identifier.

## Under the hood

Tenant context propagation in Go:

1. HTTP middleware extracts tenant ID from header, JWT claim, or subdomain.
2. Middleware calls `context.WithValue(ctx, tenantKey, tenantID)`.
3. The handler passes `ctx` to the service layer.
4. The service layer passes `ctx` to the repository layer.
5. The repository extracts `tenantID` from `ctx` and adds it to every SQL query.

A tenant-aware repository looks like:

```go
type Repository struct {
	db *sql.DB
}

func (r *Repository) GetOrders(ctx context.Context) ([]Order, error) {
	tenant := ctx.Value(tenantKey).(string)
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, total FROM orders WHERE tenant_id = $1", tenant)
	// ...
}
```

The key invariant: no repository method should accept a bare `tenantID` parameter -- it must come from `context.Context` to ensure it cannot be accidentally omitted or swapped.

## How Go uses it

- **`context.Context`**: the idiomatic Go mechanism for propagating request-scoped values. Every HTTP handler, gRPC interceptor, and background worker starts with a context that carries the tenant identifier.
- **`database/sql`**: `QueryContext`, `ExecContext`, `PrepareContext` all accept a context. The driver cancels queries if the context is cancelled.
- **`chi` / `gorilla/mux`**: middleware extracts tenant from request and injects into context.
- **`sqlx` / `pgx`**: popular Go database libraries that support `QueryContext` and prepared statements with tenant scoping.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"sync"
)

type Tenant string

type ctxKey string

const tenantKey ctxKey = "tenant"

func WithTenant(ctx context.Context, t Tenant) context.Context {
	return context.WithValue(ctx, tenantKey, t)
}

func TenantFromContext(ctx context.Context) (Tenant, bool) {
	t, ok := ctx.Value(tenantKey).(Tenant)
	return t, ok
}

type TenantStore struct {
	mu   sync.Mutex
	data map[string]map[string]string
}

func NewTenantStore() *TenantStore {
	return &TenantStore{data: make(map[string]map[string]string)}
}

func (s *TenantStore) Set(ctx context.Context, key, val string) error {
	tenant, ok := TenantFromContext(ctx)
	if !ok {
		return fmt.Errorf("no tenant in context")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data[string(tenant)] == nil {
		s.data[string(tenant)] = make(map[string]string)
	}
	s.data[string(tenant)][key] = val
	return nil
}

func (s *TenantStore) Get(ctx context.Context, key string) (string, error) {
	tenant, ok := TenantFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("no tenant in context")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[string(tenant)][key]
	if !ok {
		return "", fmt.Errorf("key not found for %s", tenant)
	}
	return val, nil
}

type TenantConfig struct {
	Tenant         Tenant
	MaxConnections int
}

type ConfigManager struct {
	mu      sync.RWMutex
	configs map[Tenant]TenantConfig
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{configs: make(map[Tenant]TenantConfig)}
}

func (m *ConfigManager) Set(cfg TenantConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configs[cfg.Tenant] = cfg
}

func (m *ConfigManager) Get(ctx context.Context) (TenantConfig, error) {
	tenant, ok := TenantFromContext(ctx)
	if !ok {
		return TenantConfig{}, fmt.Errorf("no tenant in context")
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	cfg, ok := m.configs[tenant]
	if !ok {
		return TenantConfig{}, fmt.Errorf("no config for %s", tenant)
	}
	return cfg, nil
}

func main() {
	ctx1 := WithTenant(context.Background(), Tenant("acme"))
	ctx2 := WithTenant(context.Background(), Tenant("globex"))

	store := NewTenantStore()
	store.Set(ctx1, "theme", "dark")
	store.Set(ctx2, "theme", "light")
	store.Set(ctx2, "timezone", "UTC")

	theme1, _ := store.Get(ctx1, "theme")
	theme2, _ := store.Get(ctx2, "theme")
	tz2, _ := store.Get(ctx2, "timezone")
	fmt.Printf("acme theme: %s\nglobex theme: %s\nglobex tz: %s\n", theme1, theme2, tz2)

	cfg := NewConfigManager()
	cfg.Set(TenantConfig{Tenant: "acme", MaxConnections: 50})
	cfg.Set(TenantConfig{Tenant: "globex", MaxConnections: 10})
	if c, _ := cfg.Get(ctx1); c.MaxConnections == 50 {
		fmt.Println("acme config verified")
	}
}
```

## Step-by-step execution

For `store.Set(ctx1, "theme", "dark")`:

1. `TenantFromContext(ctx1)` extracts `"acme"` from context.
2. `s.mu.Lock()` acquires the mutex.
3. Checks `s.data["acme"]` -- nil, so allocates the map.
4. Stores `data["acme"]["theme"] = "dark"`.
5. Releases mutex.

For `store.Get(ctx1, "theme")`:

1. `TenantFromContext(ctx1)` extracts `"acme"`.
2. Acquires mutex.
3. Returns `data["acme"]["theme"]` → `"dark"`.
4. If `ctx2` were used instead, it would return `"light"` from a completely separate map entry.

The mutex ensures concurrent reads/writes are safe. The tenant-keyed map guarantees isolation: tenant data never overlaps.

## Common mistakes

- **Passing tenant as a separate parameter**: every function that needs the tenant should receive it via context, not as a function argument. A bare `tenantID string` parameter is easy to forget or swap.
- **Forgetting the tenant filter in queries**: a single query without `WHERE tenant_id = ?` leaks all tenants' data. Use automated query analysis or a database proxy that injects the filter.
- **Storing tenant-scoped data in a global cache**: a shared `sync.Map` or Redis instance that is not keyed by tenant mixes data. Always prefix cache keys with the tenant identifier.
- **Connection pooling across tenants**: if tenant A exhausts the connection pool, tenant B cannot connect. Use per-tenant connection pools or configure connection limits per tenant in the database.
- **Assuming one model fits all**: a startup with 5 tenants needs shared isolation. An enterprise with 500 tenants each with 10 TB of data needs silent isolation. Re-evaluate the model as the product grows.

## Debugging walkthrough

A customer reports seeing another company's orders:

```go
func GetOrders(ctx context.Context) ([]Order, error) {
	rows, err := db.Query("SELECT id, total FROM orders")
	if err != nil {
		return nil, err
	}
	// scan rows...
}
```

**Symptom**: The `GetOrders` function returns all orders from all tenants.

**Investigation**: Search for database queries in the repository layer. Find the one missing the tenant filter. In this case, there is no `WHERE tenant_id = ?` clause.

**Root cause**: The developer forgot to extract the tenant from context and add it to the query. The function signature accepts `ctx context.Context` but never uses it for tenant filtering.

**Fix**:

```go
func GetOrders(ctx context.Context) ([]Order, error) {
	tenant := ctx.Value(tenantKey).(string)
	rows, err := db.QueryContext(ctx,
		"SELECT id, total FROM orders WHERE tenant_id = $1", tenant)
	// ...
}
```

**Prevention**: Write a linter that enforces tenant filtering on every repository method. Write integration tests that insert data for two tenants and verify cross-tenant isolation.

## Production notes

- **Tenant resolution**: extract the tenant identifier from the JWT `sub` claim, the request path (`/api/:tenant/...`), or the Host header. Never trust the client to provide the tenant -- derive it from authentication.
- **Silent model operations**: with database-per-tenant, schema migrations must run N times. Use a migration tool that supports tenant iteration (e.g., `golang-migrate` with a tenant loop).
- **Audit logging**: every query should log the tenant ID. If a cross-tenant leak is suspected, the audit log helps trace which query served which tenant's data.
- **Rate limiting**: apply per-tenant rate limits to prevent a noisy tenant from degrading service for others.
- **Shared model scaling**: as shared database grows, consider read replicas or sharding by tenant group. Citus (PostgreSQL) supports distributed tables with tenant colocation.

## Performance implications

- **Silent model**: each tenant has dedicated resources. No noisy neighbor. Connection pool per tenant. Cost scales linearly with tenant count.
- **Shared model**: one database serves all tenants. Connection pool contention is a risk. Query performance degrades as the table grows. Index the `tenant_id` column -- every query filters on it, so a composite index starting with `(tenant_id, ...)` is essential.
- **Hybrid model**: routing logic adds microseconds per request. The application must maintain a tenant-to-database mapping (usually cached). This is negligible compared to the query execution time.

## Practice task

Add a third tenant "tenant-c" and extend the tenant store with a `Delete(ctx, key)` method. Write a function `ListKeys(ctx) ([]string, error)` that returns all keys for the current tenant. Use context scoping to ensure that `ListKeys` for tenant A does not return keys from tenant B.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/13-multi-tenancy-architecture
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/13-multi-tenancy-architecture
```

The tests verify tenant context propagation, data isolation across tenants, missing-tenant errors, and config per-tenant scoping.

## Review questions

1. What are the three tenant isolation models, and what tradeoffs does each make?
2. Why should the tenant identifier be propagated through `context.Context` rather than as a separate function parameter?
3. What happens if a shared-model query omits the `WHERE tenant_id = ?` clause?
4. When would you choose the silent model over the shared model?
5. How would you prevent a noisy tenant from exhausting the shared connection pool?

## NEXT UP

Caching architecture.
