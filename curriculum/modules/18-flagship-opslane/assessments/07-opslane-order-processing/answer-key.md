# Opslane Milestone 07 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) To ensure safe cleanup: if any error, panic, or early return occurs before `tx.Commit()`, the transaction rolls back cleanly, preventing dangling locks or uncommitted partial state.**  
*Explanation:* If `tx.Commit()` succeeds, a subsequent call to `tx.Rollback()` is a harmless no-op returning `sql.ErrTxDone`. Deferring `tx.Rollback()` guarantees that if an error or panic occurs anywhere in the function, all uncommitted locks and mutations are immediately rolled back.

**2. b) `StatusCompleted` and `StatusCancelled`.**  
*Explanation:* In `AllowedTransitions`, `StatusCompleted: {}` and `StatusCancelled: {}` map to empty slices. Once an order is completed or cancelled, it has reached a terminal state and cannot transition further.

**3. b) `sql.LevelSerializable`.**  
*Explanation:* Serializable isolation is the strictest ACID level, ensuring that concurrent transactions execute as if they occurred serially, preventing phantom reads and concurrent update anomalies.

**4. c) It remains in `StatusPending` because the transaction was rolled back.**  
*Explanation:* Because the intermediate `PROCESSING` update was executed inside the transaction and `ProcessPayment` returned an error, the function returned early without committing, triggering the deferred rollback and leaving the order in its original `StatusPending` state.

**5. a) To allow transient failures (e.g. temporary network timeouts or expired card updates) to be retried safely.**  
*Explanation:* Unlike a terminal state, a failure might be transient. Allowing transitions from `FAILED` to `PROCESSING` enables automated retries or user-initiated checkout attempts with updated card details.

**6. b) To guarantee atomic consistency: either both the order status change and the audit log record succeed, or neither does, preventing desynchronized audit trails.**  
*Explanation:* In compliance-sensitive systems, audit logs cannot be out of sync with actual database mutations. Writing both within the same transaction guarantees that an audit record is created if and only if the order status update commits.

---

## Short Answer (2 points each)

**7. Trade-offs of calling external APIs inside database transactions:**
- **Risk:** External network calls can stall for seconds due to latency or network packet loss. Holding open a database transaction while waiting on a network call holds row locks and database connections in the connection pool, rapidly causing lock contention and connection pool exhaustion across the entire application.
- **Precautions:** Always enforce strict, aggressive `context.Context` timeouts (e.g. 2–3 seconds) on the external call. Alternatively, use an Outbox/Saga pattern where status is committed locally as `PROCESSING`, payment is executed asynchronously outside the transaction, and status is finalized in a second transaction upon payment callback.
*Grading: 1 pt for identifying connection pool/lock exhaustion; 1 pt for timeout enforcement or asynchronous saga pattern.*

**8. Finite State Machine (FSM) architecture:**
- **What is an FSM:** A mathematical model of computation consisting of a finite set of states, an initial state, and transition rules that dictate which states can be reached from the current state.
- **Why Centralized Matrix is Superior:** Scattered `if/else` checks across multiple controllers and services inevitably lead to missed edge cases, conflicting logic, and accidental illegal transitions (such as refunding an already cancelled order). Centralizing transition rules in an authoritative matrix (`AllowedTransitions`) creates a single source of truth that is easy to audit, test, and maintain.
*Grading: 1 pt for definition of FSM; 1 pt for central single source of truth vs scattered conditionals.*

**9. Serializable isolation against concurrent race conditions:**
If two concurrent checkout requests for the same order arrive simultaneously, lower isolation levels (Read Committed) might permit both requests to read `status = PENDING` simultaneously, execute dual payment charges, and overwrite each other. `LevelSerializable` detects concurrent read/write conflicts on the same rows and aborts the conflicting transaction with a serialization failure, guaranteeing that only one checkout transaction can proceed.
*Grading: 1 pt for explaining concurrent read-then-write race condition (double charge); 1 pt for explaining serializable conflict detection/abort.*

**10. Idempotency keys in payment processing:**
- **Mechanism:** The client generates a unique UUID (idempotency key) for the checkout attempt and sends it to the payment gateway.
- **Network Failure Mitigation:** If the payment gateway successfully charges the card but the network connection drops before the client receives the 200 OK response, the client retries with the *same* idempotency key. The payment gateway recognizes the key, suppresses the second charge, and returns the cached result of the original transaction, guaranteeing that the customer is never charged twice.
*Grading: 1 pt for defining unique client-generated key; 1 pt for explaining duplicate suppression upon network retry.*
