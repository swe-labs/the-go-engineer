# Opslane event bus

## Mission

Understand and apply Opslane event bus in the context of professional Go software engineering.

## Prerequisites

- opslane-08

## Mental Model

An event bus is the application's announcement system — instead of calling every department individually when something happens, you make an announcement on the PA system and let interested departments respond at their own pace.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, the in-memory event bus uses a map of topics to slices of channel subscribers. Each subscriber gets its own goroutine and buffered channel. NATS-backed bus uses NATS subjects and queue groups for load-balanced subscription.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/09-opslane-event-bus
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Tightly coupling services by making direct HTTP calls for every cross-service action.
- Losing events when the application restarts (no persistent event storage).
- Not handling event replay for new subscribers that need historical data.

## In Production

Event-driven architectures power modern cloud services — AWS SNS/SQS, Google Pub/Sub, Kafka, and NATS all implement the pub/sub pattern. Opslane uses the same architectural pattern at application scale.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-10`.
