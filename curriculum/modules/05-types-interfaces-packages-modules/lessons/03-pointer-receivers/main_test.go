package main

import "testing"

func TestCounterValue(t *testing.T) {
	c := Counter{}
	if c.Value() != 0 {
		t.Errorf("expected 0, got %d", c.Value())
	}
}

func TestCounterIncrement(t *testing.T) {
	c := Counter{}
	c.Increment()
	if c.Value() != 1 {
		t.Errorf("expected 1, got %d", c.Value())
	}
}

func TestCounterMultipleIncrements(t *testing.T) {
	c := Counter{}
	for i := 0; i < 5; i++ {
		c.Increment()
	}
	if c.Value() != 5 {
		t.Errorf("expected 5, got %d", c.Value())
	}
}

func TestCounterAdd(t *testing.T) {
	c := Counter{}
	c.Add(10)
	if c.Value() != 10 {
		t.Errorf("expected 10, got %d", c.Value())
	}
}

func TestCounterReset(t *testing.T) {
	c := Counter{}
	c.Add(100)
	c.Reset()
	if c.Value() != 0 {
		t.Errorf("expected 0 after reset, got %d", c.Value())
	}
}

func TestNilCounter(t *testing.T) {
	var c *Counter
	if c.Value() != 0 {
		t.Errorf("expected 0 from nil counter, got %d", c.Value())
	}
}

func TestBankAccountDeposit(t *testing.T) {
	a := BankAccount{}
	a.Deposit(50)
	if a.Balance() != 50 {
		t.Errorf("expected 50, got %.2f", a.Balance())
	}
}

func TestBankAccountWithdraw(t *testing.T) {
	a := BankAccount{}
	a.Deposit(100)
	err := a.Withdraw(40)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if a.Balance() != 60 {
		t.Errorf("expected 60, got %.2f", a.Balance())
	}
}

func TestBankAccountInsufficientFunds(t *testing.T) {
	a := BankAccount{}
	a.Deposit(10)
	err := a.Withdraw(20)
	if err == nil {
		t.Fatal("expected insufficient funds error")
	}
}

func TestBankAccountString(t *testing.T) {
	a := BankAccount{}
	a.Deposit(100.50)
	if a.Balance() != 100.50 {
		t.Errorf("expected 100.50, got %.2f", a.Balance())
	}
}
