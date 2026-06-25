# Payment workflow design

## Learning objective

Design a payment workflow using the saga pattern with compensating transactions, implement a workflow state machine in Go, and handle failures gracefully without leaving the system in an inconsistent state.

## Why this matters

Payment processing touches money. Every failure mode -- network outage, insufficient funds, fraud detection, timeout -- must be handled without losing money or corrupting state. A payment workflow cannot use distributed transactions (two-phase commit) across heterogeneous systems because databases and payment gateways do not support it. The saga pattern is the industry standard for coordinating multi-step payment workflows across services. Every fintech company, e-commerce platform, and subscription service uses this pattern. Understanding sagas is essential for any Go engineer building transaction-processing systems.

## Mental model

A payment workflow is a state machine. Each transition performs an action and moves to a new state. If a later step fails, earlier steps must be undone -- this is the saga pattern.

Think of booking a flight + hotel + car rental as a single trip. You book the flight (step 1), then the hotel (step 2), then the car (step 3). If the car rental fails, you need to cancel the hotel and the flight. Each cancellation is a compensating transaction. You cannot simply "roll back" because the airline's system is separate from the hotel's system. The saga orchestrates the forward steps and, on failure, runs the compensating steps in reverse order.

## Core idea

**Saga pattern**: a sequence of local transactions where each step has a compensating action that undoes it. There are two coordination models:

| Model | Description |
|---|---|
| Choreography | Each service publishes events; others react. No central coordinator. |
| Orchestration | A central coordinator tells each service what to do and handles failures. |

For payment workflows, orchestration is preferred because the coordinator can track state, enforce ordering, and run compensations reliably.

**State machine**: a payment progresses through states. Each state defines which transitions are valid:

```
Initiated -> Reserved -> Charged -> Completed
                 |            |
                 v            v
              Failed       Refunded
```

**Compensating transaction**: an action that semantically undoes a previous step. Canceling a hold on funds compensates a reservation. Issuing a refund compensates a charge.

## Under the hood

The orchestration engine (often called a workflow engine) maintains the saga log -- an ordered record of every step executed. When a step fails, the engine reads the log in reverse and runs each step's compensation. The saga log must be durable (written to a database) so that recovery survives a crash.

A payment pipeline typically has these steps:

1. **Validate**: check inputs, fraud screening.
2. **Reserve**: place a hold on the payment method (authorization).
3. **Capture**: settle the charge (capture the authorized amount).
4. **Fulfill**: deliver the product or service.
5. **Notify**: send receipt email.

Each of steps 2-4 must have a compensating action: release hold (void), refund, cancel order.

## How Go uses it

- **Temporal.io**: a Go-native workflow engine used by Netflix, Stripe, and Snapchat. Workflows are written as Go functions that can execute for days, survive process restarts, and automatically retry on failure.
- **Camunda Cloud / Zeebe**: supports Go clients for orchestration-based sagas.
- **In-house saga orchestrators**: many Go shops build lightweight orchestrators using a state machine library (`looplab/fsm`) and a database for the saga log.
- **`database/sql` transactions**: the simplest saga is a database transaction. Each step writes to the database. If any step fails, the earlier writes are rolled back. This only works within a single database, not across services.

## Go example

```go
package main

import (
	"errors"
	"fmt"
)

type State int

const (
	StateInitiated State = iota
	StateReserved
	StateCharged
	StateRefunded
	StateFailed
)

type PaymentWorkflow struct {
	ID     string
	Amount int
	State  State
	log    []string
}

func NewPaymentWorkflow(id string, amount int) *PaymentWorkflow {
	w := &PaymentWorkflow{ID: id, Amount: amount, State: StateInitiated}
	w.log = append(w.log, fmt.Sprintf("initiated %s for %d", id, amount))
	return w
}

func (w *PaymentWorkflow) Reserve() error {
	if w.State != StateInitiated {
		return fmt.Errorf("cannot reserve from %s", w.State)
	}
	w.State = StateReserved
	w.log = append(w.log, "reserved")
	return nil
}

func (w *PaymentWorkflow) Charge() error {
	if w.State != StateReserved {
		return fmt.Errorf("cannot charge from %s", w.State)
	}
	if w.Amount > 10000 {
		return errors.New("exceeds limit")
	}
	w.State = StateCharged
	w.log = append(w.log, "charged")
	return nil
}

func (w *PaymentWorkflow) Refund() error {
	if w.State != StateCharged {
		return fmt.Errorf("cannot refund from %s", w.State)
	}
	w.State = StateRefunded
	w.log = append(w.log, "refunded")
	return nil
}

func (w *PaymentWorkflow) Compensate() {
	switch w.State {
	case StateReserved:
		w.State = StateInitiated
		w.log = append(w.log, "compensated: released hold")
	case StateCharged:
		w.State = StateRefunded
		w.log = append(w.log, "compensated: refunded")
	}
}

type SagaStep struct {
	Name       string
	Execute    func() error
	Compensate func()
}

type Saga struct {
	steps []SagaStep
}

func (s *Saga) AddStep(name string, exec func() error, comp func()) {
	s.steps = append(s.steps, SagaStep{Name: name, Execute: exec, Compensate: comp})
}

func (s *Saga) Run() error {
	for i, step := range s.steps {
		if err := step.Execute(); err != nil {
			for j := i - 1; j >= 0; j-- {
				s.steps[j].Compensate()
			}
			return fmt.Errorf("step %q failed: %w", step.Name, err)
		}
	}
	return nil
}

func main() {
	w := NewPaymentWorkflow("ord-42", 7500)
	w.Reserve()
	w.Charge()
	fmt.Println("state:", w.State)

	w2 := NewPaymentWorkflow("ord-43", 7500)
	w2.Reserve()
	w2.Compensate()
	fmt.Println("after compensate:", w2.State)

	saga := &Saga{}
	var inventory int
	saga.AddStep("reserve_inventory", func() error {
		inventory = 50
		return nil
	}, func() { inventory = 0 })
	saga.AddStep("process_payment", func() error {
		return errors.New("card declined")
	}, func() {})
	saga.AddStep("confirm", func() error { return nil }, func() {})
	if err := saga.Run(); err != nil {
		fmt.Println("saga failed:", err)
		fmt.Println("inventory after rollback:", inventory)
	}
}
```

## Step-by-step execution

For a payment of 7500 cents:

1. `NewPaymentWorkflow("ord-42", 7500)` → state Initiated, log: `"initiated ord-42 for 7500"`.
2. `Reserve()`: state must be Initiated → set to Reserved, log: `"reserved"`.
3. `Charge()`: state must be Reserved → set to Charged, log: `"charged"`.
4. State is now Charged. The payment is captured.

For a saga with a failing step:

1. `reserve_inventory`: executes, sets inventory=50.
2. `process_payment`: returns error `"card declined"`.
3. Saga reverses: runs `reserve_inventory.Compensate()` → inventory=0.
4. Saga returns error. The system is in its original state.

## Common mistakes

- **Forgetting compensations**: every mutating step MUST have a compensating action. If a step is inherently non-reversible (e.g., "sent email"), the compensation should log the failure for manual reconciliation.
- **Synchronous saga in an HTTP handler**: running a saga that makes external API calls inside an HTTP request handler blocks the request for seconds. Use an asynchronous workflow engine instead.
- **Idempotency in saga steps**: if a step is retried (e.g., after a timeout), it must be idempotent. A charge step cannot charge twice. Use idempotency keys.
- **Not persisting saga state**: if the coordinator crashes, in-memory saga state is lost. The saga may be partially executed with no way to recover. Persist the saga log to a database.
- **Mixing orchestration and choreography**: do not have some compensations triggered by events and others by the coordinator. Pick one model and stick with it.

## Debugging walkthrough

A payment saga leaves the order in "charged" state but the inventory was already decremented:

```go
func processOrder(amount int) error {
	var inv int
	saga := &Saga{}
	saga.AddStep("decrement_inventory", func() error {
		inv = getInventory() - 1
		return nil
	}, func() { inv++ })
	saga.AddStep("charge_payment", func() error {
		return charge(amount) // this succeeds
	}, func() { refund(amount) })
	saga.AddStep("send_email", func() error {
		return sendEmail() // this fails
	}, func() {})
	return saga.Run()
}
```

**Symptom**: Inventory is decremented, payment is charged, but the email fails. The saga runs compensations: `charge_payment` is compensated (refund), `decrement_inventory` is compensated (increment). But the refund takes 5-7 business days. The customer sees the charge and then the refund. This is correct behavior, but the customer is confused.

**Root cause**: The saga worked correctly -- it rolled back all steps. But the user experience of a refund is worse than preventing the charge in the first place. Consider making the email step non-critical: if it fails, do not trigger a full rollback.

**Fix**: Mark `send_email` as a "best-effort" step that does not trigger compensation:

```go
func (s *Saga) RunNonCritical() error {
	for _, step := range s.steps {
		if err := step.Execute(); err != nil {
			// log and continue; do not compensate
		}
	}
	return nil
}
```

## Production notes

- **Saga log durability**: write each step completion and compensation to a database table. This allows recovery on restart.
- **Timeout handling**: each step should have a timeout. If a step does not respond within 10 seconds, mark it as failed and run compensations.
- **Monitoring**: track saga execution time, failure rate per step, and compensation frequency. A high compensation rate indicates a recurring failure in the pipeline.
- **Manual intervention**: some failures cannot be compensated automatically (e.g., a partial charge that succeeded on the gateway but the database write failed). Build a manual reconciliation dashboard.

## Performance implications

- **In-memory saga**: microseconds per step, but state is lost on crash.
- **Database-backed saga log**: adds ~5-10ms per log write. Acceptable for most payment workflows (sub-second total).
- **External API calls in saga steps**: 100-500ms per call (payment gateway, inventory service). The total saga time is the sum of all steps.
- **Concurrent sagas**: each saga is isolated. There is no contention unless two sagas touch the same inventory item. Use row-level locking for inventory decrements.

## Practice task

Extend the `PaymentWorkflow` with a new state `StatePendingFraudReview` and a transition `MarkFraudReview()` that requires state Initiated. Add a compensating transition `ReleaseHold()` that goes from `StateReserved` back to `StateInitiated`. Write a saga that includes a fraud check step that randomly fails, and verify compensation restores the state.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/12-payment-workflow-design
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/12-payment-workflow-design
```

The tests verify valid state transitions, invalid transition errors, amount limits, compensation behavior, and saga rollback.

## Review questions

1. What is the difference between a compensating transaction and a rollback?
2. Why is two-phase commit unsuitable for payment workflows across microservices?
3. In the saga pattern, what happens if a compensating transaction itself fails?
4. Why should saga state be persisted to a database rather than kept in memory?
5. What is the difference between orchestration and choreography sagas, and when would you use each?

## NEXT UP

Multi-tenancy architecture.
