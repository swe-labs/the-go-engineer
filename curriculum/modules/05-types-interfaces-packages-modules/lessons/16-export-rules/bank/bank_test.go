package bank

import "testing"

func TestNewAccount(t *testing.T) {
	a := NewAccount("Bob", 50.0)
	if a.Balance() != 50.0 {
		t.Errorf("expected balance 50.0, got %.2f", a.Balance())
	}
}

func TestDeposit(t *testing.T) {
	a := NewAccount("Bob", 0)
	if err := a.Deposit(100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Balance() != 100.0 {
		t.Errorf("expected balance 100.0, got %.2f", a.Balance())
	}
}

func TestDepositNegative(t *testing.T) {
	a := NewAccount("Bob", 0)
	if err := a.Deposit(-10); err == nil {
		t.Error("expected error for negative deposit")
	}
}

func TestWithdraw(t *testing.T) {
	a := NewAccount("Bob", 100)
	if err := a.Withdraw(40); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Balance() != 60.0 {
		t.Errorf("expected balance 60.0, got %.2f", a.Balance())
	}
}

func TestWithdrawInsufficient(t *testing.T) {
	a := NewAccount("Bob", 10)
	if err := a.Withdraw(20); err == nil {
		t.Error("expected error for insufficient balance")
	}
}

func TestWithdrawNegative(t *testing.T) {
	a := NewAccount("Bob", 100)
	if err := a.Withdraw(-5); err == nil {
		t.Error("expected error for negative withdrawal")
	}
}
