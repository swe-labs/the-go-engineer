# Opslane caching layer

## Mission

Understand and apply Opslane caching layer in the context of professional Go software engineering.

## Prerequisites

- opslane-10

## Mental Model

The caching layer is the application's short-term memory — frequently accessed information is kept at the front of your mind (in cache) for instant recall, while rarely used information stays in long-term storage (database).

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, the in-memory cache uses a map[string]*entry where entry contains the value, expiry time, and optional size. A background goroutine sweeps expired entries periodically. Redis-backed cache uses SETEX/GET with serialized JSON values.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/11-opslane-caching-layer
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Caching data without setting a TTL, serving stale data forever.
- Using a cache-aside pattern but forgetting to populate the cache on writes.
- Caching large data sets in a single in-memory map without eviction.

## In Production

Caching is fundamental to web performance — CDNs cache static assets, DNS caches IP lookups, database engines cache query results. Opslane applies caching at the application layer where it has the most impact on API latency.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-12`.
