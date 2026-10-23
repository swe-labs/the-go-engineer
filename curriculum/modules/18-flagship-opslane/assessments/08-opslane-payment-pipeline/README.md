# Opslane payment pipeline

## Mission

Engineer an idempotent, fault-tolerant payment pipeline for Opslane, implementing in-flight deduplication locks, cached receipt replay, and exponential backoff retry mechanics for transient gateway failures.

## Prerequisites

- `opslane-07`

## Mental Model

In financial payment processing, network timeouts are ambiguous: did the payment gateway execute the charge before the connection dropped, or did it fail before reaching the gateway? Resending a raw payment charge blindly risks double-charging the customer. Opslane mitigates this with an Idempotent Payment Pipeline: every charge request includes a unique client-generated `idempotencyKey`. If an identical request arrives while the first is `IN_PROGRESS`, it is rejected with `ErrDuplicateInFlight`. Once completed, the original `PaymentReceipt` is cached; subsequent replays return the cached receipt immediately without touching the external payment gateway. If a transient network error occurs, the pipeline automatically retries using exponential backoff up to a configured threshold.

## Visual Model

```text
Client Charge Request (idempotencyKey, orderID, amount)
  │
  ▼
┌─────────────────────────────────────────────────────────────┐
│ PaymentPipeline.Charge                                      │
│   Check records[idempotencyKey]                             │
└──────────────┬───────────────────────────────┬──────────────┘
               │                               │
       Record Exists?                  Record Absent
       ├── Status == "COMPLETED"       │
       │   └── Return Cached Receipt   │ Reserve: Status = "IN_PROGRESS"
       └── Status == "IN_PROGRESS"     ▼
           └── ErrDuplicateInFlight    ┌─────────────────────────────┐
                                       │ Loop (attempt <= maxRetries)│
                                       │   gateway.ExecuteCharge()   │
                                       └──────────────┬──────────────┘
                                                      │
                     ┌────────────────────────────────┴──────────────┐
                     │ Success                                       │ Transient Error
                     ▼                                               ▼
     ┌──────────────────────────────┐              ┌─────────────────────────────┐
     │ Cache Status = "COMPLETED"   │              │ Wait: delay *= 2            │
     │ Return New PaymentReceipt    │              │ Retry if attempts remain    │
     └──────────────────────────────┘              └─────────────────────────────┘
```

## Machine View

Under the hood:
1. `pp.records` maps the idempotency key string to an `*IdempotencyRecord`. Access is serialized via `pp.mu.Lock()`.
2. Reserving the key with `Status: "IN_PROGRESS"` acts as an in-memory distributed lock against concurrent duplicate calls.
3. For errors classified as `ErrTransientNetwork`, the loop backs off using `time.After(delay)` while respecting `ctx.Done()`. If the caller's context cancels or times out, the key reservation is cleaned up (`delete(pp.records, idempotencyKey)`) and `ctx.Err()` is returned.
4. Non-transient errors (such as `ErrCardDeclined`) break out immediately on attempt 0 without retrying, preventing wasted cycles and credit card fraud flags.

## Run Instructions

```bash
# Run payment pipeline test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/08-opslane-payment-pipeline

# Run standalone demo
go run ./curriculum/modules/18-flagship-opslane/assessments/08-opslane-payment-pipeline
```

## Code Walkthrough

1. **Sentinels & Structs (lines 15–38)**: Declares errors (`ErrDuplicateInFlight`, `ErrCardDeclined`, `ErrTransientNetwork`), `PaymentReceipt`, `IdempotencyRecord`, and `PaymentGatewayAdapter` interface.
2. **`PaymentPipeline` Constructor (lines 40–56)**: Initializes mutex, records map, retry ceiling, and base delay.
3. **`Charge()` Execution (lines 58–136)**:
   - Evaluates key cache: returns cached receipt if completed, rejects if in progress.
   - Reserves key as `IN_PROGRESS`.
   - Executes charge loop: checks for transient vs permanent error; backs off exponentially (`delay *= 2`); cleans up on context cancelation.
   - On success: marks record `COMPLETED`, stores receipt, and returns receipt.
4. **`FlakyGatewaySimulator` (lines 138–160)**: Provides atomic call tracking and simulates transient network timeouts before succeeding.

## Try It

1. Run the test suite: verify that transient errors are retried and eventually succeed.
2. Submit two identical charge requests with the same idempotency key: verify the second returns the identical transaction ID with `gateway.CallCount == 1`.
3. Test permanent card decline: verify that `ExecuteCharge` is called only once and does not retry.

Common mistakes:
- Retrying permanent errors (card declined, invalid card number), which wastes gateway fees and angers users.
- Forgetting to clean up in-progress keys on context cancelation, leaving the idempotency key locked permanently.
- Retrying without exponential backoff, which causes thundering herd problems on struggling downstream providers.

## In Production

At cloud scale, idempotency records are stored in durable, distributed key-value caches like Redis with a TTL (e.g. 24 hours) or directly inside relational database tables (`payments.idempotency_key UNIQUE`). Exponential backoff should incorporate randomized jitter to prevent synchronized retry storms across concurrent clients.

## Thinking Questions

1. Why should non-transient errors (like card declines) never be retried with exponential backoff?
2. What happens if a network partition occurs *after* the gateway charges the card but *before* the server records the receipt?
3. How does randomized jitter improve standard exponential backoff in distributed systems?

## Next Step

Proceed to [Milestone 09 — Event Bus](../09-opslane-event-bus/README.md) to implement asynchronous in-memory event publishing and topic subscription.
