# GS.3 Shutdown Capstone

## Mission

Coordinate readiness, HTTP draining, worker draining, and final cleanup in one production-style
shutdown flow.

This surface is the graceful-shutdown track output for Section 14.

## Files

- [main.go](./main.go): production-style shutdown capstone with readiness and worker coordination

## Run Instructions

```bash
go run ./14-application-architecture/graceful-shutdown/3-capstone
```

Then press `Ctrl+C` or send `SIGTERM` and watch the shutdown sequence.

## Success Criteria

You should be able to:

- explain why readiness must flip before draining begins
- describe the shutdown dependency order across HTTP, workers, and shared resources
- show how `errgroup` helps coordinate a multi-component shutdown

## Next Step

After `GS.3`, continue to [Section 15: Code Generation](../../../15-code-generation) or back to the
[Graceful Shutdown track](../README.md).
