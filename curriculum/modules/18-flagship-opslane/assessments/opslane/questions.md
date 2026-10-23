# Questions — Project Assessment: Opslane Flagship Capstone Defense

Answer each question with concrete reasoning and evidence.

---

### Question 1 (Multiple Choice)
During the Opslane capstone defense, a staff engineer asks: "What guarantees that a transient network failure between the Opslane service and PostgreSQL during an order creation request will not result in a corrupted or half-written order?"
- A) Go variables are stored in non-volatile RAM.
- B) The order creation repository wraps all line-item inserts and ledger balance updates within a single ACID transaction (`tx.BeginTx`), ensuring atomicity (all-or-nothing execution) and rolling back on any error.
- C) HTTP routers automatically retry POST requests until the database responds.
- D) Docker containers preserve database memory across power cycles.

---

### Question 2 (Multiple Choice)
Why does Opslane separate the read model from the write model (CQRS-lite) in order processing?
- A) Because Go does not support reading and writing to the same SQL table.
- B) To allow the read path to serve 95% of queries directly from Redis cache and database read replicas with sub-millisecond latency, while the write path focuses exclusively on strict ACID transaction validation and outbox event publishing.
- C) Because PostgreSQL requires different drivers for `SELECT` vs `INSERT`.
- D) It reduces the compiled binary size of the service.

---

### Question 3 (Multiple Choice)
How does the Opslane payment pipeline ensure that simultaneous requests with the same idempotency key do not execute concurrent credit card charges?
- A) By acquiring a database advisory lock or inserting a row with a unique constraint on `(tenant_id, idempotency_key)` before contacting the payment gateway, causing the concurrent request to fail or wait.
- B) By sleeping for 5 seconds before every payment call.
- C) By relying on the client's web browser to prevent multiple clicks.
- D) By disabling HTTP keep-alive on the payment endpoint.

---

### Question 4 (Multiple Choice)
What is the primary operational advantage of implementing an internal event bus with bounded worker pools over calling external microservices synchronously?
- A) It eliminates the need for unit testing.
- B) It shields the core API from downstream latency degradation, isolates panics, enables asynchronous backpressure, and decouples domain event publishing from third-party system availability.
- C) It allows Go code to bypass operating system memory limits.
- D) It converts relational SQL databases into document stores automatically.

---

### Question 5 (Multiple Choice)
When running Opslane in a container orchestration environment (Kubernetes), why must the application listen for both `SIGINT` and `SIGTERM` signals?
- A) `SIGINT` is sent by Kubernetes during rolling updates, and `SIGTERM` is only sent during system crashes.
- B) `SIGTERM` is the standard POSIX termination signal sent by Kubernetes before pod eviction, while `SIGINT` is sent during local development (Ctrl+C); handling both guarantees uniform graceful drain behavior across environments.
- C) Go runtime panics trigger `SIGTERM` automatically.
- D) Docker containers ignore `SIGTERM` by default.

---

### Question 6 (Multiple Choice)
How does Opslane prevent memory leaks in the caching layer when storing transient session and configuration data?
- A) By setting mandatory Time-To-Live (TTL) expiration on all Redis keys and configuring an allkeys-lru (Least Recently Used) eviction policy.
- B) By restarting the Redis server every 60 seconds via cron.
- C) By storing all cache keys in local global slices without capacity limits.
- D) Redis automatically deletes keys when the client HTTP connection closes.

---

### Question 7 (Oral Defense — End-to-End Request Flow)
Walk the defense panel through the complete lifecycle of a `POST /api/v1/orders` request in Opslane:
1. Ingress, middleware chain (correlation ID, logging, tenant context, authentication).
2. Domain validation, inventory check, and database transaction commit with outbox event.
3. Cache invalidation for the customer's order list.
4. Asynchronous event bus emission and worker pool consumption.
5. Response payload formatting and HTTP status code return.

---

### Question 8 (Oral Defense — Resilience Under Cascading Failures)
Suppose the Opslane Redis cluster crashes while the payment gateway API experiences a 5-second latency spike:
1. What happens to incoming order and payment requests?
2. How do circuit breakers and timeouts prevent database connection pool exhaustion?
3. How does the system degrade gracefully while maintaining data integrity?

---

### Question 9 (Oral Defense — Multi-Tenant Data Security)
Prove to the panel that Opslane is immune to cross-tenant data leaks:
1. Show how tenant context is propagated through Go's `context.Context`.
2. Explain what happens if an attacker attempts to fetch `GET /api/v1/orders/1234` where order `1234` belongs to Tenant B, but the attacker's JWT belongs to Tenant A.
3. How do automated integration tests verify tenant isolation boundaries?

---

### Question 10 (Oral Defense — Senior Engineering Trade-Offs)
Reflect on the design of Opslane:
1. What was the most difficult architectural trade-off made during development?
2. If Opslane needed to scale from 100,000 orders/day to 100,000,000 orders/day, what architectural bottlenecks would arise first, and how would you evolve the system?
