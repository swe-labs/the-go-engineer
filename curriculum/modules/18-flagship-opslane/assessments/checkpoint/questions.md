# Questions — Module 18 Checkpoint: Flagship Opslane Architecture & Capstone Defense

Answer each question with concrete reasoning and evidence.

---

### Question 1 (Multiple Choice)
In the Opslane multi-tenant architecture, what is the primary risk of relying solely on the application code (e.g. manually appending `WHERE tenant_id = $1`) to enforce tenant isolation rather than combining it with database-level mechanisms (like PostgreSQL Row-Level Security)?
- A) Manual SQL string concatenation causes PostgreSQL to disable index scans.
- B) A single developer oversight or missed query in a large codebase can cause cross-tenant data leaks, exposing private customer data to unauthorized tenants.
- C) Multi-tenant indexes cannot be created without Row-Level Security enabled.
- D) PostgreSQL limits connection pools to one tenant per database connection.

---

### Question 2 (Multiple Choice)
Why does Opslane decouple order payment processing from order creation using an internal event bus rather than executing payment calls synchronously inside the HTTP handler transaction?
- A) Synchronous payment API calls exceed the maximum call stack depth of the Go runtime.
- B) External payment gateways have variable network latency and failure rates; executing them inside the HTTP transaction holds database row locks and connections open, degrading throughput and causing cascading timeouts across the platform.
- C) Go HTTP handlers cannot invoke third-party HTTPS APIs.
- D) PostgreSQL requires all external network I/O to be performed through database triggers.

---

### Question 3 (Multiple Choice)
In Opslane Milestone 11 (Caching Layer), how does the system prevent the "thundering herd" (cache stampede) problem when a heavily accessed tenant configuration record expires?
- A) By restarting the Redis server whenever a cache miss occurs.
- B) By using single-flight request coalescing (`golang.org/x/sync/singleflight`) or mutex-guarded background refresh, ensuring only one upstream database query executes while concurrent requests await the single result.
- C) By setting TTL to infinite and never expiring cache keys.
- D) By disabling Redis and storing all cached data in unbuffered Go channels.

---

### Question 4 (Multiple Choice)
During Opslane graceful shutdown (Milestone 13), in what order must subsystems be terminated to guarantee zero data loss and zero dropped HTTP requests?
- A) Close database connection pool -> Terminate worker pools -> Stop HTTP server listener -> Flush logs.
- B) Stop HTTP server listener and wait for in-flight requests to complete -> Drain event bus and finish queued tasks in worker pools -> Close database connection pool -> Flush telemetry/logs.
- C) Kill all worker goroutines via `runtime.Goexit` -> Terminate database connections -> Exit immediately.
- D) Close Redis connection -> Shutdown database -> Drop incoming HTTP connections.

---

### Question 5 (Multiple Choice)
How does Opslane guarantee distributed trace context propagation across the asynchronous event bus and worker pools?
- A) By passing the global HTTP request pointer to the worker goroutine.
- B) By serializing the OpenTelemetry `traceparent` header into the event payload metadata, and extracting it into a new child context (`otel.GetTextMapPropagator().Extract`) inside the consuming worker.
- C) OpenTelemetry automatically traces all goroutines globally without context propagation.
- D) By storing the trace ID in a global static variable shared across all goroutines.

---

### Question 6 (Multiple Choice)
Why does Opslane enforce idempotency keys on payment pipeline execution (Milestone 08)?
- A) Because Go HTTP clients automatically compress JSON payloads.
- B) Because transient network timeouts can cause a client to retry an order; idempotency keys ensure duplicate requests do not process duplicate credit card charges or double-spend balances.
- C) Because PostgreSQL requires unique keys on every incoming network socket.
- D) Idempotency keys are only required for read queries (`GET`).

---

### Question 7 (System Architecture Defense — Dual-Write Problem)
In an oral capstone defense, the panel asks:
*"In your Order Processing milestone, what happens if the database transaction commits successfully, but the network crashes before your code publishes the `OrderCreated` event to the message bus? How does Opslane resolve this dual-write problem?"*
Provide a senior-level defense:
1. Explain the architectural flaw of naive dual-writes (DB commit + network publish).
2. Detail the Transactional Outbox Pattern implementation in Opslane.
3. Show how a background poller or CDC pipeline guarantees at-least-once event delivery.

---

### Question 8 (Multi-Tenant Security Defense)
Explain how Opslane prevents Cross-Tenant Data Leaks across the entire request lifecycle:
1. Extraction: How tenant claims in JWT tokens are parsed and validated by middleware.
2. Context binding: How a type-safe context key prevents collision with other context values.
3. Database enforcement: How repositories inject the tenant ID into parameterized queries and how PostgreSQL Row-Level Security (RLS) acts as a secondary defense layer.
4. Caching key namespacing: How cache keys are scoped to `tenant:{id}:{resource}:{id}` to prevent cross-tenant cache contamination.

---

### Question 9 (Concurrency & Backpressure Defense)
During a flash-sale event, Opslane receives 10,000 orders per second:
1. Explain how the Bounded Worker Pool (Milestone 10) prevents goroutine explosion and memory exhaustion.
2. What happens when the worker pool task queue is full? Detail your backpressure strategy (HTTP 429 Too Many Requests, rejected promises, or shed load).
3. How do you monitor worker saturation using Prometheus metrics (active workers, queue depth, task execution latency)?

---

### Question 10 (Production Incident & Disaster Recovery Defense)
Describe an incident simulation in Opslane where PostgreSQL database connection pool (`pgxpool`) was completely exhausted:
1. What were the root causes (unindexed queries, uncommitted transactions, slow third-party API calls held under lock)?
2. How would you identify connection pool starvation using `pprof` goroutine dumps and database connection stats?
3. What architectural remediation did you institute in Opslane to make connection exhaustion impossible?
