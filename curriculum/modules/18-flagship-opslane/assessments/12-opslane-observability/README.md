# Opslane observability

## Mission

Understand and apply Opslane observability in the context of professional Go software engineering.

## Prerequisites

- opslane-11

## Mental Model

Observability is the application's instrument panel — without it, you're flying blind. Logs tell you what happened, metrics tell you how the system is behaving, and traces tell you the path a request took through the system.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, zerolog writes JSON log entries with zero heap allocations. Prometheus metrics use a map of metric name to *Desc structs with atomic value updates. OpenTelemetry traces propagate context via the context.Context through service boundaries.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/12-opslane-observability
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Logging everything at the same level (no INFO/WARN/ERROR distinction).
- Only logging errors without enough context (no request ID, no stack trace).
- Not exporting metrics — operating the application without visibility.

## In Production

Prometheus, Grafana, and OpenTelemetry are the industry standard for Go observability — used by Kubernetes, Docker, CoreDNS, and thousands of production Go services.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-13`.
