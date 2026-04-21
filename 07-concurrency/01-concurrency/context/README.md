# Track B: Context

## Mission

This track teaches you how Go propagates cancellation, deadlines, and request-scoped metadata
through function chains so I/O work can stop cleanly when the caller is done.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `CT.1` | Lesson | [background](./1-background) | Introduces root contexts and the `Context` interface. | entry |
| `CT.2` | Lesson | [with cancel](./2-with-cancel) | Shows manual cancellation and goroutine cleanup. | `CT.1` |
| `CT.3` | Lesson | [with timeout](./3-with-timeout) | Adds deadline-based cancellation for bounded work. | `CT.1`, `CT.2` |
| `CT.4` | Lesson | [with value](./4-with-value) | Explains request-scoped metadata and the limits of context values. | `CT.1`, `CT.2`, `CT.3` |
| `CT.5` | Exercise | [timeout-aware API client](./5-timeout-client) | Combines timeout control with a real HTTP request boundary. | `CT.1`, `CT.2`, `CT.3`, `CT.4` |

## Suggested Order

1. Work through `CT.1` to `CT.4` in order.
2. Complete `CT.5` as the live context milestone.

## Track Milestone

`CT.5` is the current context track milestone.

If you can complete it and explain:

- why `context.WithTimeout` should wrap outbound HTTP calls
- why `http.NewRequestWithContext` is the real transport boundary
- why context values should carry request metadata instead of general dependencies

then the context part of Stage 07 is doing its job.

## Next Step

After `CT.5`, continue to the [Stage 07 overview](../README.md) or move into the
[Time and Scheduling track](../time-and-scheduling).
