# Opslane Milestone 13 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) `SIGINT` (Ctrl+C) and `SIGTERM` (Kubernetes pod termination signal).**  
*Explanation:* `SIGINT` is sent when a user presses Ctrl+C in a local terminal. `SIGTERM` is the standard termination signal sent by orchestrators (Kubernetes, Docker, systemd) to request graceful application shutdown.

**2. b) No, `SIGKILL` is handled directly by the operating system kernel and immediately terminates the process without invoking application code.**  
*Explanation:* Under POSIX standards, `SIGKILL` (signal 9) and `SIGSTOP` cannot be caught, blocked, or ignored by application processes. They force instant termination by the OS kernel.

**3. b) To instruct client HTTP connections (and load balancers) not to reuse the TCP socket for subsequent keep-alive requests, forcing them to reconnect to a healthy instance.**  
*Explanation:* HTTP/1.1 persistent connections (Keep-Alive) allow clients to send multiple sequential requests over a single TCP socket. Emitting `Connection: close` tells the client or proxy to close the socket after the current response, preventing further requests from hitting the terminating instance.

**4. b) Closes the listener port to refuse new TCP connections, closes idle connections, and waits for active in-flight requests to complete or until `ctx` deadline expires.**  
*Explanation:* `server.Shutdown(ctx)` is Go's standard library method for graceful draining. It stops accepting new traffic immediately and allows active HTTP handler goroutines to finish normally within the context timeout.

**5. a) Because `signal.Notify` does not block when delivering a signal; a buffered channel prevents missing signals if the process is briefly occupied when the signal arrives.**  
*Explanation:* If `signal.Notify` attempts to send a signal to an unbuffered channel and no receiver is actively waiting at that exact microsecond, the signal is dropped. A buffer capacity of at least 1 ensures the signal is safely queued.

**6. b) `http.ErrServerClosed`.**  
*Explanation:* When `Shutdown()` or `Close()` closes the underlying listener, `server.ListenAndServe()` immediately exits with `http.ErrServerClosed` to indicate intentional termination rather than an unexpected failure.

---

## Short Answer (2 points each)

**7. Why abrupt process termination causes HTTP 502 Bad Gateway:**
When a process terminates abruptly with in-flight requests active, the operating system kernel sends a TCP RST (reset) packet on all open sockets. Upstream reverse proxies (NGINX, Envoy, AWS ALB) waiting for responses from these sockets receive an unexpected connection reset, cannot parse a response, and return an HTTP 502 Bad Gateway error to the client. Graceful draining ensures all in-flight responses are delivered before closing sockets.
*Grading: 1 pt for explaining TCP connection reset; 1 pt for explaining reverse proxy 502 translation.*

**8. Kubernetes `iptables` race condition and drain delay:**
- **The Race:** When Kubernetes terminates a pod, two events happen asynchronously: (1) Kubelet sends `SIGTERM` to the container process; (2) Endpoints controller removes the pod's IP from the Service endpoint list and kube-proxy updates node `iptables`. Because `iptables` updates take several seconds to propagate across cluster nodes, external load balancers continue routing new client traffic to the pod *after* `SIGTERM` has been received.
- **Drain Delay Resolution:** If the application begins shutting down immediately upon `SIGTERM`, it rejects traffic that kube-proxy hasn't removed yet. Adding a 3–5 second sleep or failing readiness probes while keeping the server open allows `iptables` to update cluster-wide before the server stops accepting TCP connections.
*Grading: 1 pt for explaining the asynchronous propagation delay between kubelet and kube-proxy; 1 pt for explaining the sleep/drain delay mitigation.*

**9. Bounded `context.WithTimeout` in shutdown:**
- **Purpose:** Active requests may be stuck waiting on external API calls, slow clients, or deadlocked database queries. Passing a bounded context (e.g. 5–10s) guarantees that the shutdown process will terminate within a known time window, rather than hanging indefinitely.
- **Consequence of Unbounded Context:** If an unbounded `context.Background()` was used, a single stalled connection would prevent the container from ever exiting. The orchestrator would eventually lose patience and send an ungraceful `SIGKILL` anyway.
*Grading: 1 pt for explaining hung request mitigation; 1 pt for preventing indefinite hang / avoiding SIGKILL.*

**10. Teardown ordering (HTTP Server -> Worker Pools -> Database Pools):**
1. **HTTP Server (First):** Must be shut down first so that no new incoming user requests are accepted and in-flight HTTP requests complete.
2. **Background Worker Pools (Second):** Must be drained next so that queued background jobs (e.g. notifications, outbox publishing) are executed to completion.
3. **Database Connection Pools (Last):** In-flight HTTP requests and background workers both require database connections to commit transactions. If the database pool was closed first, all in-flight HTTP requests and background workers would fail with "database closed" errors.
*Grading: 1 pt for correct ordering; 1 pt for explaining dependency of active requests/workers on database connections.*
