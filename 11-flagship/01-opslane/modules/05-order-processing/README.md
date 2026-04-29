# OPSL.5 Order Processing

## Mission

Move Opslane from CRUD-shaped handlers into a real business workflow with explicit state transitions.

## What This Module Builds

- order service layer
- validation rules for order creation and transitions
- inventory coordination seam
- idempotent order workflow entry

## You Are Here If

- `OPSL.4` is complete
- you can explain why handlers should stop at request parsing and response shaping
- you are ready to model state transitions instead of direct table writes

## Proof Surface

- `go test ./11-flagship/01-opslane/internal/services/...`
- state transition tests cover valid and invalid order movement
- idempotency tests prove retries do not create duplicate orders

Implemented code surface:

- [SURFACE.md](./SURFACE.md)

Primary files:

- `internal/services/order.go`
- `internal/services/inventory.go`
- `internal/services/validation.go`

## Required Files and Boundaries

The order service should own workflow rules.
Handlers should call into it, and repositories should stay focused on persistence.

## Engineering Questions

- What happens when 1,000 orders arrive together?
- What happens when inventory is reserved but payment fails?
- What happens when work stops halfway through a transition?

## Next Module

Continue to [OPSL.6 Payment Pipeline](../06-payment-pipeline/README.md) after the order state machine is provable.
