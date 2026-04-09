# Track C: Graceful Shutdown

## Mission

This track teaches you how to keep services correct during deploys, restarts, and termination
signals instead of treating shutdown as an afterthought.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `GS.1` | Lesson | [signal context](./1-signal-context) | Introduces `signal.NotifyContext` and signal-aware workers. | entry |
| `GS.2` | Lesson | [HTTP graceful drain](./2-http-server) | Uses `http.Server.Shutdown` to drain in-flight requests. | `GS.1` |
| `GS.3` | Capstone | [shutdown capstone](./3-capstone) | Wires signals, readiness, workers, and drain order together. | `GS.1`, `GS.2` |

## Suggested Order

1. Start with `GS.1` to understand signal delivery and cancellation.
2. Move into `GS.2` to see how request draining works for HTTP services.
3. Finish with `GS.3` once you can reason about shutdown order across multiple resources.

## Track Milestone

`GS.3` is the current graceful-shutdown output.

If you can explain:

- why `signal.NotifyContext` is the right entry point for modern Go shutdown handling
- why `http.Server.Shutdown` is not the same as killing the server
- why readiness, HTTP drain, worker drain, and resource close order must be coordinated

then the graceful-shutdown part of Section 14 is doing its job.

## Next Step

After `GS.3`, continue back to the [Section 14 overview](../README.md) or move to
[Section 15: Code Generation](../../15-code-generation).
