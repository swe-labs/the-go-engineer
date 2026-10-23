# Opslane Milestone 08 Assessment — Payment Pipeline: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** What is the purpose of an idempotency key in payment systems?
- a) To encrypt credit card numbers in transit.
- b) To uniquely identify a payment transaction attempt, ensuring that network retries or duplicate client submissions do not result in multiple credit card charges.
- c) To calculate sales tax.
- d) To generate user passwords.

**2.** What occurs if a second request arrives with an idempotency key that is currently marked as `IN_PROGRESS`?
- a) The server charges the card a second time.
- b) The server rejects the request immediately with `ErrDuplicateInFlight`, preventing concurrent duplicate processing.
- c) The server cancels the original request.
- d) The server deletes the user's account.

**3.** Why does the payment pipeline check `if !errors.Is(err, ErrTransientNetwork)` before retrying?
- a) Non-transient errors (like `ErrCardDeclined` due to insufficient funds) cannot be resolved by immediate retries; retrying permanent errors wastes provider fees and causes unnecessary latency.
- b) Go's `for` loop cannot execute more than once on errors.
- c) Transient errors are automatically resolved by the operating system.
- d) Payment gateways forbid retries on network timeouts.

**4.** When a previously completed charge request is replayed with the same idempotency key, what does `PaymentPipeline.Charge` do?
- a) Executes a refund transaction.
- b) Returns the original cached `PaymentReceipt` immediately without invoking the external gateway adapter.
- c) Returns an HTTP 400 Bad Request.
- d) Decrements the user's account balance again.

**5.** In the retry loop, how does `delay` change between retry attempts?
- a) It remains constant at 10 milliseconds.
- b) It doubles after each attempt (`delay *= 2`), implementing exponential backoff.
- c) It decreases linearly.
- d) It is randomly set to zero.

**6.** What happens to the in-progress idempotency reservation if the caller's `context.Context` is cancelled during backoff sleep?
- a) The reservation remains locked in memory forever.
- b) The key reservation is safely deleted from `pp.records` and `ctx.Err()` is returned, allowing future retries.
- c) The entire process terminates.
- d) The database rolls back all historical payments.

---

## Short Answer (2 points each)

**7.** Explain the difference between transient errors (network timeouts, 503 Service Unavailable) and permanent errors (400 Bad Request, 402 Card Declined). Why must retry logic strictly differentiate between them?

**8.** What is a "Thundering Herd" or "Retry Storm" problem in distributed systems? How does adding randomized jitter to exponential backoff prevent it?

**9.** Why is thread synchronization (`pp.mu.Lock()`) required around `pp.records` in `PaymentPipeline`? What race condition would occur in its absence?

**10.** In an enterprise production deployment, why should the in-memory `records` map be replaced with a distributed store like Redis or PostgreSQL? What happens to an in-memory map when a service replica crashes or restarts?
