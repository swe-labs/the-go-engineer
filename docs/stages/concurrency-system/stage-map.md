# Concurrency System Stage Map

## Stage Goal

This stage teaches learners how to structure concurrent work so it can start, coordinate, cancel,
and stop without turning the program into hidden chaos.

## Public Stage Shape

| Part | Source Surface | Role In Stage |
| --- | --- | --- |
| Foundations | [Section 11: Concurrency](../../../11-concurrency/) | teaches goroutines, context, timers, and the vocabulary of controlled concurrent work |
| Patterns | [Section 12: Concurrency Patterns](../../../12-concurrency-patterns/) | teaches `errgroup`, bounded pipelines, and pressure-aware concurrency patterns |

## Live Milestone Backbone

| ID | Surface | Why It Matters |
| --- | --- | --- |
| `GC.7` | concurrent downloader | proves goroutine coordination and bounded fan-out |
| `CT.5` | timeout-aware API client | proves cancellation and deadline propagation through I/O |
| `TM.7` | console reminder | proves timer/ticker lifecycle management |
| `CP.5` | URL health checker | proves `errgroup`, sibling cancellation, and bounded concurrent work |

## Recommended Order

1. `GC.1` through `GC.7`
2. `CT.1` through `CT.5`
3. `TM.1` through `TM.7`
4. `CP.1` through `CP.5`

## Reference Surfaces That Still Matter

These are still useful, but they are not the current public beta backbone:

- `GC.8`
- `GC.9`
- `GC.10`
- `TM.4`
- `TM.5`
- `TM.6`
- `6-worker-pool`
