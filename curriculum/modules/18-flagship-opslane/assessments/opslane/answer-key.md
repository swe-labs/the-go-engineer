# Answer Key — Project Assessment: Opslane Flagship Capstone Defense

## Multiple Choice Answers

### Question 1: B
**Explanation:** Wrapping all write operations in a database ACID transaction (`tx.BeginTx`) guarantees Atomicity, Consistency, Isolation, and Durability. If a network blip, query failure, or system error occurs at any point before `tx.Commit()`, the database rolls back all partial writes, returning the database to its clean pre-transaction state without data corruption.

### Question 2: B
**Explanation:** In high-throughput platforms, 95% or more of all traffic consists of read queries (viewing orders, checking status, fetching tenant configs). Separating read models allows caching hot entities in Redis with sub-millisecond responses and routing complex queries to PostgreSQL read replicas. Meanwhile, the write path remains lean and focused on transactional ACID integrity, preventing read traffic from competing with write locks.

### Question 3: A
**Explanation:** Enforcing idempotency keys via an atomic unique constraint (`UNIQUE(tenant_id, idempotency_key)`) in PostgreSQL or an atomic `SET NX EX` in Redis guarantees that when concurrent requests arrive with the same key, only the first request succeeds in acquiring the lock. The second request either waits or immediately receives an "in-flight / duplicate" error response, preventing duplicate credit card charges.

### Question 4: B
**Explanation:** Synchronous RPC or HTTP calls between microservices create tight temporal coupling: if a downstream service slows down from 20ms to 2000ms, upstream handlers block, exhausting connection pools and crashing the ingress server. An internal event bus with bounded worker pools decouples event emission from consumption, isolates downstream panics, and allows the system to buffer tasks or apply backpressure during spikes.

### Question 5: B
**Explanation:** In production environments such as Kubernetes or ECS, orchestrators send `SIGTERM` when scaling down, restarting pods, or performing rolling deployments. In local development or CLI execution, developers use `SIGINT` (Ctrl+C). Handling both signals in Go via `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)` ensures seamless graceful draining across both development and production.

### Question 6: A
**Explanation:** Without TTL expiration and memory eviction bounds, an in-memory cache will grow indefinitely as new keys are added, eventually exhausting host memory and triggering Linux OOM killer termination. Combining TTL expirations with Redis `allkeys-lru` guarantees that expired keys are reclaimed and least-recently-used keys are evicted if memory limits are reached.

---

## Technical & Oral Defense Answers

### Question 7: End-to-End Request Flow Walkthrough
1. **Ingress & Middleware**:
   - The Go 1.22 HTTP router receives `POST /api/v1/orders`.
   - Middleware assigns a unique correlation ID (`X-Correlation-ID`), starts an OpenTelemetry span, logs the request with `slog`, authenticates the JWT bearer token, and injects `tenant_id` into `context.Context`.
2. **Domain Execution & Transactional Outbox**:
   - The handler decodes the JSON body and validates domain invariants (positive quantities, valid item IDs).
   - Starts a PostgreSQL transaction (`tx, _ := db.BeginTx(ctx, nil)`).
   - Verifies customer balance / inventory and inserts the `orders` and `order_items` records with `tenant_id`.
   - Inserts an `OrderCreated` record into the `outbox` table within the same transaction.
   - Commits the transaction (`tx.Commit()`).
3. **Cache Invalidation & Asynchronous Publishing**:
   - Deletes cached order lists for the tenant: `redis.Del(ctx, fmt.Sprintf("tenant:%s:orders", tenantID))`.
   - The outbox processor emits `OrderCreated` to the internal event bus.
   - Bounded worker pools receive the event and asynchronously execute notifications and analytics tracking.
4. **Response**: Returns `HTTP 201 Created` with the JSON order representation.

### Question 8: Resilience Under Cascading Failures
1. **Redis Cluster Outage**:
   - Opslane's cache-aside client wraps Redis calls in a circuit breaker or timeout (50ms).
   - If Redis fails, the cache layer logs a warning metric (`cache_errors_total.Inc()`) and falls back to querying PostgreSQL directly.
2. **Payment Gateway Latency Spike**:
   - Payment requests route through an asynchronous queue or dedicated worker pool with strict context timeouts (2 seconds).
   - Because payments are processed asynchronously via workers rather than synchronously in the order creation HTTP handler, HTTP API latency remains under 10ms and database connections are freed immediately.
   - If payment gateway timeouts exceed thresholds, circuit breakers trip to open, queuing payment retries with exponential backoff and jitter.

### Question 9: Multi-Tenant Data Security Proof
1. **Context Scoping**: The tenant ID is parsed from the cryptographically verified JWT and bound to a private struct context key.
2. **Access Control**: When an attacker in Tenant A requests `GET /api/v1/orders/1234` (where order 1234 belongs to Tenant B):
   - The repository query executes: `SELECT * FROM orders WHERE id = '1234' AND tenant_id = 'Tenant-A'`.
   - The query returns `sql.ErrNoRows`, and the API responds with `HTTP 404 Not Found`.
   - The system never reveals whether order 1234 exists in another tenant.
3. **Automated Verification**: The test suite includes explicit multi-tenant cross-tenant test cases verifying that Tenant A credentials cannot read, update, or delete Tenant B resources.

### Question 10: Senior Engineering Trade-Offs & Scaling to 100M Orders/Day
1. **Difficult Trade-Off**:
   - Balancing strict ACID transactional consistency with low-latency throughput. Choosing the Transactional Outbox Pattern introduced slight eventual consistency for downstream notifications, but was essential to eliminate dual-write data loss and decouple payment latencies from API response times.
2. **Scaling to 100M Orders/Day (~1,200 orders/sec steady, 10,000 orders/sec peak)**:
   - **Database Sharding**: Transition from a single PostgreSQL instance to tenant-sharded databases (using Citus or application-level tenant routing).
   - **Outbox CDC**: Replace polling workers with Change Data Capture (Debezium + Apache Kafka) streaming PostgreSQL WAL events directly to distributed consumer worker clusters.
   - **Distributed Rate Limiting**: Deploy lock-striped Envoy / Go edge rate limiters with local token buckets to protect the platform against traffic spikes.
