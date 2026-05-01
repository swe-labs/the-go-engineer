# OPSL.6 Payment Pipeline

## Mission

Model payment work as a reliability problem, not just another insert statement.

## What This Module Builds

- payment service workflow
- gateway integration seam
- retry and timeout behavior
- reconciliation-safe duplicate protection

## You Are Here If

- `OPSL.5` is complete
- order state transitions exist
- payment handling can move into a dedicated workflow layer

## Proof Surface

This module is implemented in the current tree.

Run:

```bash
go test ./11-flagship/01-opslane/internal/payment/...
go test ./11-flagship/01-opslane/internal/services/...
go test ./11-flagship/01-opslane/internal/handlers/...
go run ./11-flagship/01-opslane/scripts/progress.go
```

The proof surface now covers:

- gateway timeout behavior
- worker job processing boundaries
- provider-reference duplicate protection
- payment reconciliation after a pending/timeout case
- HTTP handler integration through the payment service

Implemented files:

- `internal/payment/gateway.go`
- `internal/payment/worker.go`
- `internal/services/payment.go`

Implementation map: [SURFACE.md](./SURFACE.md)

## Required Files and Boundaries

Keep gateway behavior behind an explicit boundary.
Do not scatter payment retry logic across handlers and repositories.

## Machine View

The learner-facing request is small:

```text
POST /api/v1/payments
```

Behind that request the machine now does more work:

```text
handler
  -> reads trusted tenant identity from auth middleware
  -> builds payment.Job
  -> calls PaymentService

PaymentService
  -> checks tenant-scoped order exists
  -> creates one pending payment row
  -> moves order to processing
  -> calls gateway through a timeout context
  -> updates payment status
  -> moves order to paid or failed when gateway result is final
```

The important idea is that the handler does not decide gateway, retry, or order-state rules.
It only translates HTTP into a workflow request.

## Diagram

```text
Client request
    |
    v
HTTP handler
    |
    v
PaymentService
    |
    +--> payments table: pending row
    |
    +--> OrderService: pending -> processing
    |
    +--> Gateway: charge with timeout
    |
    +--> payments table: settled/failed
    |
    +--> OrderService: processing -> paid/failed
```

## Try It

Change the simulated gateway in a test from `settled` to `failed`.
Observe that the payment status changes to `failed` and the order transitions to `failed`.

Then change the gateway to return `ErrGatewayTimeout`.
Observe that the payment stays `pending` and the order stays `processing`, because the system does
not know whether the external provider eventually charged the customer.

## Engineering Questions

- What happens on gateway timeout?
- How do you stop duplicate payment attempts from becoming duplicate charges?
- What should happen when success is observed locally but the follow-up callback never arrives?

## Next Step

Next: `OPSL.7` -> `11-flagship/01-opslane/modules/07-event-workers`

Open `11-flagship/01-opslane/modules/07-event-workers/README.md` to continue.
