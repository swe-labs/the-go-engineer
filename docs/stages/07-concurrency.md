# 07 Concurrency

## Purpose

`07 Concurrency` teaches how Go coordinates work across goroutines, cancellation, time, and
bounded patterns without turning programs into chaos.

## Who This Is For

- learners who already understand backend-shaped code
- developers who want stronger judgment around coordination, deadlines, and pressure

## Mental Model

Concurrency is about ownership, cancellation, and limits.
The goal is not to add goroutines everywhere; the goal is to structure work so it can start, stop,
and fail in controlled ways.

## What You Should Learn Here

- goroutine and channel fundamentals
- context and cancellation
- time and scheduling workflows
- bounded concurrency patterns like pipelines and worker pools
- failure-aware coordination under load

## Current Source Content

- [10-concurrency](../../10-concurrency/)
- [11-concurrency-patterns](../../11-concurrency-patterns/)

## Stage Support Docs

- [Concurrency support index](./concurrency-system/README.md)
- [Stage map](./concurrency-system/stage-map.md)
- [Milestone guidance](./concurrency-system/milestone-guidance.md)

## Finish This Stage When

- you can explain when concurrency helps and when it hurts
- you can propagate cancellation through real work
- you can build bounded worker or pipeline patterns
- you can spot common leak, deadlock, and coordination mistakes

## Next Stage

Move to [08 Quality & Test](./08-quality-test.md).
