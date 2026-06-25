package main

import (
	"testing"
)

func TestNewAccount_Valid(t *testing.T) {
	acc, err := NewAccount("acc-1", "Alice", 500.00)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acc.Balance != 500.00 {
		t.Errorf("expected 500.00, got %.2f", acc.Balance)
	}
}

func TestNewAccount_EmptyID(t *testing.T) {
	_, err := NewAccount("", "Alice", 100)
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestNewAccount_NegativeInitialBalance(t *testing.T) {
	_, err := NewAccount("acc-1", "Alice", -10)
	if err == nil {
		t.Fatal("expected error for negative initial balance")
	}
}

func TestAccount_Deposit(t *testing.T) {
	acc, _ := NewAccount("acc-1", "Alice", 100)
	err := acc.Deposit(50)
	if err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}
	if acc.Balance != 150 {
		t.Errorf("expected 150, got %.2f", acc.Balance)
	}
}

func TestAccount_Deposit_Negative(t *testing.T) {
	acc, _ := NewAccount("acc-1", "Alice", 100)
	err := acc.Deposit(-50)
	if err == nil {
		t.Fatal("expected error for negative deposit")
	}
}

func TestAccount_Withdraw(t *testing.T) {
	acc, _ := NewAccount("acc-1", "Alice", 100)
	err := acc.Withdraw(40)
	if err != nil {
		t.Fatalf("Withdraw failed: %v", err)
	}
	if acc.Balance != 60 {
		t.Errorf("expected 60, got %.2f", acc.Balance)
	}
}

func TestAccount_Withdraw_Insufficient(t *testing.T) {
	acc, _ := NewAccount("acc-1", "Alice", 100)
	err := acc.Withdraw(200)
	if err == nil {
		t.Fatal("expected error for insufficient balance")
	}
	if acc.Balance != 100 {
		t.Errorf("balance should remain unchanged: %.2f", acc.Balance)
	}
}

func TestTransfer_Success(t *testing.T) {
	alice, _ := NewAccount("acc-1", "Alice", 500)
	bob, _ := NewAccount("acc-2", "Bob", 100)

	err := Transfer(alice, bob, 200)
	if err != nil {
		t.Fatalf("Transfer failed: %v", err)
	}
	if alice.Balance != 300 {
		t.Errorf("expected Alice balance 300, got %.2f", alice.Balance)
	}
	if bob.Balance != 300 {
		t.Errorf("expected Bob balance 300, got %.2f", bob.Balance)
	}
}

func TestTransfer_NegativeAmount(t *testing.T) {
	alice, _ := NewAccount("acc-1", "Alice", 500)
	bob, _ := NewAccount("acc-2", "Bob", 100)
	err := Transfer(alice, bob, -50)
	if err == nil {
		t.Fatal("expected error for negative transfer")
	}
}

func TestTransfer_SameAccount(t *testing.T) {
	acc, _ := NewAccount("acc-1", "Alice", 500)
	err := Transfer(acc, acc, 100)
	if err == nil {
		t.Fatal("expected error for self-transfer")
	}
}

func TestTransfer_InsufficientBalance(t *testing.T) {
	alice, _ := NewAccount("acc-1", "Alice", 50)
	bob, _ := NewAccount("acc-2", "Bob", 100)
	err := Transfer(alice, bob, 200)
	if err == nil {
		t.Fatal("expected error for insufficient balance")
	}
}
