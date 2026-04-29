# 10 Production Operations

## Mission

Learn how to turn working Go services into operable software. This section covers structured logging, graceful shutdown, configuration, observability, containerization, deployment, and code generation.

## Section Map

| Track | Surface | Mission |
| --- | --- | --- |
| `SL.1-SL.5` | [Structured Logging](./01-structured-logging) | teach structured telemetry with `slog` |
| `GS.1-GS.3` | [Graceful Shutdown](./02-graceful-shutdown) | teach signals, deadlines, cleanup, and service shutdown |
| `CFG.1-CFG.5` | [Configuration](./04-configuration) | teach environment-driven startup discipline and validation |
| `OPS.1-OPS.5` | [Observability](./05-observability) | teach health checks, metrics, tracing, and operational visibility |
| `DOCKER.1-DOCKER.3` | [Containerization](./03-docker-and-deployment) | teach minimal, reproducible, and secure container images |
| `DEPLOY.1-DEPLOY.3` | [Deployment](./03-docker-and-deployment) | teach deployment-aware service behavior and release checks |
| `CG.1-CG.3` | [Code Generation](./06-code-generation) | teach `go:generate` and build-time automation |

## Why This Section Matters

Writing the code is only half the battle. Running that code reliably in a cluster requires a different set of skills.

1. **Visibility**: Structured logging (Track SL) and metrics (Track OPS) allow you to debug problems in production without logging into the server.
2. **Availability**: Graceful shutdown (Track GS) ensures that when you deploy a new version, the old version finishes its work before exiting, resulting in zero dropped requests.
3. **Security & Speed**: Modern containerization (Track DOCKER) ensures your artifacts are small, fast to deploy, and have a minimal attack surface.

## Suggested Learning Flow

1. **Start with Track SL & GS**: These are the fundamental runtime behaviors every production Go service must have.
2. **Continue to Track CFG & OPS**: Build the "Control Panel" for your application.
3. **Finish with Tracks DOCKER, DEPLOY, and CG**: Turn your code into a reproducible delivery artifact and verify it before release.

## Next Step

After learning production operations, continue to [Section 11 Flagship](../11-flagship) or return to the [Curriculum Overview](../README.md).
