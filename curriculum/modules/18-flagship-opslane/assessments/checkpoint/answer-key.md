# Answer Key — Module 18 Checkpoint: Flagship Opslane Architecture & Capstone Defense

## Multiple Choice Answers

### Question 1: B
**Explanation:** Enforcing tenant isolation exclusively through manual `WHERE tenant_id = $1` filters in application SQL relies on 100% human infallibility across every query, join, and subquery. A single omission creates an Insecure Direct Object Reference (IDOR) or leaks confidential customer data across tenants. Implementing defense-in-depth via PostgreSQL Row-Level Security (RLS) guarantees that even if application logic misses a clause, the database kernel strictly isolates rows based on session variables.

### Question 2: B
**Explanation:** External payment providers and third-party APIs have variable network latencies (hundreds of milliseconds to seconds) and can experience intermittent network stalls. Performing synchronous HTTP calls inside a database transaction keeps database connections and row-level locks open for the entire duration of the external call. Under moderate load, this rapidly exhausts the database connection pool (`pgxpool`), causing connection starvation and cascading failures across the entire system.

### Question 3: B
**Explanation:** When a popular cached record expires, hundreds of concurrent requests will experience a cache miss simultaneously and attempt to query the database to regenerate the cache. Using single-flight request coalescing (`singleflight.Group`) guarantees that only one database query is executed, while all concurrent requests wait for that single call to return and share the resulting value, completely eliminating the thundering herd spike.

### Question 4: B
**Explanation:** Graceful shutdown requires an orderly drain of upstream to downstream dependencies:
1. Stop the HTTP server listener to prevent new ingress connections, and allow active in-flight requests to complete.
2. Signal the internal event bus and worker pools to process and finish all queued background tasks.
3. Close the PostgreSQL connection pool once all database operations from active requests and background workers have finished.
4. Flush telemetry buffers and structured logger output to ensure zero observability data is lost.

### Question 5: B
**Explanation:** Goroutines and asynchronous event buses do not automatically share HTTP request contexts. When an event is published, the publisher serializes the OpenTelemetry W3C trace context (`traceparent`, `tracestate`) into the event payload metadata. Consuming workers extract this metadata into a new child context using `otel.GetTextMapPropagator().Extract()`, preserving complete end-to-end distributed trace spans across asynchronous boundaries.

### Question 6: B
**Explanation:** Network requests can fail after an order is processed by the server but before the client receives the 200 OK response. Clients automatically retry failed requests. By requiring an `Idempotency-Key` header and storing the key with its result inside an atomic database table, duplicate requests return the cached result of the original execution without recharging the customer or creating duplicate orders.

---

## Technical & Architectural Answers

### Question 7: System Architecture Defense (The Dual-Write Problem & Transactional Outbox)
1. **The Dual-Write Problem**:
   - In naive systems, code executes: `tx.Commit()` followed by `eventBus.Publish(event)`. If the application crashes, network fails, or event bus is temporarily down immediately after the DB commit, the database state has changed, but the event is permanently lost, causing inconsistent downstream state (e.g., payment processed, but inventory never deducted and confirmation email never sent).
2. **The Transactional Outbox Pattern**:
   - Inside the same local PostgreSQL transaction that writes the order record, insert an event record into an `outbox` table (`id`, `aggregate_type`, `aggregate_id`, `payload`, `created_at`, `status`).
   - Because both the order and outbox record write within the same ACID transaction, either both are committed atomically or neither is.
3. **Guaranteed Delivery via CDC or Poller**:
   - A dedicated background worker polls the `outbox` table (`SELECT ... FOR UPDATE SKIP LOCKED`), publishes events to the event bus or Kafka, and updates the status to `processed`. Alternatively, a Change Data Capture (CDC) tool like Debezium reads the PostgreSQL WAL directly. This guarantees at-least-once event delivery under all failure scenarios.

### Question 8: Multi-Tenant Security Defense
1. **JWT Extraction**: The authentication middleware parses the signed JWT, verifies cryptographic signatures, and extracts the `tenant_id` claim.
2. **Type-Safe Context Binding**:
   ```go
   type tenantContextKey struct{}
   ctx := context.WithValue(r.Context(), tenantContextKey{}, tenantID)
   ```
   An unexported struct key prevents collisions with other packages.
3. **Repository Layer Enforcement**:
   - Repository methods extract `tenantID` from `ctx`. If absent, the query immediately returns an error.
   - All SQL statements use parameterized queries: `SELECT * FROM orders WHERE id = $1 AND tenant_id = $2`.
   - PostgreSQL RLS session variables (`SET LOCAL app.current_tenant_id = $1`) enforce row isolation at the database kernel level.
4. **Cache Key Namespacing**:
   - Redis keys are formatted as `tenant:{tenant_id}:{entity}:{id}`, ensuring that cache keys can never collide or leak between tenants.

### Question 9: Concurrency & Backpressure Defense
1. **Bounded Worker Pool**:
   - Fixes the number of worker goroutines (e.g., 32 workers) reading from a buffered task channel (e.g., capacity 1000). Goroutines are never spawned dynamically per request, keeping heap allocation and CPU context-switching constant regardless of ingress rate.
2. **Backpressure Strategy**:
   - When the queue buffer reaches capacity, the producer `select` block handles saturation:
     ```go
     select {
     case pool.taskQueue <- task:
         // Successfully queued
     default:
         // Queue full: trigger backpressure
         metrics.TasksRejectedTotal.Inc()
         w.WriteHeader(http.StatusTooManyRequests)
         json.NewEncoder(w).Encode(map[string]string{"error": "system busy, please retry later"})
     }
     ```
3. **Prometheus Telemetry**:
   - `opslane_worker_pool_active_workers`: Gauge tracking currently busy workers.
   - `opslane_worker_pool_queue_depth`: Gauge tracking buffered items.
   - `opslane_task_execution_duration_seconds`: Histogram measuring task latency. Alerts trigger if queue depth exceeds 80% for 60 seconds.

### Question 10: Production Incident & Disaster Recovery Defense
1. **Root Causes of Pool Starvation**:
   - Missing database index on high-volume queries causing sequential table scans.
   - Unhandled transactions lacking `defer tx.Rollback()`, causing abandoned locks and leased connections.
   - Calling slow third-party HTTP endpoints inside active database transaction locks.
2. **Diagnostic Workflow**:
   - Goroutine dump via `pprof`: Thousands of goroutines blocked on `pgxpool.Acquire()`.
   - Database telemetry: `pg_stat_activity` showing connections in `idle in transaction` state for extended durations.
3. **Permanent Remediation**:
   - Strict connection timeouts: `pgxpool.Config.MaxConnLifetime`, `MaxConnIdleTime`, and `AcquireTimeout` set to 3 seconds.
   - External API calls moved entirely outside database transaction boundaries.
   - Static analysis linting ensuring all transaction blocks have immediate `defer tx.Rollback()`.
