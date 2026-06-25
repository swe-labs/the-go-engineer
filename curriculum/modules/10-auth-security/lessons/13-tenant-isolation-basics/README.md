# Tenant isolation basics

## Learning objective

Implement tenant isolation in multi-tenant Go services using database-level row filtering, context-based tenant propagation, and middleware enforcement to prevent cross-tenant data leakage.

## Why this matters

Every SaaS platform is multi-tenant: a single instance serves multiple customers (tenants) while keeping each tenant's data completely isolated. A cross-tenant data leak is a catastrophic security failure — one missing `WHERE tenant_id = ?` clause exposes all customers' data. The 2020 SolarWinds breach and numerous SaaS data leaks trace back to tenant isolation failures. In Go services, tenant isolation must be enforced at multiple layers: application middleware, database queries, and infrastructure.

## Mental model

A multi-tenant database is an apartment building. Each tenant has their own locked apartment door (tenant_id filter). The building manager (application) can enter any apartment for maintenance (admin override), but tenants cannot enter each other's apartments. There are three floor plans: everyone in their own building (separate database), their own floor (separate schema), or adjacent apartments with locks (shared table with tenant_id). The safest is each building, but the most practical is shared tables with row-level security.

## Core idea

Three tenant isolation strategies:

| Strategy | Isolation | Complexity | Cost | Best for |
|---|---|---|---|---|
| Separate database | Strongest | High (connection management) | Highest | Enterprise, compliance-heavy |
| Separate schema | Strong | Medium (schema per tenant) | Medium | Mid-market |
| Shared table (discriminator column) | Weakest | Low (WHERE tenant_id) | Lowest | Most SaaS apps |

The shared table approach uses a `tenant_id` column on every table:

```sql
CREATE TABLE documents (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_documents_tenant ON documents(tenant_id);
```

Every query MUST include `WHERE tenant_id = $1`. Missing this filter leaks data.

Database-level enforcement via PostgreSQL Row-Level Security (RLS):

```sql
ALTER TABLE documents ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON documents
    USING (tenant_id = current_setting('app.tenant_id')::UUID);
```

RLS ensures that even if the application forgets the filter, the database enforces it.

## Under the hood

Tenant context propagation in Go uses `context.Context`:

1. Middleware extracts tenant ID from the request (subdomain, header, JWT claim).
2. Stores tenant ID in context: `context.WithValue(ctx, tenantKey, tenantID)`.
3. Database layer reads tenant ID from context and adds it to every query.
4. Repository functions always filter by tenant_id.

```go
type contextKey string
const tenantKey contextKey = "tenant_id"

func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tenantID := extractTenant(r)
        ctx := context.WithValue(r.Context(), tenantKey, tenantID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (r *Repo) GetDocuments(ctx context.Context) ([]Document, error) {
    tenantID := ctx.Value(tenantKey).(string)
    rows, _ := r.db.QueryContext(ctx,
        "SELECT id, title FROM documents WHERE tenant_id = $1", tenantID)
    // ...
}
```

## How Go uses it

Go services implement tenant isolation at multiple layers:

- **Middleware layer**: extracts tenant from subdomain, header (`X-Tenant-ID`), or JWT claims.
- **Repository layer**: all database queries include `WHERE tenant_id = $1`.
- **Database layer**: PostgreSQL RLS as defense-in-depth.
- **Cache layer**: Redis keys prefixed with tenant ID: `tenant:{id}:session:{user_id}`.
- **File storage**: S3 keys prefixed with tenant ID: `tenants/{id}/uploads/{filename}`.

The Go `database/sql` package's `QueryContext` and `ExecContext` enable passing context with tenant data to database calls.

## Go example

```go
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type tenantKeyType string

const tenantKey tenantKeyType = "tenant_id"

// Document is a tenant-scoped resource.
type Document struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"-"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// TenantStore simulates a shared database with tenant isolation.
type TenantStore struct {
	mu        sync.RWMutex
	documents []Document
}

func NewTenantStore() *TenantStore {
	return &TenantStore{}
}

func (ts *TenantStore) AddDocument(ctx context.Context, title, content string) (*Document, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 16)
	rand.Read(b)
	doc := Document{
		ID:        hex.EncodeToString(b),
		TenantID:  tenantID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}
	ts.mu.Lock()
	ts.documents = append(ts.documents, doc)
	ts.mu.Unlock()
	return &doc, nil
}

func (ts *TenantStore) GetDocuments(ctx context.Context) ([]Document, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	var result []Document
	for _, d := range ts.documents {
		if d.TenantID == tenantID {
			result = append(result, d)
		}
	}
	return result, nil
}

// GetDocument retrieves a single document with tenant isolation.
func (ts *TenantStore) GetDocument(ctx context.Context, docID string) (*Document, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	for _, d := range ts.documents {
		if d.ID == docID && d.TenantID == tenantID {
			return &d, nil
		}
	}
	return nil, errors.New("document not found")
}

// BUGGY: GetDocumentWithoutTenantFilter retrieves a document without tenant check.
// This demonstrates a cross-tenant leak vulnerability.
func (ts *TenantStore) GetDocumentWithoutTenantFilter(docID string) (*Document, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	for _, d := range ts.documents {
		if d.ID == docID {
			return &d, nil
		}
	}
	return nil, errors.New("document not found")
}

func getTenantID(ctx context.Context) (string, error) {
	tid, ok := ctx.Value(tenantKey).(string)
	if !ok || tid == "" {
		return "", errors.New("tenant not found in context")
	}
	return tid, nil
}

// TenantMiddleware extracts tenant ID and stores it in context.
func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tenantID string
		// Try header first, then subdomain, then JWT claim.
		tenantID = r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			// Extract from subdomain: tenant.example.com -> "tenant"
			parts := strings.Split(r.Host, ".")
			if len(parts) >= 2 {
				tenantID = parts[0]
			}
		}
		if tenantID == "" {
			http.Error(w, `{"error":"tenant not identified"}`, http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), tenantKey, tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TenantIsolationMiddleware logs cross-tenant access attempts.
func TenantIsolationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		// Log any request that tries to access without tenant context.
		if tenantID == "" && r.URL.Path != "/health" {
			log.Printf("[WARN] Request without tenant context: %s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	store := NewTenantStore()

	// Pre-populate data for two tenants.
	ctx1 := context.WithValue(context.Background(), tenantKey, "acme-corp")
	ctx2 := context.WithValue(context.Background(), tenantKey, "globex-inc")

	store.AddDocument(ctx1, "Acme Strategy 2026", "Confidential Acme plans")
	store.AddDocument(ctx1, "Acme Product Roadmap", "Q3-Q4 features")
	store.AddDocument(ctx2, "Globex Financial Report", "Q2 earnings")

	// Demonstrate proper isolation.
	fmt.Println("=== Correct tenant isolation ===")
	docs1, _ := store.GetDocuments(ctx1)
	fmt.Printf("Acme-Corp documents (%d):\n", len(docs1))
	for _, d := range docs1 {
		fmt.Printf("  - %s\n", d.Title)
	}

	docs2, _ := store.GetDocuments(ctx2)
	fmt.Printf("Globex-Inc documents (%d):\n", len(docs2))
	for _, d := range docs2 {
		fmt.Printf("  - %s\n", d.Title)
	}

	// Demonstrate cross-tenant leak (BUG).
	fmt.Println("\n=== Cross-tenant leak (bug without tenant filter) ===")
	leaked, _ := store.GetDocumentWithoutTenantFilter(docs1[0].ID)
	fmt.Printf("Acme document retrieved without tenant check: '%s'\n", leaked.Title)
	fmt.Println("SECURITY ISSUE: Any tenant can access any document!")

	// HTTP handlers with tenant isolation.
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/api/documents", TenantMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			docs, err := store.GetDocuments(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(docs)
		case "POST":
			var req struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			doc, err := store.AddDocument(r.Context(), req.Title, req.Content)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(doc)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	var h http.Handler = mux
	h = TenantIsolationMiddleware(h)

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}

// Ensure errors is used.
var _ = errors.New
```

## Step-by-step execution

For `GET /api/documents` with `X-Tenant-ID: acme-corp`:

1. `TenantIsolationMiddleware` logs the request (no warning — tenant ID is present).
2. `TenantMiddleware` extracts `acme-corp` from the header, stores it in context.
3. Handler calls `store.GetDocuments(ctx)`.
4. `GetDocuments` calls `getTenantID(ctx)` which returns `"acme-corp"`.
5. Loop filters documents: only those with `TenantID == "acme-corp"` are returned.
6. Two documents (Acme Strategy, Acme Product Roadmap) are returned.

For `GET /api/documents` with `X-Tenant-ID: globex-inc`:

1-4. Same flow, tenant ID is `"globex-inc"`.
5. Only the Globex Financial Report is returned.

For `GET /api/documents` without the `X-Tenant-ID` header:

1. `TenantMiddleware` finds no header and no subdomain. Returns 400 Bad Request.

## Common mistakes

- Using `SELECT *` from shared tables without a `WHERE tenant_id = ?` clause — one missing filter leaks all tenants' data.
- Implementing tenant isolation in application code without database-level enforcement — application bugs leak data. Use PostgreSQL RLS as a safety net.
- Storing `tenant_id` as a string without validation — path traversal-like attacks on tenant context could access other tenants.
- Sharing connection pools across tenants without context propagation — a goroutine from tenant A might execute a query with tenant B's credentials.
- Hardcoding default tenant IDs — forgetting to set the tenant context results in empty queries or wrong data.
- Not indexing `tenant_id` — queries without a tenant_id index cause full table scans on a shared table.

## Debugging walkthrough

Consider this repository function:

```go
func (r *Repo) GetUserByID(ctx context.Context, userID string) (*User, error) {
    row := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", userID)
    // BUG: No tenant_id filter! Any tenant can query any user.
    var u User
    err := row.Scan(&u.ID, &u.Name, &u.Email)
    return &u, err
}
```

**Symptom**: A support agent from tenant A can look up user details for tenant B's users by guessing their user IDs.

**Investigation**: The query does not include `AND tenant_id = $1`. The `ctx` carries the tenant ID, but it is not used in the query. Any user ID returns the user record regardless of tenant.

**Root cause**: Missing `tenant_id` filter in the SQL query. The developer forgot to extract and pass the tenant ID.

**Fix**:

```go
func (r *Repo) GetUserByID(ctx context.Context, userID string) (*User, error) {
    tenantID, err := getTenantID(ctx)
    if err != nil {
        return nil, err
    }
    row := r.db.QueryRowContext(ctx,
        "SELECT id, name, email FROM users WHERE id = $1 AND tenant_id = $2",
        userID, tenantID)
    var u User
    err = row.Scan(&u.ID, &u.Name, &u.Email)
    return &u, err
}
```

## Production notes

- Always use parameterized queries with tenant_id. String concatenation for tenant_id is a SQL injection vector.
- Enable PostgreSQL Row-Level Security on all tenant-scoped tables. RLS ensures that even if a query misses the tenant filter, the database enforces isolation.
- Use a consistent naming convention for the tenant context key across all services (e.g., `"x-tenant-id"`).
- For shared database deployments, monitor query patterns for missing tenant filters. Use database audit logs.
- Cache the tenant context in a `context.Context` value. Do not store it in global variables — goroutines would share it across tenants.
- For file storage, prefix object keys with tenant ID: `tenants/{tenant_id}/uploads/{file_id}`. This enables tenant-scoped delete and S3 bucket policies.

## Performance implications

- Tenant isolation adds a `WHERE tenant_id = $1` clause to every query. With an index on `tenant_id`, this is O(log n) on the index plus O(m) on the result set.
- PostgreSQL RLS adds a small overhead per query (~5-10 microseconds) for policy evaluation.
- Context propagation is free — `context.WithValue` is O(1).
- Separate database per tenant: connection pool overhead scales with tenant count. Use connection pooling middleware.
- Shared table with tenant_id index: excellent performance for up to millions of rows per tenant. At billions of rows, consider partitioning by tenant_id.

## Practice task

Write Go functions `NewTenantStore()`, `CreateDocument(ctx context.Context, title, content string) (*Document, error)`, and `GetDocuments(ctx context.Context) ([]Document, error)` where:
- `CreateDocument` requires a valid tenant ID in context.
- `GetDocuments` filters by tenant ID from context.
- Also write `GetDocumentsNoFilter(ctx context.Context) ([]Document, error)` that returns all documents (intentionally buggy — for testing).
- Then write a `main()` that creates documents for two tenants, demonstrates proper isolation, and demonstrates the cross-tenant leak via `GetDocumentsNoFilter`.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/13-tenant-isolation-basics
go test ./curriculum/modules/10-auth-security/lessons/13-tenant-isolation-basics
```

The test file `main_test.go` contains table-driven tests that verify:
- `CreateDocument` returns an error when context has no tenant ID.
- `GetDocuments` returns only documents for the correct tenant.
- `GetDocuments` returns different results for different tenants.
- `GetDocumentsNoFilter` returns documents from all tenants (cross-tenant leak).

## Review questions

1. What are the three main tenant isolation strategies? Which is most common for SaaS applications and why?
2. Why is `WHERE tenant_id = $1` necessary in every query, even with PostgreSQL RLS enabled?
3. How does `context.Context` enable tenant isolation in Go middleware?
4. Describe a cross-tenant data leak vulnerability and how PostgreSQL Row-Level Security prevents it.
5. If a Go service uses a shared database with tenant_id filtering, what happens if the middleware fails to extract the tenant ID and the handler proceeds?

## NEXT UP

SQL injection prevention — parameterized queries, prepared statements, and defense-in-depth against SQL injection in Go.
