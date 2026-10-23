package main

import (
	"errors"
	"sync"
	"testing"
)

func TestCheckingAccount_DepositAndWithdrawal(t *testing.T) {
	acc, err := NewCheckingAccount("chk-1", "Alice", 10000, 5000)
	if err != nil {
		t.Fatalf("unexpected error creating account: %v", err)
	}

	// Normal Deposit
	if err := acc.Deposit(2500); err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	if acc.Balance() != 12500 {
		t.Fatalf("expected balance 12500, got %d", acc.Balance())
	}

	// Withdrawal within balance
	if err := acc.Withdraw(10000); err != nil {
		t.Fatalf("withdraw failed: %v", err)
	}
	if acc.Balance() != 2500 {
		t.Fatalf("expected balance 2500, got %d", acc.Balance())
	}

	// Withdrawal utilizing overdraft (overdraft limit is 5000)
	if err := acc.Withdraw(6000); err != nil {
		t.Fatalf("withdrawal within overdraft failed: %v", err)
	}
	if acc.Balance() != -3500 {
		t.Fatalf("expected negative balance -3500, got %d", acc.Balance())
	}

	// Withdrawal exceeding overdraft
	if err := acc.Withdraw(2000); !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("expected ErrInsufficientFunds, got %v", err)
	}

	// Negative amounts rejected
	if err := acc.Deposit(-500); !errors.Is(err, ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount on deposit, got %v", err)
	}
	if err := acc.Withdraw(-500); !errors.Is(err, ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount on withdraw, got %v", err)
	}
}

func TestSavingsAccount_MinimumBalanceAndInterest(t *testing.T) {
	acc, err := NewSavingsAccount("sav-1", "Bob", 20000, 5000)
	if err != nil {
		t.Fatalf("failed to create savings account: %v", err)
	}

	// Initial creation below min balance rejected
	_, errBad := NewSavingsAccount("sav-bad", "Bob", 4000, 5000)
	if !errors.Is(errBad, ErrBelowMinimumBalance) {
		t.Fatalf("expected ErrBelowMinimumBalance, got %v", errBad)
	}

	// Withdrawal leaving at least minimum balance
	if err := acc.Withdraw(10000); err != nil {
		t.Fatalf("valid withdrawal failed: %v", err)
	}
	if acc.Balance() != 10000 {
		t.Fatalf("expected balance 10000, got %d", acc.Balance())
	}

	// Withdrawal breaching minimum balance (10000 - 6000 = 4000 < 5000)
	if err := acc.Withdraw(6000); !errors.Is(err, ErrBelowMinimumBalance) {
		t.Fatalf("expected ErrBelowMinimumBalance, got %v", err)
	}

	// Interest accrual: 10% on 10000 = 1000
	if err := acc.AccrueInterest(10.0); err != nil {
		t.Fatalf("accrue interest failed: %v", err)
	}
	if acc.Balance() != 11000 {
		t.Fatalf("expected balance 11000 after interest, got %d", acc.Balance())
	}
}

func TestAccountPolymorphism_Transfer(t *testing.T) {
	checking, _ := NewCheckingAccount("chk-p1", "Charlie", 5000, 0)
	savings, _ := NewSavingsAccount("sav-p2", "Charlie", 15000, 2000)

	// Transfer from Savings to Checking
	if err := Transfer(savings, checking, 5000); err != nil {
		t.Fatalf("transfer failed: %v", err)
	}

	if savings.Balance() != 10000 {
		t.Fatalf("expected savings 10000, got %d", savings.Balance())
	}
	if checking.Balance() != 10000 {
		t.Fatalf("expected checking 10000, got %d", checking.Balance())
	}

	// Transfer failing due to minimum balance
	err := Transfer(savings, checking, 9000)
	if err == nil {
		t.Fatal("expected transfer failure when exceeding min balance, got nil")
	}
}

func TestConcurrentAccountAccess(t *testing.T) {
	acc, _ := NewCheckingAccount("chk-concurrent", "Dan", 100000, 0)

	const goroutines = 20
	const iterations = 50
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = acc.Deposit(10)
				_ = acc.Withdraw(10)
				_ = acc.Balance()
			}
		}(i)
	}

	wg.Wait()

	if acc.Balance() != 100000 {
		t.Fatalf("expected net zero balance change 100000, got %d", acc.Balance())
	}
}

func TestGenericTransactionBuffer(t *testing.T) {
	buf := NewTransactionBuffer[string](3)

	buf.Push("tx-1")
	buf.Push("tx-2")
	buf.Push("tx-3")

	if buf.Len() != 3 {
		t.Fatalf("expected 3 items, got %d", buf.Len())
	}

	// Eviction of oldest
	buf.Push("tx-4")
	items := buf.All()

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0] != "tx-2" || items[1] != "tx-3" || items[2] != "tx-4" {
		t.Fatalf("expected [tx-2, tx-3, tx-4], got %v", items)
	}
}

func TestNilInterfaceTrap(t *testing.T) {
	var concreteNil *CheckingAccount = nil
	var iface Account = concreteNil

	// In Go, an interface containing a typed nil pointer is NOT nil!
	if iface == nil {
		t.Fatal("expected typed nil pointer interface to be non-nil in Go")
	}
}
