# Track C: Graceful Shutdown

## Mission

This track teaches how to keep services correct during deploys, restarts, and termination signals instead of treating shutdown as an afterthought.

## Stage Ownership

This track belongs to [10 Production Operations](../README.md).

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `GS.1` | Lesson | [signal context](./1-signal-context) | Introduces `signal.NotifyContext` and signal-aware workers. | entry |
| `GS.2` | Lesson | [HTTP graceful drain](./2-http-server) | Uses `http.Server.Shutdown` to drain in-flight requests. | `GS.1` |
| `GS.3` | Capstone | [shutdown capstone](./3-capstone) | Wires signals, readiness, workers, and drain order together. | `GS.1`, `GS.2` |

## Next Step

After `GS.3`, return to the [10 Production Operations overview](../README.md) or continue to [Code Generation](../06-code-generation).
