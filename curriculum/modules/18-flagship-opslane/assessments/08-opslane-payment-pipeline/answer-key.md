# Opslane Milestone 08 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) To uniquely identify a payment transaction attempt, ensuring that network retries or duplicate client submissions do not result in multiple credit card charges.**  
*Explanation:* An idempotency key guarantees that repeated execution of an operation yields the exact same outcome as the first execution, preventing multiple billing events during network retries.

**2. b) The server rejects the request immediately with `ErrDuplicateInFlight`, preventing concurrent duplicate processing.**  
*Explanation:* Marking the key as `IN_PROGRESS` locks the operation. If a double-click or simultaneous duplicate request arrives before the first completes, it is rejected immediately to prevent concurrent double charging.

**3. a) Non-transient errors (like `ErrCardDeclined` due to insufficient funds) cannot be resolved by immediate retries; retrying permanent errors wastes provider fees and causes unnecessary latency.**  
*Explanation:* A declined credit card or invalid CVV will not change state within 10 milliseconds. Retrying permanent errors wastes gateway resources, incurs extra transaction fees, and degrades user experience.

**4. b) Returns the original cached `PaymentReceipt` immediately without invoking the external gateway adapter.**  
*Explanation:* Once an idempotency record is `COMPLETED`, all subsequent identical requests receive the cached receipt immediately without touching the external payment gateway, guaranteeing zero additional charges.

**5. b) It doubles after each attempt (`delay *= 2`), implementing exponential backoff.**  
*Explanation:* Exponential backoff progressively widens the delay between retries (e.g. 10ms -> 20ms -> 40ms -> 80ms), giving overloaded downstream services time to recover.

**6. b) The key reservation is safely deleted from `pp.records` and `ctx.Err()` is returned, allowing future retries.**  
*Explanation:* If the client disconnects or times out before the payment executes, the in-progress reservation must be deleted so the client can retry the payment later with the same idempotency key.

---

## Short Answer (2 points each)

**7. Transient vs. Permanent errors in payment pipelines:**
- **Transient Errors:** Temporary issues that are expected to self-resolve within milliseconds to seconds (e.g. network timeouts, connection resets, 503 Service Unavailable, 429 Too Many Requests). Retrying these with backoff is appropriate.
- **Permanent Errors:** Deterministic failures where retrying with the same parameters will always produce the same failure (e.g. 400 Bad Request, 402 Card Declined / Insufficient Funds, 401 Invalid API Key). Retrying permanent errors must be aborted immediately to avoid unnecessary costs and fraud flags.
*Grading: 1 pt for distinguishing transient vs permanent characteristics; 1 pt for explaining why permanent errors must not retry.*

**8. Thundering Herd problems and randomized jitter:**
- **Problem:** If a downstream service drops connections, hundreds of concurrent callers running standard exponential backoff will all back off for the exact same duration ($2^n \times \text{delay}$) and retry simultaneously in a synchronized burst. This avalanche of synchronized retries immediately knocks the struggling service offline again.
- **Randomized Jitter Mitigation:** Adding random jitter (e.g. `sleep = delay * (0.5 + rand.Float64() * 0.5)`) spreads the retry attempts uniformly across time, smoothing the traffic curve and allowing the downstream service to recover.
*Grading: 1 pt for explaining synchronized retry storm; 1 pt for explaining how random jitter desynchronizes traffic.*

**9. Concurrency race conditions in `pp.records`:**
Go maps are not thread-safe for concurrent read and write operations. If two concurrent goroutines attempt to check or insert idempotency keys into `pp.records` without `pp.mu.Lock()`, the Go runtime detects a concurrent map read/write and immediately crashes with a fatal panic (`fatal error: concurrent map writes`). Mutex synchronization ensures that key reservation and retrieval are atomic.
*Grading: 1 pt for identifying Go runtime fatal crash on concurrent map writes; 1 pt for explaining atomic key reservation.*

**10. In-memory maps vs. Distributed stores (Redis / PostgreSQL) in production:**
- **Limitation of In-Memory Maps:** An in-memory map lives exclusively inside the RAM of a single process. In modern multi-replica deployments (e.g. 3 pods in Kubernetes), a user's first request might hit Pod 1 while a retry hits Pod 2. Pod 2 has no visibility into Pod 1's memory and will re-charge the customer, breaking idempotency. Furthermore, if a pod restarts, all in-memory receipts are lost.
- **Distributed Store Solution:** Storing idempotency records in a shared Redis cluster or relational database table (`payments` with unique constraint on `idempotency_key`) guarantees that all service replicas share a single, durable source of truth that survives pod restarts.
*Grading: 1 pt for explaining multi-pod isolation in microservices; 1 pt for explaining distributed store (Redis/DB) shared durability.*
