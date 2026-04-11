# Track B: Structured Logging

## Mission

This track teaches you how to turn logs into queryable operational data instead of leaving them as
strings that only make sense after a production incident has already happened.

## Beta Stage Ownership

This track belongs to [8 Production Engineering](../../docs/stages/08-production-engineering.md).

Within the beta public shell, it is the first live learner path for that stage.
It is where learners build runtime visibility habits before moving into shutdown and deployment
behavior.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `SL.1` | Lesson | [slog basics](./1-slog-basics) | Introduces handlers, levels, attributes, and groups. | entry |
| `SL.2` | Lesson | [context logger](./2-context-logger) | Carries request-scoped fields through a call chain. | `SL.1` |
| `SL.3` | Lesson | [custom handler](./3-custom-handler) | Explains the extension point behind logging backends. | `SL.1`, `SL.2` |
| `SL.4` | Lesson | [zerolog comparison](./4-zerolog-comparison) | Shows when performance pressure justifies a different logger. | `SL.1`, `SL.3` |
| `SL.5` | Exercise | [PII redactor](./5-exercise) | Applies the track by redacting sensitive attributes automatically. | `SL.1`, `SL.2`, `SL.3`, `SL.4` |

## Suggested Order

1. Work through `SL.1` to `SL.4` in order.
2. Complete `SL.5` once you can explain the difference between log data and log formatting.

## Track Milestone

`SL.5` is the current structured-logging output.

If you can explain:

- why `slog` separates records from handlers
- why request-scoped logging depends on context propagation
- why redaction belongs inside the logger pipeline instead of inside every call site

then the structured-logging part of Section 14 is doing its job.

## Next Step

After `SL.5`, continue to the [Graceful Shutdown track](../graceful-shutdown) or back to the
[Section 14 overview](../README.md).

In the beta shell, this keeps you inside
[8 Production Engineering](../../docs/stages/08-production-engineering.md).
