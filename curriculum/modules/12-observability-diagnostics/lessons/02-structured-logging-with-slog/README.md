# Structured logging with slog

## Learning objective

Use Go's `log/slog` package to write structured log entries with levels, JSON and text handlers, and typed attributes, and select the appropriate handler and log level for a given production scenario.

## Why this matters

Unstructured log lines (`"user 42 logged in"`) cannot be searched, aggregated, or alerted on. A log aggregator like Grafana Loki or Datadog needs key-value pairs (`user_id=42 action=login`) to build dashboards and alerts. If you log `"user 42 logged in"`, the aggregator cannot compute "how many unique users logged in per hour" or "alert when login errors > 10 per minute". Structured logging makes your log data machine-parseable and queryable at scale. In Go 1.21+, `log/slog` is the standard library solution for structured logging.

## Mental model

Structured logging is the difference between a sticky note and a spreadsheet row. A sticky note says "error occurred". A spreadsheet row says:

| time | level | msg | query_time_ms | path |
|---|---|---|---|---|
| 12:00:00 | ERROR | db query failed | 5200 | /api/users |

The sticky note is human-readable in isolation but useless for search, aggregation, and alerting. The spreadsheet row is machine-parseable: a log aggregator can sum `query_time_ms > 1000`, alert on ERROR level, or filter by path.

slog's handler API is an assembly line: each log call creates a `Record` (item on the line), passes it through the handler (worker that serializes it), and writes it to output (shipping dock). The handler decides the output format (text or JSON) and the log level determines whether the record reaches the handler at all.

## Core idea

A structured log entry is a key-value mapping, not a flat string. Each attribute has a typed value (string, int, bool, duration) that the handler serializes according to its format. The schema for every log line is known at compile time (the keys are string literals) and the values are typed at runtime.

slog provides four built-in levels with increasing severity: `Debug`, `Info`, `Warn`, `Error`. Each level has two calling conventions: convenience methods (`logger.Info(msg, args...)`) and the generic `logger.LogAttrs(ctx, level, msg, attrs...)`. The convenience methods accept alternating key-value pairs; `LogAttrs` accepts typed `slog.Attr` values and avoids reflection overhead.

## Under the hood

`slog.Handler` is an interface with three methods:

- `Enabled(ctx, level) bool` — returns whether this handler would accept a record at the given level.
- `Handle(ctx, record) error` — serializes and writes the record.
- `WithAttrs(attrs []Attr) Handler` — returns a new handler with pre-attached attributes.

The default text handler (`slog.NewTextHandler`) writes `key=value` pairs separated by spaces. The JSON handler (`slog.NewJSONHandler`) writes a JSON object per line. Both accept a `HandlerOptions` struct that controls the minimum level (`Level`), whether to add source file/line (`AddSource`), and a custom level-to-string function (`Leveler`).

When you call `logger.Info("msg", "key1", val1, "key2", val2)`, the internal path is:

1. `Logger.Info` calls `Logger.log(nil, LevelInfo, "msg", args...)`.
2. `Logger.log` verifies the handler is enabled for LevelInfo.
3. It creates a `slog.Record` with the current time, LevelInfo, message, and caller PC.
4. It converts the variadic args to attrs by pairing consecutive args as key-value. If an odd number of args is provided, the last key gets a `nil` value.
5. The record is passed to the handler's `Handle` method.
6. The handler serializes the record and writes to the output writer.

`slog.Record` is optimized for zero-allocation in the common case: it stores up to 10 attrs inline (in a ring buffer) and only allocates when exceeded.

## How Go uses it

slog replaces `log.Printf` and `log.Println` in Go 1.21+ production code. Standard patterns:

```go
// Directly via the default logger
slog.Info("server starting", "port", 8080, "env", "production")

// With a custom JSON handler
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
logger.Warn("high memory usage", "memory_mb", 4096, "threshold_mb", 2048)

// With typed attributes (avoids reflection)
logger.LogAttrs(context.Background(), slog.LevelError, "db timeout",
    slog.String("query", "SELECT * FROM orders"),
    slog.Int("timeout_ms", 5000),
)

// With groups for logical organization
logger.Info("request completed",
    slog.Group("http",
        slog.String("method", "POST"),
        slog.String("path", "/api/orders"),
        slog.Int("status", 201),
    ),
)
```

## Go example

```go
package main

import (
	"log/slog"
	"os"
	"time"
)

func LogRequest(logger *slog.Logger, method, path string, status int, dur time.Duration) {
	logger.Info("http request",
		"method", method,
		"path", path,
		"status", status,
		"duration_ms", dur.Milliseconds(),
	)
}

func main() {
	textLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	LogRequest(textLogger, "GET", "/api/users", 200, 42*time.Millisecond)
	LogRequest(jsonLogger, "POST", "/api/orders", 500, 5200*time.Millisecond)
}
```

The text handler output looks like:
```
time=2026-01-15T10:30:00.000Z level=INFO msg="http request" method=GET path=/api/users status=200 duration_ms=42
```

The JSON handler output looks like:
```json
{"time":"2026-01-15T10:30:00.000Z","level":"INFO","msg":"http request","method":"POST","path":"/api/orders","status":500,"duration_ms":5200}
```

## Step-by-step execution

For `slog.Info("user login", "user_id", "u-42", "ip", "10.0.0.1")` with a JSON handler:

1. `slog.Info` is a convenience function that calls `slog.Default().Info(...)`.
2. The default logger's handler checks `Enabled` with LevelInfo → true.
3. A `slog.Record` is created: `time=now, level=INFO, msg="user login"`.
4. The variadic args `"user_id", "u-42", "ip", "10.0.0.1"` are split into key-value pairs and added to the record as attrs.
5. The handler serializes the record as a JSON object: `{"time":"...","level":"INFO","msg":"user login","user_id":"u-42","ip":"10.0.0.1"}` followed by a newline.
6. The JSON bytes are written to the handler's `io.Writer` (os.Stderr by default).
7. The function returns. The record is garbage collected.

## Common mistakes

- **Using slog for unstructured messages**: `slog.Info(fmt.Sprintf("user %d logged in", id))` puts the entire message in the `msg` field, making it opaque to structured queries. Always pass structured data as separate attributes: `slog.Info("user logged in", "user_id", id)`.

- **Mismatching key-value pairs**: `slog.Info("process", "key1", val1, "key2")` — an odd number of args causes the last key to receive a nil value. Go does not warn about this at compile time. Always use `LogAttrs` with typed `slog.Attr` when the number of args is dynamic.

- **Using slog.Default() in library code**: Libraries should accept a `*slog.Logger` parameter (or use context-based retrieval) instead of calling `slog.Default()`. If the application sets a custom handler, the library should use it too.

- **Ignoring the handler API**: Instead of implementing `slog.Handler` for custom output, developers often parse and re-format serialized log output with regex — which breaks when the output format changes. Implementing `slog.Handler` is ~50 lines and gives full control over serialization.

- **Setting the wrong minimum level in production**: `slog.LevelDebug` in production can generate millions of log lines per second from debug-level statements, overwhelming the log aggregator and increasing costs. Set the minimum level to `slog.LevelInfo` in production and only enable Debug for specific services during debugging.

## Debugging walkthrough

Consider this code that misses a key-value pair:

```go
slog.Info("creating order",
    "user_id", userID,
    "product_id", productID,
    "quantity", // missing value!
)
```

**Symptom**: The log aggregator shows `quantity=!BADKEY` for every log line.

**Root cause**: The `quantity` key has no corresponding value. slog treats the missing value as `nil`, which the handler serializes as `!BADKEY` (text handler) or `null` (JSON handler).

**Fix**: Use `LogAttrs` with typed attributes:

```go
slog.LogAttrs(context.Background(), slog.LevelInfo, "creating order",
    slog.String("user_id", userID),
    slog.String("product_id", productID),
    slog.Int("quantity", quantity),
)
```

Using `LogAttrs` makes every attribute explicit and type-safe. The compiler cannot catch missing values in variadic pairs, but `LogAttrs` eliminates the ambiguity entirely. The same pattern applies to all four log levels.

## Production notes

slog is the standard structured logging library in Go 1.21+. Kubernetes uses JSON logging natively. Datadog, Grafana Loki, and AWS CloudWatch all ingest structured JSON logs.

In production Go services, slog replaces `log.Println` with `log/slog` for every log line. The structured format enables:

- **Real-time dashboards**: error rate by endpoint, latency p50/p95/p99 by handler.
- **Historical analysis**: latency trend over the last 30 days, correlated with deploys.
- **Automated alerts**: error count > threshold for 5 minutes, or p99 latency > SLO for 10 minutes.

Production logging best practices:

- Use JSON handler in production (text is for development).
- Set `HandlerOptions.Level` to `LevelInfo` in production.
- Never log secrets, passwords, API keys, or PII (see lesson 05).
- Use structured attributes for queryable fields, not the message string.

## Performance implications

The cost of structured logging depends on the number of attributes and the handler:

- **Text handler**: ~0.5μs per attribute for serialization. 10 attributes → ~5μs.
- **JSON handler**: ~1μs per attribute (JSON encoding is more expensive than text key=value).
- **LogAttrs vs variadic**: `LogAttrs` with typed `slog.Attr` avoids the reflection step that variadic args require. For hot-path logging (inside a tight loop), use `LogAttrs` for ~30% lower allocation overhead.
- **Disabled levels**: If the handler's `Level` is higher than the log call's level, `Enabled` returns false and the cost is a single function call (~5ns). There is no serialization overhead.

For the critical path of a request handler (1–2 log calls with 5–10 attributes), the total cost is ~5–10μs. For a handler with a 50ms total latency, this is 0.01–0.02% overhead — negligible.

## Practice task

Complete the function `LogRequest` that writes a structured HTTP request log entry:

```go
func LogRequest(logger *slog.Logger, method, path string, status int, dur time.Duration) {
	// Log a structured entry with method, path, status, and duration_ms
}
```

The tests verify both text and JSON handler output, and also verify `LogAtLevel` with all four levels.

## Tests / verification

```bash
go run ./curriculum/modules/12-observability-diagnostics/lessons/02-structured-logging-with-slog
go test ./curriculum/modules/12-observability-diagnostics/lessons/02-structured-logging-with-slog
```

## Review questions

1. What is the difference between `slog.Info("msg", "key", val)` and `slog.LogAttrs(nil, slog.LevelInfo, "msg", slog.String("key", val))`?
2. When would you choose a JSON handler over a text handler? When would you choose text?
3. What happens if you pass an odd number of key-value pairs to `slog.Info`? How does the handler render the missing value?
4. Why should library code accept a `*slog.Logger` parameter instead of calling `slog.Default()`?
5. Given `slog.Info("query executed", "duration_ms", 2500)`, what attribute would a log aggregator filter on to find slow queries? How would the same query look with `fmt.Printf`?

## NEXT UP

Request-scoped logging — attaching per-request context to log entries so every log line from a single request carries the same request ID.
