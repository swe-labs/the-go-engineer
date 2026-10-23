# Opslane order processing

## Mission

Build a deterministic order lifecycle state machine and transactional orchestration engine in Go, enforcing valid status transitions, serializable ACID transaction semantics, external payment integration, and immutable audit logging.

## Prerequisites

- `opslane-06`

## Mental Model

E-commerce and financial order processing cannot tolerate undefined states, partial executions, or double charges. An order must progress strictly through allowed state transitions: an order in `COMPLETED` or `CANCELLED` status cannot be reprocessed. In Opslane, the order lifecycle is governed by an explicit state transition matrix (`AllowedTransitions`). Processing an order wraps status progression, external payment authorization, and audit log generation in a single atomic database transaction with `LevelSerializable` isolation. If payment fails or a crash occurs, `defer tx.Rollback()` ensures the database reverts to its exact prior state.

## Visual Model

```text
               ┌──────────┐
               │ PENDING  │
               └────┬─────┘
                    │
         ┌──────────┴──────────┐
         │ (ProcessOrder)      │ (CancelOrder)
         ▼                     ▼
  ┌────────────┐        ┌───────────┐
  │ PROCESSING │        │ CANCELLED │ [Terminal]
  └──────┬─────┘        └───────────┘
         │
    ┌────┴──────────────────────────┐
    │ Payment Gateway Execution     │
    ├─────────────────┬─────────────┤
    │ Success         │ Failure     │
    ▼                 ▼             │
┌───────────┐   ┌─────────┐         │
│ COMPLETED │   │ FAILED  ├─────────┘ (Retry Capability)
└───────────┘   └─────────┘
 [Terminal]
```

## Machine View

Under the hood:
1. `AllowedTransitions` enforces that `COMPLETED` and `CANCELLED` states have zero allowed outbound transitions, guaranteeing immutability once finalized.
2. `db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})` guarantees the highest degree of transaction isolation, preventing dirty reads, non-repeatable reads, and phantom reads under concurrent access.
3. If `op.gateway.ProcessPayment()` returns an error, the handler returns early without calling `tx.Commit()`. The deferred `tx.Rollback()` executes automatically, discarding the intermediate `PROCESSING` status update and ensuring the order remains in its original status.
4. An immutable audit record is inserted within the same transaction commit, guaranteeing an unforgeable trace of the order's completion.

## Run Instructions

```bash
# Run order processing test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/07-opslane-order-processing

# Run standalone demo
go run ./curriculum/modules/18-flagship-opslane/assessments/07-opslane-order-processing
```

## Code Walkthrough

1. **State Machine Definition (lines 20–50)**: Defines `OrderStatus` constants, `AllowedTransitions` map, and `CanTransition` validator.
2. **`PaymentGateway` Interface (lines 53–55)**: Inverts external payment processing dependency.
3. **`OrderProcessor` & Schema (lines 57–85)**: Sets up database tables for `orders` and `audit_logs`.
4. **`ProcessOrder` Workflow (lines 88–144)**:
   - Starts serializable transaction with `defer tx.Rollback()`.
   - Reads current status and amount.
   - Evaluates `CanTransition(currentStatus, StatusProcessing)`.
   - Updates status to `PROCESSING`.
   - Invokes external `gateway.ProcessPayment(ctx, orderID, amount)`.
   - On payment failure: aborts and rolls back transaction.
   - On payment success: updates status to `COMPLETED`, writes audit log, and commits transaction.

## Try It

1. Run the test suite: verify that an order advances from PENDING to COMPLETED.
2. Test payment failure: inject a failing gateway and verify that `ProcessOrder` rolls back, leaving the order in PENDING.
3. Test terminal state protection: attempt to reprocess a COMPLETED order and verify `ErrInvalidStateTransition` is returned.

Common mistakes:
- Not using `defer tx.Rollback()`, leaving open transactions and locks if an error returns before commit.
- Updating database status before verifying payment authorization, risking stranded "completed" orders when payments fail.
- Allowing status transitions from arbitrary strings without validating against a state machine.

## In Production

In distributed e-commerce backends, payment calls should be designed with idempotency keys so that network timeouts can be retried without double charging the customer. If an external payment provider is slow, short context timeouts prevent database connections and transaction locks from being held indefinitely.

## Thinking Questions

1. Why should `defer tx.Rollback()` always be called immediately after `db.BeginTx()`?
2. What are the risks of performing an external network RPC (like a payment gateway) inside an open database transaction?
3. How does the state machine allow failed orders to be retried while preventing completed orders from being reprocessed?

## Next Step

Proceed to [Milestone 08 — Payment Pipeline](../08-opslane-payment-pipeline/README.md) to implement idempotent payment charging, provider webhooks, and ledger transactions.
