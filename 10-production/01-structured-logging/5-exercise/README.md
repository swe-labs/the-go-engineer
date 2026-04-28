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


## 



## 



## 



## 



## 



## Try 



## 




## Prerequisites

You should be comfortable with Go syntax, basic data structures, and the control flow mechanics covered in earlier sections.

## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Solution Walkthrough

The solution demonstrates a complete implementation, proving the concept by bridging the individual requirements into a single, cohesive executable.

## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## Verification Surface

The correctness of this component is proven by its associated test suite. We verify boundaries, handle edge cases, and ensure performance constraints are met.

## In Production

Logging Personally Identifiable Information (PII) or secrets is a major compliance violation in production systems. When a user's password, credit card number, or social security number leaks into centralized logging platforms (like Datadog, Splunk, or CloudWatch), it triggers massive security incidents because logs are often accessible to many engineers and are retained for months. Asking developers to "just be careful" and manually redact fields at every `slog.Info` call site never works at scale. The `ReplaceAttr` pattern this exercise teaches pushes redaction down to the handler layer. This guarantees that no matter where a developer logs a `password` or `ssn` field in the codebase, it will be scrubbed before leaving the application. In production, this pattern is often extended to redact entire struct fields (using reflection) or scrub data based on complex compliance rules, forming the backbone of secure telemetry.

## Thinking Questions

1. Why is redacting fields in a centralized `ReplaceAttr` handler safer than relying on developers to redact fields manually before calling `slog.Info`?
2. If you need to log an API request body for debugging, but it might contain PII, what strategies can you use to ensure the logs remain secure?
3. What is the performance cost of inspecting every attribute with `ReplaceAttr` on every log line, and how can you minimize it?
4. How would you handle a situation where a developer logs sensitive data as part of a free-text string message (e.g., `slog.Info("user password is " + password)`) instead of as a structured attribute?

## Next Step

After `SL.5`, continue to [GS.1 signal context](../../02-graceful-shutdown) or back to the
[Structured Logging track](../README.md).


