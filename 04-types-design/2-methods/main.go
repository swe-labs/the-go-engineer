// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"fmt"
	"math"
)

// ============================================================================
// Section 6: Types & Interfaces — Methods
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What methods are: functions attached to a type
//   - VALUE receivers vs POINTER receivers — the most critical distinction
//   - When to use each receiver type (the golden rule)
//   - Why methods exist: they enable interfaces
//   - Method sets and how they affect interface satisfaction
//
// RUN: go run ./04-types-design/2-methods
// ============================================================================

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

func (c *Circle) Scale(factor float64) {
	c.Radius *= factor
}

type BankAccount struct {
	Owner   string
	Balance float64
}

func (a *BankAccount) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Println("  ❌ Deposit amount must be positive")
		return
	}
	a.Balance += amount
	fmt.Printf("  💰 Deposited $%.2f → Balance: $%.2f\n", amount, a.Balance)
}

func (a *BankAccount) Withdraw(amount float64) bool {
	if amount > a.Balance {
		fmt.Printf("  ❌ Insufficient funds: need $%.2f, have $%.2f\n", amount, a.Balance)
		return false
	}
	a.Balance -= amount
	fmt.Printf("  💸 Withdrew $%.2f → Balance: $%.2f\n", amount, a.Balance)
	return true
}

func (a BankAccount) Summary() string {
	return fmt.Sprintf("Account(%s, $%.2f)", a.Owner, a.Balance)
}

func main() {
	fmt.Println("=== Methods: Functions Attached to Types ===")
	fmt.Println()

	fmt.Println("--- Value Receiver Methods ---")
	c := Circle{Radius: 5.0}
	fmt.Printf("  %s\n", c)

	fmt.Printf("  Area:      %.2f\n", c.Area())
	fmt.Printf("  Perimeter: %.2f\n", c.Perimeter())
	fmt.Printf("  Radius:    %.2f (unchanged)\n", c.Radius)
	fmt.Println()

	fmt.Println("--- Pointer Receiver Methods ---")
	c.Scale(2.0)
	fmt.Printf("  After Scale(2.0): radius = %.2f\n", c.Radius)
	fmt.Printf("  New area: %.2f\n", c.Area())
	fmt.Println()

	fmt.Println("--- Bank Account (Pointer Receivers) ---")
	account := BankAccount{Owner: "Rasel", Balance: 1000.00}
	fmt.Printf("  %s\n", account.Summary())
	fmt.Println()

	account.Deposit(500.00)
	account.Withdraw(200.00)
	account.Withdraw(9999.00)
	fmt.Println()

	fmt.Printf("  Final: %s\n", account.Summary())

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Value receiver (c Circle):  works on a COPY — for read-only methods")
	fmt.Println("  - Pointer receiver (c *Circle): works on ORIGINAL — for mutation")
	fmt.Println("  - THE GOLDEN RULE: If ANY method needs a pointer, make ALL methods pointers")
	fmt.Println("  - Go auto-dereferences: c.Scale() works even if c is not a pointer")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.3 interfaces")
	fmt.Println("   Current: TI.2 (methods)")
	fmt.Println("---------------------------------------------------")
}
