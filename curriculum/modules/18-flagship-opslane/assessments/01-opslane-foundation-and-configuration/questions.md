# Opslane Milestone 01 Assessment — Foundation and Configuration: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** Why does `LoadConfig` validate that `PORT` falls between 1024 and 65535?
- a) Ports below 1024 are reserved for UDP only.
- b) Ports below 1024 are privileged root ports on POSIX systems, and ports above 65535 exceed the 16-bit TCP port range.
- c) Go HTTP routers cannot bind to 8-bit ports.
- d) Kubernetes will refuse to deploy any container exposing a port outside this range.

**2.** What is the primary purpose of `recoverMiddleware` in `Application.Routes()`?
- a) Automatically retrying failed HTTP database transactions.
- b) Intercepting unhandled panics inside HTTP handler goroutines to prevent server crash, logging error details with `slog`, and returning an HTTP 500 JSON response.
- c) Rerouting 404 requests to a static fallback page.
- d) Garbage collecting dead request contexts.

**3.** How does `InitLogger` tailor its output based on the environment?
- a) It writes to `/var/log/opslane.log` in production and `os.Stdout` in dev.
- b) It uses `slog.NewJSONHandler` for machine-parseable log aggregation in production, and human-readable `slog.NewTextHandler` in development.
- c) It disables debug logging completely in all environments.
- d) It encrypts log payloads with AES-256 before outputting to stdout.

**4.** Why does `LoadConfig` reject the default JWT secret when `ENV=production`?
- a) Because Go's crypto package prevents compiling hardcoded strings.
- b) To prevent a catastrophic security vulnerability where production instances run with publicly known credentials from the source repository.
- c) Because JWT tokens expire faster in production.
- d) To satisfy database migration constraints.

**5.** Which standard library HTTP server settings protect against Slowloris denial-of-service attacks?
- a) `MaxHeaderBytes` only.
- b) `ReadTimeout` and `WriteTimeout`.
- c) `DisableKeepAlives`.
- d) `TLSConfig.MinVersion`.

**6.** What information is exposed by `/healthz` versus `/readyz` in this milestone?
- a) `/healthz` returns uptime, service name, and status; `/readyz` indicates readiness state and worker pool capacity.
- b) `/healthz` tests database credentials; `/readyz` tests internet connectivity.
- c) Both endpoints execute identical internal logic and return duplicate payloads.
- d) `/healthz` is an unauthenticated endpoint; `/readyz` requires a bearer token.

---

## Short Answer (2 points each)

**7.** Explain the architectural advantages of loading and validating all configuration parameters at startup rather than calling `os.Getenv` on demand inside HTTP handlers.

**8.** What happens if a panic occurs in an HTTP handler goroutine when `recoverMiddleware` is omitted? Explain the Go runtime behavior.

**9.** How does structured logging with `log/slog` improve observability in Kubernetes compared to unstructured `log.Printf`?

**10.** In `main.go`, `http.ErrServerClosed` is explicitly ignored when evaluating `server.ListenAndServe()`. Why is this error expected during normal application termination?
