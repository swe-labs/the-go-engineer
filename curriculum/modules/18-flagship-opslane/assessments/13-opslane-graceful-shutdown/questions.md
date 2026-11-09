# Opslane Milestone 13 Assessment — Graceful Shutdown: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** Which operating system signals are typically intercepted to initiate graceful shutdown in Go services?
- a) `SIGKILL` and `SIGSTOP`.
- b) `SIGINT` (Ctrl+C) and `SIGTERM` (Kubernetes pod termination signal).
- c) `SIGHUP` only.
- d) `SIGUSR1`.

**2.** Can an application catch or intercept a `SIGKILL` signal?
- a) Yes, using `signal.Notify(c, syscall.SIGKILL)`.
- b) No, `SIGKILL` is handled directly by the operating system kernel and immediately terminates the process without invoking application code.
- c) Only on Windows.
- d) Only if running as the root user.

**3.** What is the purpose of setting `w.Header().Set("Connection", "close")` when responding with HTTP 503 during the shutdown sequence?
- a) To encrypt the response payload.
- b) To instruct client HTTP connections (and load balancers) not to reuse the TCP socket for subsequent keep-alive requests, forcing them to reconnect to a healthy instance.
- c) To reduce database connection pool usage.
- d) To flush browser cookies.

**4.** What does `server.Shutdown(ctx)` do when invoked?
- a) Immediately kills all running goroutines and exits with code 1.
- b) Closes the listener port to refuse new TCP connections, closes idle connections, and waits for active in-flight requests to complete or until `ctx` deadline expires.
- c) Reboots the operating system.
- d) Deletes the binary from disk.

**5.** Why is `stopChan` created with buffer capacity 1 (`make(chan os.Signal, 1)`)?
- a) Because `signal.Notify` does not block when delivering a signal; a buffered channel prevents missing signals if the process is briefly occupied when the signal arrives.
- b) Channels in Go cannot be unbuffered.
- c) To limit the server to receiving exactly one HTTP request.
- d) To conserve memory.

**6.** What error does `server.ListenAndServe()` return when `server.Shutdown()` is called?
- a) `io.EOF`.
- b) `http.ErrServerClosed`.
- c) `context.Canceled`.
- d) `syscall.ECONNREFUSED`.

---

## Short Answer (2 points each)

**7.** Explain why immediate process termination (`os.Exit(0)`) during a Kubernetes rolling update causes intermittent HTTP 502 Bad Gateway errors for clients.

**8.** What is the network race condition between Kubernetes `iptables` endpoint removal and pod `SIGTERM` signal delivery? How does adding a 3–5 second sleep/drain phase resolve it?

**9.** Explain the purpose of passing a bounded `context.WithTimeout` (e.g. 5–10 seconds) into `coordinator.GracefulShutdown(ctx)`. What would happen if an unbounded context was passed?

**10.** Explain the recommended shutdown order across HTTP servers, background worker pools, and database connection pools. Why must database pools be closed *last*?
