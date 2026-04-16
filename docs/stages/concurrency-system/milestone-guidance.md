# Concurrency System Milestone Guidance

## What Counts As Stage Completion

You should be able to finish the current live milestone path and explain the reasoning behind it,
not just run the final programs.

## Milestones

### `GC.7` concurrent downloader

This is the first proof that you can coordinate multiple goroutines without losing ownership of
results, synchronization, or shutdown.

### `CT.5` timeout-aware API client

This proves that cancellation is not just a local variable.
It has to reach real I/O boundaries cleanly.

### `TM.7` console reminder

This proves that timers and tickers are lifecycle tools, not background magic.

### `CP.5` URL health checker

This is the strongest current stage proof.
It asks you to combine cancellation, bounded concurrency, and failure-aware coordination in one
small system.

## Pressure Follow-Through

After `CP.5`, use the seeded diagnosis task to practice evidence-based failure analysis:

- [Diagnose CP.5 health checker failure](../expert-layer/tasks/diagnose-cp5-health-checker-failure.md)

## Bridge Path Check

If you are moving quickly through this stage, make sure you can still explain:

- why `WaitGroup` and channels solve different problems
- why context cancellation should reach HTTP and database calls
- why timers need explicit stop/drain behavior in long-running systems
- why bounded concurrency is safer than unbounded fan-out

## Ready To Move On

Move to [08 Quality & Test](../08-quality-test.md) when you can complete the
current path and explain the coordination choices in plain language.
