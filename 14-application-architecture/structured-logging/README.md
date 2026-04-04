# Section 23: Structured Logging

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| `slog` basics | Beginner | **Critical** | Key-value pairs, log levels, handlers |
| Context-keyed logger | Intermediate | **Critical** | Request-scoped logging, middleware |
| Custom `slog.Handler` | Advanced | High | Adapting output to any sink |
| `zerolog` comparison | Advanced | Medium | Allocation-free structured logging |

## Engineering Depth

Before `slog` (added in Go 1.21), every team invented their own logger or reached for `zap`, `zerolog`, or `logrus`. The result was three different APIs across every codebase. `slog` ends that.

The key insight: **structured logs are data, not strings**. A log like `"user 42 logged in from 192.168.1.1"` cannot be queried. A structured log with `user_id=42` and `remote_ip=192.168.1.1` can be indexed and filtered in any observability platform (Datadog, Loki, CloudWatch).

**Performance hierarchy:**
- `slog.TextHandler` — human readable, ~300 ns/op, ~3 allocs/op
- `slog.JSONHandler` — machine readable, ~350 ns/op, ~4 allocs/op
- `zerolog` — allocation-free (uses byte buffers + `sync.Pool`), ~80 ns/op, 0 allocs/op

For most services `slog` is sufficient. Reach for `zerolog` only when pprof shows logging in your hot path.

## Contents

| Directory | Topic | Level |
|-----------|-------|-------|
| `1-slog-basics/` | Text/JSON handlers, levels, groups, attrs | Beginner |
| `2-context-logger/` | Context-keyed logger, HTTP middleware | Intermediate |
| `3-custom-handler/` | Implement `slog.Handler` interface | Advanced |
| `4-zerolog-comparison/` | zerolog patterns, when to use it | Advanced |

## How to Run

```bash
go run ./23-structured-logging/1-slog-basics
go run ./23-structured-logging/2-context-logger
go run ./23-structured-logging/3-custom-handler
go run ./23-structured-logging/4-zerolog-comparison
```

## References

- [Go Blog: Structured Logging with slog](https://go.dev/blog/slog)
- [Package slog](https://pkg.go.dev/log/slog)
- [zerolog](https://github.com/rs/zerolog)
