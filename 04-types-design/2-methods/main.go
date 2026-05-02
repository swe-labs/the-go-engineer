// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Methods
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Associating functions with named types using method receivers.
//   - Differentiating between Value Receivers and Pointer Receivers.
//   - Understanding automatic pointer conversion (syntactic sugar).
//   - Applying the consistency rule for receiver design.
//
// WHY THIS MATTERS:
//   - Methods enable encapsulation and high cohesion by grouping logic
//     with the data it operates on. This is the foundation of
//     object-oriented patterns in Go and the mechanism used to satisfy
//     interfaces.
//
// RUN:
//   go run ./04-types-design/2-methods
//
// KEY TAKEAWAY:
//   - Methods define the behavioral interface of a data structure.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"fmt"
	"math"
)

// Section 04: Types & Design - Methods
//   - What methods are: functions attached to a type
//   - VALUE receivers vs POINTER receivers - the most critical distinction
//   - When to use each receiver type (the golden rule)
//   - Why methods exist: they enable interfaces
//   - Method sets and how they affect interface satisfaction
//

// Circle represents a circular geometry in 2D space.
type Circle struct {
	Radius float64
}

// Area calculates the area of the circle using its radius.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle.
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// String returns a human-readable representation of the circle.
func (c Circle) String() string {
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Scale resizes the circle by a given factor. It requires a pointer receiver to modify the original struct.
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
	fmt.Printf("  💰 Deposited $%.2f -> Balance: $%.2f\n", amount, a.Balance)
}

func (a *BankAccount) Withdraw(amount float64) bool {
	if amount > a.Balance {
		fmt.Printf("  ❌ Insufficient funds: need $%.2f, have $%.2f\n", amount, a.Balance)
		return false
	}
	a.Balance -= amount
	fmt.Printf("  💸 Withdrew $%.2f -> Balance: $%.2f\n", amount, a.Balance)
	return true
}

func (a BankAccount) Summary() string {
	return fmt.Sprintf("Account(%s, $%.2f)", a.Owner, a.Balance)
}

func main() {
	fmt.Println("=== Methods: Associating Behavior with Types ===")
	fmt.Println()

	// 1. Value receivers operate on a copy of the data.
	// Methods defined on 'T' can be called on both values and pointers.
	fmt.Println("--- Value Receiver Methods ---")
	c := Circle{Radius: 5.0}
	fmt.Printf("  %s\n", c)
	fmt.Printf("  Area:      %.2f\n", c.Area())
	fmt.Printf("  Perimeter: %.2f\n", c.Perimeter())
	fmt.Printf("  Radius:    %.2f (unchanged)\n", c.Radius)
	fmt.Println()

	// 2. Pointer receivers operate on the original memory address.
	// Go automatically takes the address of 'c' (&c) to satisfy the *Circle receiver.
	fmt.Println("--- Pointer Receiver Methods ---")
	c.Scale(2.0)
	fmt.Printf("  After Scale(2.0): radius = %.2f\n", c.Radius)
	fmt.Printf("  New area: %.2f\n", c.Area())
	fmt.Println()

	// 3. Encapsulation and state mutation.
	// Grouping logic (Deposit/Withdraw) with data (Balance) ensures consistency.
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
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.7 -> 04-types-design/7-receiver-sets")
	fmt.Println("Run    : go run ./04-types-design/7-receiver-sets")
	fmt.Println("Current: TI.2 (methods)")
	fmt.Println("---------------------------------------------------")
}
