# Opslane worker pools

## Mission

Understand and apply Opslane worker pools in the context of professional Go software engineering.

## Prerequisites

- opslane-09

## Mental Model

A worker pool is the application's assembly line — a fixed number of workers at stations, each taking a task from the conveyor belt (channel), processing it, and reaching for the next one. The belt can hold some buffer of waiting tasks.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, the work channel is a buffered channel created with make(chan Task, bufferSize). Workers are goroutines in a for-range loop over the channel — when the channel is closed, the loop exits. Panic recovery uses defer/recover with a stack trace log.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/10-opslane-worker-pools
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Spawning a new goroutine for every task without any limit.
- Not handling worker panics — a single panic kills the entire pool.
- Forgetting to drain the work channel on shutdown, losing pending tasks.

## In Production

Worker pool patterns are used by Sidekiq (Ruby), Celery (Python), and Kubernetes job controllers. Go's goroutine-based pools are the foundation of high-throughput Go services.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-11`.
