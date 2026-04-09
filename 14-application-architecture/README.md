# Section 14: Application Architecture

## Mission

This section teaches you how Go services become maintainable once the codebase stops fitting in a
single file or a single process concern.

The live v2 slice focuses on three architecture habits:

- designing packages with clear boundaries and small public surfaces
- logging with structured data instead of ad hoc strings
- shutting down services cleanly when deploys and signals hit production

Section 14 still contains additional reference material beyond this live slice. Those surfaces
remain available while the first public architecture path is migrated.

## Who Should Start Here

### Full Path

Start here after finishing Section 13 in order.

### Bridge Path

You can move faster if you already understand:

- interfaces, package boundaries, and dependency injection
- HTTP handlers and middleware basics
- context cancellation and long-running process behavior

Even on the bridge path, do not skip `PD.1`, `SL.1`, or `GS.1`.
They anchor the three architecture tracks in this slice.

### Targeted Path

This section is a multi-track v2 slice.
Choose the track that matches your immediate need:

- Package design for layout, visibility, and dependency boundaries
- Structured logging for production-friendly observability
- Graceful shutdown for reliable deploys and process lifecycle handling

## Section Map

| Track | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| Package Design | [PD.1 naming](./package-design) | `PD.3` | package naming, export rules, and project layout |
| Structured Logging | [SL.1 slog basics](./structured-logging) | `SL.5` | slog, request-scoped logging, handlers, and redaction |
| Graceful Shutdown | [GS.1 signal context](./graceful-shutdown) | `GS.3` | signals, HTTP draining, and coordinated shutdown |

## Suggested Order

1. Work through Package Design first so the section starts with boundary discipline.
2. Continue into Structured Logging once the package surface feels clear.
3. Finish with Graceful Shutdown so the section ends on service lifecycle and production behavior.

## Section Milestones

This live v2 slice has three promoted outputs:

- `PD.3` project layout
- `SL.5` PII redactor exercise
- `GS.3` shutdown capstone

The following surfaces remain available as legacy reference material for later alpha work:

- `grpc`
- `docker-and-deployment`
- `enterprise-capstone`

If you can complete the live slice and explain:

- why good package design is about boundary clarity, not folder decoration
- why structured logs are data contracts, not prettier strings
- why graceful shutdown is part of correctness, not deployment polish

then you are ready to move into code generation in Section 15.

## Pilot Role In V2

This live v2 slice keeps the current `14-application-architecture` layout intact while promoting
the main public architecture path:

- `PD.1` through `PD.3` form the live package-design track
- `SL.1` through `SL.5` form the live structured-logging track
- `GS.1` through `GS.3` form the live graceful-shutdown track
- `grpc`, `docker-and-deployment`, and `enterprise-capstone` remain legacy reference surfaces

## References

1. [Effective Go: Package Names](https://go.dev/doc/effective_go#package-names)
2. [Package slog](https://pkg.go.dev/log/slog)
3. [signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext)

## Next Step

After you finish the track or milestone you care about here, continue to
[Section 15: Code Generation](../15-code-generation).
