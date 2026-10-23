# Questions — Portfolio Review Pack & Interview Defense

Answer each question with concrete reasoning, architectural trade-offs, and operational evidence.

---

### Question 1 (Multiple Choice)
When presenting a high-throughput rate limiter in a senior technical interview, an interviewer asks why you used 64-shard lock striping rather than `sync.Map`. Which explanation demonstrates the deepest understanding of Go runtime mechanics?
- A) `sync.Map` uses mutexes internally for all operations, making it slower than striped mutexes in read-heavy workloads.
- B) `sync.Map` is optimized for stable key sets and read-dominated access; concurrent writes to keys cause dirty map promotions and global `mu` lock contention, whereas 64-shard striping isolates lock contention to `hash(key) % 64` and eliminates cache-line bouncing across CPU cores.
- C) `sync.Map` does not support storing struct values directly, requiring memory allocation on every read.
- D) Striped mutexes prevent goroutines from sleeping when waiting for a lock, using CPU spinlocks instead.

---

### Question 2 (Multiple Choice)
In an automated portfolio auditor, what is the primary risk of relying solely on `os.Stat` to verify that unit tests exist (`*_test.go`)?
- A) `os.Stat` cannot read files inside hidden or private directories.
- B) An empty `_test.go` file or a test file with `// +build ignore` will pass the file existence check but provides zero test coverage and cannot be compiled or executed by `go test`.
- C) `os.Stat` loads the entire file into heap memory, causing memory pressure during large audits.
- D) Go test binaries must be compiled prior to checking `os.Stat`.

---

### Question 3 (Multiple Choice)
When defending your PostgreSQL database implementation (Module 09) during a system design interview, why did you implement explicit database transaction rollbacks via `defer tx.Rollback()` instead of relying on the Go garbage collector or connection pool to clean up failed transactions?
- A) The garbage collector terminates transactions by sending `SIGKILL` to PostgreSQL backends.
- B) If a panic or early return occurs before `tx.Commit()`, uncommitted transactions hold row-level locks and active connections indefinitely until connection idle timeout, causing connection pool exhaustion and database deadlocks.
- C) PostgreSQL rejects new queries on a database if any transaction object is garbage collected without an explicit commit.
- D) `defer tx.Rollback()` automatically commits the transaction if no errors occurred in the enclosing function.

---

### Question 4 (Multiple Choice)
Why is hexagonal architecture (ports and adapters) preferred over standard layered MVC architecture for long-lived backend microservices subject to enterprise interview evaluation?
- A) Hexagonal architecture runs 2x faster because Go interfaces are optimized by the compiler into direct branch jumps.
- B) Hexagonal architecture enforces inward dependency flow: business logic and domain entities never import external frameworks, databases, or HTTP libraries; external concerns implement domain ports, enabling deterministic unit testing with mock adapters.
- C) Hexagonal architecture requires zero interface definitions in Go source code.
- D) Hexagonal architecture automatically generates OpenAPI and Swagger specs at compile time.

---

### Question 5 (Multiple Choice)
In your bounded worker pool implementation (Module 11), what guarantees that the server will not leak goroutines when client connections abruptly disconnect?
- A) The worker pool channel is closed immediately upon receiving a new connection.
- B) Workers listen to both the task queue and `ctx.Done()`, ensuring cancellation propagates down to active workers, which exit immediately and release pool capacity.
- C) The Go runtime automatically kills any goroutine whose parent HTTP handler has returned.
- D) Using unbuffered channels guarantees goroutines are reclaimed by runtime GC cycles.

---

### Question 6 (Multiple Choice)
Why should production container builds for Go services utilize `CGO_ENABLED=0` and target `scratch` or `distroless` images?
- A) Non-CGO binaries can run without an operating system kernel.
- B) It produces statically linked binaries with zero external libc/glibc dependencies, minimizing CVE attack surface by eliminating package managers and shells, and reducing image transfer latency to under 20MB.
- C) Go binaries compiled with CGO cannot bind to TCP ports lower than 1024.
- D) Docker containers cannot execute dynamically linked binaries regardless of base image.

---

### Question 7 (System Architecture Defense)
In your Opslane Flagship or Task API project, an interviewer challenges:
*"Why did you choose relational ACID transactions in PostgreSQL instead of an eventually consistent DynamoDB or Cassandra cluster for order processing and ledger updates?"*
Provide a senior-level defense:
1. Contrast the consistency guarantees and business risk of dirty reads / double spending.
2. Discuss read/write latency trade-offs at 10,000 orders/sec.
3. Explain how you would scale reads and writes if PostgreSQL becomes the bottleneck (e.g., read replicas, sharding, outbox pattern with Kafka).

---

### Question 8 (Concurrency & Failure Engineering Defense)
During a live interview, the interviewer asks:
*"In your Bounded Worker Pool, what happens if a worker encounters a panic during task execution? Does the entire server crash, does the worker permanently die, or does the pool leak capacity?"*
Explain the exact defensive engineering required:
1. Show the `recover()` pattern inside the worker goroutine loop.
2. Explain how pool capacity is preserved so other queued tasks continue processing.
3. Detail how error and panic telemetry is logged with stack traces and reported to monitoring systems.

---

### Question 9 (Performance & Memory Profiling Defense)
An interviewer presents a pprof heap profile showing high GC CPU usage and frequent pauses:
1. What runtime metrics (`GOGC`, alloc_space, inuse_space) and profiling tools (`go tool pprof`, `trace`) would you use to diagnose the root cause?
2. Explain how escape analysis determines stack vs heap allocation in Go.
3. Describe three concrete optimization techniques implemented in your portfolio (e.g., `sync.Pool`, slice preallocation, zero-copy buffer reuse) to eliminate allocations on hot paths.

---

### Question 10 (Production Deployment & Graceful Shutdown Defense)
Explain the sequence of events during a production deployment of your containerized service:
1. How does Kubernetes coordinate between endpoint deregistration (kube-proxy / ingress) and container process shutdown (`SIGTERM`)?
2. Why is an immediate `os.Exit(0)` upon receiving `SIGTERM` an anti-pattern?
3. Detail the step-by-step Go code implementation using `signal.NotifyContext` and `server.Shutdown(ctx)` with an active drain delay to achieve 100% zero-downtime rolling deploys.
