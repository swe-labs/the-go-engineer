// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 14: Application Architecture - Package Design: Visibility and Export Rules
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Uppercase = exported (public), lowercase = unexported (private)
//   - The internal/ directory: compiler-enforced visibility boundaries
//   - Designing minimal public APIs
//   - Why Go doesn't have private/public/protected keywords
//
// ENGINEERING DEPTH:
//   Go keeps visibility simple. Instead of parsing access modifiers like many
//   other languages, the compiler can decide exported vs unexported by checking
//   the first letter of an identifier.
//
// RUN: go run ./14-application-architecture/package-design/2-visibility
// ============================================================================

// User is an exported type - other packages can use it.
// The struct fields follow the same rule:
//
//	Name  -> exported
//	email -> unexported
type User struct {
	Name  string
	email string
}

// NewUser is an exported constructor function.
func NewUser(name, email string) User {
	return User{
		Name:  name,
		email: sanitizeEmail(email),
	}
}

// Email is an exported accessor for the unexported email field.
func (u User) Email() string {
	return u.email
}

// sanitizeEmail is an unexported helper.
func sanitizeEmail(email string) string {
	return email
}

func main() {
	fmt.Println("=== Visibility and Export Rules ===")
	fmt.Println()

	user := NewUser("Alice", "alice@example.com")
	fmt.Printf("  Name (exported):        %s\n", user.Name)
	fmt.Printf("  Email (via getter):     %s\n", user.Email())
	fmt.Println()

	fmt.Println("=== Export Rules ===")
	fmt.Println("  Uppercase first letter -> exported (visible outside the package)")
	fmt.Println("  Lowercase first letter -> unexported (private to the package)")
	fmt.Println("  There are no public/private/protected keywords in Go")
	fmt.Println()

	fmt.Println("=== The internal/ Directory ===")
	fmt.Println("  myapp/internal/auth -> importable by myapp/* only")
	fmt.Println("  myapp/internal/db   -> not importable by other modules")
	fmt.Println("  This is compiler-enforced, not just convention")
	fmt.Println()

	fmt.Println("=== Design Principle: Export the Minimum ===")
	fmt.Println("  1. Start with everything unexported")
	fmt.Println("  2. Export only what external packages need")
	fmt.Println("  3. Use constructors (NewX) to validate on creation")
	fmt.Println("  4. Use getters for controlled access to internal state")
	fmt.Println("  5. Once exported, it is part of your API - removing it is a breaking change")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: PD.3 project layout")
	fmt.Println("   Current: PD.2 (visibility)")
	fmt.Println("---------------------------------------------------")
}
