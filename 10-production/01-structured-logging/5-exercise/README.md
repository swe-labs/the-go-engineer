# SL.5 PII Redactor

## Mission

Use `slog.HandlerOptions.ReplaceAttr` to build a logger that automatically redacts sensitive
attributes before they reach the output handler.

This surface is the structured-logging track output for Stage 10.

## Files

- [main.go](./main.go): completed solution
- [_starter/main.go](./_starter/main.go): exercise starter

## Run Instructions

```bash
go run ./10-production/01-structured-logging/5-exercise
go run ./10-production/01-structured-logging/5-exercise/_starter
```

## Success Criteria

You should be able to:

- use `ReplaceAttr` to transform attributes centrally
- redact sensitive keys without changing the logging call sites
- explain why this is safer than manually editing every `slog.Info` call

## Next Step

After `SL.5`, continue to [GS.1 signal context](../../02-graceful-shutdown) or back to the
[Structured Logging track](../README.md).
