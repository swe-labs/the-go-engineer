# PD.3 Project Layout

## Mission

Decide when a Go project should stay flat, when it needs `internal/`, and when `cmd/` or other
top-level directories actually earn their place.

This surface is the package-design track output for Section 14.

## Files

- [main.go](./main.go): layout guide with small, medium, and large project examples

## Run Instructions

```bash
go run ./14-application-architecture/package-design/3-project-layout
```

## Success Criteria

You should be able to:

- explain when a flat layout is enough
- describe what `cmd/`, `internal/`, and `pkg/` are for
- call out common layout anti-patterns like `utils/`, `helpers/`, and premature folder sprawl

## Next Step

After `PD.3`, continue to [SL.1 slog basics](../../structured-logging) or back to the
[Package Design track](../README.md).
