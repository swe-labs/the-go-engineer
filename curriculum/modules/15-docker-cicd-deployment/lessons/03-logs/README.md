# Logs

## Learning objective

Write structured logs to stdout/stderr using Go's `log/slog` package, distinguish between log levels, and follow the 12-factor app logging convention required by Docker and container orchestrators.

## Why this matters

In Docker and Kubernetes, logs are not files — they are streams. Docker captures container stdout and stderr, and log aggregators (Loki, Elastic, Datadog) ingest those streams. If your Go service writes to a file, those logs are invisible to the orchestration layer. If your logs are unstructured, they cannot be searched, filtered, or alerted on. Proper logging is the foundation of observability in production.

## Mental model

Logs are the black box flight recorder for your service. Each log line is a timestamped sensor reading. Structured fields (user_id, duration_ms, error_kind) let you filter for anomalies in a dashboard. Unstructured chatter ("starting server...") is noise. The aggregation pipeline is the investigation board — if your sensors do not speak a common format (JSON), the board is useless. In Docker, stdout = flight data recorder, stderr = alarm bell.

## Core idea

The 12-factor app mandates that applications log to stdout/stderr, not to files. Docker captures these streams and routes them to the configured logging driver (json-file, journald, fluentd, etc.). Go's `log/slog` (structured logging) package outputs JSON by default when using `slog.NewJSONHandler`.

| Stream | Purpose | Docker behavior |
|---|---|---|
| stdout | Normal operational logs | Captured by `docker logs` |
| stderr | Errors and warnings | Captured by `docker logs`, often highlighted |
| File | Never used in containers | Lost on container restart, breaks aggregation |

## Under the hood

`slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})` returns a handler that encodes each log entry as a JSON object. The handler's `Handle()` method is called for every enabled log event; it calls `encoding/json.Encoder.Encode()` which writes to the `io.Writer`. Under the hood, `sync.Pool` reuses encoder buffers to reduce allocations. The `Level` type is an int; `HandlerOptions.Level` is a `Leveler` interface that can be a dynamic `atomic.Int64` for live log level changes.

When `slog.Info()` is called, the logger acquires the handler's mutex, calls `Handle()`, and returns. For async behavior, wrap the handler in a custom writer that sends entries to a channel consumed by a background goroutine.

## How Go uses it

Go 1.21 introduced `log/slog` as the standard structured logging package. Key types:

| Type | Purpose |
|---|---|
| `slog.Logger` | The top-level logger, holds a handler |
| `slog.Handler` | Interface for log output (JSONHandler, TextHandler) |
| `slog.Level` | Severity: Debug, Info, Warn, Error |
| `slog.Attr` | Key-value pair for structured data |
| `slog.Group` | Nested key-value group |

Production patterns: request-scoped loggers with `slog.With` to add trace_id, handler middleware that logs request duration and status, and sampling for high-volume health check endpoints.

## Go example

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger.Debug("starting up", "version", "1.0.0")
	logger.Info("server listening", "addr", ":8080")
	logger.Warn("high memory usage", "pct", 87.5)
	logger.Error("connection refused", "remote", "10.0.0.1:5432", "attempt", 3)

	slog.SetDefault(logger)
	slog.Info("this uses the default logger")
}
```

Output (each line is a JSON object written to stdout):

```json
{"time":"2026-06-25T10:00:00Z","level":"DEBUG","msg":"starting up","version":"1.0.0"}
{"time":"2026-06-25T10:00:00Z","level":"INFO","msg":"server listening","addr":":8080"}
{"time":"2026-06-25T10:00:00Z","level":"WARN","msg":"high memory usage","pct":87.5}
{"time":"2026-06-25T10:00:00Z","level":"ERROR","msg":"connection refused","remote":"10.0.0.1:5432","attempt":3}
```

## Step-by-step execution

1. `slog.NewJSONHandler(os.Stdout, opts)` creates a handler that writes JSON to stdout.
2. `slog.New(handler)` creates a Logger wrapping that handler.
3. `logger.Info("server listening", "addr", ":8080")` calls the handler's `Handle` method.
4. The handler creates a `slog.Record` with the time, level, message, and attrs.
5. The handler acquires its internal mutex and calls `json.Encoder.Encode`.
6. The JSON bytes are written to stdout.
7. Docker's logging driver reads stdout and forwards the line to the configured destination.

## Common mistakes

- **Logging to files instead of stdout.** Writing to `/var/log/app.log` violates the 12-factor principle. In Docker, file logs are invisible to `docker logs`, lost when the container restarts, and accumulate in the container layer.
- **Using `log.Println` for everything.** The standard `log` package has no log levels. Production systems need DEBUG, INFO, WARN, ERROR levels so operators can filter noise.
- **Calling `log.Fatal` in library code.** `log.Fatal` calls `os.Exit(1)`, which prevents deferred cleanup and closes the process without sending HTTP 500 responses. Return errors instead.
- **Logging sensitive data.** Logging request bodies, auth tokens, or SQL queries verbatim creates compliance violations. Use `slog.Any` with a custom `LogValuer` that redacts secrets.
- **Unstructured log messages.** Writing `fmt.Printf("User %s logged in", user)` instead of `slog.Info("user login", "user_id", id, "ip", remoteAddr)`. Unstructured logs cannot be parsed by log aggregators.

## Debugging walkthrough

Consider a service that logs to a file inside a Docker container:

```go
f, _ := os.Create("/var/log/app.log")
log.SetOutput(f)
log.Println("server started")
```

**Symptom:** `docker logs <container>` shows nothing. The log file exists inside the container but is not captured by Docker.

**Investigation:** Check the logging driver and container filesystem. Docker only captures stdout/stderr, not files. Run `docker exec <container> cat /var/log/app.log` to verify the file exists.

**Root cause:** The application writes logs to a file instead of stdout. Docker's logging driver reads from the stdout/stderr file descriptors of PID 1.

**Fix:** Replace file logging with stdout logging:
```go
log.SetOutput(os.Stdout)
log.Println("server started")
```
Or switch to `slog` with a JSON handler on stdout.

## Production notes

- Always set `slog.SetDefault(logger)` so that libraries using `slog.Default()` produce structured logs.
- Use `slog.HandlerOptions{AddSource: true}` in development for file:line information; disable in production for performance.
- For high-volume endpoints (health checks, metrics), consider sampling or suppressing debug logs.
- Use `slog.With("service", "myapp", "version", version)` at startup to add global attrs to every log line.
- Never log secrets, passwords, tokens, or PII. Implement a `redactedString` type that implements `slog.LogValuer` to strip sensitive data.
- Use log rotation awareness: Docker's json-file driver rotates logs automatically. Do not implement rotation in your application.

## Performance implications

- `slog` with JSON handler allocates per-log-call: ~200-500 bytes for a typical entry with 5 attrs.
- Synchronous I/O blocks the calling goroutine. For very high throughput (10k+ log/s), use an async handler that writes to a channel.
- Disabling debug logs at the handler level (`Level: slog.LevelInfo`) skips the handler call entirely — the record is not created.
- `slog` does not format the message string if the level is below the threshold, but the attrs are still evaluated. Use `slog.Debug("msg", "expensive", computeValue())` with caution.

## Practice task

Write a function `NewJSONLogger(w io.Writer, minLevel slog.Level) *slog.Logger` that creates a JSON logger with the given minimum level. Then write a function `LogRequest(logger *slog.Logger, method, path string, status int, duration time.Duration)` that logs a structured request line. Write a `main()` that creates the logger, logs a few requests, and prints the output.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/03-logs
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/03-logs
```

## Review questions

1. Why should Docker containers log to stdout/stderr instead of files?
2. What is the difference between `slog.Info` and `slog.Debug`? How do you control which levels are output?
3. How does `slog.With` help with request-scoped logging?
4. What happens to log output when you call `log.Fatal` inside a Docker container?
5. Why is JSON the preferred log format for containerized applications?

## NEXT UP

File permissions — how Unix file modes work, octal notation, and how Go's `os.FileMode` maps to Linux permissions.
