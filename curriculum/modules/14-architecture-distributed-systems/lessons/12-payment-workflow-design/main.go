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

func (s State) String() string {
	switch s {
	case StateInitiated:
		return "initiated"
	case StateReserved:
		return "reserved"
	case StateCharged:
		return "charged"
	case StateRefunded:
		return "refunded"
	case StateFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type PaymentWorkflow struct {
	ID     string
	Amount int
	State  State
	log    []string
}

func NewPaymentWorkflow(id string, amount int) *PaymentWorkflow {
	w := &PaymentWorkflow{ID: id, Amount: amount, State: StateInitiated}
	w.log = append(w.log, fmt.Sprintf("%s: initiated payment %s for %d cents", w.State, id, amount))
	return w
}

func (w *PaymentWorkflow) Reserve() error {
	if w.State != StateInitiated {
		return fmt.Errorf("cannot reserve from state %s", w.State)
	}
	w.State = StateReserved
	w.log = append(w.log, fmt.Sprintf("reserved %d cents", w.Amount))
	return nil
}

func (w *PaymentWorkflow) Charge() error {
	if w.State != StateReserved {
		return fmt.Errorf("cannot charge from state %s", w.State)
	}
	if w.Amount > 10000 {
		return errors.New("amount exceeds limit")
	}
	w.State = StateCharged
	w.log = append(w.log, fmt.Sprintf("charged %d cents", w.Amount))
	return nil
}

func (w *PaymentWorkflow) Refund() error {
	if w.State != StateCharged {
		return fmt.Errorf("cannot refund from state %s", w.State)
	}
	w.State = StateRefunded
	w.log = append(w.log, fmt.Sprintf("refunded %d cents", w.Amount))
	return nil
}

func (w *PaymentWorkflow) Fail() {
	w.State = StateFailed
	w.log = append(w.log, fmt.Sprintf("failed payment %s", w.ID))
}

func (w *PaymentWorkflow) Compensate() {
	switch w.State {
	case StateReserved:
		w.State = StateInitiated
		w.log = append(w.log, "compensated: released reservation")
	case StateCharged:
		w.State = StateRefunded
		w.log = append(w.log, fmt.Sprintf("compensated: refunded %d cents", w.Amount))
	case StateInitiated:
		w.log = append(w.log, "nothing to compensate")
	}
}

func (w *PaymentWorkflow) Log() []string { return w.log }

type Saga struct {
	steps []SagaStep
}

type SagaStep struct {
	Name       string
	Execute    func() error
	Compensate func()
}

func (s *Saga) AddStep(name string, exec func() error, comp func()) {
	s.steps = append(s.steps, SagaStep{Name: name, Execute: exec, Compensate: comp})
}

func (s *Saga) Run() error {
	completed := 0
	for i, step := range s.steps {
		if err := step.Execute(); err != nil {
			for j := i - 1; j >= 0; j-- {
				s.steps[j].Compensate()
			}
			return fmt.Errorf("step %q failed: %w", step.Name, err)
		}
		completed++
	}
	_ = completed
	return nil
}

func main() {
	w := NewPaymentWorkflow("pay-123", 5000)
	w.Reserve()
	w.Charge()
	fmt.Println("State:", w.State)
	fmt.Println("Log:")
	for _, l := range w.Log() {
		fmt.Println(" ", l)
	}

	w2 := NewPaymentWorkflow("pay-999", 5000)
	w2.Reserve()
	w2.Fail()
	w2.Compensate()
	fmt.Println("\nCompensation log:")
	for _, l := range w2.Log() {
		fmt.Println(" ", l)
	}

	saga := &Saga{}
	inv := 0
	saga.AddStep("reserve_inventory", func() error {
		inv = 100
		return nil
	}, func() { inv = 0 })
	saga.AddStep("process_payment", func() error {
		return errors.New("insufficient funds")
	}, func() {})
	saga.AddStep("confirm_order", func() error { return nil }, func() {})
	err := saga.Run()
	fmt.Println("\nSaga result:", err)
	fmt.Println("Inventory after saga:", inv)
}
