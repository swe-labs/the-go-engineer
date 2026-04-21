# Track B: Structured Logging

## Mission

This track teaches how to turn logs into queryable operational data instead of leaving them as strings that only make sense after a production incident.

## Stage Ownership

This track belongs to [10 Production Operations](../README.md).

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `SL.1` | Lesson | [slog basics](./1-slog-basics) | Introduces handlers, levels, attributes, and groups. | entry |
| `SL.2` | Lesson | [context logger](./2-context-logger) | Carries request-scoped fields through a call chain. | `SL.1` |
| `SL.3` | Lesson | [custom handler](./3-custom-handler) | Explains the extension point behind logging backends. | `SL.1`, `SL.2` |
| `SL.4` | Lesson | [zerolog comparison](./4-zerolog-comparison) | Shows when performance pressure justifies a different logger. | `SL.1`, `SL.3` |
| `SL.5` | Exercise | [PII redactor](./5-exercise) | Applies the track by redacting sensitive attributes automatically. | `SL.1`, `SL.2`, `SL.3`, `SL.4` |

## Next Step

After `SL.5`, continue to the [Graceful Shutdown track](../02-graceful-shutdown) or return to the [10 Production Operations overview](../README.md).
