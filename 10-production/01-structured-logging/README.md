# Track SL: Structured Logging

## Mission

Build a telemetry foundation for your applications. Learn how to move beyond simple string-based logging to **Structured Logging** with Go's standard library `log/slog`. This track covers key-value logging, log levels, context propagation, and building custom handlers for production-grade observability.

## Stage Ownership

This track belongs to [10 Production Operations](../README.md).

## Track Map

| ID | Type | Surface | Mission | Requires |
| --- | --- | --- | --- | --- |
| `SL.1` | Lesson | [slog Basics](./1-slog-basics) | Learn the core `log/slog` API and key-value pairs. | entry |
| `SL.2` | Lesson | [Context Logging](./2-context-logger) | Propagate Request IDs and attributes via `context`. | `SL.1` |
| `SL.3` | Lesson | [Custom Handlers](./3-custom-handler) | Build a handler for Slack, ELK, or specialized JSON. | `SL.2` |
| `SL.4` | Lesson | [Zerolog Comparison](./4-zerolog-comparison) | Understand when to use third-party libraries. | `SL.3` |
| `SL.5` | Exercise | [Telemetry Integration](./5-exercise) | Build a unified logging strategy for a service. | `SL.1-4` |

## Why This Track Matters

In production, you cannot "Step through" code with a debugger. Logs are your primary eyes and ears.

1. **Searchability**: JSON logs allow tools like ELK, Datadog, or CloudWatch to filter and aggregate your logs instantly.
2. **Context**: Structured logs allow you to attach a `RequestID` to every log line, making it possible to trace a single user's journey through multiple services.
3. **Performance**: `slog` is designed to be high-performance, ensuring that observability doesn't slow down your application.

## Next Step

After mastering structured logging, learn how to safely shut down your application without losing data. Continue to [Track GS: Graceful Shutdown](../02-graceful-shutdown).
