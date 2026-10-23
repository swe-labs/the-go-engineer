# Answer Key — Portfolio Review Pack & Interview Defense

## Multiple Choice Answers

### Question 1: B
**Explanation:** `sync.Map` is designed for cache-like workloads where keys are mostly read and rarely modified or added. In high-concurrency write scenarios (such as token bucket rate limiters), writes to new keys dirty the read map, triggering promotions under an internal mutex (`mu.Lock()`), leading to catastrophic contention. In contrast, 64-shard lock striping partitions keys across 64 independent mutexes via `hash(key) % 64`. Each lock operates on an independent memory address, reducing false sharing and cache-line bouncing across multi-core systems, achieving over 38 million operations/sec with zero allocations.

### Question 2: B
**Explanation:** An empty `_test.go` file or one with disabled build tags (`//go:build ignore`) satisfies filesystem existence (`os.Stat` returns `nil`), but provides zero verification of application behavior. A robust auditor checks non-trivial file length and runs `go test` programmatically to verify assertions execute and pass cleanly.

### Question 3: B
**Explanation:** In Go SQL drivers, acquiring a transaction reserves a dedicated connection from the connection pool. If an unhandled error, panic, or early return happens before `tx.Commit()`, the connection remains leased and uncommitted locks remain held inside the database engine. In PostgreSQL, this blocks concurrent updates to the same rows and eventually starves the connection pool (`max_open_conns`). Calling `defer tx.Rollback()` is safe because if `tx.Commit()` succeeds, `Rollback()` returns `sql.ErrTxDone` and is a harmless no-op.

### Question 4: B
**Explanation:** Hexagonal architecture (ports and adapters) isolates core business logic inside pure Go packages that have zero dependencies on HTTP routers, ORMs, or third-party SDKs. The domain defines interfaces (ports) for repositories and external services, which infrastructure adapters implement. This architecture allows swapping PostgreSQL for an in-memory test database, or replacing HTTP with gRPC, without touching domain business invariants.

### Question 5: B
**Explanation:** In Go, goroutines do not terminate when a caller function returns or when an HTTP client disconnects unless they explicitly listen to a cancellation signal. Passing `context.Context` down to worker tasks and selecting on `ctx.Done()` ensures that when an HTTP request aborts or a timeout expires, workers immediately stop processing, return to the pool, and release held resources.

### Question 6: B
**Explanation:** Statically linked Go binaries compiled with `CGO_ENABLED=0` do not require dynamic linkers (`ld-linux.so`) or C standard libraries (`glibc`). Placing the binary in a `scratch` or `distroless` container eliminates hundreds of potential CVEs from system packages, shells, and utilities, while shrinking container images from hundreds of megabytes to under 20MB.

---

## Technical & Architecture Defense Answers

### Question 7: System Architecture Defense (PostgreSQL ACID vs DynamoDB)
1. **Consistency Guarantees & Double-Spend Prevention**:
   - Order processing, billing, and balance accounting require strict consistency (serializable or read-committed with row-level locks). In a relational database, `SELECT ... FOR UPDATE` or transactional decrementing (`UPDATE accounts SET balance = balance - 100 WHERE id = 1 AND balance >= 100`) guarantees atomic balance checks.
   - In an eventually consistent NoSQL store (DynamoDB or Cassandra), concurrent requests can read stale balances, creating catastrophic double-spending conditions without complex conditional writes or distributed transactions.
2. **Latency Trade-offs & Throughput**:
   - PostgreSQL handles thousands of transactional writes per second on modern NVMe storage with tuned connection pooling (`pgxpool`), connection pooling proxies (`PgBouncer`), and write-ahead log (WAL) tuning.
3. **Scaling Strategy Beyond Single-Node Limits**:
   - Read scaling: Deploy streaming read replicas for reporting, read-heavy catalog queries, and dashboard traffic.
   - Write scaling: Implement tenant-based sharding or application-level partitioning where tenant orders map to independent databases.
   - Asynchronous decoupled processing: Use the Transactional Outbox Pattern to record order events in a local Postgres outbox table within the same ACID transaction, followed by a CDC pipeline (Debezium / Kafka) publishing to asynchronous event workers.

### Question 8: Concurrency & Failure Recovery Defense (Worker Pool Panics)
1. **Worker Panic Recovery Pattern**:
   ```go
   func (w *Worker) run() {
       go func() {
           for task := range w.taskQueue {
               func() {
                   defer func() {
                       if r := recover(); r != nil {
                           slog.Error("worker recovered from panic",
                               "err", r,
                               "stack", string(debug.Stack()),
                           )
                           metrics.WorkerPanicsTotal.Inc()
                       }
                   }()
                   task.Execute()
               }()
           }
       }()
   }
   ```
2. **Pool Preservation**:
   - Placing `defer recover()` inside the per-task loop closure ensures that a panicking task is caught, logged with its stack trace, and metrics are incremented without killing the long-running worker goroutine. The loop advances to the next task in the channel without leaking pool concurrency.

### Question 9: Memory Profiling & Allocation Elimination Defense
1. **Diagnostic Workflow**:
   - Run `go tool pprof http://localhost:6060/debug/pprof/heap` with `alloc_space` to identify where total memory is allocated, and `inuse_space` to identify persistent memory footprint.
   - Run `go tool trace` to visualize GC pause frequency, STW latency, and goroutine scheduling delays.
2. **Escape Analysis**:
   - The Go compiler (`go build -gcflags="-m"`) checks if a pointer created within a function outlives the stack frame. If a variable is returned as a pointer, stored in an interface, or passed to a function that escapes, it escapes to the heap, requiring garbage collection.
3. **Three Production Optimization Techniques**:
   - **`sync.Pool`**: Reusing byte buffers and temporary structs across HTTP requests eliminates GC churn on hot parsing paths.
   - **Slice Preallocation**: Initializing slices with known capacity (`make([]Item, 0, expectedLen)`) prevents multiple backing array reallocation and copying cycles during `append`.
   - **Avoiding Interface Boxing on Hot Paths**: Passing concrete types rather than `any` / `interface{}` prevents runtime heap allocations required for type descriptor wrapping.

### Question 10: Production Deployment & Graceful Shutdown Defense
1. **Kubernetes Termination Coordination**:
   - When a pod is deleted, the Kubernetes API removes the pod IP from the Service endpoints asynchronously while simultaneously sending `SIGTERM` to the container process. If the container terminates immediately, requests still in-flight in kube-proxy or ingress load balancers receive TCP RST and HTTP 502.
2. **Graceful Drain Sequence in Go**:
   ```go
   ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
   defer stop()

   <-ctx.Done()
   slog.Info("shutdown signal received, draining traffic...")

   // Step 1: Wait briefly for load balancer / kube-proxy iptables propagation
   time.Sleep(5 * time.Second)

   // Step 2: Shutdown HTTP listener, allowing active in-flight requests to complete
   drainCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
   defer cancel()

   if err := server.Shutdown(drainCtx); err != nil {
       slog.Error("server forced to shutdown", "err", err)
   }
   slog.Info("server exited gracefully")
   ```
   - This guarantees 100% zero-downtime rolling deploys.
