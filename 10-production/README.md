# 10 Production Operations

## Mission

This stage combines runtime behavior, configuration, observability, delivery mechanics, and build-time automation into one production-oriented path. It prepares learners to deploy Go services responsibly.

By the end of this stage, a learner should be able to:

- implement structured logging with `slog`
- orchestrate graceful shutdown sequences for zero-downtime deploys
- validate configuration robustly on boot
- containerize Go applications effectively with Docker
- integrate code generation workflows (like `sqlc`)

## Stage Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `SL.1-5` | Structured Logging | Build a telemetry foundation with `slog`. |
| `GS.1-3` | Graceful Shutdown | Coordinate signals, HTTP draining, and cleanup. |
| `CFG.1-5` & `OPS.1-5` | Config & Ops | Add startup discipline, validation, metrics, and alerting habits. |
| `DOCKER.1-3` & `DEPLOY.1-3` | Delivery | Turn the service into a deployable container artifact. |
| `CG.1-3` | Code Generation | Introduce generation (e.g., `sqlc`) in a CI context. |

## Why This Stage Exists Now

The learner already knows:

- how to build secure, well-architected APIs
- how to test and profile them

That is enough to start asking engineering questions like:

- how do we deploy this service?
- how do we know if it is broken in production?
- how do we shut it down without dropping user requests?

## Suggested Learning Flow

1. Structured logging and graceful shutdown form the operational runtime foundation.
2. Configuration and Ops add startup discipline and observability habits.
3. Docker and Deploy turn the service into a deliverable runtime artifact.
4. Code Generation introduces tooling in the context where teams usually adopt it.

## Next Step

After this section, continue to [11 Flagship](../11-flagship).
