// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "testing"

// ============================================================================
// Tests for: Bank Account (Composition & Embedding)
// ============================================================================
//
// These tests verify Deposit, Withdraw, interest, and overdraft logic.
// They also prove that method shadowing works correctly on embedded structs.
//
// RUN: go test -v ./06-composition-and-embedding/3-bank-account
// ============================================================================

func TestDeposit(t *testing.T) {
	acc := &Account{AccountNumber: "TEST-001", Balance: 100.0, OwnerName: "Test"}

	err := acc.Deposit(50.0)
	if err != nil {
		t.Fatalf("Deposit(50) returned unexpected error: %v", err)
	}
	if acc.Balance != 150.0 {
		t.Errorf("Balance after deposit = %.2f, want 150.00", acc.Balance)
	}
}

func TestDepositNegativeAmount(t *testing.T) {
	acc := &Account{AccountNumber: "TEST-002", Balance: 100.0, OwnerName: "Test"}

	err := acc.Deposit(-10.0)
	if err == nil {
		t.Error("Deposit(-10) should have returned an error, got nil")
	}
}

func TestWithdraw(t *testing.T) {
	acc := &Account{AccountNumber: "TEST-003", Balance: 200.0, OwnerName: "Test"}

	err := acc.Withdraw(75.0)
	if err != nil {
		t.Fatalf("Withdraw(75) returned unexpected error: %v", err)
	}
	if acc.Balance != 125.0 {
		t.Errorf("Balance after withdraw = %.2f, want 125.00", acc.Balance)
	}
}

func TestWithdrawInsufficientFunds(t *testing.T) {
	acc := &Account{AccountNumber: "TEST-004", Balance: 50.0, OwnerName: "Test"}

	err := acc.Withdraw(100.0)
	if err == nil {
		t.Error("Withdraw(100) with balance 50 should have returned an error, got nil")
	}
	if acc.Balance != 50.0 {
		t.Errorf("Balance should remain 50.00 after failed withdraw, got %.2f", acc.Balance)
	}
}

func TestSavingsAccountInterest(t *testing.T) {
	sa := &SavingsAccount{
		Account:      Account{AccountNumber: "SAV-001", Balance: 1000.0, OwnerName: "Test"},
		InterestRate: 0.05,
	}

	sa.AddInterest()

	// Interest = 1000 * 0.05 = 50.  New balance = 1050.
	want := 1050.0
	if sa.Balance != want {
		t.Errorf("Balance after AddInterest = %.2f, want %.2f", sa.Balance, want)
	}
}

func TestOverdraftWithinLimit(t *testing.T) {
	oa := &OverdraftAccount{
		Account:        Account{AccountNumber: "OVR-001", Balance: 100.0, OwnerName: "Test"},
		OverdraftLimit: 200.0,
	}

	// Withdraw 250 — within limit (100 balance + 200 overdraft = 300 available)
	err := oa.Withdraw(250.0)
	if err != nil {
		t.Fatalf("Withdraw(250) within overdraft limit returned error: %v", err)
	}
	if oa.Balance != -150.0 {
		t.Errorf("Balance after overdraft withdraw = %.2f, want -150.00", oa.Balance)
	}
}

func TestOverdraftExceedsLimit(t *testing.T) {
	oa := &OverdraftAccount{
		Account:        Account{AccountNumber: "OVR-002", Balance: 100.0, OwnerName: "Test"},
		OverdraftLimit: 200.0,
	}

	// Withdraw 400 — exceeds limit (100 + 200 = 300 available)
	err := oa.Withdraw(400.0)
	if err == nil {
		t.Error("Withdraw(400) exceeding overdraft limit should have returned an error")
	}
	if oa.Balance != 100.0 {
		t.Errorf("Balance should remain 100.00 after failed overdraft, got %.2f", oa.Balance)
	}
}
