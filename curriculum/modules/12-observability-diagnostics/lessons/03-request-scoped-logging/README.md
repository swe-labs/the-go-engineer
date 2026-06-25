# Request-scoped logging

## Learning objective

Attach per-request fields (request ID, user ID, trace ID) to every log line within a request's scope by storing a derived logger in `context.Context` and retrieving it in handlers and downstream functions.

## Why this matters

Without request-scoped logging, log lines from 1000 concurrent requests are interleaved in a single stream with no way to group them by request. When investigating a failed order for `user-42`, you search for `user-42` in the log aggregator and get 10,000 results across 1000 different requests — you cannot tell which log lines belong to the failed order. Request-scoped logging stamps every log line with a unique request ID, making it possible to extract a complete timeline for any single request by searching for that ID.

## Mental model

Request-scoped logging is like putting a unique label on every file in a project folder. Each request gets a folder (context) with a unique label (request_id). Every log line from that request goes into the folder with the label automatically stamped on it. After the request completes, you search for the label and find every log line in chronological order — a complete trace of what happened during that request.

The flow is: middleware extracts/injects a request ID → stores a derived logger (with request ID attached) in context → every handler and downstream function retrieves the logger from context → all log lines from that request include the request ID.

## Core idea

Go's `context.Context` is the carrier for request-scoped data. Store a `*slog.Logger` in context (derived from the base logger with request-specific fields via `With`). Any function that receives the context can retrieve the logger and log with the request-scoped attributes automatically included.

`slog.With(key, val)` returns a new logger that prepends the given attributes to every subsequent log call. When the derived logger is stored in context, every log call using that logger includes the request-scoped fields without the caller explicitly passing them.

## Under the hood

`slog.Logger.With` returns a new `*slog.Logger` that wraps the original handler in a new handler returned by `handler.WithAttrs(attrs)`. The new handler stores the attributes and prepends them to every `Record` before calling `Handle`. This means:

- The stored attrs are evaluated once when `With` is called, not per-log-call.
- The attrs share a single attribute slice across all log calls using that logger, reducing allocation.
- Log calls using the derived logger do not need to specify the request-scoped fields — they are automatically injected by the handler.

When storing the logger in context, use an unexported context key type to prevent collisions between packages. The standard pattern is:

```go
type contextKey string
const loggerKey contextKey = "logger"

func AttachLogger(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
        return logger
    }
    return slog.Default() // fallback
}
```

## How Go uses it

slog in Go 1.21+ introduced `slog.FromContext(ctx)` and `slog.WithContext(ctx, logger)` built into the standard library. These use the same pattern internally and are the preferred API for new code.

Middleware pattern with slog context:

```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-Id")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        logger := slog.With("request_id", requestID)
        ctx := slog.WithContext(r.Context(), logger)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

Any handler downstream can call `slog.FromContext(ctx).Info("processing")` and the `request_id` attribute is automatically included.

## Go example

```go
package main

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const loggerKey contextKey = "logger"

func AttachLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func LogWithCtx(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	LoggerFromContext(ctx).LogAttrs(ctx, level, msg, attrs...)
}

func HandleRequest(ctx context.Context, requestID string) {
	requestLogger := LoggerFromContext(ctx).With("request_id", requestID)
	ctx = AttachLogger(ctx, requestLogger)
	LogWithCtx(ctx, slog.LevelInfo, "processing request")
	processInner(ctx)
}

func processInner(ctx context.Context) {
	LogWithCtx(ctx, slog.LevelInfo, "inner processing step")
}

func main() {
	baseLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := AttachLogger(context.Background(), baseLogger)
	HandleRequest(ctx, "req-001")
	HandleRequest(ctx, "req-002")
}
```

Run with `go run .` to see that every log line is tagged with the correct `request_id=req-001` or `request_id=req-002`, and that both requests use the same base logger but produce separate log lines with distinct IDs.

## Step-by-step execution

For a single call `HandleRequest(ctx, "req-001")`:

1. `LoggerFromContext(ctx)` retrieves the base logger from context. If none exists, it returns `slog.Default()`.
2. `baseLogger.With("request_id", "req-001")` creates a new logger whose handler prepends `request_id=req-001` to every record.
3. `AttachLogger(ctx, requestLogger)` stores the derived logger in a new context (with `loggerKey`).
4. `LogWithCtx(ctx, slog.LevelInfo, "processing request")` retrieves the logger from context and calls `logger.LogAttrs(ctx, LevelInfo, "processing request")`.
5. The handler serializes the record: `time=... level=INFO msg="processing request" request_id=req-001`.
6. `processInner(ctx)` receives the context with the derived logger and logs `"inner processing step"`, which also includes `request_id=req-001`.
7. The function returns. The derived logger and its context go out of scope.

## Common mistakes

- **Creating the logger before the middleware**: If the logger is created at package initialization time (`var logger = slog.With("request_id", "...")`) instead of per-request inside middleware, every request shares the same logger with the same request_id. Always derive the logger inside the middleware, not before it.

- **Using context.Background() instead of request context**: If you create a derived logger but attach it to `context.Background()` instead of the request's `context.Context`, the request-scoped fields are lost. Always use the request's context from `r.Context()`.

- **Forgetting to propagate context**: If a downstream function creates a new context with `context.Background()`, it loses the stored logger. Always pass the incoming context through the call chain.

- **Storing a *slog.Logger in a global variable instead of context**: A global logger cannot have per-request fields because concurrent requests overwrite each other's values. Use context to carry the per-request logger.

- **Not providing a fallback**: If `LoggerFromContext` returns nil when no logger is stored, calling methods on nil causes a panic. Always return `slog.Default()` as a fallback.

## Debugging walkthrough

Consider this code that fails to propagate the request-scoped logger:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    logger := slog.With("request_id", r.Header.Get("X-Request-Id"))
    ctx := slog.WithContext(r.Context(), logger)
    processData(ctx)
}

func processData(ctx context.Context) {
    // BUG: ignoring ctx and using background context
    slog.Info("processing data") // no request_id!
}
```

**Symptom**: Log aggregator shows `processing data` without a `request_id` field. Cannot correlate this log line to any specific request.

**Root cause**: `processData` ignores the context parameter and uses the default logger via `slog.Info`, which has no request-scoped attributes.

**Fix**: Use `slog.FromContext(ctx)` to retrieve the request-scoped logger:

```go
func processData(ctx context.Context) {
    slog.FromContext(ctx).Info("processing data")
}
```

Or use the context-aware convenience functions introduced in Go 1.21:

```go
func processData(ctx context.Context) {
    slog.InfoContext(ctx, "processing data")
}
```

`slog.InfoContext` internally calls `slog.FromContext` and includes the stored attributes. This eliminates the error of using the wrong logger.

## Production notes

Request-scoped logging is universal in production Go services. Every HTTP framework (chi, gin, echo) provides middleware for request-scoped loggers. The `request_id` is the primary dimension for grouping log lines into a single request view in observability platforms.

Production best practices:

- Use `slog.WithContext` and `slog.FromContext` (Go 1.21+) instead of a custom context key.
- Use `slog.InfoContext(ctx, ...)` and friends instead of retrieving the logger manually.
- Add the request ID, trace ID, user ID (if authenticated), and handler name as request-scoped fields.
- Never store mutable data in context — the logger is immutable and derived loggers are cheap to create.

Without request-scoped logging, each incident investigation starts with "find the request_id for that failed request" — and you cannot find it if it was never logged. With request-scoped logging, the first step is "search for the request_id from the user's error report" and every relevant log line appears immediately.

## Performance implications

`slog.Logger.With` creates a new logger that wraps the handler. The cost is negligible:

- `With` calls `handler.WithAttrs(attrs)` which returns a new handler with the attrs stored.
- No serialization happens at `With` time — attrs are serialized only when a log record passes through the handler.
- The derived logger shares the underlying handler's writer, mutex, and configuration. Only the attrs slice is new.

For middleware that does this per-request, the allocation is one handler wrapper + one attrs slice, totaling ~200 bytes per request. For a service handling 1000 requests/second, this is ~200KB/second — negligible.

The per-log-call cost is also tiny: the handler prepends the stored attrs to the record's attrs, which is a memcpy of a few pointers.

## Practice task

Implement `AttachLogger`, `LoggerFromContext`, and `LogWithCtx` functions described in the example, then write a handler chain that uses them:

1. `AttachLogger(ctx, logger)` stores a logger in context.
2. `LoggerFromContext(ctx)` retrieves it (or returns `slog.Default()`).
3. `LogWithCtx(ctx, level, msg, attrs...)` logs using the context logger.
4. Use `HandleRequest(ctx, requestID)` to demonstrate a two-step handler with request-scoped fields.

The tests verify that log output contains `request_id` and the correct ID value.

## Tests / verification

```bash
go run ./curriculum/modules/12-observability-diagnostics/lessons/03-request-scoped-logging
go test ./curriculum/modules/12-observability-diagnostics/lessons/03-request-scoped-logging
```

## Review questions

1. Why must request-scoped loggers be stored in `context.Context` rather than in a global variable?
2. What does `slog.Logger.With("request_id", id)` do at the handler level? When are the stored attributes evaluated?
3. What happens if a downstream function uses `context.Background()` instead of the incoming request context? How would you detect this in logs?
4. Name three fields that should be added as request-scoped attributes in production middleware.
5. Why is `slog.InfoContext(ctx, "msg")` safer than retrieving the logger manually and calling `logger.Info`?

## NEXT UP

Correlation IDs — propagating a single identifier across service boundaries so a distributed request can be traced end-to-end.
