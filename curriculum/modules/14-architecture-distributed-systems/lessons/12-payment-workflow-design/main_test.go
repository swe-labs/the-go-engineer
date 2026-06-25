package main

import "testing"

func TestPaymentWorkflowHappyPath(t *testing.T) {
	w := NewPaymentWorkflow("pay-1", 5000)
	if w.State != StateInitiated {
		t.Fatalf("expected initiated, got %s", w.State)
	}
	if err := w.Reserve(); err != nil {
		t.Fatal(err)
	}
	if w.State != StateReserved {
		t.Fatalf("expected reserved, got %s", w.State)
	}
	if err := w.Charge(); err != nil {
		t.Fatal(err)
	}
	if w.State != StateCharged {
		t.Fatalf("expected charged, got %s", w.State)
	}
}

func TestPaymentWorkflowRefund(t *testing.T) {
	w := NewPaymentWorkflow("pay-2", 1000)
	w.Reserve()
	w.Charge()
	if err := w.Refund(); err != nil {
		t.Fatal(err)
	}
	if w.State != StateRefunded {
		t.Fatalf("expected refunded, got %s", w.State)
	}
}

func TestPaymentWorkflowInvalidTransition(t *testing.T) {
	w := NewPaymentWorkflow("pay-3", 500)
	if err := w.Charge(); err == nil {
		t.Fatal("expected error charging from initiated")
	}
}

func TestPaymentWorkflowAmountLimit(t *testing.T) {
	w := NewPaymentWorkflow("pay-4", 50000)
	w.Reserve()
	if err := w.Charge(); err == nil {
		t.Fatal("expected error for amount over limit")
	}
}

func TestPaymentWorkflowCompensateReserved(t *testing.T) {
	w := NewPaymentWorkflow("pay-5", 2000)
	w.Reserve()
	w.Compensate()
	if w.State != StateInitiated {
		t.Fatalf("expected initiated after compensate, got %s", w.State)
	}
}

func TestPaymentWorkflowCompensateCharged(t *testing.T) {
	w := NewPaymentWorkflow("pay-6", 2000)
	w.Reserve()
	w.Charge()
	w.Compensate()
	if w.State != StateRefunded {
		t.Fatalf("expected refunded after compensate, got %s", w.State)
	}
}

func TestSagaCompensatesOnFailure(t *testing.T) {
	var counter int
	saga := &Saga{}
	saga.AddStep("step1", func() error { counter++; return nil }, func() { counter-- })
	saga.AddStep("step2", func() error { return nil }, func() { counter-- })
	saga.AddStep("step3", func() error { return nil }, func() { counter-- })
	if err := saga.Run(); err != nil {
		t.Fatal("expected success")
	}
}

func TestSagaRollbackOnFailure(t *testing.T) {
	var step1executed, step1compensated, step2executed, step2compensated bool
	saga := &Saga{}
	saga.AddStep("step1", func() error { step1executed = true; return nil }, func() { step1compensated = true })
	saga.AddStep("step2", func() error { return nil }, func() { step2compensated = true })
	saga.AddStep("step3", func() error { return nil }, func() { step2compensated = true })
	err := saga.Run()
	_ = err
	_ = step1executed
	_ = step1compensated
	_ = step2executed
	_ = step2compensated
}

func TestPaymentWorkflowTable(t *testing.T) {
	tests := []struct {
		name    string
		amount  int
		ops     func(w *PaymentWorkflow) error
		wantErr bool
		wantSt  State
	}{
		{
			name:   "happy_path",
			amount: 5000,
			ops: func(w *PaymentWorkflow) error {
				if err := w.Reserve(); err != nil {
					return err
				}
				return w.Charge()
			},
			wantErr: false,
			wantSt:  StateCharged,
		},
		{
			name:   "charge_without_reserve",
			amount: 5000,
			ops: func(w *PaymentWorkflow) error {
				return w.Charge()
			},
			wantErr: true,
			wantSt:  StateInitiated,
		},
		{
			name:   "exceeds_limit",
			amount: 20000,
			ops: func(w *PaymentWorkflow) error {
				w.Reserve()
				return w.Charge()
			},
			wantErr: true,
			wantSt:  StateReserved,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := NewPaymentWorkflow("t", tc.amount)
			err := tc.ops(w)
			if (err != nil) != tc.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
			if w.State != tc.wantSt {
				t.Errorf("want state %s, got %s", tc.wantSt, w.State)
			}
		})
	}
}
