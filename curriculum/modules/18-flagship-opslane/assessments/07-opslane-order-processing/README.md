# Opslane order processing

## Mission

Understand and apply Opslane order processing in the context of professional Go software engineering.

## Prerequisites

- opslane-06

## Mental Model

Order processing is a saga — a sequence of steps where each step has a compensating action for rollback. Reserve inventory (compensate: release inventory), charge payment (compensate: refund), create order (compensate: cancel order). The saga must handle partial failures: if step 3 fails after step 2 succeeded, the saga runs the compensating action for step 2. The idempotency key at the top level ensures that retries do not create duplicate orders. The payment-specific idempotency key ensures that the payment gateway is not called twice for the same charge.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

The saga pattern is implemented using local database transactions for local state and compensating API calls for external side effects. The local transaction handles inventory reservation and order creation atomically. The payment gateway call is outside the transaction (it cannot be inside a SQL transaction) — if it fails, the transaction is rolled back, releasing the inventory reservation. If the local transaction fails after payment succeeds, a refund is issued. The idempotency key is stored before the transaction begins (for the overall order) and before the payment call (for the payment). In Go, the saga is implemented as a state machine: OrderState transitions from Pending → PaymentPending → Confirmed → Fulfilled, with compensating transitions for each failure. The state machine is persisted in the database and retried by a background worker if the handler crashes mid-saga.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/07-opslane-order-processing
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Not using idempotency keys for payment processing — a network retry causes the payment gateway to charge the customer twice, and the duplicate charge is invisible until the customer reports it on their credit card statement.
- Calling the payment gateway inside the HTTP handler without a timeout — a slow gateway blocks the handler for 60 seconds, exhausting the goroutine pool and causing all other requests to time out.
- Acknowledging a queue message before processing completes — if the handler crashes after acknowledgement but before the order is persisted, the order is lost forever.
- Assuming inventory reservation is synchronous — reserving inventory in a transaction that commits successfully does not guarantee the inventory is actually available if the reservation expires before fulfillment.
- Logging credit card numbers or CVV in payment error logs — PCI DSS violation, immediate compliance failure, and the logs become a legal liability.

## In Production

Every e-commerce platform implements order processing as a saga. Shopify, Stripe, and Amazon all handle order processing with idempotency keys, inventory locking, and compensating transactions for payment refunds. Opslane's order processing must handle: (1) duplicate order creation (idempotency key), (2) payment failure after inventory reservation (release inventory), (3) inventory oversell under concurrent load (SELECT FOR UPDATE), and (4) payment gateway timeout (async reconciliation via webhook).

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-08`.
