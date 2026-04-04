// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 6: Composition & Embedding — Bank Account (Exercise Starter)
// Level: Intermediate
// ============================================================================
//
// EXERCISE: Build a Bank Account System with Embedded Types
//
// REQUIREMENTS:
//  1. [ ] Define an `Account` struct with AccountNumber, Balance, OwnerName
//  2. [ ] Implement Deposit(amount float64) error — reject negative amounts
//  3. [ ] Implement Withdraw(amount float64) error — reject insufficient funds
//  4. [ ] Define `SavingsAccount` embedding Account + InterestRate field
//  5. [ ] Implement AddInterest() on SavingsAccount
//  6. [ ] Define `OverdraftAccount` embedding Account + OverdraftLimit field
//  7. [ ] Shadow the Withdraw method on OverdraftAccount to allow overdraft
//  8. [ ] Test all account types in main()
//
// HINTS:
//   - Use pointer receivers (*Account) for methods that modify Balance
//   - Embedding: type SavingsAccount struct { Account; InterestRate float64 }
//   - Embedded fields are "promoted" — sa.Balance accesses Account.Balance
//
// RUN: go run ./06-composition-and-embedding/3-bank-account/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define Account struct with Deposit, Withdraw, GetBalance, String methods

// TODO: Define SavingsAccount with embedded Account

// TODO: Define OverdraftAccount with shadowed Withdraw

func main() {
	fmt.Println("=== Bank Account Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your bank account system!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
