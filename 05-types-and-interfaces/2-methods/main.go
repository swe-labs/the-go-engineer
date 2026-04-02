// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"math"
)

// ============================================================================
// Section 5: Types & Interfaces — Methods
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What methods are: functions attached to a type
//   - VALUE receivers vs POINTER receivers — the most critical distinction
//   - When to use each receiver type (the golden rule)
//   - Why methods exist: they enable interfaces (Section 05/3)
//   - Method sets and how they affect interface satisfaction
//
// ANALOGY:
//   Think of a TV remote. The remote (struct) has state: volume level,
//   current channel, power status. The BUTTONS on the remote are methods:
//     remote.VolumeUp()     ← modifies state (pointer receiver)
//     remote.CurrentChannel() ← reads state (value receiver)
//   Methods are functions that BELONG to a type, just like buttons belong
//   to a specific remote.
//
// RUN: go run ./05-types-and-interfaces/2-methods
// ============================================================================

// Circle represents a geometric circle with a given radius.
// We'll attach methods to this type to compute area and perimeter.
type Circle struct {
	Radius float64 // The distance from center to edge
}

// Area is a VALUE RECEIVER method.
// The receiver "(c Circle)" is a COPY of the original Circle.
// Changes to "c" inside this method do NOT affect the original.
//
// SYNTAX BREAKDOWN:
//
//	func (c Circle) Area() float64
//	^^^^  ^^^^^^^  ^^^^   ^^^^^^^
//	func  receiver  name   return type
//
// The receiver is like an implicit first parameter.
// c.Area() is equivalent to Area(c) — but reads better.
//
// USE VALUE RECEIVERS WHEN:
//   - The method only READS data (doesn't modify the struct)
//   - The struct is small (< 3-4 fields, cheap to copy)
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius // πr²
}

// Perimeter is another value receiver — it only reads c.Radius.
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius // 2πr
}

// String makes Circle implement the fmt.Stringer interface.
// When you pass a Circle to fmt.Println, Go calls this method automatically.
func (c Circle) String() string {
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Scale is a POINTER RECEIVER method.
// The receiver "(c *Circle)" is a POINTER to the original Circle.
// Changes to "c" inside this method MODIFY the original struct.
//
// SYNTAX: func (c *Circle) Scale(factor float64)
//
//	The * before Circle means "pointer to Circle"
//
// USE POINTER RECEIVERS WHEN:
//   - The method MODIFIES the struct (mutates state)
//   - The struct is large (expensive to copy)
//   - You need consistency (if ANY method uses pointer, ALL should)
func (c *Circle) Scale(factor float64) {
	c.Radius *= factor // Modifies the ORIGINAL circle
}

// BankAccount demonstrates why pointer receivers are essential for mutation.
// Without a pointer receiver, deposits would modify a COPY and be lost.
type BankAccount struct {
	Owner   string  // Account holder's name
	Balance float64 // Current balance in dollars
}

// Deposit adds money to the account.
// MUST be a pointer receiver — otherwise the balance change is lost!
func (a *BankAccount) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Println("  ❌ Deposit amount must be positive")
		return
	}
	a.Balance += amount // Modifies the ORIGINAL account
	fmt.Printf("  💰 Deposited $%.2f → Balance: $%.2f\n", amount, a.Balance)
}

// Withdraw removes money from the account.
// Returns an error-like bool to indicate insufficient funds.
func (a *BankAccount) Withdraw(amount float64) bool {
	if amount > a.Balance {
		fmt.Printf("  ❌ Insufficient funds: need $%.2f, have $%.2f\n", amount, a.Balance)
		return false
	}
	a.Balance -= amount // Modifies the ORIGINAL account
	fmt.Printf("  💸 Withdrew $%.2f → Balance: $%.2f\n", amount, a.Balance)
	return true
}

// Summary is a value receiver — it only reads data, doesn't change anything.
func (a BankAccount) Summary() string {
	return fmt.Sprintf("Account(%s, $%.2f)", a.Owner, a.Balance)
}

func main() {
	fmt.Println("=== Methods: Functions Attached to Types ===")
	fmt.Println()

	// --- VALUE RECEIVER METHODS (Read-Only) ---
	fmt.Println("--- Value Receiver Methods ---")
	c := Circle{Radius: 5.0}
	fmt.Printf("  %s\n", c) // Calls c.String() automatically

	// Call value receiver methods — these read but don't modify
	fmt.Printf("  Area:      %.2f\n", c.Area())      // πr² = 78.54
	fmt.Printf("  Perimeter: %.2f\n", c.Perimeter()) // 2πr = 31.42
	fmt.Printf("  Radius:    %.2f (unchanged)\n", c.Radius)
	fmt.Println()

	// --- POINTER RECEIVER METHODS (Mutation) ---
	fmt.Println("--- Pointer Receiver Methods ---")
	c.Scale(2.0) // Doubles the radius — modifies the ORIGINAL Circle

	// NOTE: Go is smart. Even though Scale takes *Circle, you can call
	// it on a Circle value. Go automatically takes the address: (&c).Scale(2.0)
	// This is syntactic sugar — you don't need to write (&c).Scale(2.0).
	fmt.Printf("  After Scale(2.0): radius = %.2f\n", c.Radius) // 10.0
	fmt.Printf("  New area: %.2f\n", c.Area())                  // 314.16
	fmt.Println()

	// --- REAL-WORLD EXAMPLE: Bank Account ---
	fmt.Println("--- Bank Account (Pointer Receivers) ---")
	account := BankAccount{Owner: "Rasel", Balance: 1000.00}
	fmt.Printf("  %s\n", account.Summary())
	fmt.Println()

	account.Deposit(500.00)   // Balance → $1500.00
	account.Withdraw(200.00)  // Balance → $1300.00
	account.Withdraw(9999.00) // Fails: insufficient funds
	fmt.Println()

	fmt.Printf("  Final: %s\n", account.Summary())

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Value receiver (c Circle):  works on a COPY — for read-only methods")
	fmt.Println("  - Pointer receiver (c *Circle): works on ORIGINAL — for mutation")
	fmt.Println("  - THE GOLDEN RULE: If ANY method needs a pointer, make ALL methods pointers")
	fmt.Println("  - Go auto-dereferences: c.Scale() works even if c is not a pointer")
}
