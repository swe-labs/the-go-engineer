# Opslane Milestone 01 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) Ports below 1024 are privileged root ports on POSIX systems, and ports above 65535 exceed the 16-bit TCP port range.**  
*Explanation:* Ports 1–1023 require superuser (root) privileges to bind on Linux/macOS. TCP port numbers are 16-bit unsigned integers ($2^{16}-1 = 65535$).

**2. b) Intercepting unhandled panics inside HTTP handler goroutines to prevent server crash, logging error details with `slog`, and returning an HTTP 500 JSON response.**  
*Explanation:* In Go's `net/http`, an unrecovered panic inside a handler drops the TCP connection. Middleware using `defer recover()` catches the panic, records structured telemetry, and gracefully sends an internal server error response to the client.

**3. b) It uses `slog.NewJSONHandler` for machine-parseable log aggregation in production, and human-readable `slog.NewTextHandler` in development.**  
*Explanation:* JSON formatted logs allow daemon collectors (FluentBit, Vector, Datadog) to parse keys into indexed fields without CPU-heavy regex scanning. Text handlers prioritize developer readability in local terminal sessions.

**4. b) To prevent a catastrophic security vulnerability where production instances run with publicly known credentials from the source repository.**  
*Explanation:* Hardcoded repository secrets are exposed to anyone with code access. If a production system booted with a default secret, attackers could forge valid administrative JWT tokens trivially.

**5. b) `ReadTimeout` and `WriteTimeout`.**  
*Explanation:* Slowloris attacks hold HTTP connections open indefinitely by sending headers very slowly. Setting `ReadTimeout` forces the server to terminate connections that do not deliver complete request headers within the specified window.

**6. a) `/healthz` returns uptime, service name, and status; `/readyz` indicates readiness state and worker pool capacity.**  
*Explanation:* Liveness probes (`/healthz`) tell the orchestrator whether the process is alive; readiness probes (`/readyz`) indicate whether the instance is prepared to handle inbound requests.

---

## Short Answer (2 points each)

**7. Architectural advantages of startup configuration parsing:**
1. **Fail-Fast Reliability:** If a required environment variable is missing or malformed, the process aborts immediately upon launch with a descriptive error before registering with load balancers, rather than failing on an arbitrary client request hours later.
2. **Immutable Runtime State:** Storing configuration in an immutable struct avoids repeated OS system calls (`os.Getenv`), prevents race conditions from concurrent environment mutations, and enables seamless dependency injection into components.
*Grading: 1 pt for fail-fast explanation; 1 pt for immutability/dependency injection benefits.*

**8. Runtime behavior on unhandled panics:**
In Go's `net/http` standard library server, each inbound HTTP connection runs in its own goroutine. The standard server wraps handler execution in an internal recovery mechanism that terminates the connection and logs the panic to `Server.ErrorLog`. However, without custom application recovery middleware, the client receives a truncated TCP reset (or empty 500 without a clean JSON body), contextual attributes (`trace_id`, request path) are lost, and application metrics are not incremented.
*Grading: 1 pt for describing goroutine termination / TCP reset; 1 pt for noting the loss of structured telemetry and clean JSON response.*

**9. Structured logging in Kubernetes:**
`slog` outputs standardized JSON objects containing timestamp, severity level, message, and explicit typed attributes (`"port": 8080`, `"env": "production"`). When container runtime engines (containerd) ship stdout to centralized aggregators (Elasticsearch, Loki, CloudWatch), the JSON payload is automatically parsed into indexed fields. Engineers can filter logs instantly (`env="production" AND port=8080`) rather than running expensive full-text regex scans on unstructured text.
*Grading: 1 pt for explaining JSON attribute indexing vs raw text; 1 pt for mentioning log aggregator ingestion.*

**10. Handling `http.ErrServerClosed`:**
When `server.Shutdown(ctx)` or `server.Close()` is called to initiate graceful shutdown, the server stops accepting new connections and immediately causes `server.ListenAndServe()` to return `http.ErrServerClosed`. This error signals that the server shut down as commanded, rather than failing due to an unexpected network or OS error. Ignoring it prevents `main()` from incorrectly treating a successful shutdown as a fatal crash.
*Grading: 1 pt for identifying it as the return value of graceful shutdown; 1 pt for explaining why it should not be logged as a crash.*
