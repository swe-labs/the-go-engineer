// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 6: Composition and Embedding - Bank Account (Exercise Starter)
// Level: Intermediate
// ============================================================================
//
// EXERCISE: Build a bank account model with embedded types
//
// REQUIREMENTS:
//  1. [ ] Define an `Account` struct with AccountNumber, Balance, and OwnerName
//  2. [ ] Implement Deposit(amount float64) error to reject non-positive amounts
//  3. [ ] Implement Withdraw(amount float64) error to reject insufficient funds
//  4. [ ] Define `SavingsAccount` embedding `Account` plus an InterestRate field
//  5. [ ] Implement AddInterest() using promoted state and behavior from Account
//  6. [ ] Define `OverdraftAccount` embedding `Account` plus an OverdraftLimit field
//  7. [ ] Shadow Withdraw on OverdraftAccount to allow overdraft behavior
//  8. [ ] Demonstrate both account types clearly in main()
//  9. [ ] Make `go test ./05-composition/05-composition-and-embedding/3-bank-account` pass
//
// HINTS:
//   - Use pointer receivers (*Account) for methods that modify Balance
//   - Embedding looks like: type SavingsAccount struct { Account; InterestRate float64 }
//   - Promoted fields mean `sa.Balance` accesses `sa.Account.Balance`
//   - OverdraftAccount should reuse the embedded Account fields while changing Withdraw behavior
//
// RUN: go run ./05-composition/05-composition-and-embedding/3-bank-account/_starter
// TEST: go test ./05-composition/05-composition-and-embedding/3-bank-account
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define Account with Deposit, Withdraw, GetBalance, and String methods

// TODO: Define SavingsAccount with embedded Account and AddInterest

// TODO: Define OverdraftAccount with embedded Account and a shadowed Withdraw

func main() {
	fmt.Println("=== Bank Account Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement the bank account model described above.")
	fmt.Println("Use the tests to confirm your embedding and shadowing behavior.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
