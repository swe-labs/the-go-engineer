// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Bank Account Project
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Combining named-field composition and embedding in a domain model.
//   - Promoting method sets from base types to specialized variants.
//   - Implementing method shadowing to override behavior (Overdraft vs Standard).
//   - Managing shared state through pointer receivers in composed hierarchies.
//
// WHY THIS MATTERS:
//   - In a production banking or financial system, you must balance
//     code reuse (base account logic) with strict specialized
//     behaviors (savings interest, overdraft limits). Go's composition
//     model allows you to "promote" common logic while explicitly
//     "shadowing" methods that require differentiation. This approach
//     avoids the "Fragile Base Class" problem seen in inheritance
//     languages, where a change in a parent class can break
//     unrelated specialized types.
//
// RUN:
//   go run ./04-types-design/18-bank-account-project
//
// KEY TAKEAWAY:
//   - Embedding allows specialization through selective method shadowing.
// ============================================================================

// Commercial use is prohibited without permission.

package main

//

import (
	"errors"
	"fmt"
)

// Account (Struct) encapsulates the base state and logic for any financial account.
// Account (Struct): (Struct) encapsulates the base state and logic for any financial account.
type Account struct {
	AccountNumber string
	Balance       float64
	OwnerName     string
}

// Deposit (Method) modifies account state by adding funds; uses a pointer receiver for persistence.
// Account.Deposit (Method): (Method) modifies account state by adding funds; uses a pointer receiver for persistence.
func (acc *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	acc.Balance += amount
	return nil
}

// Withdraw (Method) implements standard fund removal with strict balance checks.
// Account.Withdraw (Method): (Method) implements standard fund removal with strict balance checks.
func (acc *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if acc.Balance < amount {
		return fmt.Errorf("insufficient funds in %s", acc.AccountNumber)
	}
	acc.Balance -= amount
	return nil
}

// GetBalance (Method) provides read-only access to the internal balance state.
// Account.GetBalance (Method): (Method) provides read-only access to the internal balance state.
func (acc *Account) GetBalance() float64 {
	return acc.Balance
}

// String (Method) implements the fmt.Stringer interface for domain-specific summary output.
// Account.String (Method): (Method) implements the fmt.Stringer interface for domain-specific summary output.
func (acc *Account) String() string {
	return fmt.Sprintf("Account [%s] Owner: %s, Balance: $%.2f",
		acc.AccountNumber, acc.OwnerName, acc.Balance)
}

// SavingsAccount (Struct) specializes the Account type by embedding it to promote base financial behavior.
// SavingsAccount (Struct): (Struct) specializes the Account type by embedding it to promote base financial behavior.
type SavingsAccount struct {
	Account
	InterestRate float64
}

// AddInterest (Method) implements interest calculation logic and utilizes the promoted Deposit method.
// SavingsAccount.AddInterest (Method): (Method) implements interest calculation logic and utilizes the promoted Deposit method.
func (sa *SavingsAccount) AddInterest() {
	interest := sa.Balance * sa.InterestRate
	fmt.Printf("Processing interest: $%.2f (rate: %.1f%%)\n", interest, sa.InterestRate*100)
	_ = sa.Deposit(interest) // Using promoted Deposit method
}

// OverdraftAccount (Struct) specializes the Account type and provides custom withdrawal behavior.
// OverdraftAccount (Struct): (Struct) specializes the Account type and provides custom withdrawal behavior.
type OverdraftAccount struct {
	Account
	OverdraftLimit float64
}

// Withdraw (Method) shadows the base Account.Withdraw method to implement specialized overdraft logic.
// OverdraftAccount.Withdraw (Method): (Method) shadows the base Account.Withdraw method to implement specialized overdraft logic.
func (oa *OverdraftAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	// Specialized logic: allow balance to drop below zero up to the overdraft limit.
	if (oa.Balance + oa.OverdraftLimit) < amount {
		return fmt.Errorf("exceeds overdraft limit of $%.2f", oa.OverdraftLimit)
	}
	oa.Balance -= amount
	return nil
}

func main() {
	fmt.Println("=== Case Study: Financial Type Design ===")
	fmt.Println()

	// 1. Specialized Embedding (Savings).
	// The SavingsAccount embeds Account, promoting Deposit() and Balance access.
	fmt.Println("--- Savings Account Flow ---")
	savAcc := SavingsAccount{
		Account: Account{
			AccountNumber: "SAV-001",
			Balance:       1000.00,
			OwnerName:     "Alice Saver",
		},
		InterestRate: 0.02,
	}
	savAcc.AddInterest()
	fmt.Printf("  Final Savings State: %s\n", savAcc.String())
	fmt.Println()

	// 2. Method Shadowing (Overdraft).
	// The OverdraftAccount shadows the standard Withdraw() method to
	// allow negative balances up to a specific limit.
	fmt.Println("--- Overdraft Account Flow ---")
	ovdAcc := OverdraftAccount{
		Account: Account{
			AccountNumber: "OVD-002",
			Balance:       100.00,
			OwnerName:     "Bob Spender",
		},
		OverdraftLimit: 200.00,
	}

	// This call dispatches to OverdraftAccount.Withdraw, not Account.Withdraw.
	if err := ovdAcc.Withdraw(250.00); err != nil {
		fmt.Printf("  Withdraw Failed: %v\n", err)
	} else {
		fmt.Printf("  Withdraw Success. New Balance: $%.2f\n", ovdAcc.Balance)
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ST.1 -> 04-types-design/19-strings")
	fmt.Println("Run    : go run ./04-types-design/19-strings")
	fmt.Println("Current: CO.3 (bank-account-project)")
	fmt.Println("---------------------------------------------------")
}
