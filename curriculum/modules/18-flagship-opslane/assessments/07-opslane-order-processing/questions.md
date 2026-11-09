# Opslane Milestone 07 Assessment — Order Processing: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** Why does `OrderProcessor.ProcessOrder` use `defer tx.Rollback()` immediately after opening a database transaction?
- a) Because Go transactions always roll back automatically even if committed.
- b) To ensure safe cleanup: if any error, panic, or early return occurs before `tx.Commit()`, the transaction rolls back cleanly, preventing dangling locks or uncommitted partial state.
- c) To clear the database connection pool.
- d) Because SQLite requires rollback before committing.

**2.** In `AllowedTransitions`, which of the following states are terminal (allowing zero outbound transitions)?
- a) `StatusPending` and `StatusProcessing`.
- b) `StatusCompleted` and `StatusCancelled`.
- c) `StatusFailed` only.
- d) All states can transition to any other state.

**3.** What transaction isolation level is used by `ProcessOrder`?
- a) `sql.LevelReadUncommitted`.
- b) `sql.LevelSerializable`.
- c) `sql.LevelReadCommitted`.
- d) `sql.LevelDefault`.

**4.** If `gateway.ProcessPayment` fails due to insufficient funds, what is the final state of the order in the database?
- a) `StatusCompleted`.
- b) `StatusProcessing`.
- c) It remains in `StatusPending` because the transaction was rolled back.
- d) The order row is deleted.

**5.** Why does the state machine allow transitions from `StatusFailed` back to `StatusProcessing`?
- a) To allow transient failures (e.g. temporary network timeouts or expired card updates) to be retried safely.
- b) Because failed orders must always be refunded.
- c) To bypass the payment gateway.
- d) To change the order amount.

**6.** Why are the order status update and the audit log insert executed in the *same* database transaction?
- a) To save memory in Go slices.
- b) To guarantee atomic consistency: either both the order status change and the audit log record succeed, or neither does, preventing desynchronized audit trails.
- c) Because audit logs cannot be written in separate queries.
- d) To bypass database write locks.

---

## Short Answer (2 points each)

**7.** Explain the trade-offs of calling an external network API (e.g. a payment gateway) inside an active database transaction. What precautions must an engineer take to mitigate transaction lock exhaustion?

**8.** What is a Finite State Machine (FSM), and why is implementing an explicit state transition matrix superior to scattered `if/else` checks throughout service code?

**9.** How does serializable transaction isolation (`sql.LevelSerializable`) protect against concurrent race conditions during order checkout?

**10.** Explain how idempotency keys prevent double charging customers during network failures between an order processing service and an external payment provider.
