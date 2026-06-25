package main

import "testing"

func TestBankDepositAndBalance(t *testing.T) {
	b := NewBank()
	b.Deposit("Alice", 100)
	bal, ok := b.Balance("Alice")
	if !ok || bal != 100 {
		t.Fatalf("Balance('Alice') = %d, %v; want 100, true", bal, ok)
	}
}

func TestBankWithdraw(t *testing.T) {
	b := NewBank()
	b.Deposit("Alice", 100)

	if ok := b.Withdraw("Alice", 30); !ok {
		t.Fatal("Withdraw(30) should succeed")
	}
	bal, _ := b.Balance("Alice")
	if bal != 70 {
		t.Fatalf("Balance after withdraw = %d; want 70", bal)
	}

	if ok := b.Withdraw("Alice", 200); ok {
		t.Fatal("Withdraw(200) should fail (insufficient funds)")
	}
}

func TestBankMissingAccount(t *testing.T) {
	b := NewBank()
	_, ok := b.Balance("Bob")
	if ok {
		t.Fatal("Balance('Bob') should return false for missing account")
	}
}

func TestNilBank(t *testing.T) {
	var b *Bank
	bal, ok := b.Balance("Alice")
	if ok || bal != 0 {
		t.Errorf("nil Bank Balance should return 0, false; got %d, %v", bal, ok)
	}
	if b.Withdraw("Alice", 10) {
		t.Error("nil Bank Withdraw should return false")
	}
}
